package main

import (
	"grpc-zmq-sse/app/db"
	grpc_server "grpc-zmq-sse/app/grpc-server"
	sse_server "grpc-zmq-sse/app/sse-server"
	zmq_local "grpc-zmq-sse/app/zmq-local"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic("ERROR: Could not load .env")
	}

	db.Connect()

	pub := zmq_local.Publisher()
	defer pub.Close()

	sub := zmq_local.Subscriber()
	defer sub.Close()

	go func() {
		grpc_server.StartTodo()
	}()

	sse_server.Start()
}
