package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	db "go-worker/database"
	"go-worker/message"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type Worker struct {
	ID           int
	RoutineType  string
	Active       bool
	ctx          context.Context
}

type WorkerPool struct {
	workers        []*Worker
	taskQueue      chan string
	numWorkers     int
	rdb            *redis.Client
	ctx            context.Context
	priority       string
	db             *sql.DB
	producer       Producer
}

func NewWorkerPool(numWorkers int, rdb *redis.Client, ctx context.Context, priority string, db *sql.DB, producer Producer) *WorkerPool {
	pool := &WorkerPool{
		workers:    make([]*Worker, numWorkers),
		taskQueue:  make(chan string, 100),
		numWorkers: numWorkers,
		rdb:        rdb,
		ctx:        ctx,
		priority:   priority,
		db:         db,
		producer:   producer,
	}

	for i := 0; i < numWorkers; i++ {
		worker := &Worker{
			ID:          i + 1,
			RoutineType: "Standard",
			Active:      true,
			ctx:         ctx,
		}
		pool.workers[i] = worker
		go pool.workerRoutine(worker)
	}

	return pool
}

func (wp *WorkerPool) workerRoutine(worker *Worker) {
	for worker.Active {
		select {
		case task := <-wp.taskQueue:
			wp.processMessage(task)
		case <-worker.ctx.Done():
			fmt.Printf("Worker %d stopped due to context cancellation\n", worker.ID)
			return
		}
	}
}

func (wp *WorkerPool) processMessage(msgString string) {
	fmt.Printf("Processing message: %s\n", msgString)
	var msg message.Message

	msgBytes := []byte(msgString)

	if err := json.Unmarshal(msgBytes, &msg); err != nil {
		log.Printf("Error unmarshalling message: %v\n", err)
		return
	}

	fmt.Printf("Processed message: %v\n", msg)

	if msg.EventTime.Before(time.Now().Add(-6 * time.Second)) {
		fmt.Printf("msg %v is older than 6 seconds\n", msg)
		msg.Status = "message unresolved"
		db.AddEventMessage(wp.db, 1, msg)
		return
	}

	msg.Status = fmt.Sprintf("message completed by %s priority worker", msg.Priority)
	if err := wp.producer.SendMessage(msg); err != nil {
		fmt.Printf("Failed to send message to producer: %v\n", err)
	}

	db.AddEventMessage(wp.db, 1, msg)
}

func (wp *WorkerPool) Run() {
	fmt.Println("Worker pool started")
	for {
		timeout := time.Duration(5) * time.Second
		message, err := wp.rdb.BRPop(wp.ctx, timeout, wp.priority).Result()
		if err != nil {
			fmt.Printf("Waiting for messages in Redis queue\n")
			continue
		}

		if len(message) < 2 {
			fmt.Println("Invalid message format")
			continue
		}

		wp.taskQueue <- message[1]
	}
}