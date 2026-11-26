// Kapan atomic tidak cukup: invarian yang melibatkan DUA nilai (transfer saldo).
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const total = 1000

type AtomicSaja struct{ a, b atomic.Int64 }

// setiap operasi atomik, tetapi PASANGANNYA tidak atomik
func (s *AtomicSaja) Transfer(n int64) {
	s.a.Add(-n)
	s.b.Add(n)
}
func (s *AtomicSaja) Jumlah() int64 { return s.a.Load() + s.b.Load() }

type DenganMutex struct {
	mu   sync.Mutex
	a, b int64
}

func (s *DenganMutex) Transfer(n int64) { s.mu.Lock(); s.a -= n; s.b += n; s.mu.Unlock() }
func (s *DenganMutex) Jumlah() int64    { s.mu.Lock(); defer s.mu.Unlock(); return s.a + s.b }

func main() {
	var as AtomicSaja
	as.a.Store(total)
	var dm DenganMutex
	dm.a = total

	var wg sync.WaitGroup
	var salahA, salahM atomic.Int64
	stop := make(chan struct{})
	wg.Add(2)
	go func() { // pengamat
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
				if as.Jumlah() != total {
					salahA.Add(1)
				}
				if dm.Jumlah() != total {
					salahM.Add(1)
				}
			}
		}
	}()
	go func() { // penulis
		defer wg.Done()
		for i := 0; i < 2_000_000; i++ {
			as.Transfer(1)
			dm.Transfer(1)
		}
		close(stop)
	}()
	wg.Wait()
	fmt.Println("invarian a+b == 1000 dilanggar (hanya atomic):", salahA.Load(), "kali")
	fmt.Println("invarian a+b == 1000 dilanggar (mutex)       :", salahM.Load(), "kali")
}
