package fproxy

import (
	"bufio"
	"crypto/tls"
	"net"
)

func NewHttpsListener(addr string, cert tls.Certificate) (net.Listener, error) {
	c := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	if len(c.NextProtos) == 0 {
		c.NextProtos = []string{"http/1.1"}
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &splitListener{Listener: ln, config: c}, nil
}

type splitListener struct {
	net.Listener
	config   *tls.Config
	Protocol string
}

func (l *splitListener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	// buffer reads on our conn
	bc := &conn{Conn: c, buf: bufio.NewReader(c)}

	// inspect the first few bytes
	hdr, err := bc.buf.Peek(1)
	if err != nil {
		_ = bc.Close()
		return nil, err
	}

	// iptables -t nat -A PREROUTING -i eth0 -p tcp --dport 80 -j REDIRECT --to-port 8080
	// iptables -t nat -A PREROUTING -i eth0 -p tcp --dport 443 -j REDIRECT --to-port 3000
	// https://github.com/mscdex/httpolyglot/blob/master/lib/index.js

	// TLS and HTTP connections are easy to distinguish based on the first byte sent by clients trying to connect.
	// https://github.com/mscdex/httpolyglot/issues/3#issuecomment-173680155
	// Alright, I just read through RFC 5245 - The Transport Layer Security (TLS) Protocol - Version 1.2 and
	// RFC 7230 - Hypertext Transfer Protocol (HTTP/1.1): Message Syntax and Routing which outline the structure
	// of TLS and HTTP, and therefore explain why this works.
	//
	// It boils down to two facts, the first byte of a TLS message is always 22 and
	// the first byte of an HTTP message will always be greater than 32 and less than (not equal to) 127.
	// This means that the condition in the code is asserting that any message whose first byte is
	// outside the range of a valid HTTP message must be a TLS message. And that totally works! 💃
	if firstByte := hdr[0]; firstByte < 32 || firstByte >= 127 { // tls/ssl
		l.Protocol = "https"
		return tls.Server(bc, l.config), nil
	}

	l.Protocol = "http"
	return bc, nil
}

// conn is a buffered conn for peeking into the connection
type conn struct {
	net.Conn
	buf *bufio.Reader
}

func (c *conn) Read(b []byte) (int, error) { return c.buf.Read(b) }
