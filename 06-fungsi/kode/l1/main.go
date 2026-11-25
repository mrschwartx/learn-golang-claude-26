// Latihan 1: MinMax variadic dengan error bila tanpa argumen.
package main

import (
	"errors"
	"fmt"
)

var ErrKosong = errors.New("minmax: tidak ada argumen")

func MinMax(angka ...int) (min, max int, err error) {
	if len(angka) == 0 {
		return 0, 0, ErrKosong
	}
	min, max = angka[0], angka[0]
	for _, n := range angka[1:] {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}
	return min, max, nil
}

func main() {
	if lo, hi, err := MinMax(4, -2, 9, 7); err == nil {
		fmt.Println("argumen biasa:", lo, hi)
	}

	s := []int{10, 3, 25, 8}
	lo, hi, err := MinMax(s...)
	fmt.Println("slice s...   :", lo, hi, err)

	lo, hi, err = MinMax(42)
	fmt.Println("satu argumen :", lo, hi, err)

	_, _, err = MinMax()
	fmt.Println("tanpa argumen:", err, "| errors.Is:", errors.Is(err, ErrKosong))
}
