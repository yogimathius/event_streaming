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

var (
	brokers  = []string{"kafka:29092"} 
	producer sarama.SyncProducer
)

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
		initKafkaProducer()

    brokers := []string{os.Getenv("KAFKA_BROKER")}
		topics := strings.Split(os.Getenv("KAFKA_TOPICS"), ",")
    consumer, err := sarama.NewConsumer(brokers, nil)
    if err != nil {
        log.Fatalf("Failed to start consumer: %v", err)
    }
    defer consumer.Close()
		defer func() {
			if err := producer.Close(); err != nil {
				log.Printf("Error closing Kafka producer: %v\n", err)
			}
		}()
		
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
	if message.Status == "message produced" {
		message.Status = "message consumed"
		sendMessageToKafka(message)

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
}

func delegateToHighWorker(message Message) {
	// job := rqueue.Enqueue(queueName, message)
	fmt.Printf("Delegated message %v to high-priority queue\n", message)
	message.Status = "message delegated high priority"
	sendMessageToKafka(message)
}

func delegateToMedWorker(message Message) {
	// job := rqueue.Enqueue(queueName, message)
	fmt.Printf("Delegated message %v to medium-priority queue\n", message)
	message.Status = "message delegated medium priority"
	sendMessageToKafka(message)
}

func delegateToLowWorker(message Message) {
	// job := rqueue.Enqueue(queueName, message)
	fmt.Printf("Delegated message %v to low-priority queue\n", message)
	message.Status = "message delegated low priority"
	sendMessageToKafka(message)
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
	message.EventTime = time.Now()
	
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