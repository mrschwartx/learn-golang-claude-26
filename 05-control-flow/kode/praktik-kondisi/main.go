// kondisi.go
package main

import (
	"fmt"
	"strconv"
)

func nilai(skor int) string {
	switch {
	case skor >= 90:
		return "A"
	case skor >= 80:
		return "B"
	case skor >= 70:
		return "C"
	default:
		return "D"
	}
}

func main() {
	if n, err := strconv.Atoi("42"); err != nil {
		fmt.Println("error:", err)
	} else if n%2 == 0 {
		fmt.Println(n, "genap")
	}

	for _, s := range []int{95, 82, 71, 40} {
		fmt.Println(s, nilai(s))
	}

	switch hari := "sabtu"; hari {
	case "sabtu", "minggu":
		fmt.Println("akhir pekan")
	default:
		fmt.Println("hari kerja")
	}

	switch x := 2; x {
	case 2:
		fmt.Println("dua")
		fallthrough
	case 3:
		fmt.Println("tiga (fallthrough, tanpa cek kondisi)")
	case 4:
		fmt.Println("tidak dijalankan")
	}
}
