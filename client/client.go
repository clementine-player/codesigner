package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/clementine/codesigner"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var address = flag.String("address", "localhost:5001", "Address of CodeSigner GRPC service")
var dmg = flag.String("dmg", "", "Path to unsigned DMG")
var developerID = flag.String("developer-id", "", "Developer ID to sign with")
var output = flag.String("output", "clementine.dmg", "Path to output signed dmg")
var cert = flag.String("cert", "", "Path to client TLS cert")
var key = flag.String("key", "", "Path to client TLS key")
var ca = flag.String("ca", "", "Path to CA Certificate")
var verify = flag.Bool("verify", false, "Whether to instead verify that the given DMG is signed correctly")

func main() {
	flag.Parse()
	crt, err := tls.LoadX509KeyPair(*cert, *key)
	if err != nil {
		log.Fatalf("Failed to load client TLS certificate: %v", err)
	}
	certPool := x509.NewCertPool()
	caCert, err := ioutil.ReadFile(*ca)
	if err != nil {
		log.Fatalf("Failed to load CA certificate: %v", err)
	}
	ok := certPool.AppendCertsFromPEM(caCert)
	if !ok {
		log.Fatalf("Failed to load CA certificate: %v", err)
	}
	addr, _, err := net.SplitHostPort(*address)
	if err != nil {
		log.Fatalf("Failed to parse address flag: %s", *address)
	}
	creds := credentials.NewTLS(&tls.Config{
		ServerName:   addr,
		Certificates: []tls.Certificate{crt},
		RootCAs:      certPool,
	})
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

	if *verify {
		reply, err := c.VerifyPackage(context.Background(), &codesigner.VerifyPackageRequest{
			Package: unsigned,
		})
		if err != nil {
			log.Fatalf("Failed to verify package: %v", err)
		}
		if reply.GetOk() {
			fmt.Printf("%s is signed correctly\n", *dmg)
		} else {
			fmt.Printf("%s is not signed correctly: %s\n", *dmg, reply.GetCodesignOutput())
		}
	} else {
		reply, err := c.SignPackage(context.Background(), &codesigner.SignPackageRequest{
			Package:     unsigned,
			DeveloperId: *developerID,
		})
		if err != nil {
			log.Fatalf("Failed to sign package: %v", err)
		}
		if err = ioutil.WriteFile(*output, reply.GetSignedPackage(), 0644); err != nil {
			log.Fatalf("Failed to write output dmg: %v", err)
		}
	}
}
