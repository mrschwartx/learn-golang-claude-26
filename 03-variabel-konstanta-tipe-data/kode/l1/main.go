// Latihan 1: lima cara deklarasi variabel.
package main

import "fmt"

// Level paket: hanya var yang boleh (:= tidak valid di sini).
var versi = "1.0"

var (
	host = "localhost"
	port = 8080
)

func main() {
	var umur int = 30 // 1. var dengan tipe
	var nama = "Andi" // 2. var tanpa tipe (inferensi)
	var total float64 // 3. var tanpa nilai (zero value)
	var (             // 4. var berkelompok
		x = 1
		y = 2.5
	)
	kota := "Jakarta" // 5. short declaration

	fmt.Printf("umur  %T %v\n", umur, umur)
	fmt.Printf("nama  %T %v\n", nama, nama)
	fmt.Printf("total %T %v\n", total, total)
	fmt.Printf("x     %T %v\n", x, x)
	fmt.Printf("y     %T %v\n", y, y)
	fmt.Printf("kota  %T %v\n", kota, kota)
	fmt.Printf("paket %T %v | %v:%v\n", versi, versi, host, port)
}
