package cache

import hp "container/heap"

type HeapDirection int

const (
	MIN HeapDirection = iota
	MAX HeapDirection = iota
)

type HeapNode[T any] struct {
	Weight float64
	Data   T
	Index  int
}

type HeapNodes[T any] struct {
	Direction HeapDirection
	Nodes     []*HeapNode[T]
}

func (h HeapNodes[T]) Len() int {
	return len(h.Nodes)
}

func (h HeapNodes[T]) Less(i, j int) bool {
	if h.Direction == MIN {
		return h.Nodes[i].Weight < h.Nodes[j].Weight
	}
	return h.Nodes[i].Weight > h.Nodes[j].Weight
}

func (h HeapNodes[T]) Swap(i, j int) {
	h.Nodes[i], h.Nodes[j] = h.Nodes[j], h.Nodes[i]
	h.Nodes[i].Index, h.Nodes[j].Index = h.Nodes[j].Index, h.Nodes[i].Index
}

func (h *HeapNodes[T]) Push(x any) {
	(*h).Nodes = append((*h).Nodes, x.(*HeapNode[T]))
}

func (h *HeapNodes[T]) Pop() any {
	old := (*h).Nodes
	n := len(old)
	x := old[n-1]
	(*h).Nodes = old[0 : n-1]
	return x
}

type Heap[T any] struct {
	Values *HeapNodes[T]
}

func NewMinHeap[T any]() Heap[T] {
	return Heap[T]{Values: &HeapNodes[T]{
		Direction: MIN,
	}}
}

func NewMaxHeap[T any]() Heap[T] {
	return Heap[T]{Values: &HeapNodes[T]{
		Direction: MAX,
	}}
}

func (h *Heap[T]) Push(p *HeapNode[T]) {
	hp.Push(h.Values, p)
	p.Index = h.Len() - 1
}

func (h *Heap[T]) Pop() *HeapNode[T] {
	n := hp.Pop(h.Values)
	return n.(*HeapNode[T])
}

func (h *Heap[T]) Order(elementIndex int) {
	hp.Fix(h.Values, elementIndex)
}

func (h *Heap[T]) Len() int {
	return h.Values.Len()
}
