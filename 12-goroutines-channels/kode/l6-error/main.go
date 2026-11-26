//go:build ignore

package main

func hanyaBaca(in <-chan int) {
	in <- 1 // salah arah
	close(in)
}

func hanyaKirim(out chan<- int) {
	_ = <-out // salah arah
}

func main() {}
