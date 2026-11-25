// struct.go
package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Hewan struct{ Nama string }

type Anjing struct {
	Hewan
	Ras string
}

type Produk struct {
	ID      int       `json:"id"`
	Nama    string    `json:"nama"`
	Harga   float64   `json:"harga,omitempty"`
	Rahasia string    `json:"-"`
	Dibuat  time.Time `json:"dibuat,omitzero"`
}

func main() {
	a := Anjing{Hewan: Hewan{Nama: "Rex"}, Ras: "Labrador"}
	fmt.Println(a.Nama, a.Hewan.Nama)

	b, _ := json.Marshal(Produk{ID: 1, Nama: "Kopi", Rahasia: "x"})
	fmt.Println(string(b))

	var p Produk
	_ = json.Unmarshal([]byte(`{"id":2,"nama":"Teh","harga":12.5}`), &p)
	fmt.Printf("%+v\n", p)
}
