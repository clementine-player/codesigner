package main

import (
  "context"
  "flag"
  "io/ioutil"
  "log"

  "github.com/clementine/codesigner"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials"
)

var address = flag.String("address", "localhost:5001", "Address of CodeSigner GRPC service")
var dmg = flag.String("dmg", "", "Path to unsigned DMG")
var developerID = flag.String("developer-id", "", "Developer ID to sign with")
var password = flag.String("password", "", "Password for keychain containing developer ID")
var output = flag.String("output", "clementine.dmg", "Path to output signed dmg")
var cert = flag.String("cert", "", "Path to server TLS cert")

func main() {
  flag.Parse()
  creds, err := credentials.NewClientTLSFromFile(*cert, "")
  if err != nil {
    log.Fatalf("Failed to load TLS cert")
  }
  conn, err := grpc.Dial(*address, grpc.WithTransportCredentials(creds), grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(100*1024*1024), grpc.MaxCallRecvMsgSize(100*1024*1024)))
  if err != nil {
    log.Fatalf("Could not connect to GRPC service: %v", err)
  }
  defer conn.Close()

  unsigned, err := ioutil.ReadFile(*dmg)
  if err != nil {
    log.Fatalf("Failed to read DMG file: %v", err)
  }
  c := codesigner.NewCodeSignerClient(conn)
  reply, err := c.SignPackage(context.Background(), &codesigner.SignPackageRequest{
    Package: unsigned,
    DeveloperId: *developerID,
    Password: *password,
  })
  if err != nil {
    log.Fatalf("Failed to sign package: %v", err)
  }
  if err = ioutil.WriteFile(*output, reply.GetSignedPackage(), 0644); err != nil {
    log.Fatalf("Failed to write output dmg: %v", err)
  }
}
