// Latihan 3: frekuensi kata + 3 teratas. Map tidak bisa diurutkan -> pindahkan ke slice.
package main

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
	"unicode"
)

type Frek struct {
	Kata   string
	Jumlah int
}

func main() {
	paragraf := `Go itu sederhana. Go itu cepat, dan Go itu menyenangkan!
Belajar Go berarti belajar menulis kode yang sederhana dan mudah dibaca.`

	kata := strings.FieldsFunc(strings.ToLower(paragraf), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	hitung := make(map[string]int)
	for _, k := range kata {
		hitung[k]++
	}

	daftar := make([]Frek, 0, len(hitung))
	for k, n := range hitung {
		daftar = append(daftar, Frek{k, n})
	}
	// urut jumlah turun; seri diputus alfabetis agar hasil deterministik
	slices.SortFunc(daftar, func(a, b Frek) int {
		if c := cmp.Compare(b.Jumlah, a.Jumlah); c != 0 {
			return c
		}
		return strings.Compare(a.Kata, b.Kata)
	})

	fmt.Println("3 kata teratas:")
	for i, f := range daftar[:3] {
		fmt.Printf("%d. %-10s %d\n", i+1, f.Kata, f.Jumlah)
	}
}
