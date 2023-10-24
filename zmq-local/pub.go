package zmq_local

import (
	"log"
	"os"

	zmq "github.com/pebbe/zmq4"
)

var GlobalPublisher *zmq.Socket

func Publisher() *zmq.Socket {
	publisher, err := zmq.NewSocket(zmq.PUB)
	if err != nil {
		log.Printf("Error: %s\n", err)
	}
	publisher.Bind("tcp://*:" + os.Getenv("ZMQ_PUB_PORT"))
	publisher.Bind("ipc://weather.ipc")
	GlobalPublisher = publisher

	log.Println("ZMQ Publisher is up at :" + os.Getenv("ZMQ_PUB_PORT"))
	return publisher
}
