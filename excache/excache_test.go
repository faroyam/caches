package excache_test

import (
	"testing"
	"time"

	"github.com/faroyam/caches/excache"
)

func TestCache_New(t *testing.T) {
	_, err := excache.New(0)
	if err == nil {
		t.Errorf("expected error")
	}

	_, err = excache.New(-1)
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestCache_Put(t *testing.T) {
	cache, _ := excache.New(1)
	cache.Put("key", "value", time.Second)

	if cache.Len() != 1 {
		t.Errorf("cache len %v, want %v", cache.Len(), 1)
	}
}

func TestCache_Get(t *testing.T) {
	cache, _ := excache.New(1)
	cache.Put("key", "value", time.Second)

	if value, ok := cache.Get("key"); !ok || value != "value" {
		t.Errorf("cached value %v, want %v", value, "value")
	}

	if value, ok := cache.Get("non-existing-key"); ok {
		t.Errorf("cached value %v, want %v", value, "nil")
	}
}

func TestCache_Expire(t *testing.T) {
	cache, _ := excache.New(2)

	cache.Put("key1", "value1", time.Millisecond)
	cache.Put("key2", "value2", time.Millisecond)

	if cache.Len() != 2 {
		t.Errorf("cache len %v, want %v", cache.Len(), 2)
	}

	time.Sleep(time.Millisecond * 10)

	if cache.Len() != 0 {
		t.Errorf("cache len %v, want %v", cache.Len(), 0)
	}
}

func TestCache_Get_ResetsTTL(t *testing.T) {
	cache, _ := excache.New(1)
	cache.Put("key", "value", time.Millisecond*100)

	time.Sleep(time.Millisecond * 70)

	if value, ok := cache.Get("key"); !ok || value != "value" {
		t.Errorf("cached value %v, want %v", value, "value")
	}

	time.Sleep(time.Millisecond * 70)

	if value, ok := cache.Get("key"); !ok || value != "value" {
		t.Errorf("cached value %v, want %v", value, "value")
	}
}

func TestCache_Delete(t *testing.T) {
	cache, _ := excache.New(1)
	cache.Put("key", "value", time.Second)

	cache.Delete("key")
	cache.Delete("non-existing-key")

	if cache.Len() != 0 {
		t.Errorf("cache len %v, want %v", cache.Len(), 0)
	}
}

func TestCache_Clear(t *testing.T) {
	cache, _ := excache.New(10)
	cache.Put("key1", "value1", time.Second)
	cache.Put("key2", "value2", time.Second)
	cache.Clear()

	if cache.Len() != 0 {
		t.Errorf("cache len %v, want %v", cache.Len(), 0)
	}
}

func TestReplace(t *testing.T) {
	cache, _ := excache.New(10)
	cache.Put("key", "value1", time.Second)

	if value, ok := cache.Get("key"); !ok || value != "value1" {
		t.Errorf("cached value %v, want %v", value, "value1")
	}

	cache.Put("key", "value2", time.Second)

	if value, ok := cache.Get("key"); !ok || value != "value2" {
		t.Errorf("cached value %v, want %v", value, "value2")
	}

	if cache.Len() != 1 {
		t.Errorf("cache len %v, want %v", cache.Len(), 1)
	}
}

func TestPutMoreThanCap(t *testing.T) {
	cache, _ := excache.New(2)
	key1 := "key1"
	value1 := "value1"
	key2 := "key2"
	value2 := "value2"
	key3 := "key3"
	value3 := "value3"

	cache.Put(key1, value1, time.Second*3)
	if value, ok := cache.Get(key1); !ok || value != value1 {
		t.Errorf("cached value %v, want %v", value, value1)
	}

	cache.Put(key2, value2, time.Second*2)
	if value, ok := cache.Get(key2); !ok || value != value2 {
		t.Errorf("cached value %v, want %v", value, value2)
	}

	cache.Put(key3, value3, time.Second)
	if value, ok := cache.Get(key3); !ok || value != value3 {
		t.Errorf("cached value %v, want %v", value, value3)
	}

	if value, ok := cache.Get(key2); ok {
		t.Errorf("cached value %v, want %v", value, nil)
	}

	if cache.Len() != 2 {
		t.Errorf("cache len %v, want %v", cache.Len(), 2)
	}
}
