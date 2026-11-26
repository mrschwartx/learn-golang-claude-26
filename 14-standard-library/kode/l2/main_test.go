package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestJalankan(t *testing.T) {
	tes := []struct {
		args   []string
		kode   int
		outHas string
		errHas string
	}{
		{[]string{"tambah", "-nama", "kopi", "-jumlah", "3", "x"}, 0, `tambah "kopi" x3`, ""},
		{[]string{"hapus", "-id", "7", "-paksa"}, 0, "hapus id=7 paksa=true", ""},
		{[]string{"tambah"}, 2, "", "-nama wajib"},
		{[]string{"hapus", "-id", "abc"}, 2, "", "invalid value"},
		{[]string{"hapus", "-tidakada"}, 2, "", "flag provided but not defined"},
		{[]string{"hapus", "-h"}, 0, "", "Usage of hapus"},
		{[]string{"ubah"}, 2, "", "tidak dikenal"},
		{nil, 2, "", "pemakaian"},
		// flag milik subcommand lain tidak berlaku
		{[]string{"hapus", "-nama", "x"}, 2, "", "not defined: -nama"},
	}
	for _, c := range tes {
		var o, e bytes.Buffer
		if got := Jalankan(c.args, &o, &e); got != c.kode {
			t.Errorf("%v: kode %d, mau %d", c.args, got, c.kode)
		}
		if !strings.Contains(o.String(), c.outHas) || !strings.Contains(e.String(), c.errHas) {
			t.Errorf("%v: out=%q err=%q", c.args, o.String(), e.String())
		}
	}
}
