package sse_server

import (
	"fmt"
	"grpc-zmq-sse/db"
	zmq_local "grpc-zmq-sse/zmq-local"
	"os"

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

		msgCh := make(chan string)
		go func() {
			for {
				msg, err := zmq_local.GlobalSubscriber.Recv(0)
				if err != nil {
					fmt.Printf("Error: %s", err)
					continue
				}
				fmt.Println("ZMQ SUB received: " + msg)

				err = db.GlobalConnection.Create(&db.Dump{Message: msg}).Error
				if err != nil {
					fmt.Printf("Error: %s", err)
					continue
				}

				msgCh <- msg
			}
		}()

		for {
			msg := <-msgCh
			c.SSEvent("message", msg)
			fmt.Println("SSE Sent: " + msg)
			c.Writer.Flush()
		}
	})

	r.Run(":" + os.Getenv("SSE_PORT")) // listen and serve on 0.0.0.0:8080
}
