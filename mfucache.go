package cache

type MFUCache[K comparable, V any] struct {
	baseCache[K, V]
}

func NewMFUCache[K comparable, V any](capacity int, defaultEmptyValue V) MFUCache[K, V] {
	return MFUCache[K, V]{
		newBaseCache[K, V](newMFUPolicy[K, V](), capacity, defaultEmptyValue),
	}
}
