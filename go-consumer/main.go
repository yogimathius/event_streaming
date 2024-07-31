package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

type Message struct {
	Id string `json:"id"`
	EventType string    `json:"event_type"`
	EventTime time.Time `json:"event_time"`
	Priority  string    `json:"priority"`
	Description string `json:"description"`
	Status string `json:"status"`
}

func main() {
		// connStr := "postgres://postgres:password@localhost:5432/event_streaming"
    // db, err := sql.Open("postgres", connStr)
    // if err != nil {
		// 	log.Fatal(err)
		// }
		// defer db.Close()

		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer db.Close()
		// teamName := os.Getenv("TEAM")
		// teamId := os.Getenv("TEAM_ID")
		// teamID, err := strconv.Atoi(teamId)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// workers, err := fetchWorkerData(db, teamID)
    // if err != nil {
    //     log.Fatal(err)
    // }
		// log.Printf("Workers for team %d: %s\n", teamID, workers)
    brokers := []string{os.Getenv("KAFKA_BROKER")}
		topics := strings.Split(os.Getenv("KAFKA_TOPICS"), ",")
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

    select {}
}


func consumePartition(consumer sarama.Consumer, topic string, partition int32) {
	partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
			log.Fatalf("Failed to start partition consumer for topic %s partition %d: %v", topic, partition, err)
	}
	defer partitionConsumer.Close()

	for msg := range partitionConsumer.Messages() {
			processMessage(msg)
	}
}

func processMessage(msg *sarama.ConsumerMessage) {
	var message Message
	if err := json.Unmarshal(msg.Value, &message); err != nil {
			log.Printf("Error unmarshalling message: %v\n", err)
			return
	}

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
	// job := rqueue.Enqueue(queueName, message)
	fmt.Printf("Delegated message %v to high-priority queue\n", message)
}

func delegateToMedWorker(message Message) {
	// job := rqueue.Enqueue(queueName, message)
	fmt.Printf("Delegated message %v to medium-priority queue\n", message)
}

func delegateToLowWorker(message Message) {
	// job := rqueue.Enqueue(queueName, message)
	fmt.Printf("Delegated message %v to low-priority queue\n", message)
}

func fetchWorkerData(db *sql.DB, teamID int) (string, error) {
	var routineType string
	query := `SELECT workers FROM event_streaming WHERE team_id = $1`

	row := db.QueryRow(query, teamID)
	err := row.Scan(&routineType)
	if err != nil {
			return "", err
	}

	return routineType, nil
}
