package excache

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

// Cache represents safe for concurrent use passive expiring cache.
// Passively expires old records.
// Uses heap.Interface under the hood.
type Cache struct {
	m        *sync.Mutex
	capacity int

	expireQueue expireQueue
	cache       map[string]*record
}

// New returns an initialized cache instance
func New(capacity int) (*Cache, error) {
	if capacity <= 0 {
		return nil, fmt.Errorf("capacity can't be negative")
	}
	return &Cache{
		m:        &sync.Mutex{},
		capacity: capacity,

		expireQueue: make(expireQueue, 0, capacity),
		cache:       make(map[string]*record, capacity),
	}, nil
}

// Get returns (value, true) or (nil, false) for a given key.
// Resets TTL.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.m.Lock()
	defer c.m.Unlock()

	c.expire()

	e, ok := c.cache[key]
	if !ok {
		return nil, false
	}

	c.expireQueue.update(e, e.value, e.ttl, time.Now().Add(e.ttl).UnixNano())

	return e.value, true
}

// Put inserts new record into the cache
func (c *Cache) Put(key string, value interface{}, ttl time.Duration) {
	c.m.Lock()
	defer c.m.Unlock()

	c.expire()

	if len(c.cache) >= c.capacity {
		if c.expireQueue[0].key != key {
			r := c.expireQueue.Pop().(*record)
			delete(c.cache, r.key)
		}
	}

	r, ok := c.cache[key]
	if !ok {
		r = &record{
			key:             key,
			value:           value,
			ttl:             ttl,
			expireTimeStamp: time.Now().Add(ttl).UnixNano(),
		}

		heap.Push(&c.expireQueue, r)
		c.cache[key] = r

		return
	}

	c.expireQueue.update(r, value, ttl, time.Now().Add(ttl).UnixNano())
}

// Delete removes the record associated with the specified key from the cache
func (c *Cache) Delete(key string) {
	c.m.Lock()
	defer c.m.Unlock()

	c.expire()

	r, ok := c.cache[key]
	if !ok {
		return
	}

	heap.Remove(&c.expireQueue, r.index)
	delete(c.cache, key)
}

// Clear removes all saved records
func (c *Cache) Clear() {
	c.m.Lock()
	defer c.m.Unlock()

	c.expireQueue = make(expireQueue, 0, c.capacity)
	c.cache = make(map[string]*record, c.capacity)
}

// Len returns the number of records in the cache
func (c *Cache) Len() int {
	c.m.Lock()
	defer c.m.Unlock()

	c.expire()

	return len(c.cache)
}

// Expire removes old records
func (c *Cache) Expire() {
	c.m.Lock()
	defer c.m.Unlock()

	c.expire()
}

func (c *Cache) expire() {
	now := time.Now().UnixNano()

	for c.expireQueue.Len() > 0 && now >= c.expireQueue[0].expireTimeStamp {
		r := c.expireQueue.Pop().(*record)
		delete(c.cache, r.key)
	}
}

type record struct {
	key   string
	value interface{}

	ttl             time.Duration
	expireTimeStamp int64

	index int
}

type expireQueue []*record

func (q expireQueue) Len() int { return len(q) }

func (q expireQueue) Less(i, j int) bool {
	return q[i].expireTimeStamp > q[j].expireTimeStamp
}

func (q expireQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *expireQueue) Push(x interface{}) {
	n := len(*q)
	item := x.(*record)
	item.index = n
	*q = append(*q, item)
}

func (q *expireQueue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*q = old[0 : n-1]
	return item
}

func (q *expireQueue) update(oldRecord *record, value interface{}, ttl time.Duration, expireTimeStamp int64) {
	oldRecord.value = value
	oldRecord.ttl = ttl
	oldRecord.expireTimeStamp = expireTimeStamp
	heap.Fix(q, oldRecord.index)
}
