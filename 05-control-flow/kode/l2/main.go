// Latihan 2: FizzBuzz sebagai fungsi, versi switch dan versi if-else.
package main

import (
	"fmt"
	"slices"
	"strconv"
)

func FizzBuzzSwitch(n int) []string {
	hasil := make([]string, 0, n)
	for i := 1; i <= n; i++ {
		switch {
		case i%15 == 0:
			hasil = append(hasil, "FizzBuzz")
		case i%3 == 0:
			hasil = append(hasil, "Fizz")
		case i%5 == 0:
			hasil = append(hasil, "Buzz")
		default:
			hasil = append(hasil, strconv.Itoa(i))
		}
	}
	return hasil
}

func FizzBuzzIf(n int) []string {
	hasil := make([]string, 0, n)
	for i := 1; i <= n; i++ {
		if i%15 == 0 {
			hasil = append(hasil, "FizzBuzz")
		} else if i%3 == 0 {
			hasil = append(hasil, "Fizz")
		} else if i%5 == 0 {
			hasil = append(hasil, "Buzz")
		} else {
			hasil = append(hasil, strconv.Itoa(i))
		}
	}
	return hasil
}

func main() {
	a, b := FizzBuzzSwitch(15), FizzBuzzIf(15)
	fmt.Println(a)
	fmt.Println("kedua versi sama:", slices.Equal(a, b))
}
