package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

type Message struct {
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
	Priority  string    `json:"priority"`
	Source    string    `json:"source"`
	Location  string    `json:"location"`
}

func main() {
    brokers := []string{os.Getenv("KAFKA_BROKER")}
		topics := strings.Split(os.Getenv("KAFKA_TOPICS"), ",")
		// team := os.Getenv("TEAM")
    consumer, err := sarama.NewConsumer(brokers, nil)
    if err != nil {
        log.Fatalf("Failed to start consumer: %v", err)
    }
    defer consumer.Close()

    for _, topic := range topics {
				partitions, err := consumer.Partitions(topic)
				if err != nil {
					log.Fatalf("Failed to get partitions for topic %s: %v", topic, err)
				}
		
				for _, partition := range partitions {
					log.Printf("Consuming messages from topic %s partition %d", topic, partition)
					go consumePartition(consumer, topic, partition)
				}
    }

    // Keep the main goroutine running
    select {}
}


func consumePartition(consumer sarama.Consumer, topic string, partition int32) {
	partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
			log.Fatalf("Failed to start partition consumer for topic %s partition %d: %v", topic, partition, err)
	}
	defer partitionConsumer.Close()

	for msg := range partitionConsumer.Messages() {
			log.Printf("Received message from topic %s partition %d: %s", topic, partition, string(msg.Value))
			processMessage(msg)
	}
}

func processMessage(msg *sarama.ConsumerMessage) {
	fmt.Printf("Message received: key=%s value=%s\n", string(msg.Key), string(msg.Value))

	// Unmarshal the JSON message
	var message Message
	if err := json.Unmarshal(msg.Value, &message); err != nil {
			log.Printf("Error unmarshalling message: %v\n", err)
			return
	}

	// Process the message based on priority
	switch message.Priority {
	case "High":
			delegateToHighWorker(message)
	case "Medium":
			delegateToMedWorker(message)
	case "Low":
			delegateToLowWorker(message)
	default:
			log.Printf("Unknown priority: %s\n", message.Priority)
	}
}

func delegateToHighWorker(message Message) {
	// Code to enqueue the message to the appropriate worker queue
	// For example, using Redis Queue (RQ):
	// job := rqueue.Enqueue(queueName, message)
	fmt.Printf("Delegated message %v to high-priority queue\n", message)
}

func delegateToMedWorker(message Message) {
	// Code to enqueue the message to the appropriate worker queue
	// For example, using Redis Queue (RQ):
	// job := rqueue.Enqueue(queueName, message)
	fmt.Printf("Delegated message %v to medium-priority queue\n", message)
}

func delegateToLowWorker(message Message) {
	// Code to enqueue the message to the appropriate worker queue
	// For example, using Redis Queue (RQ):
	// job := rqueue.Enqueue(queueName, message)
	fmt.Printf("Delegated message %v to low-priority queue\n", message)
}

