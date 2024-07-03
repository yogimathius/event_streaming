package main

import (
    "fmt"
    "log"
    "os"
    "time"
		"strings"
    "github.com/IBM/sarama"
)

func main() {
    brokers := []string{os.Getenv("KAFKA_BROKER")}
		topics := strings.Split(os.Getenv("KAFKA_TOPICS"), ",")
		team := os.Getenv("TEAM")
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
					go consumePartition(consumer, topic, partition, team)
				}
    }

    // Keep the main goroutine running
    select {}
}

func processMessage(msg []byte) {
    // Process the message based on its type and priority
    fmt.Printf("Processing message: %s\n", msg)
    // Simulate processing time based on priority logic
    time.Sleep(2 * time.Second) // Adjust sleep time based on priority logic
}

func consumePartition(consumer sarama.Consumer, topic string, partition int32, team string) {
	partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
			log.Fatalf("Failed to start partition consumer for topic %s partition %d: %v", topic, partition, err)
	}
	defer partitionConsumer.Close()

	for msg := range partitionConsumer.Messages() {
			log.Printf("Team %s consumed message from topic %s partition %d: %s", team, msg.Topic, msg.Partition, string(msg.Value))
	}
}