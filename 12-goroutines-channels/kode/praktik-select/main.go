// select.go
package main

import (
	"fmt"
	"time"
)

func main() {
	ping, pong := make(chan int), make(chan int)
	go func() {
		for n := range ping {
			pong <- n + 1
		}
		close(pong)
	}()
	for i := 0; i < 3; i++ {
		ping <- i * 10
		fmt.Println("balasan:", <-pong)
	}
	close(ping)

	// timeout
	lambat := make(chan string)
	go func() {
		time.Sleep(200 * time.Millisecond)
		lambat <- "selesai"
	}()
	select {
	case s := <-lambat:
		fmt.Println(s)
	case <-time.After(50 * time.Millisecond):
		fmt.Println("timeout")
	}

	// non-blocking
	kosong := make(chan int)
	select {
	case v := <-kosong:
		fmt.Println(v)
	default:
		fmt.Println("belum ada data")
	}
}
