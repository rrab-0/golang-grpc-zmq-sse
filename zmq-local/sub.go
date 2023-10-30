package zmq_local

import (
	"encoding/json"
	"log"
	"strings"

	"grpc-zmq-sse/db"
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

	err = subscriber.Connect("tcp://" + os.Getenv("ZMQ_SUB_HOST") + ":" + os.Getenv("ZMQ_SUB_PORT"))
	if err != nil {
		log.Printf("Error: %s\n", err)
	}

	GlobalSubscriber = subscriber
	log.Println("ZMQ Subscriber is up at " + os.Getenv("ZMQ_SUB_HOST") + ":" + os.Getenv("ZMQ_SUB_PORT"))

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
					continue
				}
			}

			// Remove topic from message
			if msgs := strings.Fields(msg); len(msgs) > 1 {
				var jsonMsg interface{}
				err = json.Unmarshal([]byte(msgs[1]), &jsonMsg)
				if err != nil {
					log.Printf("ZMQ SUB Error: %s\n", err)
					continue
				}
				log.Println("ZMQ SUB received: " + msgs[1])

				switch {
				case strings.Contains(msgs[1], `"status":"created"`):
					// TODO: Create
					err = db.GlobalConnection.Create(&db.Dump{Message: msgs[1]}).Error
					if err != nil {
						log.Printf("Error: %s\n", err)
						continue
					}
				case strings.Contains(msgs[1], `"status":"updated"`):
					// TODO: Update
				case strings.Contains(msgs[1], `"status":"deleted"`):
					// TODO: Delete
				}

				log.Println("PostgreSQL at sse-handler received: " + msgs[1])
				sse_server.GlobalChannelSSE <- msgs[1]
			}
		}
	}()

	return subscriber
}
