package config

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis() {

	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	dbString := os.Getenv("REDIS_DB")

	if host == "" {
		host = "localhost"
	}

	if port == "" {
		port = "6379"
	}

	db := 0

	if dbString != "" {
		if parsedDB, err := strconv.Atoi(dbString); err == nil {
			db = parsedDB
		}
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		log.Println("⚠️ Redis is not available:", err)
		return
	}

	log.Println("✅ Redis Connected Successfully")
}
