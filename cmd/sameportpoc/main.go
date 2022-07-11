package main

// Serving http and https on the same port?
// https://go.dev/play/p/5M2V9GeTZ-
// https://groups.google.com/g/golang-nuts/c/4oZp1csAm2o
import (
	"crypto/tls"
	"github.com/bingoohuang/fproxy"
	"log"
	"net/http"
	"os"
)

func main() {
	cert, err := tls.LoadX509KeyPair("root.pem", "root.key")
	if err != nil {
		log.Fatal(err)
	}
	l, err := fproxy.NewHttpsListener(":8080", cert)
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		ErrorLog: log.New(os.Stderr, "", log.LstdFlags),
	}
	srv.Serve(l)
}
