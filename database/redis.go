package database

import (
	"context"
	"strings"

	"github.com/olabanji12-ojo/church-backend/config"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var RedisClient *redis.Client

// ConnectRedis initializes the Redis connection
func ConnectRedis() *redis.Client {
	redisURL := config.GetEnv("REDIS_URI", "redis://localhost:6379")
	var opt *redis.Options
	var err error

	if redisURL != "" && (strings.HasPrefix(redisURL, "redis://") || strings.HasPrefix(redisURL, "rediss://")) {
		opt, err = redis.ParseURL(redisURL)
		if err != nil {
			logrus.Fatalf("❌ Failed to parse Redis URL: %v", err)
		}
	} else {
		opt = &redis.Options{
			Addr: redisURL,
		}
	}

	client := redis.NewClient(opt)

	// Ping Redis
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		logrus.Fatalf("❌ Failed to connect to Redis: %v", err)
	}

	logrus.Println("✅ Successfully connected to Redis")
	RedisClient = client
	return client
}
