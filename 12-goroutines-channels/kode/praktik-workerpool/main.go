// workerpool.go
package main

import (
	"fmt"
	"slices"
	"sync"
	"time"
)

type Hasil struct {
	Job, Worker, Kuadrat int
}

func worker(id int, jobs <-chan int, results chan<- Hasil, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		time.Sleep(10 * time.Millisecond) // simulasi kerja
		results <- Hasil{Job: j, Worker: id, Kuadrat: j * j}
	}
}

func main() {
	jobs := make(chan int, 9)
	results := make(chan Hasil, 9)

	var wg sync.WaitGroup
	for id := 1; id <= 3; id++ {
		wg.Add(1)
		go worker(id, jobs, results, &wg)
	}

	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	close(jobs)

	wg.Wait()
	close(results)

	var semua []Hasil
	for h := range results {
		semua = append(semua, h)
	}
	slices.SortFunc(semua, func(a, b Hasil) int { return a.Job - b.Job })
	for _, h := range semua {
		fmt.Printf("job %d -> %d\n", h.Job, h.Kuadrat)
	}
}
