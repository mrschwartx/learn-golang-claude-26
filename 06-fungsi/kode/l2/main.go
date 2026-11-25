// Latihan 2: closure Penjumlah dengan state terpisah.
package main

import "fmt"

func Penjumlah(awal int) func(int) int {
	total := awal
	return func(n int) int {
		total += n
		return total
	}
}

func main() {
	a := Penjumlah(0)
	b := Penjumlah(100)

	fmt.Println("a(5)  =", a(5))
	fmt.Println("a(10) =", a(10))
	fmt.Println("b(1)  =", b(1))
	fmt.Println("a(1)  =", a(1), "(a tidak terpengaruh b)")
	fmt.Println("b(1)  =", b(1), "(b tidak terpengaruh a)")
}
