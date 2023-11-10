package matching

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func Test_removeElement(t *testing.T) {
	tests := []struct {
		name  string
		arr   []Item
		index int
		want  []Item
	}{
		{"no elements", []Item{}, 0, []Item{}},
		{"no elements[10]", []Item{}, 10, []Item{}},
		{"no elements[-1]", []Item{}, -1, []Item{}},
		{"1 element", []Item{{Item: 1}}, 0, []Item{}},
		{"1 element[10]", []Item{{Item: 1}}, 10, []Item{{Item: 1}}},
		{"1 element[-1]", []Item{{Item: 1}}, -1, []Item{{Item: 1}}},
		{"3 element[0]", []Item{{Item: 1}, {Item: 2}, {Item: 3}}, 0, []Item{{Item: 2}, {Item: 3}}},
		{"3 element[1]", []Item{{Item: 1}, {Item: 2}, {Item: 3}}, 1, []Item{{Item: 1}, {Item: 3}}},
		{"3 element[2]", []Item{{Item: 1}, {Item: 2}, {Item: 3}}, 2, []Item{{Item: 1}, {Item: 2}}},
		{"3 element[3]", []Item{{Item: 1}, {Item: 2}, {Item: 3}}, 3, []Item{{Item: 1}, {Item: 2}, {Item: 3}}},
		{"3 element[10]", []Item{{Item: 1}, {Item: 2}, {Item: 3}}, 10, []Item{{Item: 1}, {Item: 2}, {Item: 3}}},
		{"3 element[-1]", []Item{{Item: 1}, {Item: 2}, {Item: 3}}, -1, []Item{{Item: 1}, {Item: 2}, {Item: 3}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeElement(tt.arr, tt.index); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeElement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatchingQueue_add_with_remove(t *testing.T) {
	tests := []struct {
		name    string
		dst     []Item
		src     []Item
		data    any
		inspect func(a, b []Item) error
		want    []Item
	}{
		{
			"not found",
			[]Item{{Item: 1}, {Item: 2}},
			[]Item{{Item: 4}, {Item: 5}},
			3,
			func(a, b []Item) error {
				if reflect.DeepEqual(a, []Item{{Item: 1}, {Item: 2}, {Item: 3}}) && reflect.DeepEqual(b, []Item{{Item: 4}, {Item: 5}}) {
					return nil
				}
				return fmt.Errorf("bad a,b: %v, %v", a, b)
			},
			nil,
		},
		{
			"not found long",
			[]Item{{Item: 1}, {Item: 2}, {Item: 6}, {Item: 7}, {Item: 8}},
			[]Item{{Item: 4}, {Item: 5}},
			3,
			func(a, b []Item) error {
				if reflect.DeepEqual(a, []Item{{Item: 2}, {Item: 6}, {Item: 7}, {Item: 8}, {Item: 3}}) && reflect.DeepEqual(b, []Item{{Item: 4}, {Item: 5}}) {
					return nil
				}
				return fmt.Errorf("bad a,b: %v, %v", a, b)
			},
			nil,
		},
		{
			"not found longer",
			[]Item{{Item: 1}, {Item: 2}, {Item: 6}, {Item: 7}, {Item: 8}, {Item: 9}, {Item: 10}},
			[]Item{{Item: 4}, {Item: 5}},
			3,
			func(a, b []Item) error {
				if reflect.DeepEqual(a, []Item{{Item: 7}, {Item: 8}, {Item: 9}, {Item: 10}, {Item: 3}}) && reflect.DeepEqual(b, []Item{{Item: 4}, {Item: 5}}) {
					return nil
				}
				return fmt.Errorf("bad a,b: %v, %v", a, b)
			},
			nil,
		},
		{
			"found",
			[]Item{{Item: 1}, {Item: 2}},
			[]Item{{Item: 4}, {Item: 5}},
			4,
			func(a, b []Item) error {
				if reflect.DeepEqual(a, []Item{{Item: 1}, {Item: 2}}) && reflect.DeepEqual(b, []Item{{Item: 5}}) {
					return nil
				}
				return fmt.Errorf("bad a,b: %v, %v", a, b)
			},
			[]Item{{Item: 4}},
		},
		{
			"found two",
			[]Item{{Item: 1}, {Item: 2}},
			[]Item{{Item: 4}, {Item: 8}, {Item: 5}},
			4,
			func(a, b []Item) error {
				if reflect.DeepEqual(a, []Item{{Item: 1}, {Item: 2}}) && reflect.DeepEqual(b, []Item{{Item: 5}}) {
					return nil
				}
				return fmt.Errorf("bad a,b: %v, %v", a, b)
			},
			[]Item{{Item: 4}, {Item: 8}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := NewQueue(5, func(a, b any) bool { return a.(int) == b.(int) || a.(int) == b.(int)*2 || a.(int)*2 == b.(int) }, true)
			got := q.add(&tt.dst, &tt.src, tt.data, time.Time{})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MatchingQueue.add() got = %v, want %v", got, tt.want)
			}
			if err := tt.inspect(tt.dst, tt.src); err != nil {
				t.Errorf("Inspect failed: %v", err.Error())
			}
		})
	}
}

func TestMatchingQueue_add_no_remove(t *testing.T) {
	tests := []struct {
		name    string
		dst     []Item
		src     []Item
		data    any
		inspect func(a, b []Item) error
		want    []Item
	}{
		{
			"not found",
			[]Item{{Item: 1}, {Item: 2}},
			[]Item{{Item: 4}, {Item: 5}},
			3,
			func(a, b []Item) error {
				if reflect.DeepEqual(a, []Item{{Item: 1}, {Item: 2}, {Item: 3}}) && reflect.DeepEqual(b, []Item{{Item: 4}, {Item: 5}}) {
					return nil
				}
				return fmt.Errorf("bad a,b: %v, %v", a, b)
			},
			nil,
		},
		{
			"not found long",
			[]Item{{Item: 1}, {Item: 2}, {Item: 6}, {Item: 7}, {Item: 8}},
			[]Item{{Item: 4}, {Item: 5}},
			3,
			func(a, b []Item) error {
				if reflect.DeepEqual(a, []Item{{Item: 2}, {Item: 6}, {Item: 7}, {Item: 8}, {Item: 3}}) && reflect.DeepEqual(b, []Item{{Item: 4}, {Item: 5}}) {
					return nil
				}
				return fmt.Errorf("bad a,b: %v, %v", a, b)
			},
			nil,
		},
		{
			"not found longer",
			[]Item{{Item: 1}, {Item: 2}, {Item: 6}, {Item: 7}, {Item: 8}, {Item: 9}, {Item: 10}},
			[]Item{{Item: 4}, {Item: 5}},
			3,
			func(a, b []Item) error {
				if reflect.DeepEqual(a, []Item{{Item: 7}, {Item: 8}, {Item: 9}, {Item: 10}, {Item: 3}}) && reflect.DeepEqual(b, []Item{{Item: 4}, {Item: 5}}) {
					return nil
				}
				return fmt.Errorf("bad a,b: %v, %v", a, b)
			},
			nil,
		},
		{
			"found",
			[]Item{{Item: 1}, {Item: 2}},
			[]Item{{Item: 4}, {Item: 5}},
			4,
			func(a, b []Item) error {
				if reflect.DeepEqual(a, []Item{{Item: 1}, {Item: 2}, {Item: 4}}) && reflect.DeepEqual(b, []Item{{Item: 4}, {Item: 5}}) {
					return nil
				}
				return fmt.Errorf("bad a,b: %v, %v", a, b)
			},
			[]Item{{Item: 4}},
		},
		{
			"found two",
			[]Item{{Item: 1}, {Item: 2}},
			[]Item{{Item: 4}, {Item: 8}, {Item: 5}},
			4,
			func(a, b []Item) error {
				if reflect.DeepEqual(a, []Item{{Item: 1}, {Item: 2}, {Item: 4}}) && reflect.DeepEqual(b, []Item{{Item: 4}, {Item: 8}, {Item: 5}}) {
					return nil
				}
				return fmt.Errorf("bad a,b: %v, %v", a, b)
			},
			[]Item{{Item: 4}, {Item: 8}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := NewQueue(5, func(a, b any) bool { return a.(int) == b.(int) || a.(int) == b.(int)*2 || a.(int)*2 == b.(int) }, false)
			got := q.add(&tt.dst, &tt.src, tt.data, time.Time{})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MatchingQueue.add() got = %v, want %v", got, tt.want)
			}
			if err := tt.inspect(tt.dst, tt.src); err != nil {
				t.Errorf("Inspect failed: %v", err.Error())
			}
		})
	}
}
