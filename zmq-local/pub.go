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
		log.Printf("Error: %s\n", err)
	}

	publisher.Bind("tcp://*:" + os.Getenv("ZMQ_PUB_PORT"))
	GlobalPublisher = publisher

	log.Println("ZMQ Publisher is up at :" + os.Getenv("ZMQ_PUB_PORT"))
	return publisher
}
