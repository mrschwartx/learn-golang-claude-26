// Latihan 2: Bilangan (~int | ~float64) dan Rata[T Bilangan] float64.
package main

import "fmt"

type Bilangan interface {
	~int | ~float64
}

func Rata[T Bilangan](xs []T) float64 {
	if len(xs) == 0 {
		return 0
	}
	var total T
	for _, x := range xs {
		total += x
	}
	return float64(total) / float64(len(xs))
}

type Skor int

// Tanpa tilde: hanya tipe int/float64 persis
type BilanganKetat interface{ int | float64 }

func RataKetat[T BilanganKetat](xs []T) float64 {
	var total T
	for _, x := range xs {
		total += x
	}
	return float64(total) / float64(len(xs))
}

func main() {
	fmt.Println(Rata([]int{1, 2, 3, 4}))
	fmt.Println(Rata([]float64{1.5, 2.5}))
	fmt.Println(Rata([]Skor{80, 90, 100})) // Skor memenuhi ~int
	fmt.Println(Rata([]int(nil)))

	fmt.Println(RataKetat([]int{2, 4}))
	// fmt.Println(RataKetat([]Skor{80, 90})) // error: Skor does not satisfy BilanganKetat (possibly missing ~ for int in BilanganKetat)
}
