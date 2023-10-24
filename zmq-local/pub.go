package zmq_local

import (
	"os"

	zmq "github.com/pebbe/zmq4"

	"fmt"
)

var GlobalPublisher *zmq.Socket

func Publisher() {
	publisher, _ := zmq.NewSocket(zmq.PUB)
	defer publisher.Close()

	publisher.Bind("tcp://*:" + os.Getenv("ZMQ_PUB_PORT"))
	publisher.Bind("ipc://weather.ipc")
	GlobalPublisher = publisher

	fmt.Println("ZMQ Publisher is up at :" + os.Getenv("ZMQ_PUB_PORT"))

	// //  Initialize random number generator
	// rand.Seed(time.Now().UnixNano())

	// // loop for a while apparently
	// for {
	// 	//  Get values that will fool the boss
	// 	zipcode := rand.Intn(100000)
	// 	temperature := rand.Intn(215) - 80
	// 	relhumidity := rand.Intn(50) + 10

	// 	//  Send message to all subscribers
	// 	msg := fmt.Sprintf("%05d %d %d", zipcode, temperature, relhumidity)
	// 	publisher.Send(msg, 0)
	// }
}

// func SendMessage(msg string) error {
// 	_, err := GlobalPublisher.Send(msg, 0)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
