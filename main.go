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

	go func() {
		grpc_server.StartTodo()
	}()

	pub := zmq_local.Publisher()
	defer pub.Close()

	sub := zmq_local.Subscriber()
	defer sub.Close()

	sse_server.Start()
}
