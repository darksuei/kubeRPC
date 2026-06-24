package cache

import (
	"fmt"
	"path/filepath"
	"sync"
)

type MemoryCache struct {
	mu   sync.RWMutex
	data map[string]map[string]string
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{data: make(map[string]map[string]string)}
}

func (m *MemoryCache) Ping() error {
	return nil
}

func (m *MemoryCache) HSet(key, field string, value interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data[key] == nil {
		m.data[key] = make(map[string]string)
	}
	m.data[key][field] = stringify(value)
	return nil
}

func (m *MemoryCache) HGet(key, field string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	hash, ok := m.data[key]
	if !ok {
		return "", ErrNotFound
	}
	val, ok := hash[field]
	if !ok {
		return "", ErrNotFound
	}
	return val, nil
}

func (m *MemoryCache) HGetAll(key string) (map[string]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	hash := m.data[key]
	result := make(map[string]string, len(hash))
	for k, v := range hash {
		result[k] = v
	}
	return result, nil
}

func (m *MemoryCache) HDel(key string, fields ...string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var deleted int64
	for _, f := range fields {
		if _, ok := m.data[key][f]; ok {
			delete(m.data[key], f)
			deleted++
		}
	}
	return deleted, nil
}

func (m *MemoryCache) Del(key string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.data[key]; !ok {
		return 0, nil
	}
	delete(m.data, key)
	return 1, nil
}

func (m *MemoryCache) Keys(pattern string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var keys []string
	for k := range m.data {
		if matched, _ := filepath.Match(pattern, k); matched {
			keys = append(keys, k)
		}
	}
	return keys, nil
}

func (m *MemoryCache) FlushDB() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = make(map[string]map[string]string)
	return nil
}

func stringify(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case int:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
