// Package main (l4): tipe Rupiah dengan JSON "Rp 15.000".
package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Rupiah int64

func (r Rupiah) String() string {
	n := int64(r)
	neg := n < 0
	if neg {
		n = -n
	}
	s := strconv.FormatInt(n, 10)
	var parts []string
	for len(s) > 3 {
		parts = append([]string{s[len(s)-3:]}, parts...)
		s = s[:len(s)-3]
	}
	parts = append([]string{s}, parts...)
	out := "Rp " + strings.Join(parts, ".")
	if neg {
		out = "-" + out
	}
	return out
}

func (r Rupiah) MarshalJSON() ([]byte, error) { return json.Marshal(r.String()) }

func (r *Rupiah) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("rupiah harus string: %w", err)
	}
	neg := strings.HasPrefix(s, "-")
	s = strings.TrimPrefix(s, "-")
	s, ok := strings.CutPrefix(s, "Rp ")
	if !ok {
		return fmt.Errorf("rupiah %q: awalan %q tidak ada", s, "Rp ")
	}
	// validasi pemisah ribuan: kelompok pertama 1-3 digit, sisanya tepat 3 digit
	for i, g := range strings.Split(s, ".") {
		if len(g) == 0 || len(g) > 3 || (i > 0 && len(g) != 3) {
			return fmt.Errorf("rupiah %q: format ribuan salah", s)
		}
	}
	n, err := strconv.ParseInt(strings.ReplaceAll(s, ".", ""), 10, 64)
	if err != nil {
		return fmt.Errorf("rupiah %q: %w", s, err)
	}
	if neg {
		n = -n
	}
	*r = Rupiah(n)
	return nil
}

type Produk struct {
	Nama  string `json:"nama"`
	Harga Rupiah `json:"harga"`
}

func main() {
	p := Produk{"Kopi", 15000}
	b, _ := json.Marshal(p)
	fmt.Println(string(b))
	var q Produk
	fmt.Println(json.Unmarshal(b, &q), q == p, q)
	for _, bad := range []string{`{"harga":15000}`, `{"harga":"15000"}`, `{"harga":"Rp 15.00"}`, `{"harga":"Rp 1.5000"}`, `{"harga":"Rp abc"}`} {
		var x Produk
		fmt.Printf("%-24s -> %v\n", bad, json.Unmarshal([]byte(bad), &x))
	}
}
