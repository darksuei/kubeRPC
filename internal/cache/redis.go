package cache

import (
	"context"
	"os"

	"github.com/go-redis/redis"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache() *RedisCache {
	return &RedisCache{
		client: redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		}),
	}
}

func (r *RedisCache) Ping() error {
	return r.client.Ping().Err()
}

func (r *RedisCache) HSet(key, field string, value interface{}) error {
	return r.client.WithContext(context.Background()).HSet(key, field, value).Err()
}

func (r *RedisCache) HGet(key, field string) (string, error) {
	val, err := r.client.HGet(key, field).Result()
	if err == redis.Nil {
		return "", ErrNotFound
	}
	return val, err
}

func (r *RedisCache) HGetAll(key string) (map[string]string, error) {
	return r.client.HGetAll(key).Result()
}

func (r *RedisCache) HDel(key string, fields ...string) (int64, error) {
	return r.client.HDel(key, fields...).Result()
}

func (r *RedisCache) Del(key string) (int64, error) {
	return r.client.Del(key).Result()
}

func (r *RedisCache) Keys(pattern string) ([]string, error) {
	return r.client.Keys(pattern).Result()
}

func (r *RedisCache) FlushDB() error {
	return r.client.FlushDB().Err()
}
