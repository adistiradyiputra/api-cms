package config

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     ENV.RedisURL,
		Password: ENV.RedisPassword,
		DB:       ENV.RedisDB,
	})
}
