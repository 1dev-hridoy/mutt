package config

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func MustConnectRedis() {
	url := os.Getenv("UPSTASH_REDIS_URL")
	if url == "" {
		panic("PANIC :: UPSTASH_REDIS_URL is not set.")
	}

	opts, err := redis.ParseURL(url)
	if err != nil {
		panic("PANIC :: Failed to parse Redis URL: " + err.Error())
	}

	RDB = redis.NewClient(opts)

	if err := RDB.Ping(context.Background()).Err(); err != nil {
		panic("PANIC :: Failed to connect to Redis: " + err.Error())
	}
}
