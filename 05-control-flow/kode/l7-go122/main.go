package main

import (
	"fmt"
	"slices"
	"sync"
)

func main() {
	users := []string{"andi", "budi", "citra"}

	// closure (deterministik)
	var fs []func() string
	for _, u := range users {
		fs = append(fs, func() string { return u })
	}
	for _, f := range fs {
		fmt.Print(f(), " ")
	}
	fmt.Println("<- closure")

	// goroutine, persis seperti kasus tim
	var wg sync.WaitGroup
	var mu sync.Mutex
	var hasil []string
	for _, u := range users {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			hasil = append(hasil, u)
			mu.Unlock()
		}()
	}
	wg.Wait()
	slices.Sort(hasil)
	fmt.Println(hasil, "<- goroutine")
}
