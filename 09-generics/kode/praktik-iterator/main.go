// iterator.go
package main

import (
	"fmt"
	"iter"
	"maps"
	"slices"
)

func Hitung(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; i < n; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

func Filter[T any](seq iter.Seq[T], ok func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range seq {
			if ok(v) && !yield(v) {
				return
			}
		}
	}
}

func Fibonacci() iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		a, b := 0, 1
		for i := 0; ; i++ {
			if !yield(i, a) {
				return
			}
			a, b = b, a+b
		}
	}
}

func main() {
	for v := range Hitung(10) {
		if v == 4 {
			break
		}
		fmt.Print(v, " ")
	}
	fmt.Println()

	genap := Filter(Hitung(10), func(n int) bool { return n%2 == 0 })
	fmt.Println(slices.Collect(genap))

	for i, f := range Fibonacci() {
		if i == 10 {
			break
		}
		fmt.Print(f, " ")
	}
	fmt.Println()

	m := map[string]int{"b": 2, "a": 1, "c": 3}
	fmt.Println(slices.Sorted(maps.Keys(m)))

	next, stop := iter.Pull(Hitung(3))
	defer stop()
	for {
		v, ok := next()
		if !ok {
			break
		}
		fmt.Print("pull:", v, " ")
	}
	fmt.Println()
}
