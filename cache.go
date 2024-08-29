package cache

import "sync"

type EntryExtraInfo map[struct{}]any

type cacheEntry[K comparable, V any] struct {
	key       K
	value     V
	extraInfo EntryExtraInfo
}

func newCacheEntry[K comparable, V any](key K, value V) *cacheEntry[K, V] {
	entry := &cacheEntry[K, V]{
		key:       key,
		value:     value,
		extraInfo: make(EntryExtraInfo),
	}
	return entry
}

type BaseCache[K comparable, V any] struct {
	mapping           map[K]*cacheEntry[K, V]
	size              int
	capacity          int
	lock              sync.Mutex
	defaultEmptyValue V
	policy            Policy[K, V]
}

func NewBaseCache[K comparable, V any](policy Policy[K, V], capacity int, defaultEmptyValue V) BaseCache[K, V] {
	return BaseCache[K, V]{
		mapping:           make(map[K]*cacheEntry[K, V], capacity),
		size:              0,
		capacity:          capacity,
		defaultEmptyValue: defaultEmptyValue,
		policy:            policy,
	}
}

func (c *BaseCache[K, V]) removeEntry(entry *cacheEntry[K, V]) {
	clear(entry.extraInfo)
	delete(c.mapping, entry.key)
	c.size--
}

func (c *BaseCache[K, V]) addEntry(entry *cacheEntry[K, V]) {
	c.mapping[entry.key] = entry
	c.size++
}

func (c *BaseCache[K, V]) Get(key K) any {
	c.lock.Lock()
	defer c.lock.Unlock()

	// If exists, return and move to front
	if entry, ok := c.mapping[key]; ok {
		c.policy.beforeGet(entry)
		return entry.value
	}

	return nil
}

func (c *BaseCache[K, V]) Set(key K, value V) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var entry *cacheEntry[K, V]

	if e, ok := c.mapping[key]; ok {
		entry = e
	} else {
		entry = newCacheEntry(key, value)
		c.addEntry(entry)
	}

	c.policy.afterAdd(entry)

	if c.Size() > c.capacity {
		evictedEntry := c.policy.evict()
		c.removeEntry(evictedEntry)
	}
}

func (c *BaseCache[K, V]) Evict() *cacheEntry[K, V] {
	evictedEntry := c.policy.evict()
	c.removeEntry(evictedEntry)
	return evictedEntry
}

func (c *BaseCache[K, V]) Size() int {
	return c.size
}
