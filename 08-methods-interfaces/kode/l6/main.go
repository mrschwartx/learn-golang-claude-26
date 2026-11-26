// Latihan 6: io.Reader kustom HurufBesar, disalin dengan io.Copy.
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// HurufBesar membungkus io.Reader lain (decorator).
type HurufBesar struct{ R io.Reader }

func (h HurufBesar) Read(p []byte) (int, error) {
	n, err := h.R.Read(p) // baca dari sumber, lalu ubah byte yang terbaca
	copy(p[:n], bytes.ToUpper(p[:n]))
	return n, err // teruskan err (termasuk io.EOF) apa adanya
}

func main() {
	src := strings.NewReader("halo, dunia interface go!\n")
	n, err := io.Copy(os.Stdout, HurufBesar{src})
	fmt.Println("disalin:", n, "byte, err:", err)

	// bisa ditumpuk dengan reader lain: LimitReader membatasi 10 byte
	src2 := strings.NewReader("selamat belajar go\n")
	io.Copy(os.Stdout, io.LimitReader(HurufBesar{src2}, 10))
	fmt.Println()

	// atau menerima file/stdin/HTTP body yang sama tanpa mengubah HurufBesar
	var _ io.Reader = HurufBesar{}
}
