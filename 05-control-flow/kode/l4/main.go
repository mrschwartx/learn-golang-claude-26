// Latihan 4: v pada range adalah salinan; dua cara memperbaikinya.
package main

import "fmt"

type Item struct {
	Nama string
	Stok int
}

func main() {
	items := []Item{{"a", 1}, {"b", 2}, {"c", 3}}

	for _, v := range items {
		v.Stok *= 10 // hanya mengubah salinan
	}
	fmt.Println("1. salah (salinan)       :", items)

	// bukti: alamat v berbeda dari alamat elemen
	for i, v := range items[:1] {
		fmt.Printf("   &v==&items[%d]? %t\n", i, &v == &items[i])
	}

	for i := range items { // cara 1: lewat indeks
		items[i].Stok *= 10
	}
	fmt.Println("2. benar (lewat indeks)   :", items)

	ptrs := []*Item{{"x", 1}, {"y", 2}} // cara 2: slice berisi pointer
	for _, p := range ptrs {
		p.Stok *= 10 // p salinan pointer, tapi menunjuk objek yang sama
	}
	fmt.Println("3. benar (slice pointer)  :", *ptrs[0], *ptrs[1])
}
