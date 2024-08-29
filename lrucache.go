package cache

type LRUCache[K comparable, V any] struct {
	baseCache[K, V]
}

func NewLRUCache[K comparable, V any](capacity int, defaultEmptyValue V) LRUCache[K, V] {
	return LRUCache[K, V]{
		newBaseCache[K, V](newLRUPolicy[K, V](), capacity, defaultEmptyValue),
	}
}
