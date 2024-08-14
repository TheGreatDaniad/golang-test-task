package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
)

type Message struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}

func main() {
	conn, err := amqp.Dial("amqp://danial:danial@rabbitmq:5672/") 
	if err != nil {
		log.Fatalf("Failed to connect to rabbitMQ server: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open the channel: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"message_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to declare the queue: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})
	defer redisClient.Close()

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var msg Message
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				log.Println(err)
				continue
			}

			redisKey := msg.Sender + "__" + msg.Receiver
			if err := redisClient.LPush(context.Background(), redisKey, d.Body).Err(); err != nil {
				log.Println(err)
			}
		}
	}()

	<-forever
}
