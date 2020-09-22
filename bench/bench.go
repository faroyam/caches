package bench

import (
	"sync"
)

// Map represents an interface to Go built in map type
type Map struct {
	m     *sync.Mutex
	cache map[string]interface{}
}

// New returns an initialized cache instance
func New(capacity int) *Map {
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
