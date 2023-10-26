package main

import (
	"fmt"
	"log"

	zmq_local "grpc-zmq-sse/zmq-local"

	zmq "github.com/pebbe/zmq4"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic("ERROR: Could not load .env")
	}

	// db.Connect()

	pub := zmq_local.Publisher()
	defer pub.Close()

	sub := zmq_local.Subscriber()
	defer sub.Close()

	zmq_local.GlobalSubscriber.SetSubscribe("hello2")
	zmq_local.GlobalSubscriber.SetSubscribe("hello3")
	pubMsg := "i am mock1"
	for {
		// Publisher
		_, err := zmq_local.GlobalPublisher.Send("hello1 "+pubMsg, zmq.DONTWAIT)
		if err != nil {
			fmt.Printf("ZMQ PUB Error: %s\n", err)
		}
		log.Println("ZMQ PUB Sent: " + pubMsg)

		// Sub to mock2
		msgMock2, err := zmq_local.GlobalSubscriber.Recv(zmq.DONTWAIT)
		// if err != nil {
		// 	log.Printf("ZMQ SUB Mock 2 Error: %s\n", err)
		// }
		log.Println("ZMQ SUB Mock 2 received: " + msgMock2)

		// Sub to mock3
		msgMock3, err := zmq_local.GlobalSubscriber.Recv(zmq.DONTWAIT)
		// if err != nil {
		// 	log.Printf("ZMQ SUB Mock 3 Error: %s\n", err)
		// }
		log.Println("ZMQ SUB Mock 3 received: " + msgMock3)
	}

	// go func() {
	// 	grpc_server.Start()
	// }()

	// go func() {
	// 	grpc_server.StartTodo()
	// }()

	// sse_server.Start()
}

// netsh interface portproxy add v4tov4 listenport=5555 listenaddress=172.20.10.7 connectport=5555 connectaddress=172.28.13.233
// netsh interface portproxy add v4tov4 listenport=5556 listenaddress=172.20.10.7 connectport=5556 connectaddress=172.28.13.233
// netsh interface portproxy add v4tov4 listenport=5557 listenaddress=172.20.10.7 connectport=5557 connectaddress=172.28.13.233
