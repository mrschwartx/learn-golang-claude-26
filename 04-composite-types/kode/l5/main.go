// Latihan 5: var s []int vs []int{} vs make([]int, 0).
package main

import (
	"encoding/json"
	"fmt"
)

type Resp struct {
	Items []int `json:"items"`
}

func tampil(nama string, s []int) {
	b, _ := json.Marshal(Resp{s})
	fmt.Printf("%-18s nil=%-5t len=%d cap=%d json=%s\n", nama, s == nil, len(s), cap(s), b)
}

func main() {
	var a []int
	b := []int{}
	c := make([]int, 0)
	tampil("var s []int", a)
	tampil("[]int{}", b)
	tampil("make([]int, 0)", c)

	// ketiganya aman untuk append dan range
	a = append(a, 1)
	for range b {
	}
	fmt.Println("append ke nil slice:", a)

	// jebakan: perbandingan dengan nil
	fmt.Println("len(a[:0])==0 tapi nil?", a[:0] == nil)

	// decode: JSON null vs []
	var r1, r2 Resp
	_ = json.Unmarshal([]byte(`{"items":null}`), &r1)
	_ = json.Unmarshal([]byte(`{"items":[]}`), &r2)
	fmt.Println("decode null:", r1.Items == nil, "| decode []:", r2.Items == nil)
}
