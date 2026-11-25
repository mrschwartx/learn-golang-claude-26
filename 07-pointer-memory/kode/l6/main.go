// Latihan 6: dampak GOGC. Jalankan dengan GOGC=25, 100, 400 (dan GODEBUG=gctrace=1).
package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

func main() {
	// data hidup tetap ~32 MiB
	hidup := make([][]byte, 32)
	for i := range hidup {
		hidup[i] = make([]byte, 1<<20)
	}

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	gcAwal := ms.NumGC

	// churn: ~2 GiB sampah jangka pendek, potongan 64 KiB
	var puncak uint64
	var sink []byte
	mulai := time.Now()
	for i := 0; i < 32768; i++ {
		sink = make([]byte, 64<<10)
		sink[0] = byte(i)
		if i%64 == 0 {
			runtime.ReadMemStats(&ms)
			if ms.HeapAlloc > puncak {
				puncak = ms.HeapAlloc
			}
		}
	}
	_ = sink
	runtime.ReadMemStats(&ms)
	fmt.Fprintf(os.Stderr, "")
	fmt.Printf("GOGC=%-4s siklus GC=%-4d puncak heap=%4d MiB  pause total=%v  waktu=%v\n",
		os.Getenv("GOGC"), ms.NumGC-gcAwal, puncak>>20, time.Duration(ms.PauseTotalNs).Round(time.Microsecond),
		time.Since(mulai).Round(time.Millisecond))
	runtime.KeepAlive(hidup)
}
