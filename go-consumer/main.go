package main

import (
    "fmt"
    "log"
    "os"
    "time"

    "github.com/IBM/sarama"
)

func main() {
    brokers := []string{os.Getenv("KAFKA_BROKER")}
    topics := []string{os.Getenv("KAFKA_TOPICS")}

    consumer, err := sarama.NewConsumer(brokers, nil)
    if err != nil {
        log.Fatalf("Failed to start consumer: %v", err)
    }
    defer consumer.Close()

    for _, topic := range topics {
        partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
        if err != nil {
            log.Fatalf("Failed to start partition consumer: %v", err)
        }
        defer partitionConsumer.Close()

        go func(pc sarama.PartitionConsumer) {
            for message := range pc.Messages() {
                log.Printf("Consumed message from topic %s: %s", topic, string(message.Value))
                processMessage(message.Value)
            }
        }(partitionConsumer)
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
