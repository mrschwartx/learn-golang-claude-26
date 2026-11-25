// Latihan 6: map nilai siswa dicetak terurut berdasarkan nama.
package main

import (
	"fmt"
	"maps"
	"slices"
)

func main() {
	nilai := map[string]int{"Budi": 80, "Andi": 90, "Citra": 85, "Dewi": 70, "Eko": 95}

	fmt.Print("urutan iterasi (berubah antar-eksekusi, jalankan berkali-kali):")
	for k := range nilai {
		fmt.Print(" ", k)
	}
	fmt.Println()

	fmt.Println("terurut berdasarkan nama:")
	for _, k := range slices.Sorted(maps.Keys(nilai)) {
		fmt.Printf("  %-6s %d\n", k, nilai[k])
	}
}
