package pqueue

import (
	"fmt"
	"strings"
	"time"
)

type HeapNode[V any] struct {
	entry *pQueueEntry[*HeapNode[V]]
	Data  V
}

func (h *HeapNode[V]) String() string {
	return fmt.Sprintf("%v-%d", h.Data, h.entry.weight)
}

func newHeapNode[V any](data V) *HeapNode[V] {
	return &HeapNode[V]{
		Data: data,
	}
}

type MinHeap[V any] pQueue[*HeapNode[V]]

func NewMinHeap[V any]() *MinHeap[V] {
	q := newPriorityQueue[*HeapNode[V]]()
	return (*MinHeap[V])(q)
}

func (h *MinHeap[V]) Push(data V) *HeapNode[V] {
	node := newHeapNode(data)
	entry := (*pQueue[*HeapNode[V]])(h).push(node, 1)
	node.entry = entry
	return node
}

func (h *MinHeap[V]) Pop() *HeapNode[V] {
	entry := (*pQueue[*HeapNode[V]])(h).pop()
	if entry == nil {
		return nil
	}
	return entry.value
}

func (h *MinHeap[V]) Touch(node *HeapNode[V]) {
	(*pQueue[*HeapNode[V]])(h).setWeight(node.entry.index, node.entry.weight+1)
}

func (h *MinHeap[V]) Size() int {
	return (*pQueue[*HeapNode[V]])(h).size
}

// pQueue implementation

type pQueueEntry[V any] struct {
	index     int
	value     V
	weight    int
	fetchedAt time.Time
}

func (e pQueueEntry[V]) String() string {
	return fmt.Sprintf("%v-%d", e.value, e.weight)
}

func newPriorityQueueEntry[V any](value V, weight int) *pQueueEntry[V] {
	return &pQueueEntry[V]{
		index:     0,
		value:     value,
		weight:    weight,
		fetchedAt: time.Now(),
	}
}

type pQueue[V any] struct {
	q    []*pQueueEntry[V]
	size int
}

func newPriorityQueue[V any]() *pQueue[V] {
	queue := pQueue[V]{
		q:    make([]*pQueueEntry[V], 0),
		size: 0,
	}
	return &queue
}

func (q *pQueue[V]) parent(index int) int {
	return (index - 1) / 2
}

func (q *pQueue[V]) left(index int) int {
	return index*2 + 1
}

func (q *pQueue[V]) right(index int) int {
	return index*2 + 2
}

func (q *pQueue[V]) isLeaf(index int) bool {
	return index >= (q.size/2) && index <= q.size
}

func (q *pQueue[V]) peek() *pQueueEntry[V] {
	return q.q[0]
}

func (q *pQueue[V]) swap(i, j int) {
	q.q[i].index, q.q[j].index = q.q[j].index, q.q[i].index
	q.q[i], q.q[j] = q.q[j], q.q[i]
}

func (q *pQueue[V]) heapifyDown(index int) {
	l, r := q.left(index), q.right(index)
	minMaxIndex := index

	if l < q.size && q.q[l].weight < q.q[minMaxIndex].weight {
		minMaxIndex = l
	}
	if r < q.size && q.q[r].weight < q.q[minMaxIndex].weight {
		minMaxIndex = r
	}

	if minMaxIndex != index {
		q.swap(index, minMaxIndex)
		q.heapifyDown(minMaxIndex)
	}
}

func (q *pQueue[V]) heapifyUp(index int) {
	shouldSwapParent := func(childIndex, parentIndex int) bool {
		// if q.q[parentIndex].weight == q.q[childIndex].weight {
		// 	return q.q[parentIndex].fetchedAt.After(q.q[childIndex].fetchedAt)
		// }
		return q.q[parentIndex].weight > q.q[childIndex].weight
	}

	parentIndex := q.parent(index)
	for index >= 0 && shouldSwapParent(index, parentIndex) {
		q.swap(index, parentIndex)
		index = parentIndex
		parentIndex = q.parent(index)
	}
}

func (q *pQueue[V]) setWeight(index, weight int) {
	q.q[index].weight = weight
	q.heapifyDown(index)
}

func (q *pQueue[V]) pop() *pQueueEntry[V] {
	if q.size == 0 {
		return nil
	}

	if q.size == 1 {
		q.size--
		return q.q[0]
	}

	root := q.q[0]
	q.q[0] = q.q[q.size-1]
	q.q = q.q[0 : q.size-1]
	q.size--
	q.heapifyDown(0)

	return root
}

func (q *pQueue[V]) push(value V, weight int) *pQueueEntry[V] {
	q.size++
	cEntryIndex := q.size - 1
	entry := newPriorityQueueEntry(value, weight)
	entry.index = cEntryIndex
	q.q = append(q.q, entry)

	q.heapifyUp(cEntryIndex)

	return entry
}

func (q *pQueue[V]) String() string {
	var sb strings.Builder
	for _, e := range q.q {
		sb.WriteString(fmt.Sprintf("%v->(%v,%d)|", q.q[q.parent(e.index)].value, e.value, e.weight))
	}
	return sb.String()[:sb.Len()-1]
}
