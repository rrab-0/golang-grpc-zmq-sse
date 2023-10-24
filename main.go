package main

import (
	"grpc-zmq-sse/db"
	grpc_server "grpc-zmq-sse/grpc-server"
	sse_server "grpc-zmq-sse/sse-server"
	zmq_local "grpc-zmq-sse/zmq-local"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic("ERROR: Could not load .env")
	}

	db.Connect()

	zmq_local.Publisher()
	zmq_local.Subscriber()

	sse_server.Start()
	grpc_server.Start()
}
