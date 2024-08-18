package worker

import (
	"database/sql"
	"fmt"
	"time"

	db "go-consumer/database"
	"go-consumer/message"
)

type Delegator interface {
	Delegate(message.Message, Producer)
}

type WorkerDelegator struct{
	database *sql.DB;
}

type Producer interface {
	SendMessage(message.Message) error
}

func NewWorkerDelegator(database *sql.DB) WorkerDelegator {
	return WorkerDelegator{
		database: database,
	}
}

func (w WorkerDelegator) Delegate(message message.Message, producer Producer) {
	switch message.Priority {
	case "High":
		w.delegateToHighWorker(message, producer, w.database)
	case "Medium":
		w.delegateToMedWorker(message, producer, w.database)
	case "Low":
		w.delegateToLowWorker(message, producer, w.database)
	default:
		fmt.Printf("Unknown priority: %s\n", message.Priority)
	}
}

func (w WorkerDelegator) delegateToHighWorker(message message.Message, producer Producer, database *sql.DB) {
	// job := rqueue.Enqueue(queueName, message)
	fmt.Printf("Delegated message %v to high-priority queue\n", message)
	
	if message.EventTime.Before(time.Now().Add(-6 * time.Second)) {
		fmt.Printf("Message %v is older than 6 seconds\n", message)
		message.Status = "message unresolved"
		db.AddEventMessage(database, 1, message)
		return
	}

	message.Status = "message delegated high priority"
	producer.SendMessage(message)

	db.AddEventMessage(database, 1, message)
}

func (w WorkerDelegator) delegateToMedWorker(message message.Message, producer Producer, database *sql.DB) {
	// job := rqueue.Enqueue(queueName, message)
	if message.EventTime.Before(time.Now().Add(-6 * time.Second)) {
		fmt.Printf("Message %v is older than 6 seconds\n", message)
		message.Status = "message unresolved"
		db.AddEventMessage(database, 1, message)
		return
	}

	fmt.Printf("Delegated message %v to medium-priority queue\n", message)
	message.Status = "message delegated medium priority"
	producer.SendMessage(message)


	db.AddEventMessage(database, 1, message)
}

func (w WorkerDelegator) delegateToLowWorker(message message.Message, producer Producer, database *sql.DB) {
	// job := rqueue.Enqueue(queueName, message)
	fmt.Printf("Delegated message %v to low-priority queue\n", message)
	message.Status = "message delegated low priority"
	producer.SendMessage(message)

	if message.EventTime.Before(time.Now().Add(-6 * time.Second)) {
		fmt.Printf("Message %v is older than 6 seconds\n", message)
		message.Status = "message unresolved"
		db.AddEventMessage(database, 1, message)
		return
	}

	db.AddEventMessage(database, 1, message)
}
