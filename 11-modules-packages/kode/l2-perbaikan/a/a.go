package a

import "contoh.com/perbaikan/model"

func PesananDari(p model.Pengguna) []model.Pesanan {
	return []model.Pesanan{{ID: 1, Pemilik: p}}
}
