// closure.go
package main

import (
	"fmt"
	"sort"
)

func penghitung() func() int {
	n := 0
	return func() int {
		n++
		return n
	}
}

type Operasi func(int, int) int

func terapkan(op Operasi, a, b int) int { return op(a, b) }

func main() {
	a, b := penghitung(), penghitung()
	fmt.Println(a(), a(), a(), b())

	fmt.Println(terapkan(func(x, y int) int { return x * y }, 6, 7))

	hasil := func(x int) int { return x * x }(9) // langsung dipanggil
	fmt.Println(hasil)

	kata := []string{"pisang", "apel", "jeruk"}
	sort.Slice(kata, func(i, j int) bool { return len(kata[i]) < len(kata[j]) })
	fmt.Println(kata)
}
