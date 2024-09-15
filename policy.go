package cache

import (
	ds "github.com/yairp7/go-common-lib/ds/basic"
	"github.com/yairp7/gocache/pqueue"
)

type Policy[K comparable, V any] interface {
	// Called after adding the item to the map, and before evicting items if needed
	afterAdd(entry *cacheEntry[K, V])

	// Called after getting the item from the map, and before returning it
	beforeGet(entry *cacheEntry[K, V])

	// Called after adding an item in case the size has exceeded the maximum capacity
	evict() *cacheEntry[K, V]
}

type LRUPolicy[K comparable, V any] struct {
	orderList *ds.LinkedList[*cacheEntry[K, V]]
}

func newLRUPolicy[K comparable, V any]() *LRUPolicy[K, V] {
	p := &LRUPolicy[K, V]{
		orderList: ds.NewLinkedList[*cacheEntry[K, V]](),
	}
	return p
}

var lruPolicyEntryKey = struct{}{}

func (p *LRUPolicy[K, V]) afterAdd(entry *cacheEntry[K, V]) {
	if v, ok := entry.extraInfo[lruPolicyEntryKey]; ok {
		p.orderList.MoveToFront(v.(*ds.LinkedListEntry[*cacheEntry[K, V]]))
		return
	}

	listEntry := p.orderList.AddFront(entry)
	entry.extraInfo[lruPolicyEntryKey] = listEntry
}

func (p *LRUPolicy[K, V]) beforeGet(entry *cacheEntry[K, V]) {
	lruEntry := entry.extraInfo[lruPolicyEntryKey].(*ds.LinkedListEntry[*cacheEntry[K, V]])
	p.orderList.MoveToFront(lruEntry)
}

func (p *LRUPolicy[K, V]) evict() *cacheEntry[K, V] {
	evictedEntry := p.orderList.PopTail()
	if evictedEntry != nil {
		return evictedEntry.Data
	}
	return nil
}

type LFUPolicy[K comparable, V any] struct {
	heap *pqueue.MinHeap[*cacheEntry[K, V]]
}

func newLFUPolicy[K comparable, V any]() *LFUPolicy[K, V] {
	return &LFUPolicy[K, V]{
		heap: pqueue.NewMinHeap[*cacheEntry[K, V]](),
	}
}

var lfuPolicyEntryKey = struct{}{}

func (p *LFUPolicy[K, V]) incrementFreq(node *pqueue.HeapNode[*cacheEntry[K, V]]) {
	p.heap.Touch(node)
}

func (p *LFUPolicy[K, V]) afterAdd(entry *cacheEntry[K, V]) {
	if v, ok := entry.extraInfo[lfuPolicyEntryKey]; ok {
		node := v.(*pqueue.HeapNode[*cacheEntry[K, V]])
		p.incrementFreq(node)
		return
	}

	node := p.heap.Push(entry)
	entry.extraInfo[lfuPolicyEntryKey] = node
}

func (p *LFUPolicy[K, V]) beforeGet(entry *cacheEntry[K, V]) {
	node := entry.extraInfo[lfuPolicyEntryKey].(*pqueue.HeapNode[*cacheEntry[K, V]])
	p.incrementFreq(node)
}

func (p *LFUPolicy[K, V]) evict() *cacheEntry[K, V] {
	if p.heap.Size() > 0 {
		evictedEntry := p.heap.Pop()
		return evictedEntry.Data
	}
	return nil
}
