// Latihan 7: *db.Postgres memenuhi layananpesanan.Penyimpan tanpa saling mengimpor interface.
package main

import (
	"fmt"

	"contoh.com/modul8/l7/db"
	"contoh.com/modul8/l7/layananpesanan"
)

// Dicek di sisi pemakai (main merangkai keduanya):
var _ layananpesanan.Penyimpan = (*pembungkus)(nil)

// Pesanan di db dan di layananpesanan adalah dua tipe berbeda; adapter kecil menjembatani.
type pembungkus struct{ pg *db.Postgres }

func (p *pembungkus) Simpan(ps layananpesanan.Pesanan) error {
	return p.pg.Simpan(db.Pesanan{ID: ps.ID, Total: ps.Total})
}

func main() {
	svc := layananpesanan.Baru(&pembungkus{&db.Postgres{DSN: "postgres://localhost/toko"}})
	fmt.Println("Buat(1, 7500):", svc.Buat(1, 7500))
	fmt.Println("Buat(2, 0)   :", svc.Buat(2, 0))
}
