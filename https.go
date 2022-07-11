package fproxy

import (
	"bufio"
	"bytes"
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
	hdr, err := bc.buf.Peek(4)
	if err != nil {
		_ = bc.Close()
		return nil, err
	}

	// I don't remember what the TLS handshake looks like, but this works as a POC
	if bytes.Equal(hdr, []byte{22, 3, 1, 0}) {
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
