// Latihan 1: berbagi backing array antara s dan t.
package main

import "fmt"

func info(nama string, s []int) {
	fmt.Printf("%-3s %v len=%d cap=%d\n", nama, s, len(s), cap(s))
}

func main() {
	s := []int{1, 2, 3, 4, 5}
	t := s[1:4] // [2 3 4], len=3, cap=4 (dari indeks 1 sampai akhir array)
	info("s", s)
	info("t", t)

	t[0] = 99 // menulis ke s[1]
	fmt.Println("setelah t[0]=99:")
	info("s", s)

	t = append(t, 100) // len 3 < cap 4: muat -> menimpa s[4], array sama
	fmt.Println("setelah append(t, 100):")
	info("s", s)
	info("t", t)

	t = append(t, 200) // len 4 == cap 4: realokasi -> array BARU, lepas dari s
	t[0] = -1
	fmt.Println("setelah append kedua (realokasi) dan t[0]=-1:")
	info("s", s)
	info("t", t)
}
