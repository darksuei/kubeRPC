package config

import (
	"fmt"
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

	flushRedis()
}

func flushRedis() {
	flush := os.Getenv("FLUSH_DATABASE")

	if flush == "true" {
		err := Rdb.FlushDB().Err()
		if err != nil {
			fmt.Printf("Error flushing Redis: %v\n", err)
		} else {
			fmt.Println("Redis flushed successfully")
		}
	}
}
