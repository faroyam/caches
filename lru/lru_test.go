package lru

import (
	"testing"
)

func TestPut(t *testing.T) {
	cache := New(1)
	cache.Put("key", "value")

	if cache.Len() != 1 {
		t.Errorf("cache len %v, want %v", cache.Len(), 1)
	}

	if cache.records.Len() != 1 {
		t.Errorf("records list len %v, want %v", cache.records.Len(), 1)
	}

	if e := cache.cache["key"]; e.Value.(record).value != "value" {
		t.Errorf("cached value %v, want %v", e.Value.(record).value, "value")
	}
}

func TestGet(t *testing.T) {
	cache := New(1)
	cache.Put("key", "value")

	if value, ok := cache.Get("key"); !ok || value != "value" {
		t.Errorf("cached value %v, want %v", value, "value")
	}

	if value, ok := cache.Get("non-existing-key"); ok {
		t.Errorf("cached value %v, want %v", value, "nil")
	}
}

func TestReplace(t *testing.T) {
	cache := New(10)
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

	if cache.records.Len() != 1 {
		t.Errorf("records list len %v, want %v", cache.records.Len(), 1)
	}
}

func TestInsertMoreThanCap(t *testing.T) {
	cache := New(1)
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

	if cache.records.Len() != 1 {
		t.Errorf("records list len %v, want %v", cache.records.Len(), 1)
	}
}

func TestDelete(t *testing.T) {
	cache := New(1)
	cache.Put("key", "value")

	cache.Delete("key")
	cache.Delete("non-existing-key")

	if cache.Len() != 0 {
		t.Errorf("cache len %v, want %v", cache.Len(), 0)
	}

	if cache.records.Len() != 0 {
		t.Errorf("records list len %v, want %v", cache.records.Len(), 0)
	}
}

func TestClear(t *testing.T) {
	cache := New(10)
	cache.Put("key1", "value1")
	cache.Put("key2", "value2")
	cache.Clear()

	if cache.Len() != 0 {
		t.Errorf("cache len %v, want %v", cache.Len(), 0)
	}

	if cache.records.Len() != 0 {
		t.Errorf("records list len %v, want %v", cache.records.Len(), 0)
	}
}

func Test(t *testing.T) {
	cache := New(2)

	cache.Put("1", 1)
	cache.Put("2", 2)
	if value, ok := cache.Get("1"); !ok || value.(int) != 1 {
		t.Errorf("cached value %v, want %v", value, 1)
	}
	cache.Put("3", 3)
	if value, ok := cache.Get("2"); ok {
		t.Errorf("cached value %v, want %v", value, nil)
	}
	cache.Put("4", 4)
	if value, ok := cache.Get("1"); ok {
		t.Errorf("cached value %v, want %v", value, nil)
	}
	if value, ok := cache.Get("3"); !ok || value.(int) != 3 {
		t.Errorf("cached value %v, want %v", value, 3)
	}
	if value, ok := cache.Get("4"); !ok || value.(int) != 4 {
		t.Errorf("cached value %v, want %v", value, 4)
	}
}
