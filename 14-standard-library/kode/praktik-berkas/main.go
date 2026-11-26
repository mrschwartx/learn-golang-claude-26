// berkas.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dir, err := os.MkdirTemp("", "belajar-go-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "catatan.txt")
	isi := "baris satu\nbaris dua\nbaris tiga\n"
	if err := os.WriteFile(path, []byte(isi), 0o644); err != nil {
		panic(err)
	}

	data, _ := os.ReadFile(path)
	fmt.Print(string(data))

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	n := 0
	for sc.Scan() {
		n++
		fmt.Printf("%d: %s\n", n, strings.ToUpper(sc.Text()))
	}
	if err := sc.Err(); err != nil {
		panic(err)
	}

	// bufio.Writer: jangan lupa Flush
	out, _ := os.Create(filepath.Join(dir, "keluar.txt"))
	w := bufio.NewWriter(out)
	fmt.Fprintln(w, "ditulis lewat buffer")
	w.Flush()
	out.Close()

	entri, _ := os.ReadDir(dir)
	for _, e := range entri {
		fmt.Println("berkas:", e.Name())
	}

	if _, err := os.Open("tidak-ada.txt"); err != nil {
		fmt.Println(os.IsNotExist(err))
	}
}
