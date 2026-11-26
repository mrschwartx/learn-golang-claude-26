// Latihan 2: generator angkaGenap; PENGIRIM (generator) yang menutup channel.
package main

import "fmt"

func angkaGenap(n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out) // pemilik/pengirim menutup
		for i := 1; i <= n; i++ {
			out <- i * 2
		}
	}()
	return out
}

func main() {
	for v := range angkaGenap(5) { // range berhenti otomatis saat channel ditutup
		fmt.Print(v, " ")
	}
	fmt.Println()
}
