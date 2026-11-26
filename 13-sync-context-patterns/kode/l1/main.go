// Latihan 1: data race pada counter. Jalankan: go run -race ./l1
package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	hitung := 0
	for range 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			hitung++ // BACA-TAMBAH-TULIS tanpa sinkronisasi
		}()
	}
	wg.Wait()
	fmt.Println("hasil (seharusnya 1000):", hitung)
}
