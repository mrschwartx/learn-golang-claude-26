// Latihan 3: ErrHTTP custom + errors.As + asType (pengganti errors.AsType Go 1.26).
package main

import (
	"errors"
	"fmt"
)

type ErrHTTP struct {
	Kode  int
	Pesan string
}

func (e *ErrHTTP) Error() string { return fmt.Sprintf("http %d: %s", e.Kode, e.Pesan) }

func panggilAPI(path string) error {
	switch path {
	case "/hilang":
		return &ErrHTTP{404, "tidak ditemukan"}
	case "/rusak":
		return &ErrHTTP{500, "kesalahan server"}
	}
	return nil
}

func klien(path string) error {
	if err := panggilAPI(path); err != nil {
		return fmt.Errorf("klien GET %s: %w", path, err)
	}
	return nil
}

// Setara errors.AsType (Go 1.26+): tanpa variabel perantara di pemanggil.
func asType[T error](err error) (T, bool) {
	var t T
	ok := errors.As(err, &t)
	return t, ok
}

func main() {
	for _, p := range []string{"/hilang", "/rusak", "/ok"} {
		err := klien(p)
		if err == nil {
			fmt.Println(p, "-> sukses")
			continue
		}

		// cara 1: errors.As
		var eh *ErrHTTP
		if errors.As(err, &eh) {
			fmt.Printf("%s -> As: kode=%d retry=%t\n", p, eh.Kode, eh.Kode >= 500)
		}
		// cara 2: asType / errors.AsType
		if v, ok := asType[*ErrHTTP](err); ok {
			fmt.Printf("%s -> asType: %s\n", p, v.Pesan)
		}
	}
}
