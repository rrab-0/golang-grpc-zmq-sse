# Example usage of gRPC + ZeroMQ + SSE

## Run locally

1. Install protoc from https://github.com/protocolbuffers/protobuf/releases/tag/v24.4

2. Put the protoc binary in your PATH

3. Install the protocol compiler plugins for Go

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
```

```
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

4. Ensure that protoc-gen-go.exe/protoc-gen-go-grpc.exe are in your $GOPATH/bin

## some notes

-   this project uses go1.21.0
-   data flows from top to bottom:
    -   grpc client (postman/browser/etc)
    -   grpc server
    -   zmq pub
    -   zmq sub
    -   sse server
    -   sse client (postman/browser/etc)
-   google folder is a dependency for todo.proto

## generate proto

```
$ protoc --go_out=generated-proto-todo --go_opt=paths=source_relative \
> --go-grpc_out=generated-proto-todo --go-grpc_opt=paths=source_relative \
> todo.proto
```

```
protoc --go_out=../generated-proto-trying --go_opt=paths=source_relative \
> --go-grpc_out=../generated-proto-trying --go-grpc_opt=paths=source_relative \
> trying.proto
```

-   will generate to generated-proto-todo folder

## how to generate proto from a custom folder to another custom folder

-   with a dir sturcture like this:

```
.
├── generated-protos
└── protos
    └── hello-world.proto

2 directories, 1 file
```

-   we first need to cd to the protos folder

```
cd protos
```

-   then do the command there

```
protoc -I ./ --go_out=../generated-protos/ --go_opt=paths=source_relative --go-grpc_out=../generated-protos/ --go-grpc_opt=paths=source_relative hello-world.proto
```

-   so that it our folder tree will look like this

```
.
├── generated-protos
│   ├── hello-world.pb.go
│   └── hello-world_grpc.pb.go
├── go.mod
└── protos
    └── hello-world.proto
```

## proto dependencies

if we have a proto file that depends on another proto file, we can put the dependencies on the same directory as the proto file that depends on it.

-   example, todo.proto imports google/api/annotations.proto like this,

```
import "google/api/annotations.proto";
```

-   then tree dir structure should look like,

```
.
├── google
│   └── api
│       ├── annotations.proto
│       └── http.proto
├── hello-world.proto
└── todo.proto
```
