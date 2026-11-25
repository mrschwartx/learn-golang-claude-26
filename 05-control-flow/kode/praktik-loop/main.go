// loop.go
package main

import "fmt"

func main() {
	// FizzBuzz dengan range integer (Go 1.22+)
	for i := range 15 {
		n := i + 1
		switch {
		case n%15 == 0:
			fmt.Print("FizzBuzz ")
		case n%3 == 0:
			fmt.Print("Fizz ")
		case n%5 == 0:
			fmt.Print("Buzz ")
		default:
			fmt.Print(n, " ")
		}
	}
	fmt.Println()

	// while-style dan loop tak hingga
	sisa := 3
	for sisa > 0 {
		sisa--
	}
	hitung := 0
	for {
		hitung++
		if hitung == 5 {
			break
		}
	}
	fmt.Println(sisa, hitung)

	// string: byte vs rune
	s := "hé"
	for i := 0; i < len(s); i++ {
		fmt.Print(s[i], " ")
	}
	fmt.Println()
	for i, r := range s {
		fmt.Printf("%d:%c ", i, r)
	}
	fmt.Println()

	// label
luar:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if j == 2 {
				continue luar
			}
			if i == 2 {
				break luar
			}
			fmt.Printf("(%d,%d) ", i, j)
		}
	}
	fmt.Println()
}
