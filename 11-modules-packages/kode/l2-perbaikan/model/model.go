// Package model memuat tipe bersama; tidak mengimpor a maupun b.
package model

type Pengguna struct{ Nama string }

type Pesanan struct {
	ID      int
	Pemilik Pengguna
}
