// iota.go
package main

import "fmt"

type Hari int

const (
	Minggu Hari = iota
	Senin
	Selasa
)

func namaHari(h Hari) string {
	switch h {
	case Minggu:
		return "Minggu"
	case Senin:
		return "Senin"
	default:
		return "Lainnya"
	}
}

const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
	GB
)

const (
	Baca  = 1 << iota // 1
	Tulis             // 2
	Eksek             // 4
)

func main() {
	fmt.Println(Minggu, Senin, Selasa, namaHari(Senin))
	fmt.Println(KB, MB, GB)

	izin := Baca | Tulis
	fmt.Println("izin:", izin, "bisa tulis:", izin&Tulis != 0, "bisa eksek:", izin&Eksek != 0)

	var kecil uint8 = 255
	kecil++
	fmt.Println("overflow uint8:", kecil)

	x, y := 0.1, 0.2
	fmt.Println(x+y, x+y == 0.3) // variabel: dihitung float64 saat runtime
	fmt.Println(0.1+0.2 == 0.3)  // konstanta: presisi tak terbatas saat kompilasi
}
