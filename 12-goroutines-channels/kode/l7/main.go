// Latihan 7: melindungi "database" yang hanya sanggup 10 query bersamaan.
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type DB struct {
	kapasitas int32
	aktif     atomic.Int32
	puncak    atomic.Int32
	ditolak   atomic.Int32
	berhasil  atomic.Int32
}

// Query meniru database: bila lebih dari kapasitas -> "tumbang" (menolak).
func (d *DB) Query() bool {
	n := d.aktif.Add(1)
	defer d.aktif.Add(-1)
	for {
		p := d.puncak.Load()
		if n <= p || d.puncak.CompareAndSwap(p, n) {
			break
		}
	}
	if n > d.kapasitas {
		d.ditolak.Add(1)
		return false
	}
	time.Sleep(20 * time.Millisecond)
	d.berhasil.Add(1)
	return true
}

func (d *DB) Laporan(nama string, t time.Duration) {
	fmt.Printf("%-28s berhasil=%-4d ditolak=%-4d puncak koneksi=%-4d waktu=%v\n",
		nama, d.berhasil.Load(), d.ditolak.Load(), d.puncak.Load(), t.Round(10*time.Millisecond))
}

const request = 500

func main() {
	// 1) satu goroutine per request, tanpa batas
	d1 := &DB{kapasitas: 10}
	mulai := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < request; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); d1.Query() }()
	}
	wg.Wait()
	d1.Laporan("tanpa batas:", time.Since(mulai))

	// 2) worker pool 10 worker
	d2 := &DB{kapasitas: 10}
	jobs := make(chan struct{}, request)
	mulai = time.Now()
	for w := 0; w < 10; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				d2.Query()
			}
		}()
	}
	for i := 0; i < request; i++ {
		jobs <- struct{}{}
	}
	close(jobs)
	wg.Wait()
	d2.Laporan("worker pool (10):", time.Since(mulai))

	// 3) channel berbuffer sebagai semaphore
	d3 := &DB{kapasitas: 10}
	sem := make(chan struct{}, 10)
	mulai = time.Now()
	for i := 0; i < request; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}        // ambil slot (memblokir bila penuh)
			defer func() { <-sem }() // lepas slot
			d3.Query()
		}()
	}
	wg.Wait()
	d3.Laporan("semaphore berbuffer (10):", time.Since(mulai))
}
