package cache

import "errors"

var ErrNotFound = errors.New("not found")

type Cache interface {
	Ping() error
	HSet(key, field string, value interface{}) error
	HGet(key, field string) (string, error)
	HGetAll(key string) (map[string]string, error)
	HDel(key string, fields ...string) (int64, error)
	Del(key string) (int64, error)
	Keys(pattern string) ([]string, error)
	FlushDB() error
}
