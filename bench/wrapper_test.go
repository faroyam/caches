package bench_test

import (
	"sync"
	"time"

	"github.com/faroyam/caches/excache"
)

// Map represents an interface to Go built in map type
type Map struct {
	m     *sync.Mutex
	cache map[string]interface{}
}

// NewMap returns an initialized cache instance
func NewMap(capacity int) *Map {
	return &Map{
		m:     &sync.Mutex{},
		cache: make(map[string]interface{}, capacity),
	}
}

// Get returns (value, true) or (nil, false) for the given key
func (c *Map) Get(key string) (interface{}, bool) {
	c.m.Lock()
	defer c.m.Unlock()

	value, ok := c.cache[key]
	if !ok {
		return nil, false
	}

	return value, true
}

// Put inserts new record in the cache
func (c *Map) Put(key string, value interface{}) {
	c.m.Lock()
	defer c.m.Unlock()

	c.cache[key] = value
}

// ExpiringMap represents an interface to expiring map
type ExpiringMap struct {
	cache *excache.Cache
}

// NewExpiringMap returns an initialized cache instance
func NewExpiringMap(capacity int) (*ExpiringMap, error) {
	cache, err := excache.New(capacity)
	if err != nil {
		return nil, err
	}
	return &ExpiringMap{
		cache: cache,
	}, nil
}

// Get returns (value, true) or (nil, false) for the given key
func (c *ExpiringMap) Get(key string) (interface{}, bool) {
	return c.cache.Get(key)
}

// Put inserts new record in the cache
func (c *ExpiringMap) Put(key string, value interface{}) {
	c.cache.Put(key, value, time.Second)
}
