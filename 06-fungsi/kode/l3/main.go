// Latihan 3: Map, Filter, Reduce lalu digabung.
package main

import "fmt"

func Map(data []int, f func(int) int) []int {
	hasil := make([]int, 0, len(data))
	for _, v := range data {
		hasil = append(hasil, f(v))
	}
	return hasil
}

func Filter(data []int, f func(int) bool) []int {
	hasil := make([]int, 0, len(data))
	for _, v := range data {
		if f(v) {
			hasil = append(hasil, v)
		}
	}
	return hasil
}

func Reduce(data []int, awal int, f func(akum, v int) int) int {
	akum := awal
	for _, v := range data {
		akum = f(akum, v)
	}
	return akum
}

func main() {
	angka := []int{1, 2, 3, 4, 5, 6}

	kuadrat := Map(angka, func(x int) int { return x * x })
	genap := Filter(kuadrat, func(x int) bool { return x%2 == 0 })
	jumlah := Reduce(genap, 0, func(a, v int) int { return a + v })

	fmt.Println("kuadrat:", kuadrat)
	fmt.Println("genap  :", genap)
	fmt.Println("jumlah :", jumlah)

	// digabung dalam satu ekspresi
	fmt.Println("digabung:", Reduce(
		Filter(Map(angka, func(x int) int { return x * x }), func(x int) bool { return x%2 == 0 }),
		0, func(a, v int) int { return a + v }))
}
