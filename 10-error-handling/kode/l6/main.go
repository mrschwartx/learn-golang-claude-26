// Latihan 6: %w vs %v pada paket yang menyembunyikan error driver database.
package main

import (
	"errors"
	"fmt"
)

// --- "driver database" pihak ketiga (detail internal) ---
type ErrDriver struct{ Kode string }

func (e *ErrDriver) Error() string { return "driver: kode " + e.Kode }

var ErrKoneksiDriver = errors.New("driver: koneksi ditolak")

func driverQuery() error { return fmt.Errorf("%w (%w)", &ErrDriver{"23505"}, ErrKoneksiDriver) }

// --- API paket repo kita ---
var ErrDuplikat = errors.New("repo: data duplikat") // sentinel milik KITA

// Versi A: %w mengekspos error driver sebagai bagian dari API
func SimpanBocor() error {
	if err := driverQuery(); err != nil {
		return fmt.Errorf("simpan user: %w", err)
	}
	return nil
}

// Versi B: %v memutus rantai; terjemahkan ke error domain milik kita
func SimpanTertutup() error {
	if err := driverQuery(); err != nil {
		var ed *ErrDriver
		if errors.As(err, &ed) && ed.Kode == "23505" {
			return fmt.Errorf("simpan user: %w", ErrDuplikat) // kontrak stabil
		}
		return fmt.Errorf("simpan user: %v", err) // pesan ada, rantai putus
	}
	return nil
}

func main() {
	a := SimpanBocor()
	fmt.Println("A pesan:", a)
	var ed *ErrDriver
	fmt.Println("A: pemanggil bisa melihat error driver?", errors.As(a, &ed), "| Is(ErrKoneksiDriver):", errors.Is(a, ErrKoneksiDriver))

	b := SimpanTertutup()
	fmt.Println("B pesan:", b)
	fmt.Println("B: error driver terlihat?", errors.As(b, &ed), "| Is(ErrDuplikat):", errors.Is(b, ErrDuplikat))

	// kasus pemakaian %v murni
	c := fmt.Errorf("simpan user: %v", driverQuery())
	fmt.Println("C (dibungkus dengan v): pesan:", c, "| Is(ErrKoneksiDriver):", errors.Is(c, ErrKoneksiDriver))
}
