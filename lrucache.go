package cache

type LRUCache[K comparable, V any] struct {
	BaseCache[K, V]
}

func NewLRUCache[K comparable, V any](capacity int, defaultEmptyValue V) LRUCache[K, V] {
	return LRUCache[K, V]{
		NewBaseCache[K, V](newLRUPolicy[K, V](), capacity, defaultEmptyValue),
	}
}
