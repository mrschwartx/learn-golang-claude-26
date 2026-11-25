// fungsi.go
package main

import (
	"errors"
	"fmt"
)

func bagi(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("pembagi nol")
	}
	return a / b, nil
}

func ukuran(luas, sisi float64) (panjang, lebar float64) {
	panjang = luas / sisi
	lebar = sisi
	return
}

func jumlah(angka ...int) int {
	total := 0
	for _, n := range angka {
		total += n
	}
	return total
}

func main() {
	if h, err := bagi(10, 4); err == nil {
		fmt.Println("hasil:", h)
	}
	if _, err := bagi(1, 0); err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println(ukuran(20, 4))
	fmt.Println(jumlah(), jumlah(1, 2, 3), jumlah([]int{4, 5, 6}...))
}
