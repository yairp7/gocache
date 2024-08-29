package cache

import (
	"sync"
	"testing"

	"gopkg.in/stretchr/testify.v1/assert"
)

func Test_LRUCache(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		lruCache := NewLRUCache[int](10, 0)

		for i := 0; i < 10; i++ {
			lruCache.Set(i, i)
		}

		for i := 9; i >= 0; i-- {
			lruCache.Set(i, i)
		}

		for i := 9; i >= 0; i-- {
			entry := lruCache.Evict()
			assert.Equal(t, i, entry.value)
		}
	})
}

func Test_MFUCache(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		mfuCache := NewMFUCache[int](10, 0)
		for i := 0; i < 10; i++ {
			mfuCache.Set(i, i)
		}

		for i := 9; i >= 0; i-- {
			for j := 0; j < i; j++ {
				mfuCache.Set(i, i)
			}
		}

		for i := 0; i < 10; i++ {
			entry := mfuCache.Evict()
			assert.Equal(t, i, entry.value)
		}
	})
}

func TestSetRace(t *testing.T) {
	cache := NewLRUCache[int, int](100000, 0)
	waitGroup := sync.WaitGroup{}

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		for i := 0; i < 10000; i++ {
			cache.Set(i, i)
		}
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		for i := 10000; i < 20000; i++ {
			cache.Set(i, i)
		}
	}()

	waitGroup.Wait()
}

func TestGetRace(t *testing.T) {
	cache := NewLRUCache[int, int](100000, 0)
	waitGroup := sync.WaitGroup{}

	for i := 0; i < 10000; i++ {
		cache.Set(i, i)
	}

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		for i := 0; i < 10000; i++ {
			cache.Get(i)
		}
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		for i := 0; i < 10000; i++ {
			cache.Get(i)
		}
	}()

	waitGroup.Wait()
}

func Benchmark_Set(b *testing.B) {
	b.Run("LRU", func(b *testing.B) {
		b.ReportAllocs()

		lruCache := NewLRUCache[int, int](1000, 0)
		for i := 0; i < 1000; i++ {
			lruCache.Set(i, i)
		}

		for i := 0; i < b.N; i++ {
			lruCache.Set(i, i)
		}
	})

	b.Run("MFU", func(b *testing.B) {
		b.ReportAllocs()
		mfuCache := NewMFUCache[int, int](1000, 0)

		for i := 0; i < 1000; i++ {
			mfuCache.Set(i, i)
		}

		for i := 0; i < b.N; i++ {
			mfuCache.Set(i, i)
		}
	})
}
