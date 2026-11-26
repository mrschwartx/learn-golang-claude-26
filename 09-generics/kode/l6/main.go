// Latihan 6: iterator bilangan prima tak hingga, ambil 10 dengan break.
package main

import (
	"fmt"
	"iter"
)

func Prima() iter.Seq[int] {
	return func(yield func(int) bool) {
		hitung := 0 // statistik untuk demo
		defer func() { fmt.Println("  [iterator berhenti & membersihkan diri; yield dipanggil", hitung, "kali]") }()
		for n := 2; ; n++ {
			if !laPrima(n) {
				continue
			}
			hitung++
			if !yield(n) { // false: konsumen melakukan break -> WAJIB berhenti
				return
			}
		}
	}
}

func laPrima(n int) bool {
	for d := 2; d*d <= n; d++ {
		if n%d == 0 {
			return false
		}
	}
	return n >= 2
}

func main() {
	fmt.Println("10 prima pertama:")
	i := 0
	for p := range Prima() {
		fmt.Print(p, " ")
		i++
		if i == 10 {
			break
		}
	}
	fmt.Println()

	// bila iterator mengabaikan hasil yield, Go memicu panic saat runtime
	buruk := func(yield func(int) bool) {
		for i := 0; ; i++ {
			yield(i) // SALAH: hasil diabaikan
			if i > 3 {
				return
			}
		}
	}
	func() {
		defer func() { fmt.Println("iterator buruk ->", recover()) }()
		for v := range buruk {
			if v == 1 {
				break
			}
		}
	}()
}
