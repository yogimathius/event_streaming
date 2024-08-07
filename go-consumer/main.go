package main

import (
	"go-consumer/kafka"
	"go-consumer/worker"
	"log"
	"os"
	"strings"
)

func main() {
	// Initialize database connection
	// database, err := db.InitDb()
	// if err != nil {
	// 	log.Fatalf("Database initialization failed: %v", err)
	// }
	// defer database.Close()

	producer, err := kafka.NewProducer([]string{"kafka:29092"})
	if err != nil {
		log.Fatalf("Kafka producer initialization failed: %v", err)
	}
	defer producer.Close()

	brokers := []string{os.Getenv("KAFKA_BROKER")}
	topics := strings.Split(os.Getenv("KAFKA_TOPICS"), ",")
	workerDelegator := worker.WorkerDelegator{}
	consumer, err := kafka.NewConsumer(brokers, producer, workerDelegator)
	if err != nil {
		log.Fatalf("Kafka consumer initialization failed: %v", err)
	}
	defer consumer.Close()

	// Start consuming messages
	consumer.StartConsuming(topics)
}
