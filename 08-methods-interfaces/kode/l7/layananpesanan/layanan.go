// Package layananpesanan: pemakai (consumer). Mendeklarasikan interface minimal miliknya sendiri.
package layananpesanan

import (
	"errors"
	"fmt"
)

type Pesanan struct {
	ID    int
	Total int
}

// Penyimpan: hanya SATU method yang benar-benar dibutuhkan paket ini.
type Penyimpan interface {
	Simpan(p Pesanan) error
}

type Layanan struct{ store Penyimpan } // menerima interface

func Baru(s Penyimpan) *Layanan { return &Layanan{store: s} } // mengembalikan struct

var ErrTotalNol = errors.New("total pesanan harus > 0")

func (l *Layanan) Buat(id, total int) error {
	if total <= 0 {
		return ErrTotalNol
	}
	if err := l.store.Simpan(Pesanan{id, total}); err != nil {
		return fmt.Errorf("simpan pesanan %d: %w", id, err)
	}
	return nil
}
