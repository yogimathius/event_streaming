package kafka

import (
	"encoding/json"
	"fmt"
	"go-worker/pool"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumer sarama.Consumer
	pool     *pool.WorkerPool
	mu       sync.Mutex
	running  bool
}

func NewConsumer(brokers []string, pool *pool.WorkerPool) (*Consumer, error) {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start consumer: %v", err)
	}

	return &Consumer{
		consumer: consumer,
		pool:     pool,
		running:  false,
	}, nil
}

func (c *Consumer) StartConsuming(topic string) {
	partitions, err := c.consumer.Partitions(topic)
	if err != nil {
		log.Fatalf("Failed to get partitions for topic %s: %v", topic, err)
	}

	for _, partition := range partitions {
		log.Printf("Consuming topic on %s partition %d\n", topic, partition)
		go c.consumePartition(topic, partition)
	}

	select {}
}

func (c *Consumer) consumePartition(topic string, partition int32) {
	partitionConsumer, err := c.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to start partition consumer for topic %s partition %d: %v", topic, partition, err)
	}
	defer partitionConsumer.Close()

	for msg := range partitionConsumer.Messages() {
		log.Printf("Received message: %s", string(msg.Value))
		if c.shouldRunPool(msg) {
			c.runWorkerPool()
		}
	}
}
type EventCreatedMessage struct {
	EventType string    `json:"event_type"`
}

func (c *Consumer) shouldRunPool(msg *sarama.ConsumerMessage) bool {
	// Add logic here to determine if the pool should run based on the message content.
	// For example, check if the message contains a specific trigger keyword.
	triggerKeyword := "event_created"
	fmt.Println("Message value: ", string(msg.Value))
	// Message value:  {"event_type":"event_created"}
	// convert to json and check event_type
	var message EventCreatedMessage
	if err := json.Unmarshal(msg.Value, &message); err != nil {
		log.Printf("Error unmarshalling message: %v\n", err)
		return false
	}


	return message.EventType == triggerKeyword
}

func (c *Consumer) runWorkerPool() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		log.Println("Worker pool is already running, skipping start")
		return
	}

	log.Println("Starting worker pool")

	go c.pool.ActivateWorkers()
	go c.pool.Run()
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
