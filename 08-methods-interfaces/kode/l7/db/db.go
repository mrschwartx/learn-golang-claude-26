// Package db: penyedia (provider). Tidak tahu apa-apa soal layananpesanan.
package db

import "fmt"

type Pesanan struct {
	ID    int
	Total int
}

type Postgres struct{ DSN string }

func (p *Postgres) Simpan(ps Pesanan) error {
	fmt.Printf("INSERT ke %s: pesanan %d total %d\n", p.DSN, ps.ID, ps.Total)
	return nil
}

// Banyak method lain yang TIDAK dibutuhkan layananpesanan:
func (p *Postgres) Hapus(id int) error                   { return nil }
func (p *Postgres) Cari(id int) (Pesanan, error)         { return Pesanan{}, nil }
func (p *Postgres) Laporan(bulan int) ([]Pesanan, error) { return nil, nil }
