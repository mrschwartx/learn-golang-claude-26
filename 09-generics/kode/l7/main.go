// Latihan 7: Proses[T any] dengan switch any(x).(type): tanda desain keliru, dan alternatifnya.
package main

import (
	"fmt"
	"strings"
)

// ---------- KELIRU: generik, tapi di dalamnya mengecek tipe satu per satu ----------
func ProsesKeliru[T any](x T) string {
	switch v := any(x).(type) {
	case int:
		return fmt.Sprint("int ", v*2)
	case string:
		return "string " + strings.ToUpper(v)
	default:
		return "tidak didukung"
	}
}

// ---------- ALTERNATIF 1: fungsi terpisah per tipe ----------
func ProsesInt(x int) string       { return fmt.Sprint("int ", x*2) }
func ProsesString(x string) string { return "string " + strings.ToUpper(x) }

// ---------- ALTERNATIF 2: interface; tiap tipe menentukan perilakunya sendiri ----------
type Pemroses interface{ Proses() string }

type Angka int
type Teks string

func (a Angka) Proses() string { return fmt.Sprint("int ", int(a)*2) }
func (t Teks) Proses() string  { return "string " + strings.ToUpper(string(t)) }

func ProsesInterface(p Pemroses) string { return p.Proses() }

func main() {
	fmt.Println(ProsesKeliru(21), "|", ProsesKeliru("go"))
	fmt.Println(ProsesKeliru(3.14), "<- float64 lolos kompilasi, baru gagal saat RUNTIME")

	fmt.Println(ProsesInt(21), "|", ProsesString("go"))
	fmt.Println(ProsesInterface(Angka(21)), "|", ProsesInterface(Teks("go")))
	// ProsesInterface(3.14) // error KOMPILASI: float64 does not implement Pemroses
}
