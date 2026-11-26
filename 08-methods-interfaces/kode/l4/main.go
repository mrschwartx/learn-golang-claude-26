// Latihan 4: jebakan nil pada interface error.
package main

import (
	"errors"
	"fmt"
)

type ErrValidasi struct{ Field string }

func (e *ErrValidasi) Error() string { return "validasi gagal: " + e.Field }

// SALAH: mengembalikan pointer nil bertipe konkret sebagai error
func validasiSalah(ok bool) error {
	var e *ErrValidasi
	if !ok {
		e = &ErrValidasi{"nama"}
	}
	return e // interface berisi (*ErrValidasi, nil) -> TIDAK nil
}

// BENAR 1: kembalikan literal nil
func validasiBenar(ok bool) error {
	if !ok {
		return &ErrValidasi{"nama"}
	}
	return nil
}

func main() {
	err := validasiSalah(true) // sukses, tapi...
	fmt.Println("validasiSalah(true) != nil :", err != nil)
	fmt.Printf("  tipe dinamis=%T, nilai=%v\n", err, err == (*ErrValidasi)(nil))
	// err.Error() akan panic karena method memakai e.Field pada pointer nil

	func() {
		defer func() { fmt.Println("  memanggil err.Error() ->", recover()) }()
		_ = err.Error()
	}()

	err = validasiBenar(true)
	fmt.Println("validasiBenar(true) != nil :", err != nil)

	err = validasiBenar(false)
	var ev *ErrValidasi
	fmt.Println("gagal ->", err, "| errors.As:", errors.As(err, &ev), ev.Field)
}
