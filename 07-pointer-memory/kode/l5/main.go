// Latihan 5: slice kecil menahan array besar (100 MiB).
package main

import (
	"fmt"
	"runtime"
	"slices"
)

const ukuran = 100 << 20 // 100 MiB

func heapMiB() uint64 {
	runtime.GC()
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return ms.HeapAlloc >> 20
}

//go:noinline
func buatBesar() []byte { return make([]byte, ukuran) }

//go:noinline
func ambilKecilBocor() []byte {
	besar := buatBesar()
	return besar[:5] // len 5, tapi menahan 100 MiB
}

//go:noinline
func ambilKecilAman() []byte {
	besar := buatBesar()
	return slices.Clone(besar[:5]) // salinan 5 byte
}

func main() {
	fmt.Println("heap awal                       :", heapMiB(), "MiB")

	kecil := ambilKecilBocor()
	fmt.Printf("setelah ambilKecilBocor (len=%d cap=%d MiB): %d MiB <- 100 MiB masih tertahan\n",
		len(kecil), cap(kecil)>>20, heapMiB())
	runtime.KeepAlive(kecil)

	kecil = nil
	fmt.Println("setelah kecil=nil + GC          :", heapMiB(), "MiB")

	aman := ambilKecilAman()
	fmt.Printf("setelah ambilKecilAman (len=%d cap=%d): %d MiB <- array besar bebas\n",
		len(aman), cap(aman), heapMiB())
	runtime.KeepAlive(aman)
}
