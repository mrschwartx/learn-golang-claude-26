// Latihan 3: tabel perkalian 1-9, berhenti total saat hasil > 50 (break berlabel).
package main

import "fmt"

func main() {
tabel:
	for i := 1; i <= 9; i++ {
		for j := 1; j <= 9; j++ {
			hasil := i * j
			if hasil > 50 {
				fmt.Printf("\nberhenti: %d x %d = %d (> 50)\n", i, j, hasil)
				break tabel
			}
			fmt.Printf("%3d", hasil)
		}
		fmt.Println()
	}

	// tanpa label: break hanya mengakhiri loop dalam
	fmt.Println("\ntanpa label (hanya loop dalam yang berhenti):")
	for i := 1; i <= 9; i++ {
		for j := 1; j <= 9; j++ {
			if i*j > 50 {
				break
			}
			fmt.Printf("%3d", i*j)
		}
		fmt.Println()
	}
}
