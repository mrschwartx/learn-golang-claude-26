// Pembanding Latihan 2: map biasa tanpa kunci diakses bersamaan.
package main

import (
	"fmt"
	"sync"
)

func main() {
	m := map[string]string{}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				m[fmt.Sprint(i)] = "x"
			}
		}()
	}
	wg.Wait()
	fmt.Println("selesai, len =", len(m))
}
