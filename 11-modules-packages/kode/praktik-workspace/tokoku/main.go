package main

import (
	"fmt"

	"contoh.com/salam"
	"contoh.com/tokoku/diskon"
	"contoh.com/tokoku/internal/harga"
	"contoh.com/tokoku/laporan"
)

func main() {
	fmt.Println(salam.Halo("Budi"))
	total := harga.Total([]int{10000, 25000})
	akhir, err := diskon.Terapkan(total, 10)
	fmt.Println(akhir, err)

	_, err = diskon.Terapkan(total, 80)
	fmt.Println(err)

	// package ketiga: laporan
	teks, err := laporan.Cetak([]int{10000, 25000, 5000}, 20)
	fmt.Print(teks)
	fmt.Println("err:", err)

	_, err = laporan.Cetak([]int{1000}, 99)
	fmt.Println("laporan gagal:", err)
}
