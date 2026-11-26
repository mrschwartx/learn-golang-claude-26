// perbaikan.go
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Penghitung struct {
	mu sync.Mutex
	n  int
}

func (p *Penghitung) Tambah() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.n++
}

// jalankan: pengganti wg.Go (Go 1.25+) untuk toolchain yang lebih lama.
func jalankan(wg *sync.WaitGroup, f func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		f()
	}()
}

func main() {
	var wg sync.WaitGroup

	// 1. Mutex
	var pc Penghitung
	for range 1000 {
		jalankan(&wg, func() { pc.Tambah() }) // pada Go 1.25+: wg.Go(func() { pc.Tambah() })
	}

	// 2. Atomic
	var at atomic.Int64
	for range 1000 {
		jalankan(&wg, func() { at.Add(1) })
	}

	// 3. Channel: satu goroutine pemilik state
	ch := make(chan struct{})
	selesai := make(chan int)
	go func() {
		n := 0
		for range ch {
			n++
		}
		selesai <- n
	}()
	for range 1000 {
		jalankan(&wg, func() { ch <- struct{}{} })
	}
	wg.Wait()
	close(ch)

	fmt.Println(pc.n, at.Load(), <-selesai)

	var once sync.Once
	for range 3 {
		once.Do(func() { fmt.Println("inisialisasi sekali") })
	}
}
