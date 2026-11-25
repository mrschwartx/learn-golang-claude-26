// Latihan 1 - Modul 1: versi Hello yang memakai nama dan loop 5 kali.
package main

import "fmt"

func main() {
	nama := "Pembelajar Go"
	fmt.Printf("Halo, Go! Saya %s.\n", nama)

	for i := 1; i <= 5; i++ {
		fmt.Println("Iterasi ke-", i)
	}
}
