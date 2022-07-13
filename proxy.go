// Copyright (C) 2018 Betalo AB - All Rights Reserved

// Courtesy: https://medium.com/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c

// $ openssl req -newkey rsa:2048 -nodes -keyout server.key -new -x509 -sha256 -days 3650 -out server.pem

package fproxy

import (
	"encoding/base64"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Proxy is a HTTPS forward proxy.
type Proxy struct {
	Logger     *zap.Logger
	Auth       string
	Avoid      string
	Forwarding *httputil.ReverseProxy

	DestDialTimeout    time.Duration
	DestReadTimeout    time.Duration
	DestWriteTimeout   time.Duration
	ClientReadTimeout  time.Duration
	ClientWriteTimeout time.Duration
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.Logger.Info("Incoming request", zap.String("host", r.Host))

	if p.Auth != "" {
		pa := r.Header.Get("Proxy-Authorization")
		if auth, ok := parseBasicProxyAuth(pa); !ok || auth != p.Auth {
			p.Logger.Warn("Authorization attempt with invalid credentials")
			http.Error(w, http.StatusText(http.StatusProxyAuthRequired), http.StatusProxyAuthRequired)
			return
		}
	}

	if r.URL.Scheme == "http" {
		p.handleHTTP(w, r)
	} else {
		p.handleTunneling(w, r)
	}
}

func (p *Proxy) handleHTTP(w http.ResponseWriter, r *http.Request) {
	p.Logger.Debug("Got HTTP request", zap.String("host", r.Host))
	if p.Avoid != "" && strings.Contains(r.Host, p.Avoid) == true {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusMethodNotAllowed)
		return
	}
	p.Forwarding.ServeHTTP(w, r)
}

func (p *Proxy) handleTunneling(w http.ResponseWriter, r *http.Request) {
	if p.Avoid != "" && strings.Contains(r.Host, p.Avoid) == true {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusMethodNotAllowed)
		return
	}

	if r.Method != http.MethodConnect {
		p.Logger.Info("Method not allowed", zap.String("method", r.Method))
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	p.Logger.Debug("Connecting", zap.String("host", r.Host))

	destConn, err := net.DialTimeout("tcp", r.Host, p.DestDialTimeout)
	if err != nil {
		p.Logger.Error("Destination dial failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	p.Logger.Debug("Connected", zap.String("host", r.Host))
	p.Logger.Debug("Hijacking", zap.String("host", r.Host))

	// Hijacker interface allows to take over the connection. After that the caller is
	// responsible to manage such connection (HTTP library won’t do it anymore).
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		p.Logger.Error("Hijacking not supported")
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		p.Logger.Error("Hijacking failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	clientConn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))

	p.Logger.Debug("Hijacked connection", zap.String("host", r.Host))

	now := time.Now()
	clientConn.SetReadDeadline(now.Add(p.ClientReadTimeout))
	clientConn.SetWriteDeadline(now.Add(p.ClientWriteTimeout))
	destConn.SetReadDeadline(now.Add(p.DestReadTimeout))
	destConn.SetWriteDeadline(now.Add(p.DestWriteTimeout))

	d := netDirection(destConn, clientConn)
	go transfer(destConn, newDebugReadCloser(clientConn, ">>> "+d, p.Logger))
	go transfer(clientConn, newDebugReadCloser(destConn, "<<< "+d, p.Logger))
}

func transfer(dest io.WriteCloser, src io.ReadCloser) {
	defer func() {
		_ = dest.Close()
		_ = src.Close()
	}()
	_, _ = io.Copy(dest, src)
}

// parseBasicProxyAuth parses an HTTP Basic Authorization string.
// "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("Aladdin:open sesame", true).
func parseBasicProxyAuth(authz string) (auth string, ok bool) {
	const prefix = "Basic "
	if !strings.HasPrefix(authz, prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(authz[len(prefix):])
	if err != nil {
		return
	}
	return string(c), true
}

// NewForwardingHTTPProxy returns a new reverse proxy that takes an incoming
// request and sends it to another server, proxying the response back to the
// client.
//
// See: https://golang.org/pkg/net/http/httputil/#ReverseProxy
func NewForwardingHTTPProxy(logger *log.Logger) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	// TODO:(alesr) Use timeouts specified via flags to customize the default
	// transport used by the reverse proxy.
	return &httputil.ReverseProxy{
		ErrorLog: logger,
		Director: director,
	}
}
