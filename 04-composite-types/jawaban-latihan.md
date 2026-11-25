# Modul 4: Composite Types (Array, Slice, Map, Struct) — Jawaban Latihan

Seluruh kode ada di folder `kode/` (modul Go `contoh.com/modul4`, satu direktori per program). Jalankan dari dalam `kode/`:

```bash
go run ./l1        # latihan 1, dst. sampai ./l7
go run ./praktik-slice   # juga: praktik-array, praktik-map, praktik-struct
```

Semua output di bawah dijalankan sungguhan di Go 1.24.7; `go vet ./...` dan `gofmt -l` bersih.

---

## Latihan 1 — `t := s[1:4]`, ubah `t[0]`, lalu `append` (`kode/l1`)

```
s   [1 2 3 4 5] len=5 cap=5
t   [2 3 4] len=3 cap=4
setelah t[0]=99:
s   [1 99 3 4 5] len=5 cap=5
setelah append(t, 100):
s   [1 99 3 4 100] len=5 cap=5
t   [99 3 4 100] len=4 cap=4
setelah append kedua (realokasi) dan t[0]=-1:
s   [1 99 3 4 100] len=5 cap=5
t   [-1 3 4 100 200] len=5 cap=8
```

**Penjelasan rinci:**

- `s` punya `len=5, cap=5`. `t := s[1:4]` tidak menyalin data: `t` menunjuk ke `s[1]` di backing array yang sama, dengan **`len=3`** (elemen 2, 3, 4) dan **`cap=4`** (dari indeks 1 sampai akhir array).
- `t[0] = 99` menulis ke `s[1]` (alamat yang sama), sehingga **`s[1]` ikut berubah menjadi 99**.
- `append(t, 100)`: `len(t)=3 < cap(t)=4`, jadi masih muat. Elemen baru ditulis ke slot setelah `t[2]`, yaitu **`s[4]`**, yang berubah dari 5 menjadi 100. `s` sendiri tidak berubah panjangnya (tetap 5), tetapi isinya berubah. Setelah ini `t` punya `len=4, cap=4`.
- `append(t, 200)` yang kedua: `len == cap`, sehingga Go **mengalokasikan array baru** (cap menjadi 8), menyalin isi, dan `t` kini menunjuk ke array baru. Sejak saat itu `t[0] = -1` **tidak lagi memengaruhi `s`**.

Ringkasnya: elemen `s[1]` berubah karena penugasan lewat `t`, elemen `s[4]` berubah karena `append` pertama (masih dalam kapasitas), dan `s[0]`, `s[2]`, `s[3]` tidak berubah.

---

## Latihan 2 — `HapusDuplikat` (`kode/l2`)

```go
func HapusDuplikat(s []string) []string {
	terlihat := make(map[string]struct{}, len(s))
	hasil := make([]string, 0, len(s))
	for _, v := range s {
		if _, ada := terlihat[v]; ada {
			continue
		}
		terlihat[v] = struct{}{}
		hasil = append(hasil, v)
	}
	return hasil
}

func HapusDuplikatTerurut(s []string) []string {
	c := slices.Clone(s)
	slices.Sort(c)
	return slices.Compact(c)
}
```

```
input         : [b a b c a d c]
dengan map    : [b a c d]
sort+compact  : [a b c d]
input utuh    : [b a b c a d c]
input kosong  : [] 0
```

| Aspek | Map sebagai himpunan | `slices.Sort` + `slices.Compact` |
|---|---|---|
| Urutan hasil | **Urutan kemunculan pertama dipertahankan** | Urutan asli hilang (hasil terurut) |
| Waktu | O(n) rata-rata | O(n log n) |
| Memori tambahan | Map berukuran sampai n | Hanya salinan slice (atau nol jika boleh mengubah input) |
| Syarat tipe elemen | `comparable` | `cmp.Ordered` untuk `Sort` (atau `SortFunc`) |

`Compact` hanya membuang duplikat yang **bersebelahan**, itulah sebabnya slice harus diurutkan lebih dulu. Selain itu `Compact` mengubah isi slice masukan, sehingga contoh di atas memakai `slices.Clone` agar input aslinya aman. Pakai versi map bila urutan penting, dan versi sort bila hasil terurut memang diinginkan atau memori terbatas.

---

## Latihan 3 — Frekuensi kata, 3 teratas (`kode/l3`)

Paragraf uji:

> Go itu sederhana. Go itu cepat, dan Go itu menyenangkan! Belajar Go berarti belajar menulis kode yang sederhana dan mudah dibaca.

```
3 kata teratas:
1. go         4
2. itu        3
3. belajar    2
```

