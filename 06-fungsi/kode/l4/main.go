// Latihan 4: urutan defer. Tebak dulu outputnya sebelum menjalankan!
package main

import "fmt"

func demo() {
	pesan := "awal"
	defer fmt.Println("defer A: pesan saat didaftarkan =", pesan) // argumen dievaluasi SEKARANG

	defer func() {
		fmt.Println("defer B (closure): pesan saat dijalankan =", pesan) // dibaca SAAT berjalan
	}()

	for i := range 3 {
		defer fmt.Println("defer loop i =", i)
	}

	pesan = "berubah"
	fmt.Println("badan fungsi selesai, pesan =", pesan)
}

func main() {
	demo()
}
