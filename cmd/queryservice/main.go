package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type Message struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}

func main() {
	r := gin.Default()
	redisClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})
	defer redisClient.Close()

	r.GET("/message/list", func(c *gin.Context) {
		sender := c.Query("sender")
		receiver := c.Query("receiver")

		if sender == "" || receiver == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "sender and receiver are required"})
			return
		}

		redisKey := sender + "__" + receiver
		messages, err := redisClient.LRange(context.Background(), redisKey, 0, -1).Result()
		if err != nil {
			log.Printf("Failed to retrieve the messages from Redis: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
			return
		}

		var result []Message
		for _, msg := range messages {
			var message Message
			if err := json.Unmarshal([]byte(msg), &message); err != nil {
				log.Printf("Failed to unmarshal the message: %s", err)
				continue
			}
			result = append(result, message)
		}

		c.JSON(http.StatusOK, result)
	})

	r.Run(":8081")
}
