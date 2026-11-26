// generik.go
package main

import (
	"cmp"
	"fmt"
	"strconv"
)

type Angka interface {
	~int | ~int64 | ~float64
}

type Rupiah int

func Jumlah[T Angka](xs []T) T {
	var total T
	for _, x := range xs {
		total += x
	}
	return total
}

func Map[T, U any](xs []T, f func(T) U) []U {
	out := make([]U, 0, len(xs))
	for _, x := range xs {
		out = append(out, f(x))
	}
	return out
}

func Filter[T any](xs []T, ok func(T) bool) []T {
	var out []T
	for _, x := range xs {
		if ok(x) {
			out = append(out, x)
		}
	}
	return out
}

func Reduce[T, A any](xs []T, awal A, f func(A, T) A) A {
	acc := awal
	for _, x := range xs {
		acc = f(acc, x)
	}
	return acc
}

func Maks[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Cari[T comparable](xs []T, target T) int {
	for i, x := range xs {
		if x == target {
			return i
		}
	}
	return -1
}

func main() {
	fmt.Println(Jumlah([]int{1, 2, 3}), Jumlah([]float64{1.5, 2.5}), Jumlah([]Rupiah{1000, 500}))
	fmt.Println(Map([]int{1, 2, 3}, strconv.Itoa))
	genap := Filter([]int{1, 2, 3, 4, 5, 6}, func(n int) bool { return n%2 == 0 })
	fmt.Println(genap, Reduce(genap, 0, func(a, n int) int { return a + n*n }))
	fmt.Println(Maks("apel", "jeruk"), Maks(3, 7), Maks(2.5, 1.5))
	fmt.Println(Cari([]string{"a", "b", "c"}, "c"), Cari([]int{1, 2}, 9))
}
