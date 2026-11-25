// Latihan 5: strconv.Atoi dengan penanganan error ramah pengguna.
package main

import (
	"errors"
	"fmt"
	"strconv"
)

func parse(input string) {
	n, err := strconv.Atoi(input)
	if err != nil {
		var numErr *strconv.NumError
		if errors.As(err, &numErr) && numErr.Err == strconv.ErrSyntax {
			fmt.Printf("%q bukan angka yang valid. Masukkan angka bulat, misalnya 12.\n", input)
		} else if errors.Is(err, strconv.ErrRange) {
			fmt.Printf("%q terlalu besar atau terlalu kecil.\n", input)
		} else {
			fmt.Println("terjadi kesalahan:", err)
		}
		fmt.Println("  (error asli:", err, ")")
		return
	}
	fmt.Printf("Berhasil: %d (n+1 = %d)\n", n, n+1)
}

func main() {
	for _, in := range []string{"12x", "12", "", "99999999999999999999"} {
		parse(in)
	}
}
