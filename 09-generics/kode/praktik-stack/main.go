// stack.go
package main

import "fmt"

type Stack[T any] struct{ items []T }

func (s *Stack[T]) Push(x T) { s.items = append(s.items, x) }

func (s *Stack[T]) Pop() (T, bool) {
	var nol T
	if len(s.items) == 0 {
		return nol, false
	}
	x := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return x, true
}

// Generic method (func (s *Stack[T]) Transform[U any](...)) butuh Go 1.27+.
// Di Go < 1.27 gunakan fungsi generik biasa yang menerima stack sebagai argumen:
func Transform[T, U any](s *Stack[T], f func(T) U) []U {
	out := make([]U, 0, len(s.items))
	for _, x := range s.items {
		out = append(out, f(x))
	}
	return out
}

type Pair[K comparable, V any] struct {
	Key K
	Val V
}

type Cache[K comparable, V any] struct {
	data map[K]V
}

func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{data: make(map[K]V)}
}

func (c *Cache[K, V]) Set(k K, v V) { c.data[k] = v }

func (c *Cache[K, V]) Get(k K) (V, bool) {
	v, ok := c.data[k]
	return v, ok
}

func main() {
	s := &Stack[string]{}
	s.Push("a")
	s.Push("bb")
	fmt.Println(Transform(s, func(x string) int { return len(x) }))
	v, ok := s.Pop()
	fmt.Println(v, ok)

	c := NewCache[string, Pair[int, bool]]()
	c.Set("x", Pair[int, bool]{1, true})
	fmt.Println(c.Get("x"))
	fmt.Println(c.Get("y"))
}
