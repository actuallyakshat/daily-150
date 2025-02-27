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

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		fmt.Printf("Failed to parse Redis URL: %v\n", err)
		return
	}

	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		opt.Password = password
	}

	RedisClient = redis.NewClient(opt)

	ctx := context.Background()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
	} else {
		fmt.Println("Successfully connected to Redis")
	}
}
