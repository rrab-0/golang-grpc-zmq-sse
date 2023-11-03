package zmq_local

import (
	"log"
	"os"

	zmq "github.com/pebbe/zmq4"
)

var GlobalPublisher *zmq.Socket
var DefaultTopic = "default-topic"

func Publisher() *zmq.Socket {
	publisher, err := zmq.NewSocket(zmq.PUB)
	if err != nil {
		log.Printf("PUB Error: %s\n", err)
	}

	err = publisher.Bind("tcp://" + os.Getenv("ZMQ_PUB_HOST") + ":" + os.Getenv("ZMQ_PUB_PORT"))
	if err != nil {
		log.Printf("PUB Error: %s\n", err)
	}
	GlobalPublisher = publisher

	log.Println("ZMQ Publisher is up at " + os.Getenv("ZMQ_PUB_HOST") + ":" + os.Getenv("ZMQ_PUB_PORT"))
	return publisher
}
