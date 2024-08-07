package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"go-consumer/message"
	"go-consumer/worker"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumer    sarama.Consumer
	producer    Producer
	workerDelegator worker.Delegator
}

type Producer interface {
	SendMessage(message.Message) error
}

func NewConsumer(brokers []string, producer Producer, workerDelegator worker.Delegator) (*Consumer, error) {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start consumer: %v", err)
	}

	return &Consumer{
		consumer:    consumer,
		producer:    producer,
		workerDelegator: workerDelegator,
	}, nil
}

func (c *Consumer) StartConsuming(topics []string, producer Producer) {
	for _, topic := range topics {
		partitions, err := c.consumer.Partitions(topic)
		if err != nil {
			log.Fatalf("Failed to get partitions for topic %s: %v", topic, err)
		}

		for _, partition := range partitions {
			log.Default().Printf("Consuming topic on %s partition %d\n", topic, partition)
			go c.consumePartition(topic, partition, producer)
		}
	}

	select {}
}

func (c *Consumer) consumePartition(topic string, partition int32, producer Producer) {
	partitionConsumer, err := c.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to start partition consumer for topic %s partition %d: %v", topic, partition, err)
	}
	defer partitionConsumer.Close()

	for msg := range partitionConsumer.Messages() {
		c.processMessage(msg, producer)
	}
}

func (c *Consumer) processMessage(msg *sarama.ConsumerMessage, producer Producer) {
	var message message.Message
	if err := json.Unmarshal(msg.Value, &message); err != nil {
		log.Printf("Error unmarshalling message: %v\n", err)
		return
	}
	if message.Status == "message produced" {
		message.Status = "message consumed"
		c.producer.SendMessage(message)

		c.workerDelegator.Delegate(message, producer)
	}
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
