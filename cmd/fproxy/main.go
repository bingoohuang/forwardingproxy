// Copyright (C) 2018 Betalo AB - All Rights Reserved

package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/bingoohuang/fproxy"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/acme/autocert"
)

var (
	pCaPath                  = flag.String("ca", "", "Filepath to certificate and private key, like -ca cert.pem,key.pem")
	pAddr                    = flag.String("addr", ":7777", "Server address, use :0 for random listening port (check it in the log)")
	pAuth                    = flag.String("auth", "", "Server authentication username:password")
	pAvoid                   = flag.String("avoid", "", "Site to be avoided")
	pLog                     = flag.String("log", "info", "Log level")
	pDestDialTimeout         = flag.Duration("dest.dial.timeout", 10*time.Second, "Destination dial timeout")
	pDestReadTimeout         = flag.Duration("dest.read.timeout", 5*time.Second, "Destination read timeout")
	pDestWriteTimeout        = flag.Duration("dest.write.timeout", 5*time.Second, "Destination write timeout")
	pClientReadTimeout       = flag.Duration("client.read.timeout", 5*time.Second, "Client read timeout")
	pClientWriteTimeout      = flag.Duration("client.write.timeout", 5*time.Second, "Client write timeout")
	pServerReadTimeout       = flag.Duration("server.read.timeout", 30*time.Second, "Server read timeout")
	pServerReadHeaderTimeout = flag.Duration("server.read.header.timeout", 30*time.Second, "Server read header timeout")
	pServerWriteTimeout      = flag.Duration("server.write.timeout", 30*time.Second, "Server write timeout")
	pServerIdleTimeout       = flag.Duration("server.idle.timeout", 30*time.Second, "Server idle timeout")
	pLetsEncrypt             = flag.Bool("le", false, "Use letsencrypt for https")
	pLEWhitelist             = flag.String("le.whitelist", "f.cn", "Hostname to whitelist for letsencrypt")
	pLECacheDir              = flag.String("le.cache.dir", "", "Cache directory for certificates")
)

func main() {
	flag.Usage = func() {
		fmt.Print(`Usage of fproxy:
  -addr                       string  Server address (default ":0")
  -auth                       string  Server authentication username:password
  -avoid                      string Site to be avoided
  -log                        string   Log level (default "info")
  -ca                         string   Filepath to certificate and private key, like -ca cert.pem,key.pem (This will enable https proxy)
  -le                                  Use letsencrypt for https
  -le.cache.dir               string   Cache directory for certificates
  -le.whitelist               string   Hostname to whitelist for letsencrypt (default "localhost")
  -server.idle.timeout        duration Server idle timeout (default 30s)
  -server.read.header.timeout duration Server read header timeout (default 30s)
  -server.read.timeout        duration Server read timeout (default 30s)
  -server.write.timeout       duration Server write timeout (default 30s)
  -client.read.timeout        duration Client read timeout (default 5s)
  -client.write.timeout       duration Client write timeout (default 5s)
  -dest.dial.timeout          duration Destination dial timeout (default 10s)
  -dest.read.timeout          duration Destination read timeout (default 5s)
  -dest.write.timeout         duration Destination write timeout (default 5s)
`)
	}
	flag.Parse()

	c := zap.NewProductionConfig()
	c.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	level, err := zapcore.ParseLevel(*pLog)
	if err != nil {
		log.Fatalf("Error: failed to parse log level: %v", err)
	}
	c.Level.SetLevel(level)

	logger, err := c.Build()
	if err != nil {
		log.Fatalln("Error: failed to initiate logger")
	}
	defer logger.Sync()
	stdLogger := zap.NewStdLog(logger)

	p := &fproxy.Proxy{
		Forwarding:         fproxy.NewForwardingHTTPProxy(stdLogger),
		Logger:             logger,
		Auth:               *pAuth,
		DestDialTimeout:    *pDestDialTimeout,
		DestReadTimeout:    *pDestReadTimeout,
		DestWriteTimeout:   *pDestWriteTimeout,
		ClientReadTimeout:  *pClientReadTimeout,
		ClientWriteTimeout: *pClientWriteTimeout,
		Avoid:              *pAvoid,
	}

	s := &http.Server{
		Addr:              *pAddr,
		Handler:           p,
		ErrorLog:          stdLogger,
		ReadTimeout:       *pServerReadTimeout,
		ReadHeaderTimeout: *pServerReadHeaderTimeout,
		WriteTimeout:      *pServerWriteTimeout,
		IdleTimeout:       *pServerIdleTimeout,
		TLSNextProto:      map[string]func(*http.Server, *tls.Conn, http.Handler){}, // Disable HTTP/2
	}

	if *pLetsEncrypt {
		if *pLEWhitelist == "" {
			p.Logger.Fatal("error: no -le.whitelist flag set")
		}
		if *pLECacheDir == "" {
			*pLECacheDir, err = ioutil.TempDir("/tmp", "letsencrypt")
			p.Logger.Info("Cache temp directory for certificates", zap.String("letsEncryptCacheDir", *pLECacheDir))
		}

		m := &autocert.Manager{
			Cache:      autocert.DirCache(*pLECacheDir),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(*pLEWhitelist),
		}

		s.TLSConfig = m.TLSConfig()
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		p.Logger.Info("Server shutting down")
		if err = s.Shutdown(context.Background()); err != nil {
			p.Logger.Error("Server shutdown failed", zap.Error(err))
		}
		close(idleConnsClosed)
	}()

	var svrErr error
	if *pCaPath != "" {
		ps := strings.SplitN(*pCaPath, ",", 2)
		if len(ps) != 2 {
			log.Fatalf("invalid flags, e.g. -ca cert.pem,key.pem")
		}
		svrErr = listenAndServeTLS(s, ps[0], ps[1], p.Logger)
	} else if *pLetsEncrypt {
		svrErr = listenAndServeTLS(s, "", "", p.Logger)
	} else {
		svrErr = listenAndServe(s, p.Logger)
	}

	if svrErr != http.ErrServerClosed {
		p.Logger.Error("Listening for incoming connections failed", zap.Error(svrErr))
	}

	<-idleConnsClosed
	p.Logger.Info("Server stopped")
}

func listenAndServeTLS(srv *http.Server, certFile, keyFile string, logger *zap.Logger) error {
	l, err := fproxy.CreateTLSListener(srv.Addr, certFile, keyFile)
	if err != nil {
		return err
	}

	defer l.Close()

	logger.Info("http/https server starting", zap.String("Listening", l.Addr().String()))
	_, port, _ := net.SplitHostPort(l.Addr().String())
	logger.Info(fmt.Sprintf("settings: export http_proxy=http://127.0.0.1:%s; export https_proxy=http://127.0.0.1:%s", port, port))
	logger.Info(fmt.Sprintf("or      : export http_proxy=https://127.0.0.1:%s; export https_proxy=https://127.0.0.1:%s", port, port))
	return srv.Serve(l)
}

func listenAndServe(srv *http.Server, logger *zap.Logger) error {
	l, err := fproxy.CreateListener(srv.Addr)
	if err != nil {
		return err
	}

	defer l.Close()

	logger.Info("http server starting", zap.String("Listening", l.Addr().String()))
	_, port, _ := net.SplitHostPort(l.Addr().String())
	logger.Info(fmt.Sprintf("settings: export http_proxy=http://127.0.0.1:%s; export https_proxy=http://127.0.0.1:%s", port, port))

	return srv.Serve(l)
}
