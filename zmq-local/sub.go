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

	mockSubErr2 := subscriber.Connect("tcp://" + os.Getenv("MOCK_IP2") + ":" + os.Getenv("ZMQ_SUB_PORT"))
	if mockSubErr2 != nil {
		log.Printf("mockSubErr2 Error: %s\n", mockSubErr2)
	}

	mockSubErr3 := subscriber.Connect("tcp://" + os.Getenv("MOCK_IP3") + ":" + os.Getenv("ZMQ_SUB_PORT"))
	if mockSubErr3 != nil {
		log.Printf("mockSubErr3 Error: %s\n", mockSubErr3)
	}

	GlobalSubscriber = subscriber
	log.Println("ZMQ Subscriber is up at :" + os.Getenv("ZMQ_SUB_PORT"))

	// // Subscribe to topic 10001
	// filter := "10001 "    // zipcode, default is NYC, 10001
	// if len(os.Args) > 1 { // can set topic with cli args
	// 	filter = os.Args[1] + " "
	// }
	// subscriber.SetSubscribe("hello2")
	// subscriber.SetSubscribe("hello3")

	return subscriber
}
