// hitung_byte.go
package main

import (
	"fmt"
	"io"
	"os"
)

// PenghitungByte membungkus io.Writer dan menghitung byte yang lewat.
type PenghitungByte struct {
	io.Writer // embedded interface: Write diteruskan otomatis
	Total     int
}

func (p *PenghitungByte) Write(b []byte) (int, error) {
	n, err := p.Writer.Write(b)
	p.Total += n
	return n, err
}

func main() {
	pb := &PenghitungByte{Writer: os.Stdout}
	fmt.Fprintln(pb, "Halo interface!")
	fmt.Fprintf(pb, "angka: %d\n", 42)
	fmt.Println("total byte:", pb.Total)
}
