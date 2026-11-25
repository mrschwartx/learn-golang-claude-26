// gc.go
package main

import (
	"fmt"
	"runtime"
)

func mb(b uint64) uint64 { return b / 1024 / 1024 }

func main() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	awal := ms.NumGC

	var simpan [][]byte
	for i := 0; i < 50; i++ {
		simpan = append(simpan, make([]byte, 1<<20)) // 1 MiB
	}
	runtime.ReadMemStats(&ms)
	fmt.Println("heap saat data masih dipegang (MiB) >= 50:", mb(ms.HeapAlloc) >= 50)

	simpan = nil
	runtime.GC()
	runtime.ReadMemStats(&ms)
	fmt.Println("heap setelah GC (MiB) < 10:", mb(ms.HeapAlloc) < 10)
	fmt.Println("GC berjalan:", ms.NumGC > awal)
}
