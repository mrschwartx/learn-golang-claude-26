// Latihan 5: pipeline berhenti dini; periksa kebocoran goroutine.
package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func generator(ctx context.Context, n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 1; i <= n; i++ {
			select {
			case out <- i:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func kuadrat(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in {
			select {
			case out <- v * v:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// gabung versi BERBATAS-KONTEKS: goroutine penyalin juga memantau ctx
func gabung(ctx context.Context, chans ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	for _, c := range chans {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range c {
				select {
				case out <- v:
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	go func() { wg.Wait(); close(out) }()
	return out
}

// gabungTanpaCtx: versi dari materi (tanpa ctx pada penyalin)
func gabungTanpaCtx(chans ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	for _, c := range chans {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range c {
				out <- v // terblokir bila konsumen pergi
			}
		}()
	}
	go func() { wg.Wait(); close(out) }()
	return out
}

func jalankan(nama string, versiAman bool) {
	awal := runtime.NumGoroutine()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // dipanggil juga bila loop selesai tanpa break

	sumber := generator(ctx, 1000)
	w1, w2, w3 := kuadrat(ctx, sumber), kuadrat(ctx, sumber), kuadrat(ctx, sumber)
	var hasil <-chan int
	if versiAman {
		hasil = gabung(ctx, w1, w2, w3)
	} else {
		hasil = gabungTanpaCtx(w1, w2, w3)
	}

	diambil := 0
	for v := range hasil {
		fmt.Printf("  %s: ambil %d\n", nama, v)
		diambil++
		if diambil == 3 {
			cancel() // berhenti dini
			break
		}
	}
	time.Sleep(200 * time.Millisecond) // beri waktu goroutine berhenti
	fmt.Printf("%s: goroutine sebelum=%d  sesudah=%d  (selisih %d)\n\n", nama, awal, runtime.NumGoroutine(), runtime.NumGoroutine()-awal)
}

func main() {
	jalankan("gabung dengan ctx   ", true)
	jalankan("gabung dari materi  ", false)
}
