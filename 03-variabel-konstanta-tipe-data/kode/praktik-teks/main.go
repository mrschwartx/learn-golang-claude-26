// teks.go
package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

func main() {
	s := "héllo, 世界"
	fmt.Println("byte:", len(s), "rune:", utf8.RuneCountInString(s))

	for i, r := range "héy" {
		fmt.Printf("posisi %d: %c (%U)\n", i, r, r)
	}

	raw := `baris 1
baris 2 \n tidak diproses`
	fmt.Println(raw)

	angka, err := strconv.Atoi("123abc")
	fmt.Println(angka, err)

	n, _ := strconv.Atoi("42")
	fmt.Println(n+1, strconv.Itoa(n)+"!", string(rune(65)))

	var sb strings.Builder
	for i := 0; i < 3; i++ {
		fmt.Fprintf(&sb, "[%d]", i)
	}
	fmt.Println(sb.String())

	f := 9.99
	fmt.Println(int(f))
}
