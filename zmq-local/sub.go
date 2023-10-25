package zmq_local

import (
	"log"

	zmq "github.com/pebbe/zmq4"

	"os"
)

var GlobalSubscriber *zmq.Socket

func Subscriber() *zmq.Socket {
	subscriber, err := zmq.NewSocket(zmq.SUB)
	if err != nil {
		log.Printf("Error: %s\n", err)
	}

	subscriber.Connect("tcp://localhost:" + os.Getenv("ZMQ_SUB_PORT"))
	GlobalSubscriber = subscriber
	log.Println("ZMQ Subscriber is up at :" + os.Getenv("ZMQ_SUB_PORT"))

	// Subscribe to topic 10001
	filter := "10001 "    // zipcode, default is NYC, 10001
	if len(os.Args) > 1 { // can set topic with cli args
		filter = os.Args[1] + " "
	}
	subscriber.SetSubscribe(filter)
	return subscriber
}
