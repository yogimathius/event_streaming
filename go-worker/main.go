package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

var ctx = context.Background()

func main() {
	priority := []string{os.Getenv("PRIORITY")}

	for {
		queue := priority[0]

		message, err := rdb.BRPop(ctx, 0, queue).Result()
		if err != nil {
			log.Printf("Failed to get message from Redis queue: %v\n", err)
			time.Sleep(time.Second) 
			continue
		}

		processMessage(message[1])
	}
}

func processMessage(message string) {
	fmt.Printf("Processing message: %s\n", message)
	// Add your processing logic here
}
