// Latihan 1: Reset(*int) vs ResetNilai(int).
package main

import "fmt"

func Reset(p *int)     { *p = 0 }
func ResetNilai(n int) { n = 0 }

func main() {
	a, b := 42, 42
	ResetNilai(a)
	Reset(&b)
	fmt.Println("setelah ResetNilai(a):", a)
	fmt.Println("setelah Reset(&b)    :", b)

	fmt.Printf("alamat a=%p b=%p\n", &a, &b)
	func(n int) { fmt.Printf("alamat parameter n (salinan)  : %p\n", &n) }(a)
	func(p *int) { fmt.Printf("nilai pointer p (alamat milik b): %p\n", p) }(&b)
}
