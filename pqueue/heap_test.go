package pqueue_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/yairp7/gocache/pqueue"
	"gopkg.in/stretchr/testify.v1/assert"
)

type TestCase struct {
	input    []int
	expected []int
}

var tests = map[string][]TestCase{
	"Touch": {
		{
			input:    []int{2, 3, 4, 5, 10},
			expected: []int{2, 4, 80, 5, 10, 3},
		},
	},
}

func Test_MinHeap(t *testing.T) {
	// t.SkipNow()
	for i, testCase := range tests["Touch"] {
		t.Run(fmt.Sprintf("Touch %d", i), func(t *testing.T) {
			q := pqueue.NewMinHeap[int]()
			nodes := make([]*pqueue.HeapNode[int], len(testCase.input))
			for i, e := range testCase.input {
				nodes[i] = q.Push(e)
				time.Sleep(time.Millisecond)
			}

			q.Touch(nodes[1])
			q.Touch(nodes[1])
			q.Touch(nodes[1])
			q.Touch(nodes[1])

			q.Touch(nodes[4])
			q.Touch(nodes[4])
			q.Touch(nodes[4])

			n := q.Push(80)

			q.Touch(nodes[3])
			q.Touch(nodes[3])

			q.Touch(n)

			for _, expected := range testCase.expected {
				node := q.Pop()
				if node == nil {
					break
				}
				assert.Equal(t, expected, node.Data)
			}
		})
	}
}
