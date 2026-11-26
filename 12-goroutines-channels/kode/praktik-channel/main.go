// channel.go
package main

import (
	"fmt"
	"sync"
)

func kirimAngka(n int, out chan<- int) {
	for i := 1; i <= n; i++ {
		out <- i
	}
	close(out) // pengirim yang menutup
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("halo dari goroutine")
	}()
	wg.Wait() // tanpa ini main bisa selesai lebih dulu

	// unbuffered: kirim dan terima saling menyinkronkan
	ch := make(chan int)
	go kirimAngka(3, ch)
	for v := range ch {
		fmt.Println("terima:", v)
	}

	// buffered: kirim tidak memblokir selama buffer belum penuh
	buf := make(chan string, 2)
	buf <- "satu"
	buf <- "dua"
	fmt.Println(len(buf), cap(buf))
	fmt.Println(<-buf, <-buf)

	// channel tertutup: ok == false
	tutup := make(chan int, 1)
	tutup <- 9
	close(tutup)
	v, ok := <-tutup
	fmt.Println(v, ok)
	v, ok = <-tutup
	fmt.Println(v, ok)
}
