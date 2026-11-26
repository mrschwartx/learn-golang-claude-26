// Latihan 3: method set dan pointer receiver. Versi yang gagal ada di ../l3-error.
package main

import "fmt"

type Resetter interface{ Reset() }

type Counter struct{ n int }

func (c *Counter) Reset() { c.n = 0 }

func pakai(r Resetter) { r.Reset() }

func main() {
	c := Counter{n: 5}

	pakai(&c) // pointer: *Counter memenuhi Resetter
	fmt.Println("setelah pakai(&c):", c.n)

	c.n = 9
	c.Reset() // variabel addressable: otomatis (&c).Reset()
	fmt.Println("c.Reset() langsung :", c.n)

	// elemen map tidak addressable -> method pointer receiver tidak bisa dipanggil
	m := map[string]*Counter{"a": {n: 3}} // simpan pointer agar bisa
	m["a"].Reset()
	fmt.Println("map berisi pointer :", m["a"].n)

	// pembuktian method set lewat reflect
	// pakai(c) // <- error kompilasi: lihat l3-error
}
