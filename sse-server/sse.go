package sse_server

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var ChannelSSE = make(chan string)

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

		fmt.Fprintf(c.Writer, "data: you are connected\n\n")
		c.Writer.Flush()

		// go func() {
		// 	defer close(msgCh)
		// 	for {
		// 		msg, err := zmq_local.GlobalSubscriber.Recv(1)
		// 		if err != nil {
		// 			if err.Error() != "resource temporarily unavailable" {
		// 				log.Printf("ZMQ SUB Error: %s\n", err)
		// 			}
		// 		}

		// 		if msg != "" {
		// 			log.Println("ZMQ SUB received: " + msg)
		// 			err = db.GlobalConnection.Create(&db.Dump{Message: msg}).Error
		// 			if err != nil {
		// 				log.Printf("Error: %s\n", err)
		// 			}
		// 			log.Println("PostgreSQL at sse-handler received: " + msg)

		// 			msgCh <- msg
		// 		}
		// 	}
		// }()

		// go func() {
		for {
			select {
			case msgFromSubZMQ := <-ChannelSSE:
				// c.SSEvent("message", msgFromSubZMQ)
				fmt.Fprintf(c.Writer, "data: %s\n\n", msgFromSubZMQ)
				c.Writer.Flush()
				log.Println("SSE Sent: " + msgFromSubZMQ)
			case <-c.Writer.CloseNotify():
				log.Println("client disconnected")
				return
			}
		}
		// }()
	})

	// r.GET("/sse", func(c *gin.Context) {
	// 	c.Header("Content-Type", "text/event-stream")
	// 	c.Header("Cache-Control", "no-cache")
	// 	c.Header("Connection", "keep-alive")
	// 	c.Header("Access-Control-Allow-Origin", "*")

	// 	msgCh := make(chan string)

	// 	// Start a goroutine to listen to ZMQ messages
	// 	go func() {
	// 		defer close(msgCh)
	// 		for {
	// 			msg, err := zmq_local.GlobalSubscriber.Recv(0) // Use 0 as the flag for non-blocking receive
	// 			if err != nil {
	// 				log.Printf("ZMQ SUB Error: %s\n", err)
	// 			}

	// 			if msg != "" {
	// 				log.Println("ZMQ SUB received: " + msg)
	// 				err = db.GlobalConnection.Create(&db.Dump{Message: msg}).Error
	// 				if err != nil {
	// 					log.Printf("Error: %s\n", err)
	// 				}
	// 				log.Println("PostgreSQL at sse-handler received: " + msg)

	// 				msgCh <- msg
	// 			}
	// 		}
	// 	}()

	// 	// Send messages to the client as they arrive
	// 	for {
	// 		select {
	// 		case msgFromSubZMQ, ok := <-msgCh:
	// 			if !ok {
	// 				continue // Channel closed, client disconnected
	// 			}
	// 			fmt.Fprintf(c.Writer, "data: %s\n\n", msgFromSubZMQ)
	// 			c.Writer.Flush()
	// 			log.Println("SSE Sent: " + msgFromSubZMQ)
	// 		case <-c.Writer.CloseNotify():
	// 			log.Println("client disconnected")
	// 			continue
	// 		}
	// 	}
	// })

	// r.GET("/sse", func(c *gin.Context) {
	// 	c.Header("Content-Type", "text/event-stream")
	// 	c.Header("Cache-Control", "no-cache")
	// 	c.Header("Connection", "keep-alive")
	// 	c.Header("Access-Control-Allow-Origin", "*")

	// 	msgCh := make(chan []byte)

	// 	// Register the client
	// 	clients[msgCh] = struct{}{}
	// 	defer func() {
	// 		// Unregister the client when the connection is closed
	// 		delete(clients, msgCh)
	// 		close(msgCh)
	// 		log.Println("Client disconnected")
	// 	}()

	// 	for msg := range msgCh {
	// 		// Send messages to the client
	// 		fmt.Fprintf(c.Writer, "data: %s\n\n", msg)
	// 		c.Writer.Flush()
	// 		log.Println("SSE Sent: " + string(msg))

	// 	}

	// 	// Check if the client connection is closed
	// 	select {
	// 	case <-c.Writer.CloseNotify():
	// 		return
	// 	}
	// })

	// go func() {
	// 	for {
	// 		msg, err := zmq_local.GlobalSubscriber.Recv(1)
	// 		if err != nil {
	// 			if err.Error() != "resource temporarily unavailable" {
	// 				log.Printf("ZMQ SUB Error: %s\n", err)
	// 			}
	// 		}

	// 		payload := map[string]interface{}{
	// 			"message": msg,
	// 		}

	// 		gMsg, _ := json.Marshal(payload)

	// 		if msg != "" {
	// 			log.Println("ZMQ SUB received: " + msg)
	// 			err = db.GlobalConnection.Create(&db.Dump{Message: msg}).Error
	// 			if err != nil {
	// 				log.Printf("Error: %s\n", err)
	// 			}
	// 			log.Println("PostgreSQL at sse-handler received: " + msg)

	// 			for client := range clients {
	// 				client <- gMsg
	// 			}
	// 		}
	// 	}
	// }()

	r.Run(":" + os.Getenv("SSE_SERVER_PORT"))
}
