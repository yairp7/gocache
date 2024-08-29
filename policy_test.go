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
	node := lruPolicy.list.tail
	for node != nil {
		assert.Equal(t, entries[i].extraInfo[lruPolicyEntryKey], node)
		node = node.next
		i++
	}

	lruPolicy.afterAdd(entries[0])
	lruPolicy.afterAdd(entries[1])
	lruPolicy.afterAdd(entries[1])

	assert.Equal(t, entries[2].extraInfo[lruPolicyEntryKey], lruPolicy.list.PopTail())
	assert.Equal(t, entries[0].extraInfo[lruPolicyEntryKey], lruPolicy.list.PopTail())
	assert.Equal(t, entries[1].extraInfo[lruPolicyEntryKey], lruPolicy.list.PopTail())
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

	assert.Equal(t, entries[1].extraInfo[lruPolicyEntryKey], lruPolicy.list.PopTail())
	assert.Equal(t, entries[0].extraInfo[lruPolicyEntryKey], lruPolicy.list.PopTail())
	assert.Equal(t, entries[2].extraInfo[lruPolicyEntryKey], lruPolicy.list.PopTail())
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

func Test_MFUPolicy_afterAdd(t *testing.T) {
	mfuPolicy := newMFUPolicy[string, string]()

	entries := []*cacheEntry[string, string]{
		{key: "a", value: "a", extraInfo: make(EntryExtraInfo)},
		{key: "b", value: "b", extraInfo: make(EntryExtraInfo)},
		{key: "c", value: "c", extraInfo: make(EntryExtraInfo)},
	}

	for _, entry := range entries {
		mfuPolicy.afterAdd(entry)
	}

	for i, node := range mfuPolicy.minHeap.Values.Nodes {
		assert.Equal(t, entries[i].extraInfo[mfuPolicyEntryKey], node)
	}

	mfuPolicy.afterAdd(entries[0])
	mfuPolicy.afterAdd(entries[1])
	mfuPolicy.afterAdd(entries[1])

	node := mfuPolicy.minHeap.Pop()
	assert.Equal(t, entries[2].key, node.Data.key)
}

func Test_MFUPolicy_beforeGet(t *testing.T) {
	mfuPolicy := newMFUPolicy[string, string]()

	entries := []*cacheEntry[string, string]{
		{key: "a", value: "a", extraInfo: make(EntryExtraInfo)},
		{key: "b", value: "b", extraInfo: make(EntryExtraInfo)},
		{key: "c", value: "c", extraInfo: make(EntryExtraInfo)},
	}

	for _, entry := range entries {
		mfuPolicy.afterAdd(entry)
	}

	mfuPolicy.beforeGet(entries[0])
	mfuPolicy.beforeGet(entries[0])
	mfuPolicy.beforeGet(entries[0])
	mfuPolicy.beforeGet(entries[0])

	mfuPolicy.beforeGet(entries[2])
	mfuPolicy.beforeGet(entries[2])
	mfuPolicy.beforeGet(entries[2])

	assert.Equal(t, entries[1].key, mfuPolicy.minHeap.Pop().Data.key)
	assert.Equal(t, entries[2].key, mfuPolicy.minHeap.Pop().Data.key)
	assert.Equal(t, entries[0].key, mfuPolicy.minHeap.Pop().Data.key)
}

func Test_MFUPolicy_evict(t *testing.T) {
	mfuPolicy := newMFUPolicy[string, string]()

	entries := []*cacheEntry[string, string]{
		{key: "a", value: "a", extraInfo: make(EntryExtraInfo)},
		{key: "b", value: "b", extraInfo: make(EntryExtraInfo)},
		{key: "c", value: "c", extraInfo: make(EntryExtraInfo)},
	}

	for _, entry := range entries {
		mfuPolicy.afterAdd(entry)
	}

	mfuPolicy.beforeGet(entries[0])
	mfuPolicy.beforeGet(entries[0])
	mfuPolicy.beforeGet(entries[0])
	mfuPolicy.beforeGet(entries[0])

	mfuPolicy.beforeGet(entries[2])
	mfuPolicy.beforeGet(entries[2])
	mfuPolicy.beforeGet(entries[2])

	assert.Equal(t, entries[1].key, mfuPolicy.evict().key)
	assert.Equal(t, entries[2].key, mfuPolicy.evict().key)
	assert.Equal(t, entries[0].key, mfuPolicy.evict().key)
}
