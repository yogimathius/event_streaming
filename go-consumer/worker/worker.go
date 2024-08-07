package worker

import (
	"fmt"

	"go-consumer/message"
)

type Delegator interface {
	Delegate(message.Message)
}

type WorkerDelegator struct{}

func (w WorkerDelegator) Delegate(message message.Message) {
	switch message.Priority {
	case "High":
		w.delegateToHighWorker(message)
	case "Medium":
		w.delegateToMedWorker(message)
	case "Low":
		w.delegateToLowWorker(message)
	default:
		fmt.Printf("Unknown priority: %s\n", message.Priority)
	}
}

func (w WorkerDelegator) delegateToHighWorker(message message.Message) {
	// job := rqueue.Enqueue(queueName, message)
	fmt.Printf("Delegated message %v to high-priority queue\n", message)
	message.Status = "message delegated high priority"
}

func (w WorkerDelegator) delegateToMedWorker(message message.Message) {
	// job := rqueue.Enqueue(queueName, message)
	fmt.Printf("Delegated message %v to medium-priority queue\n", message)
	message.Status = "message delegated medium priority"
}

func (w WorkerDelegator) delegateToLowWorker(message message.Message) {
	// job := rqueue.Enqueue(queueName, message)
	fmt.Printf("Delegated message %v to low-priority queue\n", message)
	message.Status = "message delegated low priority"
}
