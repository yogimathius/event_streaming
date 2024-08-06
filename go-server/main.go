package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

type Message struct {
	Id string `json:"id"`
	EventType string    `json:"event_type"`
	EventTime time.Time `json:"event_time"`
	Priority  string    `json:"priority"`
	Description string `json:"description"`
	Status string `json:"status"`
}

var (
	brokers  = []string{"kafka:29092"} 
	producer sarama.SyncProducer
)

func main() {
	initKafkaProducer()

	router := gin.Default()

	router.POST("/message", handleMessage)

	fmt.Println("Server listening on port 8080...")
	router.Run(":8080")

	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf("Error closing Kafka producer: %v\n", err)
		}
	}()
}

func handleMessage(c *gin.Context) {
	var msg Message
	msg.Status = "message produced"
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sendMessageToKafka(msg)

	c.JSON(http.StatusOK, gin.H{"message": "Message sent to Kafka successfully"})
}

func initKafkaProducer() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Return.Successes = true

	var err error
	producer, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}

	fmt.Println("Kafka producer initialized")
}


func sendMessageToKafka(message Message) {
	jsonMessage, err := json.Marshal(message)
	
	if err != nil {
		log.Printf("Failed to marshal message to JSON: %v\n", err)
		return
	}
	partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic: message.EventType,
		Value: sarama.StringEncoder(jsonMessage),
	})
	fmt.Printf("Message sent to partition %d at offset %d: %s\n", partition, offset, jsonMessage)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v\n", err)
		return
	}
}
