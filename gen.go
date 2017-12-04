//go:generate protoc --go_out=plugins=grpc:. -I$GOPATH/src/github.com/clementine/signer-service codesigner_service.proto
package codesigner
