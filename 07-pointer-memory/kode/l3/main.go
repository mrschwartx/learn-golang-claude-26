// Latihan 3: mengembalikan pointer ke variabel lokal itu aman.
// Lihat keputusan escape analysis dengan:  go build -gcflags=-m ./l3
package main

import "fmt"

type User struct {
	Nama string
	Umur int
}

func NewUser(nama string) *User {
	u := User{Nama: nama, Umur: 20} // variabel lokal...
	return &u                       // ...alamatnya dikembalikan: u "lolos" ke heap
}

func hitungTanpaLolos() int {
	v := User{"lokal", 1} // tidak lolos -> tetap di stack
	return v.Umur
}

func main() {
	a := NewUser("Andi")
	b := NewUser("Budi")
	fmt.Println(*a, *b, a != b, hitungTanpaLolos())
}
