package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

type Message struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}

func main() {
	r := gin.Default()

	r.POST("/message", func(c *gin.Context) {
		var msg Message
		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("1")
		conn, err := amqp.Dial("amqp://danial:danial@rabbitmq:5672/") // ideally I won't hardcode this of course, just a bit short in time :)
		if err != nil {
			log.Println(err)
		}
		defer conn.Close()
		fmt.Println("1")

		ch, err := conn.Channel()
		if err != nil {
			log.Println(err)
		}
		defer ch.Close()
		fmt.Println("1")

		q, err := ch.QueueDeclare(
			"message_queue",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Println(err)
		}

		body, err := json.Marshal(msg)
		if err != nil {
			log.Println(err)
		}

		err = ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})
		if err != nil {
			log.Println(err)
		}

		c.JSON(http.StatusOK, gin.H{"status": "message sent"})
	})

	r.Run(":8080")
}
