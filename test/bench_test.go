package test_test

import (
	"strconv"
	"testing"

	"github.com/faroyam/caches/test"

	"github.com/faroyam/caches/lfu"
	"github.com/faroyam/caches/lru"
)

const (
	size100k = 1000_000
	size1kk  = 1_000_000
)

type cache interface {
	Get(key string) (interface{}, bool)
	Put(key string, value interface{})
}

var result interface{}

func BenchmarkMapGet100k(b *testing.B) { benchmarkGet(initMap(size100k, size100k), size100k, b) }
func BenchmarkLRUGet100k(b *testing.B) { benchmarkGet(initLRUCache(size100k), size100k, b) }
func BenchmarkLFUGet100k(b *testing.B) { benchmarkGet(initLFUCache(size100k), size100k, b) }
func BenchmarkMapGet1kk(b *testing.B)  { benchmarkGet(initMap(size1kk, size1kk), size1kk, b) }
func BenchmarkLRUGet1kk(b *testing.B)  { benchmarkGet(initLRUCache(size1kk), size1kk, b) }
func BenchmarkLFUGet1kk(b *testing.B)  { benchmarkGet(initLFUCache(size1kk), size1kk, b) }

func BenchmarkMapPut100k(b *testing.B) { benchmarkPut(initMap(size100k, size100k*10), size100k, b) }
func BenchmarkLRUPut100k(b *testing.B) { benchmarkPut(initLRUCache(size100k), size100k, b) }
func BenchmarkLFUPut100k(b *testing.B) { benchmarkPut(initLFUCache(size100k), size100k, b) }
func BenchmarkMapPut1kk(b *testing.B)  { benchmarkPut(initMap(size1kk, size1kk*10), size1kk, b) }
func BenchmarkLRUPut1kk(b *testing.B)  { benchmarkPut(initLRUCache(size1kk), size1kk, b) }
func BenchmarkLFUPut1kk(b *testing.B)  { benchmarkPut(initLFUCache(size1kk), size1kk, b) }

func benchmarkGet(cache cache, size int, b *testing.B) {
	var v interface{}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v, _ = cache.Get(strconv.Itoa(n % size))
	}
	result = v
}

func benchmarkPut(cache cache, size int, b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		key := strconv.Itoa(n + size)
		cache.Put(key, key)
	}
	result, _ = cache.Get(strconv.Itoa(size))
}

func initMap(size, cap int) *test.Map {
	m := test.New(cap)
	for i := 0; i < size; i++ {
		key := strconv.Itoa(i)
		m.Put(key, key)
	}
	return m
}

func initLRUCache(size int) *lru.Cache {
	c, _ := lru.New(size)
	for i := 0; i < size; i++ {
		key := strconv.Itoa(i)
		c.Put(key, key)
	}
	return c
}

func initLFUCache(size int) *lfu.Cache {
	c, _ := lfu.New(size)
	for i := 0; i < size; i++ {
		key := strconv.Itoa(i)
		c.Put(key, key)
	}
	return c
}
