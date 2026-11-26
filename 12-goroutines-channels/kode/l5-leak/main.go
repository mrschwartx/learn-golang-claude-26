// Kegagalan 3: goroutine leak, dan perbaikannya dengan jalur keluar (done).
package main

import (
	"fmt"
	"runtime"
	"time"
)

// BOCOR: konsumen berhenti membaca setelah 1 nilai; produsen tertahan selamanya.
func bocor() int {
	out := make(chan int)
	go func() {
		for i := 0; ; i++ {
			out <- i // terblokir selamanya setelah konsumen pergi
		}
	}()
	return <-out
}

// PERBAIKAN: produsen memantau done.
func aman() int {
	out := make(chan int)
	done := make(chan struct{})
	go func() {
		for i := 0; ; i++ {
			select {
			case out <- i:
			case <-done:
				return
			}
		}
	}()
	v := <-out
	close(done) // konsumen memberi tahu: berhenti
	return v
}

func main() {
	awal := runtime.NumGoroutine()
	for i := 0; i < 100; i++ {
		bocor()
	}
	time.Sleep(50 * time.Millisecond)
	fmt.Println("BOCOR : goroutine awal =", awal, "-> sesudah 100 panggilan =", runtime.NumGoroutine())

	for i := 0; i < 100; i++ {
		aman()
	}
	time.Sleep(50 * time.Millisecond)
	fmt.Println("AMAN  : 100 panggilan tambahan, goroutine sekarang =", runtime.NumGoroutine(), "(tidak bertambah; selisih 100 berasal dari bagian BOCOR)")
}
