// pipeline.go
package main

import (
	"context"
	"fmt"
	"slices"
	"sync"
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

func gabung(chans ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	for _, c := range chans {
		jalankan(&wg, func() {
			for v := range c {
				out <- v
			}
		})
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sumber := generator(ctx, 10)
	// fan-out: 3 worker membaca dari channel yang sama
	w1, w2, w3 := kuadrat(ctx, sumber), kuadrat(ctx, sumber), kuadrat(ctx, sumber)

	var hasil []int
	for v := range gabung(w1, w2, w3) { // fan-in
		hasil = append(hasil, v)
	}
	slices.Sort(hasil)
	fmt.Println(hasil)
}
