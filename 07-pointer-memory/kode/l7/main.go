// Latihan 7: GOMEMLIMIT sebagai batas lunak. Jalankan dengan dan tanpa GOMEMLIMIT.
package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
)

func main() {
	// "layanan" dengan 40 MiB data hidup, GOGC dilonggarkan agar GC jarang (mensimulasikan GOGC tinggi)
	debug.SetGCPercent(1000)

	hidup := make([][]byte, 40)
	for i := range hidup {
		hidup[i] = make([]byte, 1<<20)
	}

	var ms runtime.MemStats
	var puncak uint64
	var sink []byte
	for i := 0; i < 8192; i++ {
		sink = make([]byte, 64<<10)
		sink[0] = byte(i)
		if i%64 == 0 {
			runtime.ReadMemStats(&ms)
			if ms.HeapSys > puncak {
				puncak = ms.HeapSys
			}
		}
	}
	_ = sink
	runtime.ReadMemStats(&ms)
	fmt.Printf("GOMEMLIMIT=%-6q siklus GC=%-3d puncak HeapSys=%d MiB\n",
		os.Getenv("GOMEMLIMIT"), ms.NumGC, puncak>>20)
	runtime.KeepAlive(hidup)
}
