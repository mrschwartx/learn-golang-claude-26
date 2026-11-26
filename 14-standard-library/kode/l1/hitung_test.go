package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHitung(t *testing.T) {
	p := filepath.Join(t.TempDir(), "a.txt")
	os.WriteFile(p, []byte("Go go GO\nhalo dunia\n\nçé ok\n"), 0o644)
	f, _ := os.Open(p)
	defer f.Close()
	h, err := Hitung(f)
	if err != nil {
		t.Fatal(err)
	}
	if h.Baris != 4 || h.Kata != 7 || h.Frek["go"] != 3 {
		t.Fatalf("salah: %+v", h)
	}
	// "Go go GO\n"=9, "halo dunia\n"=11, "\n"=1, "çé ok\n"=6 (rune, bukan byte)
	if h.Rune != 27 {
		t.Fatalf("rune=%d, mau 27", h.Rune)
	}
}
