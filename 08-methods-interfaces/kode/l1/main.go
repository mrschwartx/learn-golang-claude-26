// Latihan 1: tambah Segitiga tanpa mengubah totalLuas.
package main

import (
	"fmt"
	"math"
)

type Bentuk interface {
	Luas() float64
	Keliling() float64
}

type Persegi struct{ Sisi float64 }
type Lingkaran struct{ R float64 }

// Tipe baru: cukup mendefinisikan method. Tidak ada "implements".
type Segitiga struct{ A, B, C float64 }

func (p Persegi) Luas() float64       { return p.Sisi * p.Sisi }
func (p Persegi) Keliling() float64   { return 4 * p.Sisi }
func (l Lingkaran) Luas() float64     { return math.Pi * l.R * l.R }
func (l Lingkaran) Keliling() float64 { return 2 * math.Pi * l.R }

func (s Segitiga) Keliling() float64 { return s.A + s.B + s.C }
func (s Segitiga) Luas() float64 { // rumus Heron
	p := s.Keliling() / 2
	return math.Sqrt(p * (p - s.A) * (p - s.B) * (p - s.C))
}

var _ Bentuk = Segitiga{} // pengecekan saat kompilasi

// TIDAK diubah sama sekali sejak sebelum Segitiga ada.
func totalLuas(daftar []Bentuk) float64 {
	total := 0.0
	for _, b := range daftar {
		total += b.Luas()
	}
	return total
}

func main() {
	daftar := []Bentuk{Persegi{2}, Lingkaran{1}, Segitiga{3, 4, 5}}
	for _, b := range daftar {
		fmt.Printf("%-24T luas=%.4f keliling=%.4f\n", b, b.Luas(), b.Keliling())
	}
	fmt.Printf("total luas: %.4f (4 + pi + 6)\n", totalLuas(daftar))
}
