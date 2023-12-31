package zmq_local

import (
	"encoding/json"
	"grpc-zmq-sse/app/db"
	"log"
	"strings"

	sse_server "grpc-zmq-sse/app/sse-server"

	"github.com/google/uuid"
	zmq "github.com/pebbe/zmq4"

	"os"
)

var GlobalSubscriber *zmq.Socket

func Subscriber() *zmq.Socket {
	subscriber, err := zmq.NewSocket(zmq.SUB)
	if err != nil {
		log.Printf("SUB Error: %s\n", err)
	}

	err = subscriber.Connect("tcp://" + os.Getenv("ZMQ_SUB_HOST") + ":" + os.Getenv("ZMQ_SUB_PORT"))
	if err != nil {
		log.Printf("SUB Error: %s\n", err)
	}

	GlobalSubscriber = subscriber
	log.Println("ZMQ Subscriber is up at " + os.Getenv("ZMQ_SUB_HOST") + ":" + os.Getenv("ZMQ_SUB_PORT"))

	err = subscriber.SetSubscribe(DefaultTopic)
	if err != nil {
		log.Printf("SUB Error: %s\n", err)
	}

	go func() {
		for {
			msg, err := subscriber.Recv(zmq.DONTWAIT)
			if err != nil {
				if err.Error() != "resource temporarily unavailable" {
					log.Printf("ZMQ SUB Error: %s\n", err)
					continue
				}
			}

			// Remove topic prefix and send to SSE while also create/update/delete in PostgreSQL
			msgNoTopic := strings.TrimPrefix(msg, DefaultTopic)
			if msgNoTopic != "" {
				var jsonMsg db.SubMessage
				err = json.Unmarshal([]byte(msgNoTopic), &jsonMsg)
				if err != nil {
					log.Printf("ZMQ SUB Error: %s\n", err)
					continue
				}
				log.Println("ZMQ SUB received: " + msgNoTopic)

				todoId := jsonMsg.ID
				todoUUID, _ := uuid.Parse(todoId)
				var dbTodo db.Todo
				dbTodo.ID = todoUUID
				dbTodo.Title = jsonMsg.Title
				dbTodo.Description = jsonMsg.Description
				if jsonMsg.Completed == "true" {
					dbTodo.Completed = true
				}
				dbTodo.Completed = false

				switch {
				case jsonMsg.Status == "created":
					err = db.GlobalConnection.Create(&dbTodo).Error
					if err != nil {
						log.Printf("Error: %s\n", err)
						continue
					}
				case jsonMsg.Status == "updated":
					if err := db.GlobalConnection.Save(&dbTodo).Error; err != nil {
						log.Printf("Error: %s\n", err)
						continue
					}
				case jsonMsg.Status == "deleted":
					var deletedTodo db.Todo
					if err := db.GlobalConnection.Where("id = ?", todoId).Delete(&deletedTodo).Error; err != nil {
						log.Printf("Error: %s\n", err)
						continue
					}
				}

				log.Println("PostgreSQL at sse-handler received: " + msgNoTopic)
				sse_server.GlobalChannelSSE <- msgNoTopic
			}
		}
	}()

	return subscriber
}
