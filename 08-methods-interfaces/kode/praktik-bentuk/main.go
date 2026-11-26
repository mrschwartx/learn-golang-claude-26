// bentuk.go
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

func (p Persegi) Luas() float64       { return p.Sisi * p.Sisi }
func (p Persegi) Keliling() float64   { return 4 * p.Sisi }
func (l Lingkaran) Luas() float64     { return math.Pi * l.R * l.R }
func (l Lingkaran) Keliling() float64 { return 2 * math.Pi * l.R }

var _ Bentuk = Persegi{}
var _ Bentuk = (*Lingkaran)(nil) // *T juga memenuhi interface bila T memenuhinya

func totalLuas(daftar []Bentuk) float64 {
	total := 0.0
	for _, b := range daftar {
		total += b.Luas()
	}
	return total
}

func deskripsi(b Bentuk) string {
	switch v := b.(type) {
	case Persegi:
		return fmt.Sprintf("persegi sisi %.1f", v.Sisi)
	case Lingkaran:
		return fmt.Sprintf("lingkaran jari-jari %.1f", v.R)
	default:
		return "bentuk tidak dikenal"
	}
}

func main() {
	daftar := []Bentuk{Persegi{2}, Lingkaran{1}}
	fmt.Printf("total luas: %.2f\n", totalLuas(daftar))
	for _, b := range daftar {
		fmt.Println(deskripsi(b))
	}

	var x any = "go"
	if s, ok := x.(string); ok {
		fmt.Println("string:", s)
	}
	if _, ok := x.(int); !ok {
		fmt.Println("bukan int")
	}
}
