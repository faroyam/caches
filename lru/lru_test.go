package lru_test

import (
	"testing"

	"github.com/faroyam/caches/lru"
)

func TestCache_New(t *testing.T) {
	_, err := lru.New(0)
	if err == nil {
		t.Errorf("expected error")
	}

	_, err = lru.New(-1)
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestCache_Put(t *testing.T) {
	cache, _ := lru.New(1)
	cache.Put("key", "value")

	if cache.Len() != 1 {
		t.Errorf("cache len %v, want %v", cache.Len(), 1)
	}
}

func TestCache_Get(t *testing.T) {
	cache, _ := lru.New(1)
	cache.Put("key", "value")

	if value, ok := cache.Get("key"); !ok || value != "value" {
		t.Errorf("cached value %v, want %v", value, "value")
	}

	if value, ok := cache.Get("non-existing-key"); ok {
		t.Errorf("cached value %v, want %v", value, "nil")
	}
}

func TestCache_LRU(t *testing.T) {
	cache, _ := lru.New(1)
	if key, ok := cache.LRU(); ok {
		t.Errorf("lru %v, want %v", key, "")
	}

	cache.Put("1", "1")
	cache.Put("2", "2")
	cache.Put("2", "2`")

	if key, _ := cache.LRU(); key != "2" {
		t.Errorf("lru %v, want %v", key, "2")
	}
}

func TestCache_Delete(t *testing.T) {
	cache, _ := lru.New(1)
	cache.Put("key", "value")

	cache.Delete("key")
	cache.Delete("non-existing-key")

	if cache.Len() != 0 {
		t.Errorf("cache len %v, want %v", cache.Len(), 0)
	}

	if key, ok := cache.LRU(); ok {
		t.Errorf("lru %v, want %v", key, "''")
	}
}

func TestCache_Clear(t *testing.T) {
	cache, _ := lru.New(10)
	cache.Put("key1", "value1")
	cache.Put("key2", "value2")
	cache.Clear()

	if cache.Len() != 0 {
		t.Errorf("cache len %v, want %v", cache.Len(), 0)
	}
}

func TestReplace(t *testing.T) {
	cache, _ := lru.New(10)
	cache.Put("key", "value1")

	if value, ok := cache.Get("key"); !ok || value != "value1" {
		t.Errorf("cached value %v, want %v", value, "value1")
	}

	cache.Put("key", "value2")

	if value, ok := cache.Get("key"); !ok || value != "value2" {
		t.Errorf("cached value %v, want %v", value, "value2")
	}

	if cache.Len() != 1 {
		t.Errorf("cache len %v, want %v", cache.Len(), 1)
	}
}

func TestPutMoreThanCap(t *testing.T) {
	cache, _ := lru.New(1)
	key1 := "key1"
	value1 := "value1"
	key2 := "key2"
	value2 := "value2"

	cache.Put(key1, value1)
	if value, ok := cache.Get(key1); !ok || value != value1 {
		t.Errorf("cached value %v, want %v", value, value1)
	}

	cache.Put(key2, value2)
	if value, ok := cache.Get(key2); !ok || value != value2 {
		t.Errorf("cached value %v, want %v", value, value2)
	}

	if value, ok := cache.Get(key1); ok {
		t.Errorf("cached value %v, want %v", value, nil)
	}

	if cache.Len() != 1 {
		t.Errorf("cache len %v, want %v", cache.Len(), 1)
	}
}

func Test(t *testing.T) {
	cache, _ := lru.New(2)

	cache.Put("1", 1)
	cache.Put("2", 2)

	if value, ok := cache.Get("1"); !ok || value != 1 {
		t.Errorf("cached value %v, want %v", value, 1)
	}

	if key, _ := cache.LRU(); key != "2" {
		t.Errorf("lru key %v, want %v", key, "2")
	}

	// keys: 1 -> 2

	cache.Put("3", 3)

	if value, ok := cache.Get("2"); ok {
		t.Errorf("cached value %v, want %v", value, nil)
	}

	if key, _ := cache.LRU(); key != "1" {
		t.Errorf("lru key %v, want %v", key, "1")
	}

	// keys: 3 -> 1

	cache.Put("4", 4)

	if value, ok := cache.Get("1"); ok {
		t.Errorf("cached value %v, want %v", value, nil)
	}

	if key, _ := cache.LRU(); key != "3" {
		t.Errorf("lru key %v, want %v", key, "3")
	}

	// keys: 4 -> 3

	if value, ok := cache.Get("4"); !ok || value != 4 {
		t.Errorf("cached value %v, want %v", value, 4)
	}

	if value, ok := cache.Get("3"); !ok || value != 3 {
		t.Errorf("cached value %v, want %v", value, 3)
	}

	if key, _ := cache.LRU(); key != "4" {
		t.Errorf("lru key %v, want %v", key, "4")
	}

	// keys: 3 -> 4
}
