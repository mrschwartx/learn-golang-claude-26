# Modul 6: Fungsi

[← Modul 5: Control Flow (Conditionals & Loops)](05-control-flow.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 7: Pointer & Memory Management →](07-pointer-memory.md)

**Isi modul ini**

- [Teori](#teori)
  - [6.1 Dasar Fungsi](#61-dasar-fungsi)
  - [6.2 Multiple Return Values](#62-multiple-return-values)
  - [6.3 Named Return Values](#63-named-return-values)
  - [6.4 Variadic Functions](#64-variadic-functions)
  - [6.5 Anonymous Functions dan Closures](#65-anonymous-functions-dan-closures)
  - [6.6 Call by Value](#66-call-by-value)
  - [6.7 `defer`](#67-defer)
  - [6.8 Function Types dan Higher-Order Functions](#68-function-types-dan-higher-order-functions)
  - [6.9 Rekursi](#69-rekursi)
- [Praktik](#praktik)
  - [Multiple Return, Named Return, dan Variadic](#multiple-return-named-return-dan-variadic)
  - [Closure dan Fungsi sebagai Nilai](#closure-dan-fungsi-sebagai-nilai)
  - [`defer` dan Call by Value](#defer-dan-call-by-value)
- [Latihan](#latihan)

---


## Teori

### 6.1 Dasar Fungsi

**Fungsi** dideklarasikan dengan keyword `func`, diikuti nama, parameter, dan tipe kembalian. Tipe ditulis **setelah** nama, dan parameter berurutan dengan tipe sama boleh digabung.

```go
func tambah(a, b int) int {
    return a + b
}
```

Ciri fungsi di Go yang berbeda dari banyak bahasa lain:

| Ciri | Penjelasan |
|---|---|
| **Tidak ada overloading** | Dua fungsi dalam satu paket tidak boleh bernama sama, meski parameternya berbeda |
| **Tidak ada default argument** | Semua parameter wajib diisi; pola alternatifnya adalah *functional options* atau struct konfigurasi |
| **First-class** | Fungsi adalah nilai: bisa disimpan di variabel, dikirim sebagai argumen, dan dikembalikan dari fungsi lain |
| **Tipe fungsi** | `type Operasi func(int, int) int` — bisa diberi nama dan method |
| **`main` dan `init`** | `main()` adalah titik masuk program; `init()` dijalankan otomatis sebelum `main` (tiap file boleh punya beberapa) |

### 6.2 Multiple Return Values

Fungsi Go dapat mengembalikan **lebih dari satu nilai**. Pola paling dikenal adalah **`(hasil, error)`** — pengganti exception di bahasa lain ([Modul 10](10-error-handling.md)).

```go
func bagi(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("pembagi nol")
    }
    return a / b, nil
}

hasil, err := bagi(10, 2)
```

Nilai yang tidak dibutuhkan diabaikan dengan **blank identifier**: `hasil, _ := bagi(10, 2)`. Mengabaikan `error` dengan cara ini harus benar-benar disengaja.

### 6.3 Named Return Values

Nilai kembalian boleh diberi nama. Variabel itu otomatis dideklarasikan dengan **zero value**, dan `return` tanpa argumen (*naked return*) mengembalikan nilai terakhirnya.

```go
func bagiAman(a, b int) (hasil int, err error) {
    if b == 0 {
        err = errors.New("pembagi nol")
        return
    }
    hasil = a / b
    return
}
```

| Penggunaan | Penilaian |
|---|---|
| Menamai hasil sebagai **dokumentasi** (`(lebar, tinggi int)`) | Baik |
| Memodifikasi hasil di dalam **`defer`** (misalnya membungkus error) | Sah, dan kadang satu-satunya cara |
| *Naked return* di fungsi panjang | **Hindari** — pembaca sulit melacak nilai yang dikembalikan |

### 6.4 Variadic Functions

Parameter terakhir boleh bertipe `...T` sehingga menerima **nol atau lebih argumen**. Di dalam fungsi, ia berperilaku sebagai **slice** `[]T`.

```go
func jumlah(angka ...int) int {
    total := 0
    for _, n := range angka {
        total += n
    }
    return total
}

jumlah(1, 2, 3)
daftar := []int{4, 5, 6}
jumlah(daftar...)   // meneruskan slice dengan ...
```

Jika dipanggil tanpa argumen, `angka` bernilai `nil` slice. Contoh variadic terkenal: `fmt.Println(a ...any)` dan `append(s, elemen...)`.

### 6.5 Anonymous Functions dan Closures

**Fungsi anonim** adalah fungsi tanpa nama yang dapat langsung dipanggil atau disimpan di variabel. Jika fungsi tersebut merujuk variabel dari scope luarnya, ia menjadi **closure** — variabel yang ditangkap **tetap hidup** selama closure masih dipakai, dan ditangkap **secara referensi**, bukan salinan.

```go
func penghitung() func() int {
    n := 0
    return func() int {
        n++
        return n
    }
}

c := penghitung()
c() // 1
c() // 2
```

Setiap panggilan `penghitung()` menghasilkan closure dengan `n` miliknya sendiri. Pola ini dasar bagi *middleware* HTTP, *functional options*, dan *callback*. Soal variabel loop yang ditangkap closure, lihat [Modul 5](05-control-flow.md) (semantik per-iterasi sejak Go 1.22).

### 6.6 Call by Value

**Go selalu memakai pass-by-value**: setiap argumen **disalin** ke parameter fungsi. Namun yang disalin berbeda-beda tergantung tipenya:

| Tipe | Yang disalin | Akibatnya |
|---|---|---|
| `int`, `string`, `bool`, array, struct | Seluruh nilai | Perubahan di dalam fungsi **tidak** memengaruhi pemanggil |
| **Pointer** | Alamat (pointer-nya) | Fungsi bisa mengubah data yang ditunjuk ([Modul 7](07-pointer-memory.md)) |
| **Slice** | Header (pointer, len, cap) | Mengubah elemen **terlihat** di pemanggil; `append` yang mengubah `len` **tidak** |
| **Map**, **channel** | Referensi ke struktur internal | Perubahan isi **terlihat** di pemanggil |

Struct besar yang disalin pada setiap pemanggilan bisa mahal; untuk kasus tersebut gunakan pointer ([Modul 7](07-pointer-memory.md)).

### 6.7 `defer`

`defer` menjadwalkan sebuah pemanggilan fungsi untuk dijalankan **tepat sebelum fungsi pembungkusnya selesai** (baik lewat `return` biasa maupun `panic`).

| Aturan | Penjelasan |
|---|---|
| **LIFO** | Beberapa `defer` dijalankan dari yang terakhir didaftarkan ke yang pertama |
| **Argumen dievaluasi saat `defer` dideklarasikan** | Bukan saat dijalankan |
| Dipakai untuk **pembersihan** | `defer file.Close()`, `defer mu.Unlock()`, `defer cancel()` |
| Bisa memodifikasi **named return** | Berguna untuk membungkus error |

Jebakan: **`defer` di dalam loop** baru berjalan saat fungsi selesai, bukan di akhir iterasi — memproses ribuan file dalam satu loop dengan `defer f.Close()` menahan semua file tetap terbuka. Pindahkan isi loop ke fungsi tersendiri.

### 6.8 Function Types dan Higher-Order Functions

Karena fungsi adalah nilai, kita bisa menulis fungsi yang **menerima atau mengembalikan fungsi** (*higher-order function*):

```go
type Filter func(int) bool

func Saring(data []int, f Filter) []int { /* ... */ }
```

Contoh dalam standard library: `sort.Slice(s, func(i, j int) bool {...})`, `slices.SortFunc`, dan `http.HandlerFunc`. Pola *middleware* (`func(http.Handler) http.Handler`, [Modul 16](16-cli-web-development.md)) dan *functional options* untuk konfigurasi juga dibangun di atas konsep ini.

### 6.9 Rekursi

Fungsi boleh memanggil dirinya sendiri. Go **tidak melakukan tail-call optimization**, tetapi stack goroutine **tumbuh dinamis** (batas default 1 GB di 64-bit), sehingga rekursi dalam lebih aman dibanding bahasa dengan stack tetap. Tetap pertimbangkan versi iteratif untuk kedalaman yang tidak terbatas.

## Praktik

### Multiple Return, Named Return, dan Variadic

```go
// fungsi.go
package main

import (
	"errors"
	"fmt"
)

func bagi(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("pembagi nol")
	}
	return a / b, nil
}

func ukuran(luas, sisi float64) (panjang, lebar float64) {
	panjang = luas / sisi
	lebar = sisi
	return
}

func jumlah(angka ...int) int {
	total := 0
	for _, n := range angka {
		total += n
	}
	return total
}

func main() {
	if h, err := bagi(10, 4); err == nil {
		fmt.Println("hasil:", h)
	}
	if _, err := bagi(1, 0); err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println(ukuran(20, 4))
	fmt.Println(jumlah(), jumlah(1, 2, 3), jumlah([]int{4, 5, 6}...))
}
```

### Closure dan Fungsi sebagai Nilai

```go
// closure.go
package main

import (
	"fmt"
	"sort"
)

func penghitung() func() int {
	n := 0
	return func() int {
		n++
		return n
	}
}

type Operasi func(int, int) int

func terapkan(op Operasi, a, b int) int { return op(a, b) }

func main() {
	a, b := penghitung(), penghitung()
	fmt.Println(a(), a(), a(), b())

	fmt.Println(terapkan(func(x, y int) int { return x * y }, 6, 7))

	hasil := func(x int) int { return x * x }(9) // langsung dipanggil
	fmt.Println(hasil)

	kata := []string{"pisang", "apel", "jeruk"}
	sort.Slice(kata, func(i, j int) bool { return len(kata[i]) < len(kata[j]) })
	fmt.Println(kata)
}
```

### `defer` dan Call by Value

```go
// defer.go
package main

import (
	"errors"
	"fmt"
)

type Titik struct{ X int }

func ubahNilai(t Titik)    { t.X = 100 }
func ubahPointer(t *Titik) { t.X = 100 }
func ubahIsi(s []int)      { s[0] = 100 }
func tambahSlice(s []int)  { s = append(s, 999) }

func bungkus() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("bungkus: %w", err)
		}
	}()
	return errors.New("gagal asli")
}

func main() {
	for i := 0; i < 3; i++ {
		defer fmt.Println("defer ke-", i)
	}

	x := 1
	defer fmt.Println("x saat defer didaftarkan:", x)
	x = 2

	t := Titik{1}
	ubahNilai(t)
	fmt.Println("setelah ubahNilai:", t)
	ubahPointer(&t)
	fmt.Println("setelah ubahPointer:", t)

	s := []int{1, 2, 3}
	ubahIsi(s)
	tambahSlice(s)
	fmt.Println("slice:", s, "len:", len(s))

	fmt.Println(bungkus())
}
```

Perhatikan urutan output `defer`: LIFO, dan nilai `x` yang dicetak adalah `1`, bukan `2`.

## Latihan

1. Tulis fungsi `MinMax(angka ...int) (min, max int, err error)` yang mengembalikan error jika tidak ada argumen. Panggil dengan argumen biasa dan dengan slice (`s...`).
2. Tulis fungsi `Penjumlah(awal int) func(int) int` yang mengembalikan closure yang menambah nilai akumulasi setiap dipanggil. Buat dua penjumlah dan buktikan keduanya memiliki state terpisah.
3. Implementasikan `Map`, `Filter`, dan `Reduce` untuk `[]int` memakai parameter bertipe fungsi, lalu gabungkan: kuadratkan semua angka, ambil yang genap, lalu jumlahkan.
4. Jelaskan output program berikut tanpa menjalankannya terlebih dulu, kemudian verifikasi: sebuah fungsi mendaftarkan `defer` dalam loop `for i := range 3` yang mencetak `i`, dan satu `defer` lain yang mencetak variabel yang berubah setelah `defer` dideklarasikan.
5. Buktikan dengan program bahwa Go adalah *pass-by-value*: tulis fungsi yang menerima `struct`, `*struct`, `[]int`, dan `map[string]int` lalu ubah masing-masing. Susun tabel hasil mana yang terlihat oleh pemanggil, dan jelaskan mengapa `append` di dalam fungsi tidak mengubah panjang slice milik pemanggil.
6. Tulis fungsi `BacaSemua(paths []string)` yang membuka setiap file. Mula-mula pakai `defer f.Close()` langsung di dalam loop, lalu perbaiki dengan memindahkan isi loop ke fungsi bantu. Jelaskan masalah yang muncul pada versi pertama bila `paths` berisi ribuan file.
7. Diskusikan: Go tidak mendukung default argument maupun overloading. Bandingkan dua pola pengganti — **struct konfigurasi** dan **functional options** (`func(*Config)`) — lalu tentukan kapan Anda memilih masing-masing.

---

[← Modul 5: Control Flow (Conditionals & Loops)](05-control-flow.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 7: Pointer & Memory Management →](07-pointer-memory.md)
