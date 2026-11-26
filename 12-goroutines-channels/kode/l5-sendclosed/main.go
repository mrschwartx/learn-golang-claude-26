// Kegagalan 2: send on closed channel (dua penutup / kirim setelah tutup).
package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan int, 10)
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch <- i
		}()
	}
	wg.Wait()
	close(ch)
	fmt.Println("ditutup; sekarang mengirim lagi...")
	ch <- 99 // panic: send on closed channel
}
