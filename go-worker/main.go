package main

import (
	"context"
	db "go-worker/database"
	"go-worker/kafka"
	"go-worker/pool"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	ctx := context.Background()

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
	
	priority := []string{os.Getenv("PRIORITY")}
	routine := os.Getenv("ROUTINE")

	pool := pool.NewWorkerPool(7, rdb, ctx, priority[0], database, producer, routine)
	pool.Run()
}

