//go:build ignore

// Latihan 3: sengaja TIDAK dikompilasi (build tag ignore).
package main

type Resetter interface{ Reset() }

type Counter struct{ n int }

func (c *Counter) Reset() { c.n = 0 }

func main() {
	var r Resetter = Counter{} // nilai biasa: tidak memenuhi
	_ = r
}
