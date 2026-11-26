// Latihan 4: kumpulkan semua kegagalan validasi dengan errors.Join.
package main

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNamaKosong      = errors.New("nama wajib diisi")
	ErrEmailTidakValid = errors.New("email tidak valid")
	ErrUmurDiBawah     = errors.New("umur minimal 17")
)

type Form struct {
	Nama, Email string
	Umur        int
}

func Validasi(f Form) error {
	var errs []error
	if strings.TrimSpace(f.Nama) == "" {
		errs = append(errs, ErrNamaKosong)
	}
	if !strings.Contains(f.Email, "@") {
		errs = append(errs, fmt.Errorf("%w: %q", ErrEmailTidakValid, f.Email))
	}
	if f.Umur < 17 {
		errs = append(errs, fmt.Errorf("%w (diberikan %d)", ErrUmurDiBawah, f.Umur))
	}
	return errors.Join(errs...) // nil bila errs kosong
}

func main() {
	fmt.Println("form valid  ->", Validasi(Form{"Andi", "a@b.com", 20}))

	err := Validasi(Form{"", "bukan-email", 10})
	fmt.Println("form salah  ->")
	fmt.Println(err)

	fmt.Println("\nIs(NamaKosong)   :", errors.Is(err, ErrNamaKosong))
	fmt.Println("Is(EmailTidakValid):", errors.Is(err, ErrEmailTidakValid))
	fmt.Println("Is(UmurDiBawah)  :", errors.Is(err, ErrUmurDiBawah))
	fmt.Println("Is(io lain)      :", errors.Is(err, errors.New("lain")))

	// Join mengimplementasikan Unwrap() []error
	if j, ok := err.(interface{ Unwrap() []error }); ok {
		fmt.Println("jumlah error di dalamnya:", len(j.Unwrap()))
	}
}
