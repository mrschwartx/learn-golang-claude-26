// Latihan 4: len vs jumlah rune, dan membalik string dengan []rune.
package main

import (
	"fmt"
	"unicode/utf8"
)

func balik(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func balikSalah(s string) string { // membalik BYTE: merusak UTF-8
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func main() {
	s := "Selamat pagi 🌅"
	fmt.Println("len (byte):", len(s))
	fmt.Println("rune      :", utf8.RuneCountInString(s))
	fmt.Println("dibalik   :", balik(s))
	fmt.Printf("balik byte: %q (valid UTF-8: %t)\n", balikSalah(s), utf8.ValidString(balikSalah(s)))
}
