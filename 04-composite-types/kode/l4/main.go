// Latihan 4: embedding + tag JSON.
package main

import (
	"encoding/json"
	"fmt"
)

type Alamat struct {
	Jalan string `json:"jalan"`
	Kota  string `json:"kota"`
}

// Alamat di-embed TANPA tag -> field-nya "diratakan" ke objek luar.
type Karyawan struct {
	Alamat
	Nama    string  `json:"nama"`
	Gaji    float64 `json:"gaji,omitempty"`
	Rahasia string  `json:"-"`
}

// Pembanding: embed dengan tag -> menjadi objek bersarang.
type KaryawanBersarang struct {
	Alamat `json:"alamat"`
	Nama   string `json:"nama"`
}

func main() {
	k := Karyawan{Alamat: Alamat{"Jl. Merdeka 1", "Jakarta"}, Nama: "Andi", Rahasia: "pin-1234"}
	b, _ := json.MarshalIndent(k, "", "  ")
	fmt.Println("embed tanpa tag (gaji 0 -> hilang karena omitempty):")
	fmt.Println(string(b))

	k.Gaji = 9500000
	b, _ = json.Marshal(k)
	fmt.Println("\ndengan gaji:", string(b))

	fmt.Println("\nakses terpromosi:", k.Kota, "==", k.Alamat.Kota)

	b, _ = json.Marshal(KaryawanBersarang{Alamat{"Jl. Merdeka 1", "Jakarta"}, "Andi"})
	fmt.Println("\nembed dengan tag json:\"alamat\":", string(b))

	var dari Karyawan
	_ = json.Unmarshal([]byte(`{"jalan":"Jl. Baru","kota":"Bandung","nama":"Budi","Rahasia":"x"}`), &dari)
	fmt.Printf("\ndecode: %+v\n", dari)
}
