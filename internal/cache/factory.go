package cache

import (
	"log/slog"
	"os"
)

var Store Cache

func Init() {
	backend := os.Getenv("CACHE_TYPE")
	if backend == "memory" {
		Store = NewMemoryCache()
		slog.Info("cache: using in-memory store")
	} else {
		Store = NewRedisCache()
		slog.Info("cache: using Redis store")
	}

	if err := Store.Ping(); err != nil {
		slog.Error("cache: ping failed", "error", err)
	} else {
		slog.Info("cache: connected")
	}
}
