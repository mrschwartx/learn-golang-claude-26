package layananpesanan

import (
	"errors"
	"testing"
)

// fake: implementasi kecil di dalam test, tanpa database sungguhan
type fakeStore struct {
	disimpan []Pesanan
	err      error
}

func (f *fakeStore) Simpan(p Pesanan) error {
	if f.err != nil {
		return f.err
	}
	f.disimpan = append(f.disimpan, p)
	return nil
}

func TestBuatBerhasil(t *testing.T) {
	f := &fakeStore{}
	if err := Baru(f).Buat(1, 5000); err != nil {
		t.Fatal(err)
	}
	if len(f.disimpan) != 1 || f.disimpan[0].Total != 5000 {
		t.Fatalf("hasil tak terduga: %+v", f.disimpan)
	}
}

func TestBuatValidasi(t *testing.T) {
	f := &fakeStore{}
	err := Baru(f).Buat(1, 0)
	if !errors.Is(err, ErrTotalNol) {
		t.Fatalf("seharusnya ErrTotalNol, dapat %v", err)
	}
	if len(f.disimpan) != 0 {
		t.Fatal("tidak boleh menyimpan bila validasi gagal")
	}
}

func TestBuatGagalSimpan(t *testing.T) {
	sentinel := errors.New("koneksi putus")
	f := &fakeStore{err: sentinel}
	err := Baru(f).Buat(2, 100)
	if !errors.Is(err, sentinel) {
		t.Fatalf("error asli harus terbungkus, dapat %v", err)
	}
}
