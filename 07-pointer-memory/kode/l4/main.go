// Latihan 4: mengubah field struct di dalam map, dua cara.
package main

import "fmt"

type Produk struct {
	Nama string
	Stok int
}

func main() {
	// cara 1: salin - ubah - tulis ulang
	m := map[string]Produk{"kopi": {"Kopi", 10}}
	p := m["kopi"]
	p.Stok--
	m["kopi"] = p
	fmt.Println("salin-tulis ulang :", m)

	// cara 2: map berisi pointer
	mp := map[string]*Produk{"kopi": {"Kopi", 10}}
	mp["kopi"].Stok--
	fmt.Println("map pointer       :", *mp["kopi"])

	// jebakan cara 2: key tidak ada -> nil pointer
	defer func() { fmt.Println("recover:", recover()) }()
	mp["teh"].Stok-- // panic: nil pointer dereference
}
