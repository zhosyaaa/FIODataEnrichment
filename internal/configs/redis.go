package configs

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"os"
)

var redisClient *redis.Client

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
		log.Error().Err(err).Msg("Error connecting to Redis")
		return nil
	}
	redisClient = client
	log.Info().Msg("Connected to Redis")
	return client
}

func GetFromCache(key string) (string, error) {
	val, err := redisClient.Get(context.Background(), key).Result()
	if err != nil {
		log.Error().Err(err).Msg("Error getting data from Redis cache")
		return "", err
	}
	return val, nil
}

func SetInCache(key string, value string) error {
	err := redisClient.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		log.Error().Err(err).Msg("Error setting data in Redis cache")
		return err
	}
	return nil
}
