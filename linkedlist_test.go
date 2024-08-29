package cache_test

import (
	"testing"

	cache "github.com/yairp7/gocache"
	"gopkg.in/stretchr/testify.v1/assert"
)

func Test_Linkedlist_Add(t *testing.T) {
	list := cache.NewLinkedList[int]()
	list.Add(0)
	list.Add(1)
	list.Add(2)
	list.Add(3)
	assert.Equal(t, 4, list.Len())
}

func Test_Linkedlist_Remove(t *testing.T) {
	list := cache.NewLinkedList[int]()
	list.Add(0)
	e1 := list.Add(1)
	list.Add(2)
	e3 := list.Add(3)
	list.Remove(e1)
	list.Remove(e3)
	assert.Equal(t, 2, list.Len())
}

func Test_Linkedlist_At(t *testing.T) {
	list := cache.NewLinkedList[int]()
	e0 := list.Add(0)
	e1 := list.Add(1)
	e2 := list.Add(2)
	e3 := list.Add(3)
	list.Add(4)
	assert.Equal(t, e0, list.At(0))
	assert.Equal(t, e1, list.At(1))
	assert.Equal(t, e2, list.At(2))
	assert.Equal(t, e3, list.At(3))
}

func Test_Linkedlist_HeadTail(t *testing.T) {
	list := cache.NewLinkedList[int]()
	e0 := list.Add(0)
	list.Add(1)
	list.Add(2)
	e3 := list.Add(3)
	assert.Equal(t, e3, list.Head())
	assert.Equal(t, e0, list.Tail())
}

func Test_Linkedlist_PopTail(t *testing.T) {
	list := cache.NewLinkedList[int]()
	e0 := list.Add(0)
	e1 := list.Add(1)
	list.Add(2)
	list.Add(3)
	assert.Equal(t, e0, list.PopTail())
	assert.Equal(t, 3, list.Len())
	assert.Equal(t, e1, list.PopTail())
	assert.Equal(t, 2, list.Len())
}

func Test_Linkedlist_PopHead(t *testing.T) {
	list := cache.NewLinkedList[int]()
	list.Add(0)
	list.Add(1)
	e2 := list.Add(2)
	e3 := list.Add(3)
	assert.Equal(t, e3, list.PopHead())
	assert.Equal(t, 3, list.Len())
	assert.Equal(t, e2, list.PopHead())
	assert.Equal(t, 2, list.Len())
}

func Test_Linkedlist_Swap(t *testing.T) {
	list := cache.NewLinkedList[int]()
	list.Add(0)
	e1 := list.Add(1)
	list.Add(2)
	e3 := list.Add(3)
	list.Swap(e1, e3)
	assert.Equal(t, e1, list.Head())
	assert.Equal(t, e3, list.At(1))
	list.Swap(e1, e3)
	assert.Equal(t, e3, list.Head())
	assert.Equal(t, e1, list.At(1))
}
