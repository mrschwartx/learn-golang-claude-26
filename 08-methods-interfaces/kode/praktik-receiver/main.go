// receiver.go
package main

import "fmt"

type Penghitung struct{ n int }

func (p *Penghitung) Tambah()       { p.n++ }
func (p Penghitung) Nilai() int     { return p.n }
func (p Penghitung) String() string { return fmt.Sprintf("Penghitung(%d)", p.n) }

type Penambah interface{ Tambah() }

func main() {
	var c Penghitung
	c.Tambah() // otomatis (&c).Tambah()
	c.Tambah()
	fmt.Println(c.Nilai(), c) // String() dipakai oleh fmt

	var t Penambah = &c // harus pointer: Tambah ada di method set *Penghitung
	t.Tambah()
	fmt.Println(c)

	// var t2 Penambah = c // error: Penghitung does not implement Penambah (method Tambah has pointer receiver)
}
