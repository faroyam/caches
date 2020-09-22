package lru

import (
	"container/list"
	"fmt"
	"sync"
)

// Cache represents safe for concurrent use LRU cache
type Cache struct {
	m        *sync.Mutex
	capacity int

	records *list.List
	cache   map[string]*list.Element
}

type record struct {
	key   string
	value interface{}
}

// New returns an initialized cache instance
func New(capacity int) (*Cache, error) {
	if capacity <= 0 {
		return nil, fmt.Errorf("capacity can't be negative")
	}
	return &Cache{
		m:        &sync.Mutex{},
		capacity: capacity,
		records:  list.New(),
		cache:    make(map[string]*list.Element, capacity),
	}, nil
}

// Get returns (value, true) or (nil, false) for a given key
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

// Put inserts a new record into the cache
func (c *Cache) Put(key string, value interface{}) {
	c.m.Lock()
	defer c.m.Unlock()

	if len(c.cache) >= c.capacity {
		e := c.records.Remove(c.records.Back())
		delete(c.cache, e.(record).key)
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

// LeastRecentlyUsed returns (key, true) that was not touched for the longest time.
// Returns ("", false) if there are no keys in the cache.
// Does not "use" record i.e. returning record will remain untouched.
func (c *Cache) LeastRecentlyUsed() (string, bool) {
	c.m.Lock()
	defer c.m.Unlock()

	if e := c.records.Back(); e != nil {
		return e.Value.(record).key, true
	}
	return "", false
}

// Delete removes the record associated with the specified key from the cache
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

// Len returns the number of records in the cache
func (c *Cache) Len() int {
	return len(c.cache)
}
