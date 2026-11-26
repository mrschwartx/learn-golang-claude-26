package a

import "contoh.com/melingkar/b"

type Pengguna struct{ Nama string }

func (p Pengguna) Pesanan() []b.Pesanan { return nil }
