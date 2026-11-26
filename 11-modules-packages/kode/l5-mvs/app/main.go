package main

import (
	"fmt"

	"contoh.com/deta"
	"contoh.com/detb"
	"contoh.com/lib"
)

func main() {
	fmt.Println(deta.Versi())
	fmt.Println(detb.Versi())
	fmt.Println("yang benar-benar dibangun: lib", lib.Versi)
}
