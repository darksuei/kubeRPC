package cache

import (
	"log"
	"os"
)

var Store Cache

func Init() {
	backend := os.Getenv("CACHE_TYPE")
	if backend == "memory" {
		Store = NewMemoryCache()
		log.Println("Cache: using in-memory store")
	} else {
		Store = NewRedisCache()
		log.Println("Cache: using Redis store")
	}

	if err := Store.FlushDB(); err != nil {
		log.Printf("Cache: flush failed: %v", err)
	} else {
		log.Println("Cache: flushed successfully")
	}
}
