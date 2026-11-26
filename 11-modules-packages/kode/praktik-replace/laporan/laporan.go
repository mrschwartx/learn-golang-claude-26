// Package laporan menyusun ringkasan pesanan dari package harga dan diskon.
package laporan

import (
	"fmt"
	"strings"

	"contoh.com/tokoku/diskon"
	"contoh.com/tokoku/internal/harga"
)

// Cetak mengembalikan laporan teks untuk daftar harga dengan satu persentase diskon.
func Cetak(barang []int, persen int) (string, error) {
	total := harga.Total(barang)
	akhir, err := diskon.Terapkan(total, persen)
	if err != nil {
		return "", fmt.Errorf("laporan: %w", err)
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Jumlah barang : %d\n", len(barang))
	fmt.Fprintf(&b, "Total         : Rp%d\n", total)
	fmt.Fprintf(&b, "Diskon        : %d%%\n", persen)
	fmt.Fprintf(&b, "Bayar         : Rp%d\n", akhir)
	return b.String(), nil
}
