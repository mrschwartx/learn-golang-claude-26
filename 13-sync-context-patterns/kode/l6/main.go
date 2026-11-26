// Latihan 6: semaphore (channel berbuffer) vs worker pool; maksimum 3 dari 20 tugas.
package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const (
	nTugas = 20
	batas  = 3
)

type ukur struct {
	aktif, puncak atomic.Int32
}

func (u *ukur) tugas(id int) {
	n := u.aktif.Add(1)
	for {
		p := u.puncak.Load()
		if n <= p || u.puncak.CompareAndSwap(p, n) {
			break
		}
	}
	time.Sleep(30 * time.Millisecond)
	u.aktif.Add(-1)
}

func semaphore() (int32, time.Duration, int) {
	var u ukur
	var wg sync.WaitGroup
	sem := make(chan struct{}, batas)
	mulai := time.Now()
	maksGoroutine := 0
	for i := 0; i < nTugas; i++ {
		wg.Add(1)
		go func() { // 20 goroutine dibuat sekaligus
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			u.tugas(i)
		}()
	}
	maksGoroutine = runtime.NumGoroutine()
	wg.Wait()
	return u.puncak.Load(), time.Since(mulai), maksGoroutine
}

func pool() (int32, time.Duration, int) {
	var u ukur
	var wg sync.WaitGroup
	jobs := make(chan int, nTugas)
	mulai := time.Now()
	for w := 0; w < batas; w++ { // hanya 3 goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range jobs {
				u.tugas(id)
			}
		}()
	}
	maksGoroutine := runtime.NumGoroutine()
	for i := 0; i < nTugas; i++ {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	return u.puncak.Load(), time.Since(mulai), maksGoroutine
}

func main() {
	p, d, g := semaphore()
	fmt.Printf("semaphore berbuffer : puncak serentak=%d  waktu=%v  goroutine (sesaat setelah dibuat)=%d\n", p, d.Round(10*time.Millisecond), g)
	p, d, g = pool()
	fmt.Printf("worker pool         : puncak serentak=%d  waktu=%v  goroutine (sesaat setelah dibuat)=%d\n", p, d.Round(10*time.Millisecond), g)
}
