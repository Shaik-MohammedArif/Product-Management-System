package cache

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func InitRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password by default
		DB:       0,                // Default DB
	})

	// Test Redis connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")
	return client
}

func SetCache(client *redis.Client, key string, value string, ttl time.Duration) error {
	err := client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}
	log.Printf("Cached key: %s", key)
	return nil
}

func GetCache(client *redis.Client, key string) (string, error) {
	val, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Cache miss
	} else if err != nil {
		return "", err
	}
	return val, nil
}
