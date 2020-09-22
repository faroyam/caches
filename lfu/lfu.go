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

	er, ok := c.cache[key]
	if !ok {
		return nil, false
	}

	currentRecord := er.Value.(record)
	currentNode := currentRecord.node.Value.(node)
	newFrequency := currentNode.frequency + 1
	nextNode := currentRecord.node.Next()

	if nextNode == nil || nextNode.Value.(node).frequency != newFrequency {
		nextNode = c.nodes.InsertAfter(newNode(newFrequency, list.New()), currentRecord.node)
	}

	c.removeRecord(er, false)

	er = nextNode.Value.(node).records.PushFront(newRecord(nextNode, key, currentRecord.value))
	c.cache[key] = er

	return currentRecord.value, true
}

// Put inserts new record into the cache
func (c *Cache) Put(key string, value interface{}) {
	c.m.Lock()
	defer c.m.Unlock()

	if len(c.cache) >= c.capacity {
		er, _, _ := c.lfu()
		if er.Value.(record).key != key {
			c.removeRecord(er, true)
		}
	}

	if er, ok := c.cache[key]; ok {
		currentRecord := er.Value.(record)
		currentNode := currentRecord.node.Value.(node)

		newFrequency := currentNode.frequency + 1
		nextNode := currentRecord.node.Next()

		if nextNode == nil || nextNode.Value.(node).frequency != newFrequency {
			nextNode = c.nodes.InsertAfter(newNode(newFrequency, list.New()), currentRecord.node)
		}

		c.removeRecord(er, false)

		er = nextNode.Value.(node).records.PushFront(newRecord(nextNode, key, value))
		c.cache[key] = er

		return
	}

	currentNode := c.nodes.Front()

	if currentNode == nil || currentNode.Value.(node).frequency != 1 {
		currentNode = c.nodes.PushFront(newNode(1, list.New()))
	}

	er := currentNode.Value.(node).records.PushFront(newRecord(currentNode, key, value))
	c.cache[key] = er
}

// LFU returns one of keys (key, frequency, true) that has been touched fewer times.
// Returns ("", 0, false) if there are no keys in the cache.
// Does not "use" record i.e. returning record will remain untouched.
func (c *Cache) LFU() (string, int64, bool) {
	c.m.Lock()
	defer c.m.Unlock()

	if er, frequency, ok := c.lfu(); ok {
		return er.Value.(record).key, frequency, true
	}
	return "", 0, false
}

// Delete removes the record associated with the specified key from the cache
func (c *Cache) Delete(key string) {
	c.m.Lock()
	defer c.m.Unlock()

	er, ok := c.cache[key]
	if !ok {
		return
	}

	c.removeRecord(er, true)
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
	if en := c.nodes.Front(); en != nil {
		node := en.Value.(node)
		if er := node.records.Front(); er != nil {
			return er, node.frequency, true
		}
	}
	return nil, 0, false
}

func (c *Cache) removeRecord(er *list.Element, removeFromCache bool) record {
	currentRecord := er.Value.(record)
	currentNode := currentRecord.node.Value.(node)

	removedRecord := currentNode.records.Remove(er).(record)
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
