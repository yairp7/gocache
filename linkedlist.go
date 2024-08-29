package cache

import "slices"

var linkedListEntryKey = struct{}{}

type LinkedListEntry[V any] struct {
	data     V
	next     *LinkedListEntry[V]
	previous *LinkedListEntry[V]
}

func newLinkedListEntry[V any](initialValue V) *LinkedListEntry[V] {
	return &LinkedListEntry[V]{
		data:     initialValue,
		next:     nil,
		previous: nil,
	}
}

type LinkedList[V any] struct {
	tail *LinkedListEntry[V]
	head *LinkedListEntry[V]
	size int
}

func NewLinkedList[V any]() *LinkedList[V] {
	return &LinkedList[V]{
		tail: nil,
		head: nil,
		size: 0,
	}
}

func (l *LinkedList[V]) Add(value V) *LinkedListEntry[V] {
	entry := newLinkedListEntry(value)
	return l.add(entry)
}

func (l *LinkedList[V]) add(entry *LinkedListEntry[V]) *LinkedListEntry[V] {
	if l.head != nil {
		l.head.next = entry
		entry.previous = l.head
		l.head = entry
	} else {
		l.head = entry
		l.tail = entry
	}
	entry.next = nil
	l.size++
	return entry
}

func (l *LinkedList[V]) Remove(entry *LinkedListEntry[V]) {
	isHead := entry == l.Head()
	isTail := entry == l.Tail()
	if isHead && isTail {
		l.head = nil
		l.tail = nil
	} else if isHead {
		l.head = l.head.previous
		l.head.next = nil
	} else if isTail {
		l.tail = l.tail.next
		l.tail.previous = nil
	} else {
		entry.previous.next = entry.next
		entry.next.previous = entry.previous
	}
	l.size--
}

func (l *LinkedList[V]) At(index int) *LinkedListEntry[V] {
	direction := -1
	if index < l.size-index {
		direction = 1
	}

	var node *LinkedListEntry[V]
	if direction > 0 {
		node = l.tail
	} else {
		node = l.head
		index = l.size - index - 1
	}

	currentIndex := 0
	for node != nil {
		if currentIndex == index {
			return node
		}

		if direction > 0 {
			node = node.next
		} else {
			node = node.previous
		}

		currentIndex++
	}

	return nil
}

func (l *LinkedList[V]) MoveToFront(entry *LinkedListEntry[V]) {
	if entry == l.Head() {
		return
	}

	l.Remove(entry)
	l.add(entry)
}

func (l *LinkedList[V]) Head() *LinkedListEntry[V] {
	return l.head
}

func (l *LinkedList[V]) Tail() *LinkedListEntry[V] {
	return l.tail
}

func (l *LinkedList[V]) PopHead() *LinkedListEntry[V] {
	originalHead := l.Head()
	if l.Head() != nil {
		l.Remove(l.Head())
	}
	return originalHead
}

func (l *LinkedList[V]) PopTail() *LinkedListEntry[V] {
	originalTail := l.Tail()
	if l.Tail() != nil {
		l.Remove(l.Tail())
	}
	return originalTail
}

func (l *LinkedList[V]) Swap(entryA, entryB *LinkedListEntry[V]) {
	entryA.data, entryB.data = entryB.data, entryA.data

	if entryA.next != nil {
		entryA.next.previous = entryB
	}

	if entryB.next != nil {
		entryB.next.previous = entryA
	}

	if entryA.previous != nil {
		entryA.previous.next = entryB
	}

	if entryB.previous != nil {
		entryB.previous.next = entryA
	}

	if l.Head() == entryA {
		l.head = entryB
	} else if l.Head() == entryB {
		l.head = entryA
	}

	if l.Tail() == entryA {
		l.tail = entryB
	} else if l.Tail() == entryB {
		l.tail = entryA
	}

	entryA.next, entryB.next = entryB.next, entryA.next
	entryA.previous, entryB.previous = entryB.previous, entryA.previous
}

func (l *LinkedList[V]) Sort(cmp func(a, b V) int) *LinkedList[V] {
	items := make([]*LinkedListEntry[V], l.Len())
	l.forEach(func(i int, entry *LinkedListEntry[V]) {
		items[i] = entry
	})
	slices.SortFunc(items, func(a, b *LinkedListEntry[V]) int {
		return cmp(a.data, b.data)
	})
	newList := NewLinkedList[V]()
	for _, item := range items {
		newList.add(item)
	}
	return newList
}

func (l *LinkedList[V]) Len() int {
	return l.size
}

func (l *LinkedList[V]) forEach(f func(int, *LinkedListEntry[V])) {
	i := 0
	ptr := l.Tail()
	for ptr != nil {
		f(i, ptr)
		ptr = ptr.next
		i++
	}
}
