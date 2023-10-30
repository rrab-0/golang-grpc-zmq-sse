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

				var dbTodo db.Todo
				todoId := jsonMsg.ID
				dbTodo.Title = jsonMsg.Title
				dbTodo.Description = jsonMsg.Description

				if jsonMsg.Completed == "true" {
					dbTodo.Completed = true
				}
				dbTodo.Completed = false

				switch {
				case strings.Contains(msgNoTopic, `"status":"created"`):
					err = db.GlobalConnection.Create(&dbTodo).Error
					if err != nil {
						log.Printf("Error: %s\n", err)
						continue
					}
				case strings.Contains(msgNoTopic, `"status":"updated"`):
					var woiTodo db.Todo
					if err := db.GlobalConnection.Where("id = ?", todoId).Save(&woiTodo).Error; err != nil {
						log.Printf("Error: %s\n", err)
						continue
					}
				case strings.Contains(msgNoTopic, `"status":"deleted"`):
					var woiTodo db.Todo
					if err := db.GlobalConnection.Where("id = ?", todoId).Delete(&woiTodo).Error; err != nil {
						log.Printf("Error: %s\n", err)
						continue
					}
				}

				log.Println("PostgreSQL at sse-handler received: " + msgNoTopic)
				sse_server.GlobalChannelSSE <- msgNoTopic
			}

			// // Remove topic from message
			// if msgs := strings.Fields(msg); len(msgs) > 1 {
			// 	msgResult := strings.TrimPrefix(msg, DefaultTopic)
			// 	var jsonMsg db.SubMessage
			// 	err = json.Unmarshal([]byte(msgResult), &jsonMsg)
			// 	if err != nil {
			// 		log.Printf("ZMQ SUB Error: %s\n", err)
			// 		continue
			// 	}
			// 	log.Println("ZMQ SUB received: " + msgResult)

			// 	var dbTodo db.Todo
			// 	todoId := jsonMsg.ID
			// 	dbTodo.Title = jsonMsg.Title
			// 	dbTodo.Description = jsonMsg.Description

			// 	if jsonMsg.Completed == "true" {
			// 		dbTodo.Completed = true
			// 	}
			// 	dbTodo.Completed = false

			// 	switch {
			// 	case strings.Contains(msgs[1], `"status":"created"`):
			// 		err = db.GlobalConnection.Create(&dbTodo).Error
			// 		if err != nil {
			// 			log.Printf("Error: %s\n", err)
			// 			continue
			// 		}
			// 	case strings.Contains(msgs[1], `"status":"updated"`):
			// 		if err := db.GlobalConnection.Where("id = ?", todoId).Save(&dbTodo).Error; err != nil {
			// 			log.Printf("Error: %s\n", err)
			// 			continue
			// 		}
			// 	case strings.Contains(msgs[1], `"status":"deleted"`):
			// 		if err := db.GlobalConnection.Where("id = ?", todoId).Delete(&dbTodo).Error; err != nil {
			// 			log.Printf("Error: %s\n", err)
			// 			continue
			// 		}
			// 	}

			// 	log.Println("PostgreSQL at sse-handler received: " + msgResult)
			// 	sse_server.GlobalChannelSSE <- msgResult
			// }
		}
	}()

	return subscriber
}
