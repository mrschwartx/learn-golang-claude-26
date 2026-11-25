// Latihan 5: loop indeks (byte) vs range (rune).
package main

import "fmt"

func main() {
	s := "Go 世界 🌍"
	fmt.Println("len (byte):", len(s))

	fmt.Println("\nloop indeks (per byte):")
	n := 0
	for i := 0; i < len(s); i++ {
		fmt.Printf("%d:%02x ", i, s[i])
		n++
	}
	fmt.Println("\niterasi:", n)

	fmt.Println("\nrange (per rune):")
	m := 0
	for i, r := range s {
		fmt.Printf("%d:%c(%d byte) ", i, r, len(string(r)))
		m++
	}
	fmt.Println("\niterasi:", m)
}
