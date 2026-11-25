// Latihan 7: uang dengan float64 vs integer satuan terkecil (sen/rupiah).
package main

import "fmt"

// Uang disimpan sebagai jumlah sen (int64). Tidak ada pecahan -> tidak ada galat biner.
type Uang int64

func (u Uang) String() string {
	tanda := ""
	if u < 0 {
		tanda, u = "-", -u
	}
	return fmt.Sprintf("%sRp%d,%02d", tanda, int64(u)/100, int64(u)%100)
}

// KaliPersen menghitung persentase dengan pembulatan eksplisit (half-up) pada integer.
func (u Uang) KaliPersen(p int64) Uang { return Uang((int64(u)*p + 50) / 100) }

func main() {
	// float64: galat representasi biner
	var f float64
	for i := 0; i < 10; i++ {
		f += 0.1
	}
	fmt.Println("float64 10x0.1  :", f, "| == 1.0?", f == 1.0)
	a, b := 0.1, 0.2
	fmt.Println("float64 0.1+0.2 :", a+b, "| == 0.3?", a+b == 0.3)

	var jumlah float64
	for i := 0; i < 1000; i++ {
		jumlah += 19.99
	}
	fmt.Printf("float64 1000x19.99: %.10f (seharusnya 19990)\n", jumlah)

	// integer sen
	var total Uang
	for i := 0; i < 1000; i++ {
		total += 1999
	}
	fmt.Println("int64 sen 1000x1999:", total, "| tepat 19990,00?", total == 1999000)

	harga := Uang(1999)
	fmt.Println("diskon 15% dari", harga, "=", harga.KaliPersen(15))
	fmt.Println("PPN 11% dari", harga, "=", harga.KaliPersen(11))
}
