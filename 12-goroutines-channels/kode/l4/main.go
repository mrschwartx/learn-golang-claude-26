// Latihan 4: 20 job, 1 vs 5 worker. Dua jenis beban: menunggu (I/O) dan CPU.
package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func pool(nWorker, nJob int, kerja func()) time.Duration {
	jobs := make(chan int, nJob)
	var wg sync.WaitGroup
	for w := 0; w < nWorker; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				kerja()
			}
		}()
	}
	mulai := time.Now()
	for j := 0; j < nJob; j++ {
		jobs <- j
	}
	close(jobs)
	wg.Wait()
	return time.Since(mulai)
}

var sink uint64

func cpu() { // ~ puluhan ms kerja CPU murni
	var x uint64 = 1
	for i := 0; i < 40_000_000; i++ {
		x = x*6364136223846793005 + 1442695040888963407
	}
	sink += x
}

func main() {
	fmt.Println("GOMAXPROCS =", runtime.GOMAXPROCS(0), "| NumCPU =", runtime.NumCPU())

	io := func() { time.Sleep(50 * time.Millisecond) }
	t1 := pool(1, 20, io)
	t5 := pool(5, 20, io)
	fmt.Printf("BEBAN MENUNGGU (sleep 50ms x 20 job): 1 worker=%v | 5 worker=%v | percepatan=%.1fx\n",
		t1.Round(time.Millisecond), t5.Round(time.Millisecond), float64(t1)/float64(t5))

	c1 := pool(1, 20, cpu)
	c5 := pool(5, 20, cpu)
	fmt.Printf("BEBAN CPU     (20 job):               1 worker=%v | 5 worker=%v | percepatan=%.1fx\n",
		c1.Round(time.Millisecond), c5.Round(time.Millisecond), float64(c1)/float64(c5))
}
