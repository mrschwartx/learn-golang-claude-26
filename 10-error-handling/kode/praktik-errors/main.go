// errors.go
package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

// asType: pengganti errors.AsType (Go 1.26+) untuk toolchain yang lebih lama.
func asType[T error](err error) (T, bool) {
	var t T
	ok := errors.As(err, &t)
	return t, ok
}

var ErrStokHabis = errors.New("stok habis")

type ErrValidasi struct {
	Field string
	Pesan string
}

func (e *ErrValidasi) Error() string {
	return fmt.Sprintf("validasi %s: %s", e.Field, e.Pesan)
}

func cekStok(n int) error {
	if n <= 0 {
		return ErrStokHabis
	}
	return nil
}

func buatPesanan(jumlah int, email string) error {
	if email == "" {
		return &ErrValidasi{Field: "email", Pesan: "wajib diisi"}
	}
	if err := cekStok(jumlah); err != nil {
		return fmt.Errorf("buat pesanan: %w", err)
	}
	return nil
}

func main() {
	err := buatPesanan(0, "a@b.com")
	fmt.Println(err)
	fmt.Println("stok habis?", errors.Is(err, ErrStokHabis), err == ErrStokHabis)

	err = buatPesanan(1, "")
	var ev *ErrValidasi
	if errors.As(err, &ev) {
		fmt.Println("field bermasalah:", ev.Field)
	}
	if v, ok := asType[*ErrValidasi](err); ok {
		fmt.Println("AsType:", v.Pesan)
	}

	_, err = os.Open("tidak-ada.txt")
	fmt.Println(err)
	fmt.Println(errors.Is(err, fs.ErrNotExist))

	gabungan := errors.Join(errors.New("langkah 1 gagal"), errors.New("langkah 2 gagal"))
	fmt.Println(gabungan)
}
