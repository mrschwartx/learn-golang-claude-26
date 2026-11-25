// Latihan 6: jebakan string(65). Jalankan "go vet ./l6" untuk melihat peringatannya.
package main

import (
	"fmt"
	"strconv"
)

func main() {
	n := 65
	fmt.Println(string(n))           // "A"  <- JEBAKAN (go vet memperingatkan)
	fmt.Println(string(rune(n)))     // "A"  eksplisit: jelas ini karakter
	fmt.Println(strconv.Itoa(n))     // "65" cara benar untuk teks angka
	fmt.Println(fmt.Sprint(n) + "!") // "65!" alternatif
}
