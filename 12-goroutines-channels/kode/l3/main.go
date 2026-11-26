// Latihan 3: timeout dengan select + time.After, dan nasib goroutine lambat (leak).
package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

func lambat(d time.Duration, out chan<- string) {
	time.Sleep(d)
	out <- "selesai" // pada channel unbuffered tanpa penerima: terblokir selamanya
}

// VERSI BOCOR: channel unbuffered, penerima menyerah saat timeout.
func ambilDenganTimeout(d time.Duration) (string, error) {
	ch := make(chan string)
	go lambat(d, ch)
	select {
	case s := <-ch:
		return s, nil
	case <-time.After(100 * time.Millisecond):
		return "", fmt.Errorf("timeout setelah 100ms")
	}
}

// PERBAIKAN 1: channel berbuffer 1 -> goroutine selalu bisa mengirim lalu selesai.
func ambilBufer(d time.Duration) (string, error) {
	ch := make(chan string, 1)
	go lambat(d, ch)
	select {
	case s := <-ch:
		return s, nil
	case <-time.After(100 * time.Millisecond):
		return "", fmt.Errorf("timeout setelah 100ms")
	}
}

// PERBAIKAN 2: context ikut membatalkan pekerjaan.
func ambilContext(d time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	ch := make(chan string, 1)
	go func() {
		select {
		case <-time.After(d):
			ch <- "selesai"
		case <-ctx.Done(): // berhenti lebih awal, tidak menunggu sia-sia
		}
	}()
	select {
	case s := <-ch:
		return s, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func main() {
	fmt.Println("goroutine awal:", runtime.NumGoroutine())

	s, err := ambilDenganTimeout(30 * time.Millisecond)
	fmt.Println("cepat (30ms):", s, err)

	for i := 0; i < 50; i++ {
		ambilDenganTimeout(2 * time.Second) // semua timeout
	}
	time.Sleep(2500 * time.Millisecond) // beri waktu seluruh goroutine lambat selesai tidur
	fmt.Println("BOCOR   : setelah 50 timeout & menunggu 2,5 dtk, goroutine =", runtime.NumGoroutine())
}
