package fproxy

import (
	"crypto/tls"
	"net"
)

func CreateTLSListener(addr, certFile, keyFile string) (net.Listener, error) {
	if addr == "" {
		addr = ":https"
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	return NewHttpsListener(addr, cert)
}

func CreateListener(addr string) (net.Listener, error) {
	if addr == "" {
		addr = ":http"
	}
	return net.Listen("tcp", addr)
}
