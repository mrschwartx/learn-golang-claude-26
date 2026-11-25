//go:build ignore

// Latihan 4: sengaja TIDAK dikompilasi (build tag ignore). Salin ke direktori lain untuk melihat errornya.
package main

type Produk struct{ Stok int }

func main() {
	m := map[string]Produk{"kopi": {10}}
	m["kopi"].Stok -= 1
}
