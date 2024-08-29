package cache

import ds "github.com/yairp7/go-common-lib/ds/basic"

type Policy[K comparable, V any] interface {
	afterAdd(entry *cacheEntry[K, V])
	beforeGet(entry *cacheEntry[K, V])
	evict() *cacheEntry[K, V]
}

type LRUPolicy[K comparable, V any] struct {
	list *ds.LinkedList[*cacheEntry[K, V]]
}

func newLRUPolicy[K comparable, V any]() *LRUPolicy[K, V] {
	p := &LRUPolicy[K, V]{
		list: ds.NewLinkedList[*cacheEntry[K, V]](),
	}
	return p
}

var lruPolicyEntryKey = struct{}{}

func (p *LRUPolicy[K, V]) afterAdd(entry *cacheEntry[K, V]) {
	if v, ok := entry.extraInfo[lruPolicyEntryKey]; ok {
		p.list.MoveToFront(v.(*ds.LinkedListEntry[*cacheEntry[K, V]]))
		return
	}

	listEntry := p.list.Add(entry)
	entry.extraInfo[lruPolicyEntryKey] = listEntry
}

func (p *LRUPolicy[K, V]) beforeGet(entry *cacheEntry[K, V]) {
	lruEntry := entry.extraInfo[lruPolicyEntryKey].(*ds.LinkedListEntry[*cacheEntry[K, V]])
	p.list.MoveToFront(lruEntry)
}

func (p *LRUPolicy[K, V]) evict() *cacheEntry[K, V] {
	evictedEntry := p.list.PopTail()
	if evictedEntry != nil {
		return evictedEntry.Data
	}
	return nil
}

type LFUPolicy[K comparable, V any] struct {
	minHeap Heap[*cacheEntry[K, V]]
}

func newLFUPolicy[K comparable, V any]() *LFUPolicy[K, V] {
	return &LFUPolicy[K, V]{
		minHeap: NewMinHeap[*cacheEntry[K, V]](),
	}
}

var lfuPolicyEntryKey = struct{}{}

func (p *LFUPolicy[K, V]) incrementWeight(node *HeapNode[*cacheEntry[K, V]]) {
	node.Weight++
	p.minHeap.Order(node.Index)
}

func (p *LFUPolicy[K, V]) afterAdd(entry *cacheEntry[K, V]) {
	if v, ok := entry.extraInfo[lfuPolicyEntryKey]; ok {
		node := v.(*HeapNode[*cacheEntry[K, V]])
		p.incrementWeight(node)
		return
	}

	node := &HeapNode[*cacheEntry[K, V]]{Weight: 1, Data: entry}
	p.minHeap.Push(node)
	entry.extraInfo[lfuPolicyEntryKey] = node
}

func (p *LFUPolicy[K, V]) beforeGet(entry *cacheEntry[K, V]) {
	node := entry.extraInfo[lfuPolicyEntryKey].(*HeapNode[*cacheEntry[K, V]])
	p.incrementWeight(node)
}

func (p *LFUPolicy[K, V]) evict() *cacheEntry[K, V] {
	if p.minHeap.Len() > 0 {
		evictedEntry := p.minHeap.Pop()
		return evictedEntry.Data
	}
	return nil
}
