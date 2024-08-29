package cache

type LFUCache[K comparable, V any] struct {
	BaseCache[K, V]
}

func NewLFUCache[K comparable, V any](capacity int, defaultEmptyValue V) LFUCache[K, V] {
	return LFUCache[K, V]{
		NewBaseCache(newLFUPolicy[K, V](), capacity, defaultEmptyValue),
	}
}
