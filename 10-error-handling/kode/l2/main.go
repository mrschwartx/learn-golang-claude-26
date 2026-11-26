// Latihan 2: sentinel ErrTidakDitemukan melewati dua lapisan; errors.Is vs ==.
package main

import (
	"errors"
	"fmt"
)

var ErrTidakDitemukan = errors.New("tidak ditemukan")

// lapisan 1: repository
func cariUserDB(id int) error {
	return fmt.Errorf("query user id=%d: %w", id, ErrTidakDitemukan)
}

// lapisan 2: service
func AmbilProfil(id int) error {
	if err := cariUserDB(id); err != nil {
		return fmt.Errorf("ambil profil: %w", err)
	}
	return nil
}

func main() {
	err := AmbilProfil(42)
	fmt.Println("pesan lengkap  :", err)
	fmt.Println("errors.Is      :", errors.Is(err, ErrTidakDitemukan))
	fmt.Println("err == sentinel:", err == ErrTidakDitemukan)

	// rantai bungkusan
	for e := err; e != nil; e = errors.Unwrap(e) {
		fmt.Printf("  lapis: %-45q (%T)\n", e.Error(), e)
	}

	// jika dibungkus dengan %v, rantai putus
	putus := fmt.Errorf("ambil profil: %v", ErrTidakDitemukan)
	fmt.Println("dibungkus dengan verb v, errors.Is:", errors.Is(putus, ErrTidakDitemukan))
}
