// Latihan 1: Klasifikasi suhu dengan switch tanpa ekspresi.
package main

import "fmt"

func Klasifikasi(suhu float64) string {
	switch {
	case suhu < 15:
		return "dingin"
	case suhu < 25:
		return "sejuk"
	case suhu < 32:
		return "hangat"
	default:
		return "panas"
	}
}

func main() {
	for _, s := range []float64{-5, 14.9, 15, 24.99, 25, 31.5, 32, 40} {
		fmt.Printf("%6.2f°C -> %s\n", s, Klasifikasi(s))
	}
}
