# codesigner
Mac code signing GRPC service

## Running the codesigner server

`go run server.go -cert some.host.name.crt -key some.host.name.key -password MacKeychainPassword`

## Signing a DMG

`go run client.go -cert some.host.name.crt -developer-id "Mac Developer ID" -address some.host.name:5001 -dmg path/to/some.dmg`

## Generating GRPC TLS Certificates & Keys

1. `go get -u github.com/square/certstrap`
1. `certstrap init --common-name "Some CA Name"`
1. `certstrap request-cert some.host.name`
1. `certstrap sign some.host.name --CA "Some CA Name"`
