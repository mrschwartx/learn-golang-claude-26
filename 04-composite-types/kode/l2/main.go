// Latihan 2: HapusDuplikat dengan map, dibandingkan dengan slices.Sort + slices.Compact.
package main

import (
	"fmt"
	"slices"
)

// Mempertahankan urutan kemunculan pertama. O(n) waktu, O(n) memori tambahan.
func HapusDuplikat(s []string) []string {
	terlihat := make(map[string]struct{}, len(s))
	hasil := make([]string, 0, len(s))
	for _, v := range s {
		if _, ada := terlihat[v]; ada {
			continue
		}
		terlihat[v] = struct{}{}
		hasil = append(hasil, v)
	}
	return hasil
}

// Versi sort+compact: urutan asli HILANG (jadi terurut). O(n log n), tanpa map.
func HapusDuplikatTerurut(s []string) []string {
	c := slices.Clone(s) // jangan mengubah input
	slices.Sort(c)
	return slices.Compact(c)
}

func main() {
	in := []string{"b", "a", "b", "c", "a", "d", "c"}
	fmt.Println("input         :", in)
	fmt.Println("dengan map    :", HapusDuplikat(in))
	fmt.Println("sort+compact  :", HapusDuplikatTerurut(in))
	fmt.Println("input utuh    :", in)
	fmt.Println("input kosong  :", HapusDuplikat(nil), len(HapusDuplikat(nil)))
}
