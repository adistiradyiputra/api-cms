package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	// Check if Redis URL is set
	if ENV.RedisURL == "" {
		log.Println("Warning: Redis URL not set, skipping Redis connection")
		return
	}

	RDB = redis.NewClient(&redis.Options{
		Addr:     ENV.RedisURL,
		Password: ENV.RedisPassword,
		DB:       ENV.RedisDB,
	})

	// Test connection
	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
	} else {
		log.Println("Redis connected successfully")
	}
}
