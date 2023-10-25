package sse_server

import (
	"grpc-zmq-sse/db"
	zmq_local "grpc-zmq-sse/zmq-local"
	"log"
	"os"
	"time"

	zmq "github.com/pebbe/zmq4"

	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/sse", func(c *gin.Context) {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("Access-Control-Allow-Origin", "*")

		// msgCh := make(chan string, 100) // Use a buffered channel to avoid blocking
		msgCh := make(chan string)
		go func() {
			for {
				// TODO, fix bug here:
				// sub.Recv is blocking therefore can't "connect"
				// to sse before "msg" var is populated by pub.Send
				// msg, err := zmq_local.GlobalSubscriber.Recv(zmq.DONTWAIT)
				// if err != nil && zmq.AsErrno(err) != zmq.Errno(zmq.EFSM) {
				// 	log.Printf("ZMQ SUB Error: %s\n", err)
				// 	continue
				// }
				// log.Println("ZMQ SUB received: " + msg)

				// err = db.GlobalConnection.Create(&db.Dump{Message: msg}).Error
				// if err != nil {
				// 	log.Printf("Error: %s\n", err)
				// 	continue
				// }
				// log.Println("PostgreSQL at sse-handler received: " + msg)

				// msg, err := zmq_local.GlobalSubscriber.Recv(zmq.DONTWAIT)
				// if err == nil {
				// 	log.Println("ZMQ SUB received: " + msg)
				// 	err = db.GlobalConnection.Create(&db.Dump{Message: msg}).Error
				// 	if err != nil {
				// 		log.Printf("Error: %s\n", err)
				// 	}
				// 	log.Println("PostgreSQL at sse-handler received: " + msg)
				// 	msgCh <- msg
				// } else if zmq.AsErrno(err) != zmq.Errno(zmq.ETIMEDOUT) || zmq.AsErrno(err) != zmq.Errno(zmq.EADDRNOTAVAIL) {
				// 	log.Printf("ZMQ SUB Error: %s\n", err)
				// }

				// Current solution:
				// - still send sse msgs but change msg to "ZMQ SUB: Waiting for message from Publisher..."
				msg, err := zmq_local.GlobalSubscriber.Recv(zmq.DONTWAIT)

				switch {
				case err == nil:
					err = db.GlobalConnection.Create(&db.Dump{Message: msg}).Error
					if err != nil {
						log.Printf("Error: %s\n", err)
						continue
					}

					log.Println("PostgreSQL at sse-handler received: " + msg)
					msgCh <- msg
					continue
				// case zmq.AsErrno(err) == zmq.Errno(zmq.EFSM):
				case err.Error() == "resource temporarily unavailable":
					msg = "ZMQ SUB: Waiting for message from Publisher..."
					msgCh <- msg

					log.Println(msg)
					time.Sleep(3000 * time.Millisecond)
					continue
				default:
					log.Printf("ZMQ SUB Error: %s\n", err)
					continue
				}
			}
		}()

		for {
			select {
			case msgFromSubZMQ := <-msgCh:
				c.SSEvent("message", msgFromSubZMQ)
				c.Writer.Flush()
				log.Println("SSE Sent: " + msgFromSubZMQ)
			case <-c.Writer.CloseNotify():
				log.Println("client disconnected")
				return
			}
		}
	})

	r.Run(":" + os.Getenv("SSE_SERVER_PORT")) // listen and serve on 0.0.0.0:8080
}
