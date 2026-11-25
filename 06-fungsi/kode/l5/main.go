// Latihan 5: pass-by-value untuk struct, *struct, slice, map.
package main

import "fmt"

type Data struct{ N int }

func ubahStruct(d Data)          { d.N = 100 }
func ubahPointer(d *Data)        { d.N = 100 }
func ubahSliceElem(s []int)      { s[0] = 100 }
func appendSlice(s []int)        { s = append(s, 999) }
func ubahMap(m map[string]int)   { m["k"] = 100 }
func ubahSliceNil(s []int)       { s = []int{7, 7, 7} }
func tambahMap(m map[string]int) { m["baru"] = 1 }

func main() {
	d := Data{1}
	ubahStruct(d)
	fmt.Println("struct         :", d)

	d2 := Data{1}
	ubahPointer(&d2)
	fmt.Println("*struct        :", d2)

	s := []int{1, 2, 3}
	ubahSliceElem(s)
	fmt.Println("slice elemen   :", s)

	appendSlice(s)
	fmt.Println("slice append   :", s, "len =", len(s))

	ubahSliceNil(s)
	fmt.Println("slice diganti  :", s)

	m := map[string]int{"k": 1}
	ubahMap(m)
	tambahMap(m)
	fmt.Println("map            :", m)

	// bukti pass-by-value: alamat parameter berbeda dari alamat variabel asli
	x := 1
	fmt.Printf("alamat x di main: %p\n", &x)
	func(p *int) {
		fmt.Printf("nilai pointer   : %p (sama, karena yang disalin adalah alamatnya)\n", p)
		fmt.Printf("alamat variabel p: %p (berbeda, p adalah salinan)\n", &p)
	}(&x)
}
