package main

import (
	"encoding/json"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	for _, n := range []Rupiah{0, 5, 999, 1000, 15000, 1234567, 1000000000, -2500} {
		b, err := json.Marshal(Produk{"x", n})
		if err != nil {
			t.Fatal(err)
		}
		var q Produk
		if err := json.Unmarshal(b, &q); err != nil || q.Harga != n {
			t.Errorf("%d -> %s -> %d (%v)", n, b, q.Harga, err)
		}
	}
	if s := Rupiah(15000).String(); s != "Rp 15.000" {
		t.Errorf("got %q", s)
	}
}

func TestInvalid(t *testing.T) {
	for _, in := range []string{`15000`, `"15000"`, `"Rp 15.00"`, `"Rp 1.5000"`, `"Rp "`, `"Rp a"`, `"Rp 9999999999999999999999"`} {
		var r Rupiah
		if err := json.Unmarshal([]byte(in), &r); err == nil {
			t.Errorf("%s harus error", in)
		}
	}
}

func TestSliceDanMap(t *testing.T) {
	var m map[string]Rupiah
	if err := json.Unmarshal([]byte(`{"a":"Rp 1.000","b":"Rp 20.000"}`), &m); err != nil || m["b"] != 20000 {
		t.Fatalf("%v %v", m, err)
	}
}
