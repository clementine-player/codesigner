//go:generate protoc --go_out=plugins=grpc:. -I$GOPATH/src/github.com/clementine/codesigner codesigner_service.proto
package codesigner
