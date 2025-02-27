package initialisers

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	ctx := context.Background()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
	} else {
		fmt.Println("Successfully connected to Redis")
	}
}
