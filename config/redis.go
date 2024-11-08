package config

import (
	"log"
	"os"

	"github.com/go-redis/redis"
)

var Rdb *redis.Client

func InitRedisClient() {
	log.Println("Initializing redis client..")

	Rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}
