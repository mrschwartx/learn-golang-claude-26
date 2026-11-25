// variabel.go
package main

import "fmt"

var versi = "1.0" // level paket: wajib pakai var

func main() {
	var umur int
	var harga float64
	var aktif bool
	var nama string
	var data []int
	var peta map[string]int

	fmt.Printf("%T %v\n", umur, umur)
	fmt.Printf("%T %v\n", harga, harga)
	fmt.Printf("%T %v\n", aktif, aktif)
	fmt.Printf("%T %q\n", nama, nama)
	fmt.Printf("%T %v %t\n", data, data, data == nil)
	fmt.Printf("%T %v %t\n", peta, peta, peta == nil)

	kota := "Jakarta"
	a, b := 10, 3
	fmt.Println(versi, kota, a/b, a%b, float64(a)/float64(b))
}
