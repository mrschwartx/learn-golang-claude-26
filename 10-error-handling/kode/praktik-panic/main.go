// panic.go
package main

import (
	"fmt"
	"runtime/debug"
	"strings"
)

func ambil(data []int, i int) int {
	return data[i] // panic jika i di luar rentang
}

// Aman mengubah panic menjadi error di batas sistem.
func ambilAman(data []int, i int) (hasil int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("pulih dari panic: %v", r)
			stack := string(debug.Stack())
			fmt.Println("stack berisi 'ambil':", strings.Contains(stack, "ambil"))
		}
	}()
	return ambil(data, i), nil
}

func main() {
	fmt.Println(ambilAman([]int{1, 2, 3}, 1))
	fmt.Println(ambilAman([]int{1, 2, 3}, 5))

	defer fmt.Println("defer tetap dijalankan saat panic")
	var m map[string]int
	m["x"] = 1 // panic: assignment to entry in nil map (tidak dipulihkan)
}
