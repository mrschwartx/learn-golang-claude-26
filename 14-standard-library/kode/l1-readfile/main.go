// l1-readfile: pembanding: os.ReadFile memuat seluruh berkas ke memori.
package main

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
)

func main() {
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	rss := "n/a"
	if b, err := os.ReadFile("/proc/self/status"); err == nil {
		for _, l := range strings.Split(string(b), "\n") {
			if strings.HasPrefix(l, "VmHWM:") {
				rss = strings.TrimSpace(strings.TrimPrefix(l, "VmHWM:"))
			}
		}
	}
	fmt.Printf("baris=%d ukuran=%.1f MiB heap diminta dari OS: %.1f MiB, RSS puncak: %s\n",
		bytes.Count(data, []byte("\n")), float64(len(data))/1048576, float64(ms.HeapSys)/1048576, rss)
}