Langkahnya: ubah ke huruf kecil, pisah kata dengan `strings.FieldsFunc` (membuang tanda baca), hitung dengan `map[string]int`, pindahkan ke `[]Frek`, lalu `slices.SortFunc` (jumlah menurun, seri diputus alfabetis agar hasil deterministik).

**Mengapa map tidak bisa langsung diurutkan?** Map adalah hash table: tidak punya urutan bawaan, dan urutan iterasinya **sengaja diacak** oleh runtime Go agar program tidak bergantung padanya. Tidak ada operasi "sort map" di bahasa maupun library. Cara umum: ambil key atau pasangan key-value ke slice, urutkan slice-nya. Perhatikan juga pada hasil di atas bahwa tiga kata di peringkat ketiga seri (`belajar`, `dan`, `sederhana`, masing-masing 2); tanpa aturan pemutus seri hasilnya akan berubah-ubah antar-eksekusi karena iterasi map acak.

---

## Latihan 4 — Struct `Alamat` di-embed, encode JSON (`kode/l4`)

```go
type Alamat struct {
	Jalan string `json:"jalan"`
	Kota  string `json:"kota"`
}

type Karyawan struct {
	Alamat
	Nama    string  `json:"nama"`
	Gaji    float64 `json:"gaji,omitempty"`
	Rahasia string  `json:"-"`
}
```

```
embed tanpa tag (gaji 0 -> hilang karena omitempty):
{
  "jalan": "Jl. Merdeka 1",
  "kota": "Jakarta",
  "nama": "Andi"
}

dengan gaji: {"jalan":"Jl. Merdeka 1","kota":"Jakarta","nama":"Andi","gaji":9500000}

akses terpromosi: Jakarta == Jakarta

embed dengan tag json:"alamat": {"alamat":{"jalan":"Jl. Merdeka 1","kota":"Jakarta"},"nama":"Andi"}
```

**Yang terjadi saat struct di-embed:**

- Tanpa tag pada field yang di-embed, `encoding/json` **meratakan** (flatten) field-fieldnya ke objek luar: `jalan` dan `kota` muncul sejajar dengan `nama`, bukan bersarang di dalam objek `alamat`. Ini mengikuti perilaku field yang "dipromosikan" di Go.
- Jika field embedded diberi tag nama (`` Alamat `json:"alamat"` ``), JSON menjadi **objek bersarang** `"alamat": {...}`.
- `json:"-"` membuat `Rahasia` tidak pernah muncul (baik saat encode maupun decode), dan `omitempty` membuat `gaji` hilang selama nilainya 0.
- Peringatan: bila `Alamat` dan struct luar punya field dengan nama JSON yang sama, aturan kedalaman dan konflik JSON berlaku (field di level lebih dangkal menang, dan konflik di level yang sama membuat keduanya diabaikan). Periksa kembali bila struct yang di-embed lebih dari satu.

---

## Latihan 5 — `var s []int` vs `[]int{}` vs `make([]int, 0)` (`kode/l5`)

```
var s []int        nil=true  len=0 cap=0 json={"items":null}
[]int{}            nil=false len=0 cap=0 json={"items":[]}
make([]int, 0)     nil=false len=0 cap=0 json={"items":[]}
append ke nil slice: [1]
decode null: true | decode []: false
```

| Bentuk | `== nil` | `len` | `cap` | JSON |
|---|---|---|---|---|
| `var s []int` | **true** | 0 | 0 | `null` |
| `s := []int{}` | false | 0 | 0 | `[]` |
| `s := make([]int, 0)` | false | 0 | 0 | `[]` |

`[]int{}` dan `make([]int, 0)` setara untuk tujuan praktis. Ketiganya aman dipakai dengan `append`, `range`, dan `len`, sehingga untuk memeriksa "kosong" gunakan **`len(s) == 0`**, bukan `s == nil`.

**Kapan perbedaan ini penting di API:**

- **Respons JSON:** klien (misalnya JavaScript yang memanggil `.length` atau `.map`) bisa error bila menerima `null` padahal mengharapkan `[]`. Bila kontrak API menjanjikan array, inisialisasi dengan `[]T{}` atau `make` sebelum di-encode.
- **Membedakan "tidak dikirim" dan "kosong"** pada decode: `null` menghasilkan slice `nil`, sedangkan `[]` menghasilkan slice non-nil kosong. Berguna untuk operasi *PATCH* (ubah menjadi kosong vs tidak diubah).
- **Perbandingan dan test:** `reflect.DeepEqual` serta beberapa library test membedakan `nil` dan kosong; `slices.Equal` tidak.

---

## Latihan 6 — Array disalin, slice tidak (`kode/l6`)

