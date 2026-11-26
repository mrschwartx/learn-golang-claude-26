// Latihan 1: Keys/Values generik vs maps.Keys + slices.Collect.
package main

import (
	"fmt"
	"maps"
	"slices"
)

func Keys[K comparable, V any](m map[K]V) []K {
	out := make([]K, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func Values[K comparable, V any](m map[K]V) []V {
	out := make([]V, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	return out
}

func main() {
	m := map[string]int{"d": 4, "a": 1, "c": 3, "b": 2, "e": 5}

	fmt.Println("Keys (buatan sendiri)     :", Keys(m))
	fmt.Println("Values (buatan sendiri)   :", Values(m))
	fmt.Println("maps.Keys + Collect       :", slices.Collect(maps.Keys(m)))
	fmt.Println("maps.Values + Collect     :", slices.Collect(maps.Values(m)))

	// deterministik
	fmt.Println("terurut: slices.Sorted    :", slices.Sorted(maps.Keys(m)))
	k := Keys(m)
	slices.Sort(k)
	fmt.Println("terurut: Keys + Sort      :", k)

	// bukti bahwa urutan berubah-ubah: kumpulkan beragam urutan dari 200 iterasi
	unik := map[string]bool{}
	for i := 0; i < 200; i++ {
		unik[fmt.Sprint(Keys(m))] = true
	}
	fmt.Println("urutan berbeda dari 200 kali Keys(m):", len(unik))
}
