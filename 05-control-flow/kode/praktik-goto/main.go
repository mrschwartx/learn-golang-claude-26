// goto.go
package main

import "fmt"

func main() {
	percobaan := 0
ulang:
	percobaan++
	if percobaan < 3 {
		fmt.Println("percobaan", percobaan, "gagal, ulangi")
		goto ulang
	}
	fmt.Println("berhasil pada percobaan", percobaan)
}
