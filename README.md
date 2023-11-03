# Example usage of gRPC + ZeroMQ + SSE

## Run locally

### 1. Install protoc

1. install and check that it is working

```
$ sudo apt install -y protobuf-compiler
$ protoc --version
libprotoc 24.3
```

### 2. Install ZeroMQ (only works on WSL)

1. Install `pkg-config`

```
$ sudo apt install pkg-config
```

2. Make sure `pkg-config` is available

```
$ pkg-config --version
```

3. Install `libzmq`

```
$ sudo apt-get update
$ sudo apt-get install libzmq3-dev
```

4. Make sure `libzmq` is installed

```
$ pkg-config --modversion libzmq
```

5. Make sure `CGO_ENABLED` is set to `1` in go env

```
$ go env CGO_ENABLED
1
```

### 3. Run the project

1. Clone the repository

```
git clone https://github.com/rrab-0/golang-grpc-zmq-sse.git
```

2. Make sure go version is 1.21

```
$ go version
go version go1.21.0 linux/amd64
```

3. Run the app

```
$ cd golang-grpc-zmq-sse
$ go run main.go
```

## Run on 2 different devices (2 WSLs on a network)

1. Clone the repo

2. Configure the ports in your `.env` for both of your machines, for the HOSTs, use the IP of your WSL obtained through `ifconfig` for `GRPC_TODO_HOST`, `ZMQ_PUB_HOST`, and `SSE_SERVER_HOST`. For the `ZMQ_SUB_HOST`, use the Local Network IP of your other machine (ex: machine1 uses IP of machine2, and machine2 uses IP of machine1).

3. Forward WSL port to Windows
    1. In both of your machines, open up `cmd` as admin.
    2. Run this, don’t forget to change the localport to be ports from your .env, run this for all of your ports
        ```
        netsh advfirewall firewall add rule name="Allowing LAN connections" dir=in action=allow protocol=TCP localport={YOUR_PORT}
        ```
    3. Lastly, run this to connect/bridge your WSL IP with your Local Network IP, change the `listenaddress` to be your Local Network IP and `connectaddress` to be your WSL IP, also change the `listenport` and connectport to be the ports from your `.env` (run this command 3 times, don’t run this command for your `ZMQ_SUB_HOST`)
        ```
        netsh interface portproxy add v4tov4 listenport={YOUR_PORT} listenaddress={NETWORK_IP_ADDR} connectport={YOUR_PORT} connectaddress={WSL_IP_ADDR}
        ```
4. Don’t forget to make a `PostgreSQL` database for both of your machines in WSL and also change the `DB_NAME` in `.env`

5. Run the app

### Additional Notes

#### 1. Flow of data

-   this project uses go1.21.0
-   data flows from top to bottom:
    -   grpc client (postman/browser/etc)
    -   grpc server
    -   zmq pub
    -   zmq sub
    -   sse server
    -   sse client (postman/browser/etc)
-   google folder is a dependency for todo.proto that is currently not used

#### 2. Generate proto manually

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

#### 3. How to generate proto from a custom folder to another custom folder

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

-   then do this

```
protoc -I ./ --go_out=../generated-protos/ --go_opt=paths=source_relative --go-grpc_out=../generated-protos/ --go-grpc_opt=paths=source_relative hello-world.proto
```

-   so that our folder tree will look like this

```
.
├── generated-protos
│   ├── hello-world.pb.go
│   └── hello-world_grpc.pb.go
├── go.mod
└── protos
    └── hello-world.proto
```

#### 4. proto dependencies

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
