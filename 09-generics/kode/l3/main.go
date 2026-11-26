// Latihan 3: Set[T comparable] berbasis map[T]struct{}.
package main

import (
	"fmt"
	"unsafe"
)

type Set[T comparable] struct{ m map[T]struct{} }

func NewSet[T comparable](elem ...T) *Set[T] {
	s := &Set[T]{m: make(map[T]struct{}, len(elem))}
	for _, e := range elem {
		s.Add(e)
	}
	return s
}

func (s *Set[T]) Add(x T)      { s.m[x] = struct{}{} }
func (s *Set[T]) Has(x T) bool { _, ok := s.m[x]; return ok }
func (s *Set[T]) Remove(x T)   { delete(s.m, x) }
func (s *Set[T]) Len() int     { return len(s.m) }

func main() {
	s := NewSet("a", "b")
	s.Add("c")
	s.Add("a") // duplikat diabaikan
	fmt.Println("Len:", s.Len(), "| Has(a):", s.Has("a"), "| Has(z):", s.Has("z"))
	s.Remove("a")
	fmt.Println("setelah Remove(a): Len:", s.Len(), "| Has(a):", s.Has("a"))

	angka := NewSet(1, 2, 3, 2, 1)
	fmt.Println("Set[int]:", angka.Len())

	// kenapa struct{}?
	fmt.Println("ukuran struct{} :", unsafe.Sizeof(struct{}{}), "byte")
	fmt.Println("ukuran bool     :", unsafe.Sizeof(true), "byte")
	fmt.Println("ukuran int      :", unsafe.Sizeof(0), "byte")
}
