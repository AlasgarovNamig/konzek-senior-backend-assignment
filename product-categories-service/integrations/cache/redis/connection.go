package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
)

var redisConnection *redis.Client

func SetupRedisConnection() {

	redisConnection = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: "",
		DB:       0,
	})

	_, err := redisConnection.Ping(context.Background()).Result()
	if err != nil {
		panic("failed to create a connection to Redis")
	}
}

func CloseRedisConnection() {
	err := redisConnection.Close()
	if err != nil {
		panic("failed to close connection to Redis")
	}
}
func GetRedisConnection() *redis.Client {
	return redisConnection
}
