package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"strconv"

	"github.com/clementine/codesigner"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

var port = flag.Int("port", 5001, "Port to start GRPC server on")
var cert = flag.String("cert", "", "Path to TLS certificate")
var key = flag.String("key", "", "Path to TLS key")
var ca = flag.String("ca", "", "Path to CA certificate")

func main() {
	flag.Parse()

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(*ca)
	if err != nil {
		log.Fatalf("Failed to read CA cert: %v", err)
	}
	ok := certPool.AppendCertsFromPEM(ca)
	if !ok {
		log.Fatal("Failed to add CA")
	}

	crt, err := tls.LoadX509KeyPair(*cert, *key)
	if err != nil {
		log.Fatalf("Failed to load server TLS certificate")
	}

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		log.Fatalf("Failed to listen on port: %d, %v", *port, err)
	}
	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{crt},
		ClientCAs:    certPool,
	}
	s := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)), grpc.MaxRecvMsgSize(100*1024*1024))
	codesigner.RegisterCodeSignerServer(s, &codesigner.CodeSigner{})
	reflection.Register(s)
	log.Println("Starting server...")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to start GRPC server: %v", err)
	}
}
