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
	routineType    struct {
		Working time.Duration
		Idling  time.Duration
	}
}

func NewWorkerPool(numWorkers int, rdb *redis.Client, ctx context.Context, priority string, db *sql.DB, producer Producer, routineType string) *WorkerPool {
	routinePeriods := map[string]struct{
		Working time.Duration
		Idling  time.Duration
	}{
		"Standard":    {20 * time.Second, 5 * time.Second},
		"Intermittent": {5 * time.Second, 5 * time.Second},
		"Concentrated": {60 * time.Second, 60 * time.Second},
	}

	if _, exists := routinePeriods[routineType]; !exists {
		panic("Invalid routine type")
	}

	pool := &WorkerPool{
		workers:    make([]*Worker, numWorkers),
		taskQueue:  make(chan string, 100),
		numWorkers: numWorkers,
		rdb:        rdb,
		ctx:        ctx,
		priority:   priority,
		db:         db,
		producer:   producer,
		routineType:  routinePeriods[routineType],
	}
	for i := 0; i < 2; i++ {
		worker := &Worker{
			ID:          i + 1,
			RoutineType: "Standard",
			Active:      true,
			ctx:         ctx,
		}
		pool.workers[i] = worker
		go pool.workerRoutine(worker)
	}

	// for i := 2; i < numWorkers; i++ {
	// 	worker := &Worker{
	// 		ID:          i + 1,
	// 		RoutineType: "Standard",
	// 		Active:      true,
	// 		ctx:         ctx,
	// 	}
	// 	pool.workers[i] = worker
	// 	time.Sleep(10 * time.Second)
	// 	go pool.workerRoutine(worker)
	// }

	return pool
}

func (wp *WorkerPool) workerRoutine(worker *Worker) {
	time.Sleep(wp.routineType.Working)
	
	workingTimer := time.NewTimer(wp.routineType.Idling)
	defer workingTimer.Stop()

	for worker.Active {        

		select {
		case task := <-wp.taskQueue:
			wp.processMessage(task)
		case <-worker.ctx.Done():
			fmt.Printf("Worker %d stopped due to context cancellation\n", worker.ID)
			return
		case <-workingTimer.C:
			time.Sleep(wp.routineType.Working)
			workingTimer.Reset(wp.routineType.Idling)
		}
	}
}

func (wp *WorkerPool) processMessage(msgString string) {
	var msg message.Message

	msgBytes := []byte(msgString)

	if err := json.Unmarshal(msgBytes, &msg); err != nil {
		log.Printf("Error unmarshalling message: %v\n", err)
		return
	}
	var timeToComplete int

	switch msg.Priority {
	case "High":
			timeToComplete = -5
	case "Medium":
			timeToComplete = -10
	case "Low":
			timeToComplete = -15
	default:
			timeToComplete = 0 // Default value if priority doesn't match any case
	}
	
	if msg.EventTime.Before(time.Now().Add(time.Duration(timeToComplete) * time.Second)) {
		fmt.Printf("!!!!MSG %v IS OLDER THAN %d SECONDS!!!!\n", msg.Status, timeToComplete)
		msg.Status = "message unresolved"
		db.AddEventMessage(wp.db, 1, msg)
		return
	}

	msg.Status = fmt.Sprintf("message completed by %s priority worker", wp.priority)
	fmt.Printf("msg %v is completed by %s worker\n", msg.EventType, wp.priority)
	if err := wp.producer.SendMessage(msg); err != nil {
		fmt.Printf("Failed to send message to producer: %v\n", err)
	}

	db.AddEventMessage(wp.db, 1, msg)
}

func (wp *WorkerPool) Run() {
	fmt.Println("Worker pool started")
	for {
		time := wp.routineType.Working
		message, err := wp.rdb.BRPop(wp.ctx, time, wp.priority).Result()
		if err != nil {
			continue
		}

		if len(message) < 2 {
			fmt.Println("Invalid message format")
			continue
		}

		wp.taskQueue <- message[1]
	}
}