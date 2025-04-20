package redisclient

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func GetCache(client *redis.Client, key string) (string, error) {
	return client.Get(ctx, key).Result()
}

func SetCache(client *redis.Client, key string, value string, expiration time.Duration) error {
	return client.Set(ctx, key, value, expiration).Err()
}
