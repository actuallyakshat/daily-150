package initialisers

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	currentEnv := os.Getenv("ENV")

	if currentEnv != "development" {
		redisURL := os.Getenv("REDIS_URL")
		if redisURL == "" {
			redisURL = "localhost:6379"
		}

		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			log.Printf("Failed to parse Redis URL: %v\n", err)
			return
		}

		if password := os.Getenv("REDIS_PASSWORD"); password != "" {
			opt.Password = password
		}

		RedisClient = redis.NewClient(opt)
	} else {
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
		})
	}

	ctx := context.Background()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Failed to connect to Redis: %v\n", err)
	} else {
		log.Println("Successfully connected to Redis")
	}
}
