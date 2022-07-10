package main

// Serving http and https on the same port?
// https://go.dev/play/p/5M2V9GeTZ-
// https://groups.google.com/g/golang-nuts/c/4oZp1csAm2o
import (
	"bufio"
	"bytes"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"
)

type SplitListener struct {
	net.Listener
	config   *tls.Config
	Protocol string
}

func (l *SplitListener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	// buffer reads on our conn
	bc := &conn{Conn: c, buf: bufio.NewReader(c)}

	// inspect the first few bytes
	hdr, err := bc.buf.Peek(4)
	if err != nil {
		bc.Close()
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

func main() {
	addr := ":8080"

	c := &tls.Config{}
	if len(c.NextProtos) == 0 {
		c.NextProtos = []string{"http/1.1"}
	}

	var err error
	c.Certificates = make([]tls.Certificate, 1)
	c.Certificates[0], err = tls.LoadX509KeyPair("cert.pem", "cert.key")
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	var l net.Listener = &SplitListener{Listener: ln, config: c}
	srv := &http.Server{
		ErrorLog: log.New(os.Stderr, "", log.LstdFlags),
	}
	srv.Serve(l)
}
