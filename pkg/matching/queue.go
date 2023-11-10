package matching

import (
	"sort"
	"sync"
	"time"
)

type Matcher func(a, b any) bool

type Item struct {
	Item      any
	Timestamp time.Time
}

type Queue struct {
	sync.Mutex
	first         []Item
	second        []Item
	matchingFunc  Matcher
	maxLen        int
	removeMatched bool
}

func NewQueue(maxLen int, matchingFunc Matcher, removeMatched bool) *Queue {
	ret := &Queue{
		maxLen:        maxLen,
		matchingFunc:  matchingFunc,
		removeMatched: removeMatched,
	}
	return ret
}

func removeElement(arr []Item, index int) []Item {
	if index < 0 {
		return arr
	}
	N := len(arr)
	if index >= N {
		return arr
	}
	if N == 0 {
		return arr
	}
	if N == 1 {
		return arr[:0]
	}
	if index < N-1 {
		copy(arr[index:], arr[index+1:])
	}
	return arr[:N-1]
}

func addItem(arr []Item, elem Item, maxLen int) []Item {
	N := len(arr)
	if N < maxLen {
		arr = append(arr, elem)
		return arr
	}

	copy(arr, arr[1+N-maxLen:])
	arr = arr[:maxLen]
	arr[maxLen-1] = elem
	return arr
}

func (q *Queue) add(a, b *[]Item, data any, t time.Time) []Item {
	q.Lock()
	defer q.Unlock()

	elem := Item{data, t}
	matches := make([]Item, 0, 3)
	indices := make([]int, 0, 3)
	for i, item := range *b {
		if q.matchingFunc(data, item.Item) {
			matches = append(matches, item)
			if q.removeMatched {
				indices = append(indices, i)
			}
		}
	}

	if len(matches) == 0 {
		*a = addItem(*a, elem, q.maxLen)
		return nil
	}

	if q.removeMatched {
		sort.Sort(sort.Reverse(sort.IntSlice(indices)))
		for _, idx := range indices {
			*b = removeElement(*b, idx)
		}
	}

	if !q.removeMatched {
		*a = addItem(*a, elem, q.maxLen)
	}

	return matches
}

func (q *Queue) AddFirst(data any, t time.Time) []Item {
	return q.add(&q.first, &q.second, data, t)
}

func (q *Queue) AddSecond(data any, t time.Time) []Item {
	return q.add(&q.second, &q.first, data, t)
}
