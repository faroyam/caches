package lfu

import (
	"container/list"
	"fmt"
	"sync"
)

// Cache represents safe for concurrent use Least Frequently Used cache
type Cache struct {
	m        *sync.Mutex
	capacity int

	nodes *list.List
	cache map[string]*list.Element
}

// New returns an initialized cache instance
func New(capacity int) (*Cache, error) {
	if capacity <= 0 {
		return nil, fmt.Errorf("capacity can't be negative")
	}
	return &Cache{
		m:        &sync.Mutex{},
		capacity: capacity,
		nodes:    list.New(),
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

	currentRecord := e.Value.(record)
	currentNode := currentRecord.node.Value.(node)
	newFrequency := currentNode.frequency + 1
	nextNode := currentRecord.node.Next()

	if nextNode == nil || nextNode.Value.(node).frequency != newFrequency {
		nextNode = c.nodes.InsertBefore(newNode(newFrequency, list.New()), currentRecord.node)
	}

	c.removeRecord(e, false)

	e = nextNode.Value.(node).records.PushBack(newRecord(nextNode, key, currentRecord.value))
	c.cache[key] = e

	return currentRecord.value, true
}

// Put inserts new record into the cache
func (c *Cache) Put(key string, value interface{}) {
	c.m.Lock()
	defer c.m.Unlock()

	if len(c.cache) >= c.capacity {
		e, _, _ := c.lfu()
		if e.Value.(record).key != key {
			c.removeRecord(e, true)
		}
	}

	if e, ok := c.cache[key]; ok {
		currentRecord := e.Value.(record)
		currentNode := currentRecord.node.Value.(node)

		newFrequency := currentNode.frequency + 1
		nextNode := currentRecord.node.Next()

		if nextNode == nil || nextNode.Value.(node).frequency != newFrequency {
			nextNode = c.nodes.InsertBefore(newNode(newFrequency, list.New()), currentRecord.node)
		}

		c.removeRecord(e, false)

		e = nextNode.Value.(node).records.PushBack(newRecord(nextNode, key, value))
		c.cache[key] = e

		return
	}

	frontNode := c.nodes.Front()

	if frontNode == nil || frontNode.Value.(node).frequency != 1 {
		frontNode = c.nodes.PushBack(newNode(1, list.New()))
	}

	er := frontNode.Value.(node).records.PushBack(newRecord(frontNode, key, value))
	c.cache[key] = er
}

// LFU returns one of keys (key, frequency, true) that has been touched fewer times.
// Returns ("", 0, false) if there are no keys in the cache.
// Does not "use" record i.e. returning record will remain untouched.
func (c *Cache) LFU() (string, int64, bool) {
	c.m.Lock()
	defer c.m.Unlock()

	if e, frequency, ok := c.lfu(); ok {
		return e.Value.(record).key, frequency, true
	}
	return "", 0, false
}

// Delete removes the record associated with the specified key from the cache
func (c *Cache) Delete(key string) {
	c.m.Lock()
	defer c.m.Unlock()

	e, ok := c.cache[key]
	if !ok {
		return
	}

	c.removeRecord(e, true)
}

// Clear removes all saved records
func (c *Cache) Clear() {
	c.m.Lock()
	defer c.m.Unlock()

	c.cache = make(map[string]*list.Element, c.capacity)
	c.nodes = list.New()
}

// Len returns the number of records in the cache
func (c *Cache) Len() int {
	return len(c.cache)
}

func (c *Cache) lfu() (*list.Element, int64, bool) {
	if frontNode := c.nodes.Back(); frontNode != nil {
		node := frontNode.Value.(node)
		if e := node.records.Back(); e != nil {
			return e, node.frequency, true
		}
	}
	return nil, 0, false
}

func (c *Cache) removeRecord(e *list.Element, removeFromCache bool) record {
	currentRecord := e.Value.(record)
	currentNode := currentRecord.node.Value.(node)

	removedRecord := currentNode.records.Remove(e).(record)
	if currentNode.records.Len() == 0 {
		c.nodes.Remove(currentRecord.node)
	}

	if removeFromCache {
		delete(c.cache, removedRecord.key)
	}

	return removedRecord
}

type node struct {
	frequency int64
	records   *list.List
}

func newNode(frequency int64, records *list.List) node {
	return node{
		frequency: frequency,
		records:   records,
	}
}

type record struct {
	node  *list.Element
	key   string
	value interface{}
}

func newRecord(node *list.Element, key string, value interface{}) record {
	return record{
		node:  node,
		key:   key,
		value: value,
	}
}
