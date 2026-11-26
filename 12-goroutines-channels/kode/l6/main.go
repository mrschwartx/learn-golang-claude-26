// Latihan 6: pipeline 3 tahap dengan channel terarah.
package main

import "fmt"

// generator: hanya menghasilkan -> mengembalikan <-chan int
func generator(n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 1; i <= n; i++ {
			out <- i
		}
	}()
	return out
}

// kuadrat: membaca dari <-chan, menulis ke <-chan baru
func kuadrat(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in {
			out <- v * v
		}
	}()
	return out
}

// cetak: hanya membaca
func cetak(in <-chan int) {
	for v := range in {
		fmt.Println("hasil:", v)
	}
}

func main() {
	cetak(kuadrat(generator(5)))
}
