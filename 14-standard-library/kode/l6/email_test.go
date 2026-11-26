package main

import "testing"

func TestPecah(t *testing.T) {
	n, d, ok := Pecah("budi.santoso+tag@mail.example.co.id")
	if !ok || n != "budi.santoso+tag" || d != "mail.example.co.id" {
		t.Fatalf("%q %q %v", n, d, ok)
	}
	for _, s := range []string{"", "a@b", "a@@b.com", "a b@c.com", "a@-b.com", "a@b..com", "@c.com", "a@c.d"} {
		if _, _, ok := Pecah(s); ok {
			t.Errorf("%q harus ditolak", s)
		}
	}
}
