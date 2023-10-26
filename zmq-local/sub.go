package zmq_local

import (
	"grpc-zmq-sse/db"
	"log"

	sse_server "grpc-zmq-sse/sse-server"

	zmq "github.com/pebbe/zmq4"

	"os"
)

var GlobalSubscriber *zmq.Socket

func Subscriber() *zmq.Socket {
	subscriber, err := zmq.NewSocket(zmq.SUB)
	if err != nil {
		log.Printf("Error: %s\n", err)
	}

	err = subscriber.Connect("tcp://localhost:" + os.Getenv("ZMQ_SUB_PORT"))
	if err != nil {
		log.Printf("Error: %s\n", err)
	}

	GlobalSubscriber = subscriber
	log.Println("ZMQ Subscriber is up at :" + os.Getenv("ZMQ_SUB_PORT"))

	err = subscriber.SetSubscribe(DefaultTopic)
	if err != nil {
		log.Printf("Error: %s\n", err)
	}

	go func() {
		defer close(sse_server.GlobalChannelSSE)
		for {
			msg, err := subscriber.Recv(1)
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

				sse_server.GlobalChannelSSE <- msg
			}
		}
	}()

	return subscriber
}
