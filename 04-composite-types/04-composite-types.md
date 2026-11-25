# Modul 4: Composite Types (Array, Slice, Map, Struct)

[← Modul 3: Variabel, Konstanta & Tipe Data](03-variabel-konstanta-tipe-data.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 5: Control Flow (Conditionals & Loops) →](05-control-flow.md)

**Isi modul ini**

- [Teori](#teori)
  - [4.1 Array — Ukuran Tetap](#41-array--ukuran-tetap)
  - [4.2 Slice — Struktur Data Andalan](#42-slice--struktur-data-andalan)
  - [4.3 Capacity dan Growth](#43-capacity-dan-growth)
  - [4.4 `make()`, `copy`, dan Pustaka `slices`](#44-make-copy-dan-pustaka-slices)
  - [4.5 Konversi Slice dan Array](#45-konversi-slice-dan-array)
  - [4.6 Strings sebagai Urutan Byte](#46-strings-sebagai-urutan-byte)
  - [4.7 Map dan Comma-Ok Idiom](#47-map-dan-comma-ok-idiom)
  - [4.8 Struct](#48-struct)
  - [4.9 Struct Tags dan JSON](#49-struct-tags-dan-json)
  - [4.10 Embedding Structs](#410-embedding-structs)
- [Praktik](#praktik)
  - [Slice: Sharing, Capacity, dan Growth](#slice-sharing-capacity-dan-growth)
  - [Array, Konversi, dan Struct](#array-konversi-dan-struct)
  - [Map, Comma-Ok, dan Urutan Key](#map-comma-ok-dan-urutan-key)
  - [Struct, Tag JSON, dan Embedding](#struct-tag-json-dan-embedding)
- [Latihan](#latihan)

---


## Teori

### 4.1 Array — Ukuran Tetap

**Array** adalah kumpulan elemen bertipe sama dengan **panjang tetap yang menjadi bagian dari tipenya**. Artinya `[3]int` dan `[4]int` adalah dua tipe yang berbeda dan tidak bisa saling ditukar.

```go
var a [3]int                    // [0 0 0]
b := [3]string{"x", "y", "z"}
c := [...]int{10, 20, 30}      // panjang dihitung compiler
```

Sifat penting array:

- **Value semantics**: menugaskan atau mengirim array ke fungsi **menyalin seluruh isinya**.
- Array bisa dibandingkan dengan `==` jika tipe elemennya *comparable*.
- Dalam praktik, array jarang dipakai langsung — **slice** (4.2) jauh lebih umum. Array berguna untuk data berukuran tetap seperti hash (`[32]byte` untuk SHA-256) atau key map komposit.

### 4.2 Slice — Struktur Data Andalan

**Slice** adalah "jendela" dinamis ke sebuah array di belakang layar. Secara internal sebuah slice hanyalah **header kecil berisi tiga hal**:

| Field | Arti |
|---|---|
| **pointer** | Menunjuk ke elemen pertama yang dilihat slice di dalam *backing array* |
| **length** (`len`) | Jumlah elemen yang bisa diakses |
| **capacity** (`cap`) | Jumlah elemen dari posisi awal slice sampai akhir backing array |

```go
s := []int{1, 2, 3, 4, 5}
t := s[1:3]   // [2 3] — berbagi backing array dengan s
t[0] = 99     // s menjadi [1 99 3 4 5]
```

**Slicing (`s[a:b]`) tidak menyalin data** — slice baru melihat memori yang sama. Ini efisien, tetapi artinya perubahan lewat satu slice terlihat oleh slice lain yang berbagi array yang sama.

**`nil` slice vs slice kosong**: `var s []int` bernilai `nil`, sedangkan `[]int{}` tidak `nil`; keduanya memiliki `len == 0` dan aman dipakai dengan `append` dan `range`. Perbedaannya terlihat saat dicek `s == nil` dan saat di-encode ke JSON (`null` vs `[]`). Gunakan `len(s) == 0` untuk memeriksa "kosong".

### 4.3 Capacity dan Growth

`append` menambahkan elemen ke slice dan **mengembalikan slice baru** — hasilnya harus selalu ditampung: `s = append(s, x)`.

- Jika masih ada kapasitas, elemen ditulis ke backing array yang sama.
- Jika kapasitas habis, Go **mengalokasikan array baru yang lebih besar**, menyalin isi lama, lalu slice yang dikembalikan menunjuk ke array baru tersebut.

Pola pertumbuhan secara umum: kapasitas **berlipat ganda** untuk slice kecil, lalu pertumbuhannya melambat (sekitar 1,25×) untuk slice besar. **Nilai persisnya bergantung pada versi Go dan ukuran memory class**, jadi jangan menggantungkan logika pada angka kapasitas tertentu.

Jika Anda tahu perkiraan jumlah elemen, **alokasikan kapasitas di awal** untuk menghindari salin-ulang berulang:

```go
hasil := make([]string, 0, len(input)) // len 0, cap len(input)
```

### 4.4 `make()`, `copy`, dan Pustaka `slices`

| Bentuk | Fungsi |
|---|---|
| `make([]T, n)` | Slice dengan `len = n`, `cap = n`, semua elemen zero value |
| `make([]T, 0, n)` | Slice kosong dengan kapasitas awal `n` |
| `copy(dst, src)` | Menyalin elemen; jumlah yang disalin = `min(len(dst), len(src))` |
| `s[a:b:c]` | **Full slice expression** — membatasi `cap` menjadi `c - a`, mencegah `append` menimpa data milik slice lain |
| `clear(s)` | Mengisi semua elemen dengan zero value (`clear(m)` mengosongkan map) |

Sejak Go 1.21, paket **`slices`** menyediakan operasi umum yang dulu harus ditulis manual:

| Fungsi | Kegunaan |
|---|---|
| `slices.Contains(s, x)` / `slices.Index(s, x)` | Cari elemen |
| `slices.Sort(s)` / `slices.SortFunc(s, cmp)` | Urutkan |
| `slices.Reverse(s)` | Balik urutan |
| `slices.Max(s)` / `slices.Min(s)` | Nilai terbesar/terkecil |
| `slices.Equal(a, b)` | Bandingkan isi |
| `slices.Delete(s, i, j)` / `slices.Insert(s, i, v...)` | Hapus/sisipkan elemen |
| `slices.Clone(s)` | Salin ke slice baru |

### 4.5 Konversi Slice dan Array

| Arah | Sintaks | Perilaku |
|---|---|---|
| Array → slice | `arr[:]` | Slice **berbagi memori** dengan array |
| Slice → array (salinan) | `[3]int(s)` (Go 1.20+) | **Menyalin** 3 elemen pertama; **panic** jika `len(s) < 3` |
| Slice → pointer ke array | `(*[3]int)(s)` (Go 1.17+) | **Berbagi memori**; panic jika `len(s) < 3` |

### 4.6 Strings sebagai Urutan Byte

Secara konseptual string mirip slice `[]byte` yang **hanya-baca**: `len(s)` menghitung byte, `s[i]` mengambil byte, dan `s[a:b]` membuat string baru yang berbagi memori. Konversi `[]byte(s)` dan `string(b)` **menyalin data**. Detail tentang rune dan UTF-8 ada di [Modul 3](03-variabel-konstanta-tipe-data.md).

### 4.7 Map dan Comma-Ok Idiom

**Map** adalah struktur *hash table* yang memetakan **key** ke **value**.

```go
umur := map[string]int{"andi": 30, "budi": 25}
umur["citra"] = 28
delete(umur, "budi")
```

Aturan penting map:

- Key harus bertipe **comparable** (`int`, `string`, struct sederhana, array). Slice, map, dan function **tidak boleh** jadi key.
- Membaca key yang tidak ada **tidak error** — mengembalikan zero value. Karena itu ada **comma-ok idiom** untuk membedakan "key tidak ada" dari "key ada dengan nilai zero".
- **Menulis ke map `nil` menyebabkan panic**; membaca dari map `nil` aman. Selalu inisialisasi dengan `make` atau literal.
- **Urutan iterasi map sengaja diacak** — jangan pernah bergantung pada urutan tertentu.
- Map **tidak aman** dipakai bersamaan dari banyak goroutine tanpa sinkronisasi ([Modul 13](13-sync-context-patterns.md)).

```go
if v, ok := umur["andi"]; ok {
    fmt.Println("ketemu:", v)
}
```

### 4.8 Struct

**Struct** mengelompokkan beberapa field menjadi satu tipe baru — padanan paling dekat dengan "class tanpa method dan inheritance" di bahasa lain.

```go
type User struct {
    Nama  string
    Umur  int
    Email string
}

u := User{Nama: "Andi", Umur: 30}      // field yang dilewati = zero value
p := &User{Nama: "Budi"}               // pointer ke struct
anon := struct{ X, Y int }{1, 2}       // anonymous struct
```

Poin penting: **menyertakan nama field** dalam literal lebih aman daripada bentuk posisional; struct bisa dibandingkan dengan `==` jika semua field-nya comparable; Go tidak punya constructor bawaan, sehingga konvensinya adalah fungsi `NewUser(...)` yang mengembalikan `*User` atau `User`.

### 4.9 Struct Tags dan JSON

**Struct tag** adalah string metadata di samping field, dibaca lewat reflection ([Modul 20](20-advanced-topics.md)) oleh paket seperti `encoding/json`, ORM, dan validator.

| Tag | Efek pada `encoding/json` |
|---|---|
| `` `json:"nama"` `` | Nama field di JSON menjadi `nama` |
| `` `json:"harga,omitempty"` `` | Field dihilangkan jika kosong (0, `""`, `nil`, slice kosong) |
| `` `json:"dibuat,omitzero"` `` | Field dihilangkan jika bernilai **zero value** (Go 1.24+) — cocok untuk `time.Time` |
| `` `json:"-"` `` | Field **tidak pernah** di-encode/decode |

Hanya field **exported** (berhuruf kapital) yang ikut di-encode oleh `encoding/json`.

### 4.10 Embedding Structs

**Embedding** menyertakan sebuah tipe di dalam struct tanpa memberi nama field; seluruh field dan method-nya "dipromosikan" (*promoted*) sehingga bisa diakses langsung. Ini cara Go mencapai **reuse lewat komposisi**, bukan inheritance.

```go
type Hewan struct{ Nama string }

type Anjing struct {
    Hewan          // embedded
    Ras string
}

a := Anjing{Hewan: Hewan{Nama: "Rex"}, Ras: "Labrador"}
fmt.Println(a.Nama)  // dipromosikan, setara a.Hewan.Nama
```

Embedding **bukan inheritance**: `Anjing` tidak bisa dipakai di tempat yang meminta `Hewan`. Hubungan "is-a" dimodelkan dengan interface ([Modul 8](08-methods-interfaces.md)). Mulai **Go 1.27**, key pada struct literal boleh berupa field yang dipromosikan, sehingga `Anjing{Nama: "Rex", Ras: "Labrador"}` juga valid.

## Praktik

### Slice: Sharing, Capacity, dan Growth

```go
// slice.go
package main

import (
	"fmt"
	"slices"
)

func main() {
	s := []int{1, 2, 3, 4, 5}
	t := s[1:3]
	fmt.Println("t:", t, "len:", len(t), "cap:", cap(t))
	t[0] = 99
	fmt.Println("s setelah t[0]=99:", s)

	// append pada t masih muat di kapasitas -> MENIMPA s[3]
	t = append(t, 100)
	fmt.Println("s setelah append ke t:", s)

	// full slice expression: cap dibatasi, append memaksa array baru
	u := s[1:3:3]
	u = append(u, 200)
	fmt.Println("s tidak berubah:", s, "u:", u)

	var x []int
	prev := -1
	for i := 0; i < 2000; i++ {
		x = append(x, i)
		if cap(x) != prev {
			fmt.Print(cap(x), " ")
			prev = cap(x)
		}
	}
	fmt.Println()

	angka := []int{5, 2, 8, 1}
	slices.Sort(angka)
	fmt.Println(angka, slices.Contains(angka, 8), slices.Index(angka, 5), slices.Max(angka))
	angka = slices.Delete(angka, 1, 3)
	fmt.Println(angka)
}
```

Pola kapasitas yang tercetak pada Go 1.27 di mesin 64-bit kurang lebih `4 8 16 32 64 128 256 512 848 1280 1792 2560` — angka bisa berbeda di versi Go lain.

### Array, Konversi, dan Struct

```go
// array.go
package main

import "fmt"

type Titik struct{ X, Y int }

func ubah(a [3]int) { a[0] = 999 } // menerima SALINAN

func main() {
	arr := [3]int{1, 2, 3}
	ubah(arr)
	fmt.Println("array asli:", arr)

	s := arr[:]
	s[0] = 100
	fmt.Println("array berubah lewat slice:", arr)

	salinan := [3]int(s)
	salinan[1] = -1
	fmt.Println("salinan:", salinan, "asli:", arr)

	p1, p2 := Titik{1, 2}, Titik{1, 2}
	fmt.Println("sama:", p1 == p2)

	grid := map[Titik]string{{0, 0}: "asal", {1, 2}: "target"}
	fmt.Println(grid[Titik{1, 2}])
}
```

### Map, Comma-Ok, dan Urutan Key

```go
// map.go
package main

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

func main() {
	teks := "kucing anjing kucing ikan anjing kucing"
	hitung := make(map[string]int)
	for _, kata := range strings.Fields(teks) {
		hitung[kata]++
	}

	// membaca key yang tidak ada -> zero value, gunakan comma-ok
	v, ok := hitung["burung"]
	fmt.Println(v, ok)

	// urutan iterasi acak, jadi urutkan key dulu
	for _, k := range slices.Sorted(maps.Keys(hitung)) {
		fmt.Println(k, hitung[k])
	}

	var kosong map[string]int
	fmt.Println(kosong["x"], len(kosong)) // membaca map nil aman
	// kosong["x"] = 1                    // panic: assignment to entry in nil map
}
```

### Struct, Tag JSON, dan Embedding

```go
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
```

## Latihan

1. Buat slice `[]int{1,2,3,4,5}`, ambil `t := s[1:4]`, ubah `t[0]`, lalu `append` ke `t`. Jelaskan secara rinci elemen mana dari `s` yang ikut berubah dan mengapa, dengan menyebut `len` dan `cap` masing-masing.
2. Tulis fungsi `HapusDuplikat(s []string) []string` yang mempertahankan urutan kemunculan pertama. Gunakan map sebagai himpunan, lalu bandingkan dengan solusi memakai `slices.Sort` dan `slices.Compact`.
3. Tulis program penghitung frekuensi kata dari sebuah paragraf, lalu cetak 3 kata paling sering muncul. Mengapa Anda tidak bisa langsung mengurutkan sebuah `map`?
4. Buat dua struct `Alamat` dan `Karyawan` (dengan `Alamat` di-*embed*), lalu encode ke JSON memakai tag `json` (nama field huruf kecil, satu field `json:"-"`, satu `omitempty`). Apa yang terjadi pada hasil JSON saat struct di-embed?
5. Jelaskan perbedaan `var s []int`, `s := []int{}`, dan `s := make([]int, 0)` dari sisi `nil`, `len`, `cap`, dan hasil encode JSON. Kapan perbedaan ini penting di API?
6. Buktikan bahwa array dikirim sebagai salinan ke fungsi sedangkan slice tidak, dengan program kecil. Jelaskan mengapa slice yang dikirim ke fungsi tetap bisa mengubah isi slice asal, tetapi `append` di dalam fungsi tidak mengubah `len` slice milik pemanggil.
7. Diskusikan: sebuah fungsi menerima `[]byte` dari buffer yang akan dipakai ulang dan menyimpan irisan `buf[10:20]` ke struct jangka panjang. Risiko apa yang muncul dari berbagi backing array, dan bagaimana `slices.Clone` atau full slice expression membantu?

---

[← Modul 3: Variabel, Konstanta & Tipe Data](03-variabel-konstanta-tipe-data.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 5: Control Flow (Conditionals & Loops) →](05-control-flow.md)
