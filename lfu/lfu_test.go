package lfu_test

import (
	"testing"

	"github.com/faroyam/caches/lfu"
)

func TestCache_New(t *testing.T) {
	_, err := lfu.New(0)
	if err == nil {
		t.Errorf("expected error")
	}

	_, err = lfu.New(-1)
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestCache_Put(t *testing.T) {
	cache, _ := lfu.New(1)
	cache.Put("key", "value")

	if cache.Len() != 1 {
		t.Errorf("cache len %v, want %v", cache.Len(), 1)
	}
}

func TestCache_Get(t *testing.T) {
	cache, _ := lfu.New(1)
	cache.Put("key", "value")

	if value, ok := cache.Get("key"); !ok || value != "value" {
		t.Errorf("cached value %v, want %v", value, "value")
	}

	if value, ok := cache.Get("non-existing-key"); ok {
		t.Errorf("cached value %v, want %v", value, "nil")
	}
}

func TestCache_LFU(t *testing.T) {
	cache, _ := lfu.New(1)
	if key, _, ok := cache.LFU(); ok {
		t.Errorf("lfu %v, want %v", key, "")
	}

	cache.Put("1", "1")
	cache.Put("2", "2")
	cache.Put("2", "2`")

	if key, frequency, _ := cache.LFU(); key != "2" || frequency != 2 {
		t.Errorf("lfu %v, want %v", key, "2")
		t.Errorf("frequency %v, want %v", frequency, 2)
	}
}

func TestCache_Delete(t *testing.T) {
	cache, _ := lfu.New(1)
	cache.Put("key", "value")

	cache.Delete("key")
	cache.Delete("non-existing-key")

	if cache.Len() != 0 {
		t.Errorf("cache len %v, want %v", cache.Len(), 0)
	}

	if key, frequency, ok := cache.LFU(); ok {
		t.Errorf("lfu key %v, want %v", key, "''")
		t.Errorf("frequency %v, want %v", frequency, 0)
	}
}

func TestCache_Clear(t *testing.T) {
	cache, _ := lfu.New(10)
	cache.Put("key1", "value1")
	cache.Put("key2", "value2")
	cache.Clear()

	if cache.Len() != 0 {
		t.Errorf("cache len %v, want %v", cache.Len(), 0)
	}
}

func TestReplace(t *testing.T) {
	cache, _ := lfu.New(10)
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
	cache, _ := lfu.New(1)
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
	cache, _ := lfu.New(2)

	cache.Put("1", 1)
	cache.Put("2", 2)

	if value, ok := cache.Get("1"); !ok {
		t.Errorf("cached value %v, want %v", value, 1)
	}

	if key, frequency, _ := cache.LFU(); key != "2" || frequency != 1 {
		t.Errorf("lfu key %v, want %v", key, "2")
		t.Errorf("frequency %v, want %v", frequency, 1)
	}

	// key: 1, frequency: 2
	// key: 2, frequency: 1

	cache.Put("3", 3)

	if value, ok := cache.Get("1"); !ok {
		t.Errorf("cached value %v, want %v", value, 1)
	}

	if value, ok := cache.Get("2"); ok {
		t.Errorf("cached value %v, want %v", value, nil)
	}

	if value, ok := cache.Get("3"); !ok {
		t.Errorf("cached value %v, want %v", value, 3)
	}

	if key, frequency, _ := cache.LFU(); key != "3" || frequency != 2 {
		t.Errorf("lfu key %v, want %v", key, "3")
		t.Errorf("frequency %v, want %v", frequency, 2)
	}

	// key: 1, frequency: 3
	// key: 3, frequency: 2

	cache.Put("4", 4)

	if value, ok := cache.Get("1"); !ok {
		t.Errorf("cached value %v, want %v", value, 1)
	}

	if value, ok := cache.Get("3"); ok {
		t.Errorf("cached value %v, want %v", value, nil)
	}

	if value, ok := cache.Get("4"); !ok {
		t.Errorf("cached value %v, want %v", value, 4)
	}

	if key, frequency, _ := cache.LFU(); key != "4" || frequency != 2 {
		t.Errorf("lfu key %v, want %v", key, "4")
		t.Errorf("frequency %v, want %v", frequency, 2)
	}

	// key: 1, frequency: 4
	// key: 4, frequency: 2

	cache.Put("4", 40)

	if value, ok := cache.Get("4"); !ok {
		t.Errorf("cached value %v, want %v", value, 40)
	}

	// key: 1, frequency: 4
	// key: 4, frequency: 4

	cache.Put("4", 400)
	cache.Put("5", 5)

	if value, ok := cache.Get("4"); !ok {
		t.Errorf("cached value %v, want %v", value, 400)
	}

	if value, ok := cache.Get("5"); !ok {
		t.Errorf("cached value %v, want %v", value, 5)
	}

	if key, frequency, _ := cache.LFU(); key != "5" || frequency != 2 {
		t.Errorf("lfu key %v, want %v", key, "5")
		t.Errorf("frequency %v, want %v", frequency, 2)
	}

	// key: 4, frequency: 6
	// key: 5: frequency: 2
}
