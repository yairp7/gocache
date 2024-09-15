package pqueue

import (
	"fmt"
	"testing"
	"time"

	"gopkg.in/stretchr/testify.v1/assert"
)

type TestCase struct {
	input    []int
	expected []int
}

var tests = map[string][]TestCase{
	"Basic": {
		{
			input:    []int{2, 3, 4, 5, 10},
			expected: []int{13, 2, 3, 1, 4, 10, 5},
		},
	},
}

func Test_PriorityQueue(t *testing.T) {
	for i, testCase := range tests["Basic"] {
		t.Run(fmt.Sprintf("Basic %d", i), func(t *testing.T) {
			q := newPriorityQueue[int]()
			entries := make([]*pQueueEntry[int], 0)
			for _, e := range testCase.input {
				entries = append(entries, q.push(e, e))
				time.Sleep(time.Millisecond)
			}

			q.push(1, 3)
			time.Sleep(time.Millisecond)

			q.setWeight(3, 20)
			time.Sleep(time.Millisecond)

			q.push(13, 1)
			time.Sleep(time.Millisecond)

			for _, expected := range testCase.expected {
				e := q.pop()
				assert.Equal(t, expected, e.value)
			}
		})
	}
}
