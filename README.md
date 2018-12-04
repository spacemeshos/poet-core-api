poet-core-api
=====

## pcrpc

This package implements both a client and server for POET-core RPC system
which is based off of the high-performance cross-platform
[gRPC](http://www.grpc.io/) RPC framework. By default, only the Go
client+server libraries are compiled within the package. In order to compile
the client-side libraries for other supported languages, the `protoc` tool will
need to be used to generate the compiled protos for a specific language.

The following languages are supported as clients to `pcrpc`: C++, Go, Node.js,
Java, Ruby, Android Java, PHP, Python, C# and Objective-C.

## pccli

This package implements CLI for `pcrpc` client.

## pc_test.go

Integration tests for POET-core service.

####

```
go get -u github.com/kardianos/govendor
govendor sync
```

```bash
$ go test -v
```
