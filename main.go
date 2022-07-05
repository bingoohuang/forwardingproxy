// Copyright (C) 2018 Betalo AB - All Rights Reserved

package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/acme/autocert"
)

var (
	pCertPath = flag.String("cert", "", "Filepath to certificate")
	pKeyPath  = flag.String("key", "", "Filepath to private key")
	pAddr     = flag.String("addr", ":0", "Server address")
	pAuth     = flag.String("auth", "", "Server authentication username:password")
	pAvoid    = flag.String("avoid", "", "Site to be avoided")
	pLog      = flag.String("log", "info", "Log level")

	pDestDialTimeout         = flag.Duration("dest.dial.timeout", 10*time.Second, "Destination dial timeout")
	pDestReadTimeout         = flag.Duration("dest.read.timeout", 5*time.Second, "Destination read timeout")
	pDestWriteTimeout        = flag.Duration("dest.write.timeout", 5*time.Second, "Destination write timeout")
	pClientReadTimeout       = flag.Duration("client.read.timeout", 5*time.Second, "Client read timeout")
	pClientWriteTimeout      = flag.Duration("client.write.timeout", 5*time.Second, "Client write timeout")
	pServerReadTimeout       = flag.Duration("server.read.timeout", 30*time.Second, "Server read timeout")
	pServerReadHeaderTimeout = flag.Duration("server.read.header.timeout", 30*time.Second, "Server read header timeout")
	pServerWriteTimeout      = flag.Duration("server.write.timeout", 30*time.Second, "Server write timeout")
	pServerIdleTimeout       = flag.Duration("server.idle.timeout", 30*time.Second, "Server idle timeout")

	pLetsEncrypt = flag.Bool("le", false, "Use letsencrypt for https")
	pLEWhitelist = flag.String("le.whitelist", "", "Hostname to whitelist for letsencrypt")
	pLECacheDir  = flag.String("le.cache.dir", "/tmp", "Cache directory for certificates")
)

func main() {
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

	p := &Proxy{
		Forwarding:         NewForwardingHTTPProxy(stdLogger),
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
		if *pLECacheDir == "/tmp" {
			p.Logger.Info("-le.cache.dir should be set, using '/tmp' for now...")
		}

		m := &autocert.Manager{
			Cache:      autocert.DirCache(*pLECacheDir),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(*pLEWhitelist),
		}

		s.Addr = ":https"
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
	if *pCertPath != "" && *pKeyPath != "" || *pLetsEncrypt {
		svrErr = ListenAndServeTLS(s, *pCertPath, *pKeyPath, p.Logger)
	} else {
		svrErr = ListenAndServe(s, p.Logger)
	}

	if svrErr != http.ErrServerClosed {
		p.Logger.Error("Listening for incoming connections failed", zap.Error(svrErr))
	}

	<-idleConnsClosed
	p.Logger.Info("Server stopped")
}

func ListenAndServeTLS(srv *http.Server, certFile, keyFile string, logger *zap.Logger) error {
	addr := srv.Addr
	if addr == "" {
		addr = ":https"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	logger.Info("Server starting", zap.String("Listening", ln.Addr().String()))

	defer ln.Close()

	return srv.ServeTLS(ln, certFile, keyFile)
}

func ListenAndServe(srv *http.Server, logger *zap.Logger) error {
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	logger.Info("Server starting", zap.String("Listening", ln.Addr().String()))

	return srv.Serve(ln)
}
