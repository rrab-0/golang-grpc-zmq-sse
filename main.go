package main

import (
	"grpc-zmq-sse/db"
	grpc_server "grpc-zmq-sse/grpc-server"
	sse_server "grpc-zmq-sse/sse-server"
	zmq_local "grpc-zmq-sse/zmq-local"
	"log"

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

	// sub always recving msgs
	go func() {
		defer close(sse_server.ChannelSSE)
		for {
			msg, err := zmq_local.GlobalSubscriber.Recv(1)
			if err != nil {
				if err.Error() != "resource temporarily unavailable" {
					log.Printf("ZMQ SUB Error: %s\n", err)
				}
			}

			if msg != "" {
				log.Println("ZMQ SUB received: " + msg)
				err = db.GlobalConnection.Create(&db.Dump{Message: msg}).Error
				if err != nil {
					log.Printf("Error: %s\n", err)
				}
				log.Println("PostgreSQL at sse-handler received: " + msg)

				sse_server.ChannelSSE <- msg
			}
		}
	}()

	go func() {
		grpc_server.Start()
	}()

	go func() {
		grpc_server.StartTodo()
	}()

	sse_server.Start()
}
