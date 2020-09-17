package lru

import (
	"strconv"
	"testing"
)

const (
	size1kk  = 1_000_000
	size10kk = 10_000_000
)

func BenchmarkLookup10kkMap(b *testing.B) {
	m := initMap(size10kk, size10kk)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, ok := m[strconv.Itoa(n)]; ok {
		}
	}
}

func BenchmarkLookup10kkLruCache(b *testing.B) {
	c := initCache(size10kk)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, ok := c.Get(strconv.Itoa(n)); ok {
		}
	}
}

func BenchmarkLookup1kkMap(b *testing.B) {
	m := initMap(size1kk, size1kk)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, ok := m[strconv.Itoa(n)]; ok {
		}
	}
}

func BenchmarkLookup1kkLruCache(b *testing.B) {
	c := initCache(size1kk)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, ok := c.Get(strconv.Itoa(n)); ok {
		}
	}
}

func BenchmarkPut10kkMap(b *testing.B) {
	m := initMap(size10kk, size10kk*2)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		key := strconv.Itoa(n + size10kk)
		m[key] = key
	}
}

func BenchmarkPut10kkLruCache(b *testing.B) {
	c := initCache(size10kk)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		key := strconv.Itoa(n + size10kk)
		c.Put(key, key)
	}
}

func BenchmarkPut1kkMap(b *testing.B) {
	m := initMap(size1kk, size1kk*2)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		key := strconv.Itoa(n + size1kk)
		m[key] = key
	}
}

func BenchmarkPut1kkLruCache(b *testing.B) {
	c := initCache(size1kk)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		key := strconv.Itoa(n + size1kk)
		c.Put(key, key)
	}
}

func initMap(size int, cap int) map[string]interface{} {
	m := make(map[string]interface{}, cap)
	for i := 0; i < size; i++ {
		key := strconv.Itoa(i)
		m[key] = key
	}
	return m
}

func initCache(size int) *Cache {
	c := New(size)
	for i := 0; i < size; i++ {
		key := strconv.Itoa(i)
		c.Put(key, key)
	}
	return c
}
