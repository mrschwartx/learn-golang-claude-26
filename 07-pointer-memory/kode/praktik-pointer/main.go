// pointer.go
package main

import "fmt"

type User struct {
	Nama string
	Umur int
}

// ptr: pengganti new(ekspresi) untuk Go < 1.26.
func ptr[T any](v T) *T { return &v }

func tukar(a, b *int) { *a, *b = *b, *a }

func ulangTahun(u *User) { u.Umur++ }

func main() {
	x := 10
	p := &x
	*p = 20
	fmt.Println(x, *p)

	a, b := 1, 2
	tukar(&a, &b)
	fmt.Println(a, b)

	var kosong *int
	fmt.Println(kosong == nil)
	// fmt.Println(*kosong) // panic: nil pointer dereference

	n := new(int)
	*n = 7
	fmt.Println(*n)

	// Go 1.26+: umur := new(30). Toolchain di sini Go 1.24, jadi pakai variabel perantara
	// (atau helper generik) yang berefek sama.
	umur := ptr(30)
	fmt.Println(*umur)

	u := User{"Andi", 30}
	ulangTahun(&u)
	fmt.Println(u)
}
