// map.go
package main

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

func main() {
	teks := "kucing anjing kucing ikan anjing kucing"
	hitung := make(map[string]int)
	for _, kata := range strings.Fields(teks) {
		hitung[kata]++
	}

	// membaca key yang tidak ada -> zero value, gunakan comma-ok
	v, ok := hitung["burung"]
	fmt.Println(v, ok)

	// urutan iterasi acak, jadi urutkan key dulu
	for _, k := range slices.Sorted(maps.Keys(hitung)) {
		fmt.Println(k, hitung[k])
	}

	var kosong map[string]int
	fmt.Println(kosong["x"], len(kosong)) // membaca map nil aman
	// kosong["x"] = 1                    // panic: assignment to entry in nil map
}
