package zmq_local

import (
	zmq "github.com/pebbe/zmq4"

	"fmt"
	"os"
)

var GlobalSubscriber *zmq.Socket

func Subscriber() {
	subscriber, _ := zmq.NewSocket(zmq.SUB)
	defer subscriber.Close()
	subscriber.Connect("tcp://localhost:" + os.Getenv("ZMQ_SUB_PORT"))
	GlobalSubscriber = subscriber

	fmt.Println("ZMQ Subscriber is up at :" + os.Getenv("ZMQ_SUB_PORT"))

	// Subscribe to zipcode, default is NYC, 10001
	filter := "10001 "
	if len(os.Args) > 1 {
		filter = os.Args[1] + " "
	}
	subscriber.SetSubscribe(filter)

	// for {
	// 	msg, err := subscriber.Recv(0)
	// 	if err != nil {
	// 		fmt.Printf("Error: %s", err)
	// 	}

	// 	if msgs := strings.Fields(msg); len(msgs) > 1 {
	// 		fmt.Printf("data from Publisher: " + msgs[1] + "\n")
	// 	}
	// }
}

// func GetMessage() (string, error) {
// 	msg, err := GlobalSubscriber.Recv(0)
// 	if err != nil {
// 		return msg, err
// 	}

// 	return msg, nil
// }
