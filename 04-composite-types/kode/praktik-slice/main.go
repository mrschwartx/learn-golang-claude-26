// slice.go
package main

import (
	"fmt"
	"slices"
)

func main() {
	s := []int{1, 2, 3, 4, 5}
	t := s[1:3]
	fmt.Println("t:", t, "len:", len(t), "cap:", cap(t))
	t[0] = 99
	fmt.Println("s setelah t[0]=99:", s)

	// append pada t masih muat di kapasitas -> MENIMPA s[3]
	t = append(t, 100)
	fmt.Println("s setelah append ke t:", s)

	// full slice expression: cap dibatasi, append memaksa array baru
	u := s[1:3:3]
	u = append(u, 200)
	fmt.Println("s tidak berubah:", s, "u:", u)

	var x []int
	prev := -1
	for i := 0; i < 2000; i++ {
		x = append(x, i)
		if cap(x) != prev {
			fmt.Print(cap(x), " ")
			prev = cap(x)
		}
	}
	fmt.Println()

	angka := []int{5, 2, 8, 1}
	slices.Sort(angka)
	fmt.Println(angka, slices.Contains(angka, 8), slices.Index(angka, 5), slices.Max(angka))
	angka = slices.Delete(angka, 1, 3)
	fmt.Println(angka)
}
