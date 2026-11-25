// Latihan 6: array disalin ke fungsi, slice berbagi data; append tidak mengubah len milik pemanggil.
package main

import "fmt"

func ubahArray(a [3]int)   { a[0] = 999 }
func ubahSlice(s []int)    { s[0] = 999 }
func tambah(s []int) []int { return append(s, 42) }
func tambahLupa(s []int)   { s = append(s, 42); fmt.Printf("   di dalam fungsi: %v len=%d\n", s, len(s)) }

func main() {
	arr := [3]int{1, 2, 3}
	ubahArray(arr)
	fmt.Println("array setelah ubahArray:", arr, "(tidak berubah: fungsi menerima salinan)")

	sl := []int{1, 2, 3}
	ubahSlice(sl)
	fmt.Println("slice setelah ubahSlice:", sl, "(berubah: header disalin, pointer ke array sama)")

	// append di dalam fungsi tanpa mengembalikan hasil
	fmt.Println("tambahLupa(sl):")
	tambahLupa(sl)
	fmt.Printf("   di pemanggil   : %v len=%d (len TIDAK berubah)\n", sl, len(sl))

	// kasus ada kapasitas sisa: data tertulis ke array bersama, tapi len pemanggil tetap
	sisa := make([]int, 3, 10)
	tambahLupa(sisa)
	fmt.Printf("   di pemanggil   : %v len=%d | sisa[:4]=%v (elemen ke-4 TERTULIS di array bersama)\n",
		sisa, len(sisa), sisa[:4])

	sl = tambah(sl)
	fmt.Println("pola yang benar, sl = tambah(sl):", sl)
}
