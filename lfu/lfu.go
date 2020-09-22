package lfu

import (
	"container/list"
	"fmt"
	"sync"
)

// Cache represents safe for concurrent use LFU cache
type Cache struct {
	m        *sync.Mutex
	capacity int

	nodes *list.List
	cache map[string]*list.Element
}

type node struct {
	frequency int64
	records   *list.List
}

type record struct {
	node  *list.Element
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
		nextNode = c.nodes.InsertAfter(newNode(newFrequency, list.New()), currentRecord.node)
	}

	currentNode.records.Remove(e)
	if currentNode.records.Len() == 0 {
		c.nodes.Remove(currentRecord.node)
	}

	e = nextNode.Value.(node).records.PushFront(newRecord(nextNode, key, currentRecord.value))

	c.cache[key] = e

	return currentRecord.value, true
}

// Put inserts new record into the cache
func (c *Cache) Put(key string, value interface{}) {
	c.m.Lock()
	defer c.m.Unlock()

	if len(c.cache) >= c.capacity {
		leastFrequencyNode := c.nodes.Front()
		leastFrequencyRecords := leastFrequencyNode.Value.(node).records
		r := leastFrequencyRecords.Front()

		if r.Value.(record).key != key {
			e := leastFrequencyRecords.Remove(r)
			if leastFrequencyRecords.Len() == 0 {
				c.nodes.Remove(leastFrequencyNode)
			}

			delete(c.cache, e.(record).key)
		}
	}

	if e, ok := c.cache[key]; ok {
		currentRecord := e.Value.(record)
		currentNode := currentRecord.node.Value.(node)

		newFrequency := currentNode.frequency + 1
		nextNode := currentRecord.node.Next()

		if nextNode == nil || nextNode.Value.(node).frequency != newFrequency {
			nextNode = c.nodes.InsertAfter(newNode(newFrequency, list.New()), currentRecord.node)
		}

		currentNode.records.Remove(e)
		if currentNode.records.Len() == 0 {
			c.nodes.Remove(currentRecord.node)
		}

		e = nextNode.Value.(node).records.PushFront(newRecord(nextNode, key, value))

		c.cache[key] = e

		return
	}

	insertNode := c.nodes.Front()

	if insertNode == nil || insertNode.Value.(node).frequency != 1 {
		insertNode = c.nodes.PushFront(newNode(1, list.New()))
	}

	e := insertNode.Value.(node).records.PushFront(newRecord(insertNode, key, value))

	c.cache[key] = e
}

// LeastFrequentlyUsed returns one of keys (key, frequency, true) that has been touched fewer times.
// Returns ("", 0, false) if there are no keys in the cache.
// Does not "use" record i.e. returning record will remain untouched.
func (c *Cache) LeastFrequentlyUsed() (string, int64, bool) {
	c.m.Lock()
	defer c.m.Unlock()

	if e := c.nodes.Front(); e != nil {
		node := e.Value.(node)
		if r := node.records.Front(); r != nil {
			return r.Value.(record).key, node.frequency, true
		}
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

	currentRecord := e.Value.(record)
	currentNode := currentRecord.node.Value.(node)

	currentNode.records.Remove(e)
	if currentNode.records.Len() == 0 {
		c.nodes.Remove(currentRecord.node)
	}

	delete(c.cache, currentRecord.key)
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

func newNode(frequency int64, records *list.List) node {
	return node{
		frequency: frequency,
		records:   records,
	}
}

func newRecord(node *list.Element, key string, value interface{}) record {
	return record{
		node:  node,
		key:   key,
		value: value,
	}
}
