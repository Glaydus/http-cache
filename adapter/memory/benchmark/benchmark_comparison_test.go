package main

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	hcache "github.com/Glaydus/http-cache"
	"github.com/Glaydus/http-cache/adapter/memory"

	"github.com/allegro/bigcache"
)

const maxEntrySize = 256

func BenchmarkHTTPCacheMamoryAdapterSet(b *testing.B) {
	cache, expiration := initHTTPCacheMemoryAdapter(b.N)
	for i := 0; i < b.N; i++ {
		cache.Set(uint64(i), value(), expiration)
	}
}

func BenchmarkBigCacheSet(b *testing.B) {
	cache := initBigCache(b.N)
	for i := 0; i < b.N; i++ {
		err := cache.Set(strconv.Itoa(i), value())
		if err != nil {
			return
		}
	}
}

func BenchmarkHTTPCacheMamoryAdapterGet(b *testing.B) {
	b.StopTimer()
	cache, expiration := initHTTPCacheMemoryAdapter(b.N)
	for i := 0; i < b.N; i++ {
		cache.Set(uint64(i), value(), expiration)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(uint64(i))
	}
}

func BenchmarkBigCacheGet(b *testing.B) {
	b.StopTimer()
	cache := initBigCache(b.N)
	for i := 0; i < b.N; i++ {
		err := cache.Set(strconv.Itoa(i), value())
		if err != nil {
			return
		}
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := cache.Get(strconv.Itoa(i))
		if err != nil {
			return
		}
	}
}

func BenchmarkHTTPCacheMamoryAdapterSetParallel(b *testing.B) {
	cache, expiration := initHTTPCacheMemoryAdapter(b.N)
	rand.New(rand.NewSource(time.Now().Unix()))

	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(1000)
		counter := 0
		for pb.Next() {
			cache.Set(uint64(id), value(), expiration)
			counter = counter + 1
		}
	})
}

func BenchmarkBigCacheSetParallel(b *testing.B) {
	cache := initBigCache(b.N)
	rand.New(rand.NewSource(time.Now().Unix()))

	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(1000)
		counter := 0
		for pb.Next() {
			if err := cache.Set(strconv.FormatUint(uint64(id), 10), value()); err == nil {
				counter = counter + 1
			}
		}
	})
}

func BenchmarkHTTPCacheMemoryAdapterGetParallel(b *testing.B) {
	b.StopTimer()
	cache, expiration := initHTTPCacheMemoryAdapter(b.N)
	for i := 0; i < b.N; i++ {
		cache.Set(uint64(i), value(), expiration)
	}

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			cache.Get(uint64(counter))
			counter = counter + 1
		}
	})
}

func BenchmarkBigCacheGetParallel(b *testing.B) {
	b.StopTimer()
	cache := initBigCache(b.N)
	for i := 0; i < b.N; i++ {
		err := cache.Set(strconv.Itoa(i), value())
		if err != nil {
			return
		}
	}

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			if _, err := cache.Get(strconv.Itoa(counter)); err == nil {
				counter = counter + 1
			}
		}
	})
}

func value() []byte {
	return make([]byte, 100)
}

func initHTTPCacheMemoryAdapter(entries int) (hcache.Adapter, time.Time) {
	if entries < 2 {
		entries = 2
	}
	adapter, _ := memory.NewAdapter(
		memory.AdapterWithCapacity(entries),
		memory.AdapterWithAlgorithm(memory.LRU),
	)

	return adapter, time.Now().Add(1 * time.Minute)
}

func initBigCache(entriesInWindow int) *bigcache.BigCache {
	cache, _ := bigcache.NewBigCache(bigcache.Config{
		Shards:             256,
		LifeWindow:         10 * time.Minute,
		MaxEntriesInWindow: entriesInWindow,
		MaxEntrySize:       maxEntrySize,
		Verbose:            true,
	})

	return cache
}
