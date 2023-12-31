package sse_server

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var GlobalChannelSSE = make(chan string, 2)
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
		/**
		 * If flush isn't called before rcving msgs,
		 * the event will not be sent until the buffer is filled
		 * (infinite loading for client when trying to connect to /sse).
		 */
		fmt.Fprintf(c.Writer, "data: you are connected\n\n")
		c.Writer.Flush()

		for {
			select {
			case msgFromSubZMQ := <-GlobalChannelSSE:
				// c.SSEvent("message", msgFromSubZMQ)
				fmt.Fprintf(c.Writer, "data: %s\n\n", msgFromSubZMQ)
				c.Writer.Flush()
				log.Println("SSE Sent: " + msgFromSubZMQ)
			case <-c.Writer.CloseNotify():
				log.Println("client disconnected")
				return
			}
		}
	})

	r.Run(os.Getenv("SSE_SERVER_HOST") + ":" + os.Getenv("SSE_SERVER_PORT"))
}
