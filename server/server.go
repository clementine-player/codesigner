package main

import (
  "flag"
  "log"
  "net"
  "strconv"

  "github.com/clementine/codesigner"
  "google.golang.org/grpc"
  "google.golang.org/grpc/reflection"
)

var port = flag.Int("port", 5000, "Port to start GRPC server on")

func main() {
  flag.Parse()
  listener, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
  if err != nil {
    log.Fatalf("Failed to listen on port: %d, %v", *port, err)
  }
  s := grpc.NewServer()
  codesigner.RegisterCodeSignerServer(s, &codesigner.CodeSigner{})
  reflection.Register(s)
  log.Println("Starting server...")
  if err := s.Serve(listener); err != nil {
    log.Fatalf("Failed to start GRPC server: %v", err)
  }
}
