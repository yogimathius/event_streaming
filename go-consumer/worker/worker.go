package worker

import (
	"fmt"

	"go-consumer/message"
)

type Delegator interface {
	Delegate(message.Message, Producer)
}

type WorkerDelegator struct{}

type Producer interface {
	SendMessage(message.Message) error
}

func (w WorkerDelegator) Delegate(message message.Message, producer Producer) {
	switch message.Priority {
	case "High":
		w.delegateToHighWorker(message, producer)
	case "Medium":
		w.delegateToMedWorker(message, producer)
	case "Low":
		w.delegateToLowWorker(message, producer)
	default:
		fmt.Printf("Unknown priority: %s\n", message.Priority)
	}
}

func (w WorkerDelegator) delegateToHighWorker(message message.Message, producer Producer) {
	// job := rqueue.Enqueue(queueName, message)
	fmt.Printf("Delegated message %v to high-priority queue\n", message)
	message.Status = "message delegated high priority"
	producer.SendMessage(message)
}

func (w WorkerDelegator) delegateToMedWorker(message message.Message, producer Producer) {
	// job := rqueue.Enqueue(queueName, message)
	fmt.Printf("Delegated message %v to medium-priority queue\n", message)
	message.Status = "message delegated medium priority"
	producer.SendMessage(message)
}

func (w WorkerDelegator) delegateToLowWorker(message message.Message, producer Producer) {
	// job := rqueue.Enqueue(queueName, message)
	fmt.Printf("Delegated message %v to low-priority queue\n", message)
	message.Status = "message delegated low priority"
	producer.SendMessage(message)
}
