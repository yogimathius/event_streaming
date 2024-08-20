package kafka

import (
	"fmt"
	"go-worker/pool"
	"log"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumer    sarama.Consumer
	pool        *pool.WorkerPool
}

func NewConsumer(brokers []string, pool *pool.WorkerPool) (*Consumer, error) {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start consumer: %v", err)
	}

	return &Consumer{
		consumer:    consumer,
		pool:        pool,
	}, nil
}

func (c *Consumer) StartConsuming(topics []string) {
	topic := "event_created"
	partitions, err := c.consumer.Partitions(topic)
	if err != nil {
		log.Fatalf("Failed to get partitions for topic %s: %v", topic, err)
	}

	for _, partition := range partitions {
		log.Default().Printf("Consuming topic on %s partition %d\n", topic, partition)
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

	for range partitionConsumer.Messages() {
		c.pool.Run()
	}
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
