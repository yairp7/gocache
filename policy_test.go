package cache

import (
	"testing"

	"gopkg.in/stretchr/testify.v1/assert"
)

func Test_LRUPolicy_afterAdd(t *testing.T) {
	lruPolicy := newLRUPolicy[string, string]()

	entries := []*cacheEntry[string, string]{
		{key: "a", value: "a", extraInfo: make(EntryExtraInfo)},
		{key: "b", value: "b", extraInfo: make(EntryExtraInfo)},
		{key: "c", value: "c", extraInfo: make(EntryExtraInfo)},
	}

	for _, entry := range entries {
		lruPolicy.afterAdd(entry)
	}

	i := 0
	node := lruPolicy.orderList.Tail()
	for node != nil {
		assert.Equal(t, entries[i].extraInfo[lruPolicyEntryKey], node)
		node = node.Next()
		i++
	}

	lruPolicy.afterAdd(entries[0])
	lruPolicy.afterAdd(entries[1])
	lruPolicy.afterAdd(entries[1])

	assert.Equal(t, entries[2].extraInfo[lruPolicyEntryKey], lruPolicy.orderList.PopTail())
	assert.Equal(t, entries[0].extraInfo[lruPolicyEntryKey], lruPolicy.orderList.PopTail())
	assert.Equal(t, entries[1].extraInfo[lruPolicyEntryKey], lruPolicy.orderList.PopTail())
}

func Test_LRUPolicy_beforeGet(t *testing.T) {
	lruPolicy := newLRUPolicy[string, string]()

	entries := []*cacheEntry[string, string]{
		{key: "a", value: "a", extraInfo: make(EntryExtraInfo)},
		{key: "b", value: "b", extraInfo: make(EntryExtraInfo)},
		{key: "c", value: "c", extraInfo: make(EntryExtraInfo)},
	}

	for _, entry := range entries {
		lruPolicy.afterAdd(entry)
	}

	lruPolicy.beforeGet(entries[0])
	lruPolicy.beforeGet(entries[0])
	lruPolicy.beforeGet(entries[0])
	lruPolicy.beforeGet(entries[0])

	lruPolicy.beforeGet(entries[2])
	lruPolicy.beforeGet(entries[2])
	lruPolicy.beforeGet(entries[2])

	assert.Equal(t, entries[1].extraInfo[lruPolicyEntryKey], lruPolicy.orderList.PopTail())
	assert.Equal(t, entries[0].extraInfo[lruPolicyEntryKey], lruPolicy.orderList.PopTail())
	assert.Equal(t, entries[2].extraInfo[lruPolicyEntryKey], lruPolicy.orderList.PopTail())
}

func Test_LRUPolicy_evict(t *testing.T) {
	lruPolicy := newLRUPolicy[string, string]()

	entries := []*cacheEntry[string, string]{
		{key: "a", value: "a", extraInfo: make(EntryExtraInfo)},
		{key: "b", value: "b", extraInfo: make(EntryExtraInfo)},
		{key: "c", value: "c", extraInfo: make(EntryExtraInfo)},
	}

	for _, entry := range entries {
		lruPolicy.afterAdd(entry)
	}

	lruPolicy.beforeGet(entries[0])
	lruPolicy.beforeGet(entries[0])
	lruPolicy.beforeGet(entries[0])
	lruPolicy.beforeGet(entries[0])

	lruPolicy.beforeGet(entries[2])
	lruPolicy.beforeGet(entries[2])
	lruPolicy.beforeGet(entries[2])

	assert.Equal(t, entries[1].key, lruPolicy.evict().key)
	assert.Equal(t, entries[0].key, lruPolicy.evict().key)
	assert.Equal(t, entries[2].key, lruPolicy.evict().key)
}

func Test_LFUPolicy_afterAdd(t *testing.T) {
	lfuPolicy := newLFUPolicy[string, string]()

	entries := []*cacheEntry[string, string]{
		{key: "a", value: "a", extraInfo: make(EntryExtraInfo)},
		{key: "b", value: "b", extraInfo: make(EntryExtraInfo)},
		{key: "c", value: "c", extraInfo: make(EntryExtraInfo)},
	}

	for _, entry := range entries {
		lfuPolicy.afterAdd(entry)
	}

	lfuPolicy.afterAdd(entries[0])
	lfuPolicy.afterAdd(entries[1])
	lfuPolicy.afterAdd(entries[1])

	assert.Equal(t, entries[2].key, lfuPolicy.heap.Pop().Data.key)
	assert.Equal(t, entries[0].key, lfuPolicy.heap.Pop().Data.key)
	assert.Equal(t, entries[1].key, lfuPolicy.heap.Pop().Data.key)
}

func Test_LFUPolicy_beforeGet(t *testing.T) {
	lfuPolicy := newLFUPolicy[string, string]()

	entries := []*cacheEntry[string, string]{
		{key: "a", value: "a", extraInfo: make(EntryExtraInfo)},
		{key: "b", value: "b", extraInfo: make(EntryExtraInfo)},
		{key: "c", value: "c", extraInfo: make(EntryExtraInfo)},
	}

	for _, entry := range entries {
		lfuPolicy.afterAdd(entry)
	}

	lfuPolicy.beforeGet(entries[0])
	lfuPolicy.beforeGet(entries[0])
	lfuPolicy.beforeGet(entries[0])
	lfuPolicy.beforeGet(entries[0])

	lfuPolicy.beforeGet(entries[2])
	lfuPolicy.beforeGet(entries[2])
	lfuPolicy.beforeGet(entries[2])

	assert.Equal(t, entries[1].key, lfuPolicy.heap.Pop().Data.key)
	assert.Equal(t, entries[2].key, lfuPolicy.heap.Pop().Data.key)
	assert.Equal(t, entries[0].key, lfuPolicy.heap.Pop().Data.key)
}

func Test_LFUPolicy_evict(t *testing.T) {
	lfuPolicy := newLFUPolicy[string, string]()

	entries := []*cacheEntry[string, string]{
		{key: "a", value: "a", extraInfo: make(EntryExtraInfo)},
		{key: "b", value: "b", extraInfo: make(EntryExtraInfo)},
		{key: "c", value: "c", extraInfo: make(EntryExtraInfo)},
	}

	for _, entry := range entries {
		lfuPolicy.afterAdd(entry)
	}

	lfuPolicy.beforeGet(entries[0])
	lfuPolicy.beforeGet(entries[0])
	lfuPolicy.beforeGet(entries[0])
	lfuPolicy.beforeGet(entries[0])

	lfuPolicy.beforeGet(entries[2])
	lfuPolicy.beforeGet(entries[2])
	lfuPolicy.beforeGet(entries[2])

	assert.Equal(t, entries[1].key, lfuPolicy.evict().key)
	assert.Equal(t, entries[2].key, lfuPolicy.evict().key)
	assert.Equal(t, entries[0].key, lfuPolicy.evict().key)
}
