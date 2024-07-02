package main

import (
	"log"
	"os"
	"os/signal"
	"github.com/IBM/sarama"
)

var (
	brokers = []string{"kafka:9092"}
	config  = sarama.NewConfig()
)

var highPriorityTopic = "high-priority-topic"
var mediumPriorityTopic = "medium-priority-topic"
var lowPriorityTopic = "low-priority-topic"

func main() {
	// Create consumer
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalf("Error closing consumer: %v", err)
		}
	}()

	// Create partition consumers
	highPriorityPartitionConsumer, err := consumer.ConsumePartition(highPriorityTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to start high priority partition consumer: %v", err)
	}

	mediumPriorityPartitionConsumer, err := consumer.ConsumePartition(mediumPriorityTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to start medium priority partition consumer: %v", err)
	}

	lowPriorityPartitionConsumer, err := consumer.ConsumePartition(lowPriorityTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to start low priority partition consumer: %v", err)
	}

	// Handle high-priority messages
	go func() {
		for message := range highPriorityPartitionConsumer.Messages() {
			// Process high-priority message
			log.Printf("High priority message: %s", message.Value)
		}
	}()

	// Handle medium-priority messages
	go func() {
		for message := range mediumPriorityPartitionConsumer.Messages() {
			// Process medium-priority message
			log.Printf("Medium priority message: %s", message.Value)
		}
	}()

	// Handle low-priority messages
	go func() {
		for message := range lowPriorityPartitionConsumer.Messages() {
			// Process low-priority message
			log.Printf("Low priority message: %s", message.Value)
		}
	}()

	// Handle OS signals to gracefully shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan

	log.Println("Interrupt signal received, shutting down...")
}

