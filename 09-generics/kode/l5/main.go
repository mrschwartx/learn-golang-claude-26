// Latihan 5: pada toolchain Go < 1.27, method generik ditolak. Lihat ../l5-error-*.
// Di sini ditunjukkan pola pengganti yang valid di semua versi: fungsi generik biasa.
package main

import "fmt"

type Stack[T any] struct{ items []T }

func (s *Stack[T]) Push(x T) { s.items = append(s.items, x) }

// Pengganti method generik: fungsi generik dengan Stack sebagai argumen.
func Transform[T, U any](s *Stack[T], f func(T) U) []U {
	out := make([]U, 0, len(s.items))
	for _, x := range s.items {
		out = append(out, f(x))
	}
	return out
}

// Interface yang boleh: method konkret (tanpa type parameter di method)
type Pengubah interface {
	TransformInt(f func(int) string) []string
}

func (s *Stack[T]) TransformInt(f func(int) string) []string { return nil } // hanya contoh signature

func main() {
	s := &Stack[string]{}
	s.Push("a")
	s.Push("bb")
	fmt.Println(Transform(s, func(x string) int { return len(x) }))
}
