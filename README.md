# Golang gRPC bidirectional streaming

- client sends random numbers to server
- server receives number and sends it back if the number greater than all previous numbers
- both client and server handle context errors (try to close client during send)

## Requirements

- go 1.20
- protobuf installed
- go support for protobuf installed

## Installation

### Window 10

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
go get -u google.golang.org/grpc
```

Make sure ```protoc-gen-go``` added in PATH

It should create two binaries `server` and `client`

## Quick Run Project
First clone the repo then go to go-mysql-crud folder. After that build your image and run by docker. Make sure you have docker in your machine. 

```
git clone https://github.com/faithcredit/go-BiDi-streaming.git

cd go-grpc-mongo-crud
cd greet/greet_server
go run .
cd greet/greet_client
go run .
```
## Using CLI

```
go install github.com/ktr0731/evans@latest

evans -p 50051 -r
```
