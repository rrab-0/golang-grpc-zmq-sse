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

	pub := zmq_local.Publisher()
	defer pub.Close()

	// go func() {
	// 	for {
	// 		msg := "Hello from ZMQ PUB"
	// 		_, err := zmq_local.GlobalPublisher.Send(zmq_local.DefaultTopic+" "+msg, zmq.DONTWAIT)
	// 		if err != nil {
	// 			log.Printf("ZMQ PUB Error: %s\n", err)
	// 		}

	// 		log.Println("ZMQ PUB sent: " + msg)
	// 	}
	// }()

	sub := zmq_local.Subscriber()
	defer sub.Close()

	// go func() {
	// 	for {
	// 		msg, err := zmq_local.GlobalSubscriber.Recv(1)
	// 		if err != nil {
	// 			if err.Error() != "resource temporarily unavailable" {
	// 				log.Printf("ZMQ SUB Error: %s\n", err)
	// 				continue
	// 			}
	// 		}

	// 		log.Println("ZMQ SUB received: " + msg)
	// 	}
	// }()

	// sigs := make(chan os.Signal, 1)
	// signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	// <-sigs

	go func() {
		grpc_server.StartTodo()
	}()

	// go func() {
	// 	for {
	// 		log.Println("channel msg: " + <-sse_server.GlobalChannelSSE)
	// 	}
	// }()

	sse_server.Start()
}
