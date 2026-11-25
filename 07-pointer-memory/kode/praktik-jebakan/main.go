// jebakan.go
package main

import "fmt"

type Item struct{ Jumlah int }

func tambahSalah(s []int)       { s = append(s, 99) }
func tambahBenar(s []int) []int { return append(s, 99) }
func tambahPtr(s *[]int)        { *s = append(*s, 99) }

func main() {
	s := make([]int, 1, 10)
	tambahSalah(s)
	fmt.Println("salah:", s)
	s = tambahBenar(s)
	fmt.Println("benar:", s)
	tambahPtr(&s)
	fmt.Println("pointer:", s)

	m := map[string]Item{"a": {1}}
	// m["a"].Jumlah = 5  // error: cannot assign to struct field in map
	v := m["a"]
	v.Jumlah = 5
	m["a"] = v
	fmt.Println(m)

	mp := map[string]*Item{"b": {1}}
	mp["b"].Jumlah = 9 // boleh: pointer dapat dideferensikan
	fmt.Println(mp["b"].Jumlah)
}
