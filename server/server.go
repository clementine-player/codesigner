package main

import (
  "flag"
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

func main() {
  flag.Parse()
  listener, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
  if err != nil {
    log.Fatalf("Failed to listen on port: %d, %v", *port, err)
  }
  creds, err := credentials.NewServerTLSFromFile(*cert, *key)
  if err != nil {
    log.Fatalf("Failed to load server credentials: %v", err)
  }
  s := grpc.NewServer(grpc.Creds(creds), grpc.MaxRecvMsgSize(100*1024*1024))
  codesigner.RegisterCodeSignerServer(s, &codesigner.CodeSigner{})
  reflection.Register(s)
  log.Println("Starting server...")
  if err := s.Serve(listener); err != nil {
    log.Fatalf("Failed to start GRPC server: %v", err)
  }
}
