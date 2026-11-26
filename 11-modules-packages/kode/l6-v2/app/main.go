package main

import (
	"fmt"

	pustaka "contoh.com/pustaka"
	pustakav2 "contoh.com/pustaka/v2"
)

func main() {
	fmt.Println(pustaka.Sapa("Andi"))
	fmt.Println(pustakav2.Sapa("Andi", "Selamat pagi"))
}
