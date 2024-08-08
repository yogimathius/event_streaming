package main

import (
	"database/sql"
	"fmt"
	db "gin-kafka-producer/database"
	"gin-kafka-producer/kafka"
	"gin-kafka-producer/message"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Producer interface {
	SendMessage(message.Message) error
}

func main() {
	database, err := db.InitDb()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer database.Close()
	producer, err := kafka.NewProducer([]string{"kafka:29092"})
	if err != nil {
		log.Fatalf("Kafka producer initialization failed: %v", err)
	}
	defer producer.Close()
	router := gin.Default()

	router.POST("/message", func(c *gin.Context) {
		handleMessage(c, producer, database)
	})

	fmt.Println("Server listening on port 8080...")
	router.Run(":8080")

	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf("Error closing Kafka producer: %v\n", err)
		}
	}()
}

func handleMessage(c *gin.Context, producer Producer, database *sql.DB) {
	var msg message.Message
	currentEventId, err := db.FetchLatestEvent(database)
	if err != nil {
		log.Printf("Error fetching latest event: %v\n", err)
	}
	if currentEventId == 0 {
		event := db.Event{
			GuestSatisfaction: true,
			StressMarks:       0,
		}
		currentEventId, err = db.CreateEvent(database, event)
	}
	msg.Status = "message produced"
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	producer.SendMessage(msg)

	c.JSON(http.StatusOK, gin.H{"message": "Message sent to Kafka successfully"})
}
