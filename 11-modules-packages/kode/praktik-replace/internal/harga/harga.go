package harga

// Total menjumlahkan harga barang (dalam rupiah).
func Total(barang []int) int {
	total := 0
	for _, h := range barang {
		total += h
	}
	return total
}
