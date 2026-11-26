// Latihan 1: 5 goroutine + WaitGroup; urutan output tidak terjamin.
package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for id := 1; id <= 5; id++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Print("goroutine-", id, " ")
		}()
	}
	wg.Wait()
	fmt.Println()
}
