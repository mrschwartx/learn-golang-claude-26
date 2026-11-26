// Latihan 1: BacaAngka membungkus strconv.Atoi dengan %w; errors.As ke *strconv.NumError.
package main

import (
	"errors"
	"fmt"
	"strconv"
)

func BacaAngka(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("baca angka dari %q: %w", s, err)
	}
	return n, nil
}

func main() {
	for _, in := range []string{"123", "12x", "", "99999999999999999999"} {
		n, err := BacaAngka(in)
		if err == nil {
			fmt.Println("OK:", n)
			continue
		}
		fmt.Println("error:", err)

		var ne *strconv.NumError
		if errors.As(err, &ne) {
			fmt.Printf("  Func=%q Num=%q Err=%v\n", ne.Func, ne.Num, ne.Err)
			fmt.Println("  syntax?", errors.Is(err, strconv.ErrSyntax), "| range?", errors.Is(err, strconv.ErrRange))
		}
	}
}
