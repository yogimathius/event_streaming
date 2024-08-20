package main

import (
	"context"
	"fmt"
	db "go-worker/database"
	"go-worker/kafka"
	"go-worker/message"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

type Producer interface {
	SendMessage(message.Message) error
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	ctx := context.Background()

	fmt.Printf("main Context: %v\n", ctx)
	fmt.Printf("main Redis client: %v\n", rdb)

	priority := []string{os.Getenv("PRIORITY")}

	database, err := db.InitDb()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer database.Close()

	producer, err := kafka.NewProducer([]string{"kafka:29092"})
	if err != nil {
		log.Fatalf("Kafka producer initialization failed: %v", err)
	}
	defer producer.Close()
	
	pool := NewWorkerPool(5, rdb, ctx, priority[0], database, producer)
	pool.Run()
}

