// Latihan 7: bahaya menyimpan irisan buffer yang dipakai ulang.
package main

import (
	"fmt"
	"slices"
)

type Pesan struct{ Header []byte }

func main() {
	buf := []byte("0123456789HEADER-ABCxxxxxxxx")

	// BERBAHAYA: Header berbagi backing array dengan buf
	salah := Pesan{Header: buf[10:20]}
	// AMAN 1: salin
	aman1 := Pesan{Header: slices.Clone(buf[10:20])}
	// AMAN 2 (sebagian): full slice expression membatasi cap, append tak menimpa buf
	batas := Pesan{Header: buf[10:20:20]}

	fmt.Printf("sebelum  : salah=%q aman1=%q batas=%q\n", salah.Header, aman1.Header, batas.Header)

	// buffer dipakai ulang untuk data berikutnya
	copy(buf, []byte("################ZZZZZZZZZZZZ"))
	fmt.Printf("sesudah  : salah=%q aman1=%q batas=%q\n", salah.Header, aman1.Header, batas.Header)

	// append pada slice dengan cap tersisa menimpa buf
	buf2 := []byte("0123456789HEADER-ABCxxxxxxxx")
	longgar := buf2[10:16] // cap masih panjang
	_ = append(longgar, '!', '!')
	fmt.Printf("append pada irisan longgar menimpa buf2: %q\n", buf2)

	buf3 := []byte("0123456789HEADER-ABCxxxxxxxx")
	ketat := buf3[10:16:16]
	_ = append(ketat, '!', '!')
	fmt.Printf("append pada irisan ber-cap ketat      : %q (buf3 utuh)\n", buf3)
}
