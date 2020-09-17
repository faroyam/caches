package lru

import (
	"container/list"
	"sync"
)

type Cache struct {
	m        *sync.Mutex
	capacity int

	records *list.List
	cache   map[string]*list.Element
}

// New returns an initialized cache instance
func New(capacity int) *Cache {
	return &Cache{
		m:        &sync.Mutex{},
		capacity: capacity,
		records:  list.New(),
		cache:    make(map[string]*list.Element, capacity),
	}
}

// Get returns (value, true) or (nil, false) for the given key
func (c *Cache) Get(key string) (interface{}, bool) {
	c.m.Lock()
	defer c.m.Unlock()

	e, ok := c.cache[key]
	if !ok {
		return nil, false
	}

	c.records.MoveToFront(e)
	return e.Value.(record).value, true
}

// Put inserts new record in the cache
func (c *Cache) Put(key string, value interface{}) {
	c.m.Lock()
	defer c.m.Unlock()

	if len(c.cache) >= c.capacity {
		e := c.records.Remove(c.records.Back()).(record)
		delete(c.cache, e.key)
	}

	if e, ok := c.cache[key]; ok {
		c.records.Remove(e)
	}

	e := c.records.PushFront(record{
		key:   key,
		value: value,
	})
	c.cache[key] = e
}

// Delete removes a record from the cache
func (c *Cache) Delete(key string) {
	c.m.Lock()
	defer c.m.Unlock()

	e, ok := c.cache[key]
	if !ok {
		return
	}
	r := c.records.Remove(e).(record)
	delete(c.cache, r.key)
}

// Clear removes all saved records
func (c *Cache) Clear() {
	c.m.Lock()
	defer c.m.Unlock()

	c.cache = make(map[string]*list.Element, c.capacity)
	c.records = list.New()
}

// Len returns the number of records in the cache.
// The complexity is O(1).
func (c *Cache) Len() int {
	return len(c.cache)
}

type record struct {
	key   string
	value interface{}
}