```
array setelah ubahArray: [1 2 3] (tidak berubah: fungsi menerima salinan)
slice setelah ubahSlice: [999 2 3] (berubah: header disalin, pointer ke array sama)
tambahLupa(sl):
   di dalam fungsi: [999 2 3 42] len=4
   di pemanggil   : [999 2 3] len=3 (len TIDAK berubah)
   di dalam fungsi: [0 0 0 42] len=4
   di pemanggil   : [0 0 0] len=3 | sisa[:4]=[0 0 0 42] (elemen ke-4 TERTULIS di array bersama)
pola yang benar, sl = tambah(sl): [999 2 3 42]
```

**Penjelasan:**

- Array bersifat *value semantics*: argumen yang dikirim ke fungsi adalah **salinan seluruh isinya**, sehingga perubahan di dalam fungsi tidak terlihat pemanggil.
- Go selalu *pass by value*, termasuk untuk slice. Yang disalin adalah **header slice** (pointer, `len`, `cap`), bukan isinya. Karena pointer di salinan menunjuk ke backing array yang sama, `s[0] = 999` di dalam fungsi mengubah data yang juga dilihat pemanggil.
- `append` di dalam fungsi hanya mengubah **salinan header** milik fungsi (`len` salinan itu menjadi 4). Header milik pemanggil tetap `len=3`, sehingga ia tidak melihat elemen baru. Karena itu `append` harus selalu ditampung dan dikembalikan: `sl = tambah(sl)`.
- Kasus halus (baris `sisa`): bila slice masih punya kapasitas sisa, `append` di fungsi **menulis ke backing array bersama** meskipun `len` pemanggil tidak berubah. Elemen ke-4 sebenarnya sudah terisi (`sisa[:4]` menampilkannya), sebuah sumber bug yang sulit dilacak.

---

## Latihan 7 — Berbagi backing array pada buffer yang dipakai ulang (`kode/l7`)

```
sebelum  : salah="HEADER-ABC" aman1="HEADER-ABC" batas="HEADER-ABC"
sesudah  : salah="######ZZZZ" aman1="HEADER-ABC" batas="######ZZZZ"
append pada irisan longgar menimpa buf2: "0123456789HEADER!!BCxxxxxxxx"
append pada irisan ber-cap ketat      : "0123456789HEADER-ABCxxxxxxxx" (buf3 utuh)
```

**Risiko berbagi backing array:**

1. **Data berubah diam-diam.** `buf[10:20]` bukan salinan, melainkan jendela ke buffer yang sama. Saat buffer dipakai ulang untuk pesan berikutnya, `Header` di struct jangka panjang ikut tertimpa (`"HEADER-ABC"` menjadi `"######ZZZZ"`). Ini *data corruption* tanpa error atau panic.
2. **Kebocoran memori.** Irisan kecil menahan **seluruh** backing array tetap hidup (garbage collector tidak bisa membebaskannya), sehingga struct yang kecil bisa menahan buffer besar.
3. **`append` menimpa data tetangga.** Irisan yang punya kapasitas sisa menulis ke bagian buffer di sebelahnya (`HEADER-ABC` menjadi `HEADER!!BC`).
4. **Race condition** bila buffer dipakai oleh goroutine lain secara bersamaan.

**Cara mengatasinya:**

- **`slices.Clone`** (atau `append([]byte(nil), x...)`, atau `copy`) membuat salinan independen: `aman1` tidak terpengaruh sama sekali. Ini solusi yang tepat untuk data yang disimpan jangka panjang, dan sekaligus melepas buffer besar agar bisa dibebaskan GC.
- **Full slice expression** `buf[10:20:20]` membatasi `cap`, sehingga `append` dipaksa mengalokasikan array baru dan tidak menimpa tetangga (baris terakhir output). Namun cara ini **tidak** menyelesaikan masalah penimpaan isi buffer oleh pemakaian ulang (lihat kolom `batas` pada output: tetap tertimpa), jadi bukan pengganti `Clone` untuk menyimpan data jangka panjang.

---

## Lampiran: program Praktik dari materi

Keempat program Praktik (`praktik-slice`, `praktik-array`, `praktik-map`, `praktik-struct`) sudah saya jalankan dan hasilnya sesuai teori (berbagi memori, perbedaan array dan slice, comma-ok, urutan key terurut, tag JSON, embedding).

Satu catatan: pola kapasitas pada `praktik-slice` di mesin ini adalah `1 2 4 8 16 32 64 128 256 512 848 1280 1792 2560`, sedangkan materi menyebut contoh `4 8 16 ...` untuk Go 1.27. Perbedaan ini wajar, sebab materi sendiri menegaskan bahwa angka kapasitas bergantung pada versi Go dan ukuran memory class, jadi jangan menggantungkan logika program pada angka tersebut. Bagian materi tentang Go 1.27 (key struct literal boleh berupa field yang dipromosikan) tidak saya coba karena toolchain di sini Go 1.24.
