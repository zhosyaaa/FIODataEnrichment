package configs

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
)

var (
	RedisClient *redis.Client
)

func InitRedis() *redis.Client {
	redisAddr := os.Getenv("REDIS_ADDR")

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Error connecting to Redis: %v\n", err)
		return nil
	}

	RedisClient = client
	return client
}
