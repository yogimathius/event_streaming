package kafka

import (
	"encoding/json"
	"fmt"
	"time"

	"go-consumer/message"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	syncProducer sarama.SyncProducer
	brokers      []string
}

func NewProducer(brokers []string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("error creating Kafka producer: %v", err)
	}

	return &KafkaProducer{
		syncProducer: producer,
		brokers:      brokers,
	}, nil
}

func (p *KafkaProducer) SendMessage(message message.Message) error {
	jsonMessage, err := json.Marshal(message)
	message.EventTime = time.Now()

	if err != nil {
		return fmt.Errorf("failed to marshal message to JSON: %v", err)
	}

	_, _, err = p.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: message.EventType,
		Value: sarama.StringEncoder(jsonMessage),
	})

	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %v", err)
	}

	return nil
}

func (p *KafkaProducer) Close() error {
	return p.syncProducer.Close()
}
