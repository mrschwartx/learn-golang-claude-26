# Modul 5: Control Flow (Conditionals & Loops)

[← Modul 4: Composite Types (Array, Slice, Map, Struct)](04-composite-types.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 6: Fungsi →](06-fungsi.md)

**Isi modul ini**

- [Teori](#teori)
  - [5.1 `if` dan `if-else`](#51-if-dan-if-else)
  - [5.2 `switch`](#52-switch)
  - [5.3 `for` — Satu-satunya Loop](#53-for--satu-satunya-loop)
  - [5.4 `for range`](#54-for-range)
  - [5.5 Semantik Variabel Loop (Go 1.22)](#55-semantik-variabel-loop-go-122)
  - [5.6 Iterating Maps dan Strings](#56-iterating-maps-dan-strings)
  - [5.7 `break`, `continue`, dan Label](#57-break-continue-dan-label)
  - [5.8 `goto` (Tidak Disarankan)](#58-goto-tidak-disarankan)
- [Praktik](#praktik)
  - [Conditionals dan Switch](#conditionals-dan-switch)
  - [Loop, Range, dan Label](#loop-range-dan-label)
  - [Variabel Loop per Iterasi dan Salinan `range`](#variabel-loop-per-iterasi-dan-salinan-range)
  - [`goto` (hanya untuk Pengenalan)](#goto-hanya-untuk-pengenalan)
- [Latihan](#latihan)

---


## Teori

### 5.1 `if` dan `if-else`

Go tidak memakai tanda kurung di sekitar kondisi, tetapi **kurung kurawal wajib** meski hanya satu baris. Kondisi harus bertipe `bool` — tidak ada konversi implisit dari angka atau pointer.

```go
if umur >= 17 {
    fmt.Println("dewasa")
} else if umur >= 13 {
    fmt.Println("remaja")
} else {
    fmt.Println("anak")
}
```

**Init statement**: `if` boleh diawali pernyataan singkat yang variabelnya **hanya hidup di dalam if/else tersebut**. Ini pola idiomatik untuk menangani error:

```go
if n, err := strconv.Atoi(input); err != nil {
    return err
} else if n < 0 {
    return errors.New("negatif")
} else {
    fmt.Println(n)
}
```

Gaya idiomatik Go adalah **menjaga happy path tetap rata ke kiri**: tangani error lebih dulu dan `return` lebih awal, sehingga tidak perlu `else` setelah blok yang diakhiri `return`.

### 5.2 `switch`

`switch` di Go lebih aman dan fleksibel daripada di C/Java:

| Bentuk | Penjelasan |
|---|---|
| `switch x { case 1: ... }` | Mencocokkan nilai; **tidak ada fall-through otomatis**, jadi `break` tidak perlu |
| `case 1, 2, 3:` | Beberapa nilai dalam satu `case` |
| `switch { case x > 10: ... }` | **Tanpa ekspresi** — setiap `case` adalah kondisi boolean; pengganti rantai `if-else` yang rapi |
| `switch h := f(); h { ... }` | Dengan init statement |
| `fallthrough` | Memaksa lanjut ke `case` berikutnya **tanpa memeriksa kondisinya**; jarang dipakai |
| `default:` | Dijalankan jika tidak ada `case` yang cocok; posisinya bebas |

Nilai `case` tidak harus konstanta dan tipe `switch` bisa apa saja yang comparable. Ada juga **type switch** untuk memeriksa tipe dinamis interface ([Modul 8](08-methods-interfaces.md)).

### 5.3 `for` — Satu-satunya Loop

Go hanya punya keyword `for`, tetapi tiga bentuknya mencakup `for`, `while`, dan `do-while` dari bahasa lain:

| Bentuk | Contoh |
|---|---|
| Tiga komponen (C-style) | `for i := 0; i < 10; i++ { ... }` |
| Hanya kondisi (seperti `while`) | `for sisa > 0 { ... }` |
| Tak hingga | `for { ... }` — keluar dengan `break` atau `return` |
| `range` | `for i, v := range koleksi { ... }` |

### 5.4 `for range`

`range` dapat dipakai pada banyak tipe, dengan nilai iterasi berbeda:

| Tipe | Nilai per iterasi |
|---|---|
| Array / slice | `index, value` (value adalah **salinan** elemen) |
| String | `index byte, rune` — per karakter Unicode |
| Map | `key, value` — **urutan acak** |
| Channel | Satu nilai tiap penerimaan, sampai channel ditutup ([Modul 12](12-goroutines-channels.md)) |
| Integer (Go 1.22+) | `for i := range 5` — `i` dari `0` sampai `4` |
| Function iterator (Go 1.23+) | Nilai yang dihasilkan fungsi iterator ([Modul 9](09-generics.md)) |

Gunakan `_` untuk mengabaikan nilai yang tidak diperlukan; `for range 3 { ... }` (tanpa variabel) juga valid untuk mengulang sejumlah kali.

> Karena `value` adalah **salinan**, mengubah `v` di dalam `for _, v := range users` **tidak mengubah** elemen slice. Ubah lewat index: `users[i].Aktif = true`.

### 5.5 Semantik Variabel Loop (Go 1.22)

Sebelum Go 1.22, satu variabel loop dipakai bersama oleh seluruh iterasi, sehingga closure atau goroutine yang menangkapnya sering menghasilkan nilai yang sama — salah satu bug paling terkenal di Go. **Sejak Go 1.22, setiap iterasi memiliki salinan variabel sendiri**, dan bug itu hilang.

Perilaku baru ini aktif jika `go.mod` menyatakan `go 1.22` atau lebih baru. Proyek lama dengan `go` direktif lebih rendah tetap memakai semantik lama sampai direktifnya dinaikkan.

### 5.6 Iterating Maps dan Strings

**Map**: urutan iterasi sengaja diacak dan berbeda antar-eksekusi. Untuk output yang deterministik, kumpulkan key lalu urutkan:

```go
for _, k := range slices.Sorted(maps.Keys(m)) {
    fmt.Println(k, m[k])
}
```

**String**: dua cara iterasi memberi hasil berbeda.

| Cara | Hasil untuk `"hé"` |
|---|---|
| `for i := 0; i < len(s); i++ { s[i] }` | Per **byte** (3 iterasi, karena `é` = 2 byte) |
| `for i, r := range s` | Per **rune** (2 iterasi); `i` adalah posisi byte awal tiap rune |

Gunakan `range` untuk teks manusia; gunakan loop indeks hanya jika Anda memang memproses byte (protokol atau data biner).

### 5.7 `break`, `continue`, dan Label

- `break` keluar dari `for`, `switch`, atau `select` **terdalam**.
- `continue` melompat ke iterasi berikutnya dari `for`.
- **Label** memungkinkan `break` atau `continue` menargetkan loop **luar**, berguna untuk loop bersarang.

```go
luar:
for i := 0; i < 3; i++ {
    for j := 0; j < 3; j++ {
        if j == 2 { continue luar }
        if i == 2 { break luar }
    }
}
```

Jebakan: `break` di dalam `switch` yang berada di dalam `for` hanya keluar dari `switch`, bukan dari `for`. Gunakan label jika memang ingin keluar dari loop.

### 5.8 `goto` (Tidak Disarankan)

Go mendukung `goto` dengan batasan ketat: tidak boleh **melompati deklarasi variabel** dan tidak boleh **masuk ke dalam blok**. Hampir semua kasus lebih rapi diselesaikan dengan fungsi, `defer`, atau loop berlabel. `goto` kadang muncul di kode hasil generate atau parser, tetapi **hindari di kode aplikasi**.

## Praktik

### Conditionals dan Switch

```go
// kondisi.go
package main

import (
	"fmt"
	"strconv"
)

func nilai(skor int) string {
	switch {
	case skor >= 90:
		return "A"
	case skor >= 80:
		return "B"
	case skor >= 70:
		return "C"
	default:
		return "D"
	}
}

func main() {
	if n, err := strconv.Atoi("42"); err != nil {
		fmt.Println("error:", err)
	} else if n%2 == 0 {
		fmt.Println(n, "genap")
	}

	for _, s := range []int{95, 82, 71, 40} {
		fmt.Println(s, nilai(s))
	}

	switch hari := "sabtu"; hari {
	case "sabtu", "minggu":
		fmt.Println("akhir pekan")
	default:
		fmt.Println("hari kerja")
	}

	switch x := 2; x {
	case 2:
		fmt.Println("dua")
		fallthrough
	case 3:
		fmt.Println("tiga (fallthrough, tanpa cek kondisi)")
	case 4:
		fmt.Println("tidak dijalankan")
	}
}
```

### Loop, Range, dan Label

```go
// loop.go
package main

import "fmt"

func main() {
	// FizzBuzz dengan range integer (Go 1.22+)
	for i := range 15 {
		n := i + 1
		switch {
		case n%15 == 0:
			fmt.Print("FizzBuzz ")
		case n%3 == 0:
			fmt.Print("Fizz ")
		case n%5 == 0:
			fmt.Print("Buzz ")
		default:
			fmt.Print(n, " ")
		}
	}
	fmt.Println()

	// while-style dan loop tak hingga
	sisa := 3
	for sisa > 0 {
		sisa--
	}
	hitung := 0
	for {
		hitung++
		if hitung == 5 {
			break
		}
	}
	fmt.Println(sisa, hitung)

	// string: byte vs rune
	s := "hé"
	for i := 0; i < len(s); i++ {
		fmt.Print(s[i], " ")
	}
	fmt.Println()
	for i, r := range s {
		fmt.Printf("%d:%c ", i, r)
	}
	fmt.Println()

	// label
luar:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if j == 2 {
				continue luar
			}
			if i == 2 {
				break luar
			}
			fmt.Printf("(%d,%d) ", i, j)
		}
	}
	fmt.Println()
}
```

### Variabel Loop per Iterasi dan Salinan `range`

```go
// closure_loop.go
package main

import "fmt"

type User struct {
	Nama  string
	Aktif bool
}

func main() {
	var fungsi []func() int
	for i := range 3 {
		fungsi = append(fungsi, func() int { return i })
	}
	for _, f := range fungsi {
		fmt.Print(f(), " ") // Go 1.22+: 0 1 2
	}
	fmt.Println()

	users := []User{{"Andi", false}, {"Budi", false}}
	for _, u := range users {
		u.Aktif = true // hanya mengubah salinan
	}
	fmt.Println(users)
	for i := range users {
		users[i].Aktif = true
	}
	fmt.Println(users)
}
```

### `goto` (hanya untuk Pengenalan)

```go
// goto.go
package main

import "fmt"

func main() {
	percobaan := 0
ulang:
	percobaan++
	if percobaan < 3 {
		fmt.Println("percobaan", percobaan, "gagal, ulangi")
		goto ulang
	}
	fmt.Println("berhasil pada percobaan", percobaan)
}
```

Bandingkan dengan versi memakai `for` — hasilnya sama tetapi lebih mudah dibaca.

## Latihan

1. Tulis fungsi `Klasifikasi(suhu float64) string` dengan `switch` tanpa ekspresi (dingin < 15, sejuk < 25, hangat < 32, panas selebihnya), lalu uji dengan beberapa nilai.
2. Ubah FizzBuzz di atas menjadi fungsi yang mengembalikan `[]string`, lalu tulis ulang tanpa `switch` memakai rantai `if-else`. Mana yang lebih mudah dibaca, dan mengapa?
3. Cetak tabel perkalian 1–9 dengan dua loop bersarang, lalu hentikan seluruh loop begitu hasil perkalian melewati 50 menggunakan `break` berlabel.
4. Buktikan bahwa mengubah `v` pada `for _, v := range slice` tidak mengubah elemen slice. Tulis dua cara memperbaikinya (lewat index dan lewat slice berisi pointer).
5. Iterasi `"Go 世界 🌍"` dengan loop indeks lalu dengan `range`. Bandingkan jumlah iterasi dan jelaskan perbedaannya dalam kaitan dengan UTF-8.
6. Buat map nilai siswa, lalu cetak isinya **terurut berdasarkan nama** menggunakan `maps.Keys` dan `slices.Sorted`. Jelaskan mengapa Go sengaja mengacak urutan iterasi map.
7. Diskusikan: sebelum Go 1.22 sebuah tim menulis `for _, u := range users { go func() { proses(u) }() }` dan menemukan bahwa semua goroutine memproses user yang sama. Jelaskan penyebabnya, bagaimana Go 1.22 memperbaikinya, dan apa yang harus dicek tim ini di `go.mod` agar perbaikan itu aktif.

---

[← Modul 4: Composite Types (Array, Slice, Map, Struct)](04-composite-types.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 6: Fungsi →](06-fungsi.md)
