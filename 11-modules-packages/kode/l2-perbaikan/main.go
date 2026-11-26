package main

import (
	"fmt"

	"contoh.com/perbaikan/a"
	"contoh.com/perbaikan/b"
	"contoh.com/perbaikan/model"
)

func main() {
	for _, p := range a.PesananDari(model.Pengguna{Nama: "Andi"}) {
		fmt.Println("pesanan", p.ID, "milik", b.Ringkas(p))
	}
}
