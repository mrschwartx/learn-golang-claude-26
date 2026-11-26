package main

import (
	"fmt"

	"contoh.com/tokoku/diskon"
)

func main() {
	akhir, err := diskon.Terapkan(100000, 10)
	fmt.Println(akhir, err)
}
