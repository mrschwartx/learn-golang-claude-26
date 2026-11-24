# Modul 7: Pointer & Memory Management

[← Modul 6: Fungsi](06-fungsi.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 8: Methods & Interfaces →](08-methods-interfaces.md)

**Isi modul ini**

- [Teori](#teori)
  - [7.1 Apa itu Pointer?](#71-apa-itu-pointer)
  - [7.2 Pointers Basics](#72-pointers-basics)
  - [7.3 Pointers dengan Struct](#73-pointers-dengan-struct)
  - [7.4 Pointer dengan Map dan Slice](#74-pointer-dengan-map-dan-slice)
  - [7.5 Stack vs Heap](#75-stack-vs-heap)
  - [7.6 Garbage Collection](#76-garbage-collection)
  - [7.7 Memory Leak yang Tetap Bisa Terjadi](#77-memory-leak-yang-tetap-bisa-terjadi)
  - [7.8 Kapan Memakai Pointer](#78-kapan-memakai-pointer)
- [Praktik](#praktik)
  - [Dasar Pointer](#dasar-pointer)
  - [Jebakan Slice dan Map](#jebakan-slice-dan-map)
  - [Mengamati Garbage Collector](#mengamati-garbage-collector)
- [Latihan](#latihan)

---


## Teori

### 7.1 Apa itu Pointer?

**Pointer** adalah variabel yang menyimpan **alamat memori** dari variabel lain. Go punya pointer, tetapi **tanpa aritmetika pointer** seperti di C — Anda tidak bisa "menggeser" pointer ke alamat berikutnya, sehingga seluruh kelas bug akses memori sembarangan dihilangkan.

| Operator | Nama | Arti |
|---|---|---|
| `&x` | Address-of | Mengambil alamat variabel `x` |
| `*p` | Dereference | Membaca/menulis nilai yang ditunjuk pointer `p` |
| `*T` | Tipe pointer | "Pointer ke `T`", misalnya `*int` |

Zero value pointer adalah **`nil`**. Mendereferensi pointer `nil` menyebabkan **panic** (*nil pointer dereference*), salah satu penyebab crash paling umum di Go.

### 7.2 Pointers Basics

```go
x := 10
p := &x        // p bertipe *int, menyimpan alamat x
fmt.Println(*p) // 10
*p = 20         // mengubah x lewat pointer
fmt.Println(x)  // 20
```

Fungsi **`new(T)`** mengalokasikan variabel `T` dengan zero value dan mengembalikan pointer-nya. Mulai **Go 1.26**, `new` juga menerima **ekspresi** sebagai nilai awal: `new(42)` menghasilkan `*int` yang menunjuk ke `42`. Ini sangat berguna untuk field opsional bertipe pointer (misalnya pada struct JSON dan protobuf), yang sebelumnya membutuhkan variabel perantara.

Alasan utama memakai pointer:

1. **Mengubah data milik pemanggil** (karena Go selalu pass-by-value, [Modul 6](06-fungsi.md)).
2. **Menghindari penyalinan** struct besar.
3. **Menyatakan "tidak ada nilai"** dengan `nil` (misalnya field opsional `*int`).

### 7.3 Pointers dengan Struct

Go memudahkan penggunaan pointer ke struct: akses field lewat pointer **tidak perlu dereference eksplisit** — `p.Nama` otomatis berarti `(*p).Nama`.

```go
type User struct{ Nama string; Umur int }

func ulangTahun(u *User) { u.Umur++ }

u := User{"Andi", 30}
ulangTahun(&u)
```

Mengembalikan pointer ke variabel lokal **aman** di Go: compiler memindahkan variabel tersebut ke heap bila perlu (lihat 7.6). Konstruktor idiomatik biasanya `func NewUser(nama string) *User { return &User{Nama: nama} }`.

### 7.4 Pointer dengan Map dan Slice

Slice, map, dan channel sudah berisi **referensi internal**, sehingga mengirimnya ke fungsi tidak menyalin datanya. Tetap ada jebakan:

| Situasi | Perilaku |
|---|---|
| Ubah elemen slice di fungsi | Terlihat oleh pemanggil (backing array sama) |
| `append` di fungsi | `len` slice pemanggil **tidak berubah**; kembalikan slice baru atau kirim `*[]T` |
| Ubah isi map di fungsi | Terlihat oleh pemanggil |
| `m[k].Field = x` pada `map[string]T` | **Error kompilasi** — elemen map **tidak addressable** |

Solusi untuk `m[k].Field`: ambil struct ke variabel, ubah, lalu tulis kembali (`v := m[k]; v.Field = x; m[k] = v`), atau simpan pointer sebagai value (`map[string]*T`).

> Pointer ke slice (`*[]int`) dan pointer ke map jarang dibutuhkan. Lebih idiomatik mengembalikan slice baru: `s = tambah(s, x)`.

### 7.5 Stack vs Heap

Go mengalokasikan variabel di **stack** (cepat, dibersihkan otomatis saat fungsi selesai) atau **heap** (butuh garbage collector). Anda **tidak memilih secara langsung** — compiler memutuskannya lewat **escape analysis** ([Modul 20](20-advanced-topics.md)): jika sebuah variabel masih dirujuk setelah fungsi selesai (misalnya pointer-nya dikembalikan), ia "lolos" (*escape*) ke heap.

Konsekuensi praktis: alokasi heap lebih mahal karena menambah beban GC, tetapi jangan mengorbankan keterbacaan demi menghindarinya sebelum ada bukti dari profiling ([Modul 19](19-performance-deployment.md)).

### 7.6 Garbage Collection

Go menyediakan **garbage collector (GC)** otomatis. Karakteristiknya:

| Aspek | Penjelasan |
|---|---|
| **Concurrent** | GC berjalan bersamaan dengan program, dengan jeda *stop-the-world* sangat singkat |
| **Tri-color mark & sweep** | Menandai objek yang masih terjangkau dari root, lalu membersihkan sisanya |
| **Non-generational, non-compacting** | Semua objek diperlakukan setara; memori tidak dipadatkan |
| **Green Tea GC** | Diperkenalkan sebagai eksperimen di Go 1.25 dan menjadi **default di Go 1.26**; meningkatkan efisiensi penandaan objek kecil lewat locality yang lebih baik |

Dua pengaturan utama untuk mengatur perilaku GC:

| Pengaturan | Fungsi |
|---|---|
| **`GOGC`** | Persentase pertumbuhan heap sebelum GC berikutnya; default `100` (GC dipicu saat heap dua kali ukuran data hidup). Nilai lebih tinggi = GC lebih jarang, memori lebih besar |
| **`GOMEMLIMIT`** | Batas memori **lunak** total (Go 1.19+), misalnya `GOMEMLIMIT=450MiB`; GC lebih agresif saat mendekati batas — berguna di container |

Untuk mengamati GC: `GODEBUG=gctrace=1 ./app` mencetak satu baris per siklus GC, dan `runtime.ReadMemStats` memberi statistik dari dalam program. Detail internal GC dibahas di [Modul 20](20-advanced-topics.md).

### 7.7 Memory Leak yang Tetap Bisa Terjadi

GC tidak menjamin bebas kebocoran memori — GC hanya membersihkan objek yang **tidak lagi terjangkau**. Penyebab umum "leak" di Go:

| Penyebab | Penjelasan |
|---|---|
| **Goroutine leak** | Goroutine menunggu channel yang tidak akan pernah terisi atau tertutup, sehingga memorinya tidak pernah bebas |
| **Slice yang menahan array besar** | `kecil := besar[:10]` membuat seluruh backing array tetap hidup; salin dengan `slices.Clone` jika hanya butuh sebagian kecil |
| **Cache atau map global tanpa batas** | Data terus ditambah tanpa eviction |
| **Timer/Ticker tidak dihentikan** | Panggil `Stop()` (sejak Go 1.23 timer yang tidak terreferensi bisa di-GC, tetapi menghentikan secara eksplisit tetap praktik baik) |
| **`defer` dalam loop** | Menahan resource sampai fungsi selesai ([Modul 6](06-fungsi.md)) |

Go 1.27 menjadikan profil **`goroutineleak`** pada `runtime/pprof` tersedia secara umum untuk mendeteksi goroutine yang terblokir permanen ([Modul 19](19-performance-deployment.md)).

### 7.8 Kapan Memakai Pointer

| Gunakan pointer jika... | Gunakan nilai (value) jika... |
|---|---|
| Fungsi/method harus **mengubah** data | Tipe kecil dan imutabel (angka, string, `time.Time`) |
| Struct **besar** dan sering dikirim | Struct kecil yang jarang diubah |
| Perlu merepresentasikan **"tidak ada"** (`nil`) | Nilai selalu ada |
| Struct memuat **mutex** atau state yang tak boleh disalin | Ingin semantik imutabel dan aman antar-goroutine |

Prinsip: jangan memakai pointer **hanya** demi performa tanpa pengukuran. Salinan struct kecil sering lebih cepat daripada alokasi heap dan beban GC. Hindari juga *pointer ke interface* (`*io.Reader`) — hampir selalu keliru.

## Praktik

### Dasar Pointer

```go
// pointer.go
package main

import "fmt"

type User struct {
	Nama string
	Umur int
}

func tukar(a, b *int) { *a, *b = *b, *a }

func ulangTahun(u *User) { u.Umur++ }

func main() {
	x := 10
	p := &x
	*p = 20
	fmt.Println(x, *p)

	a, b := 1, 2
	tukar(&a, &b)
	fmt.Println(a, b)

	var kosong *int
	fmt.Println(kosong == nil)
	// fmt.Println(*kosong) // panic: nil pointer dereference

	n := new(int)
	*n = 7
	fmt.Println(*n)

	// Go 1.26+: new menerima ekspresi
	umur := new(30)
	fmt.Println(*umur)

	u := User{"Andi", 30}
	ulangTahun(&u)
	fmt.Println(u)
}
```

### Jebakan Slice dan Map

```go
// jebakan.go
package main

import "fmt"

type Item struct{ Jumlah int }

func tambahSalah(s []int)       { s = append(s, 99) }
func tambahBenar(s []int) []int { return append(s, 99) }
func tambahPtr(s *[]int)        { *s = append(*s, 99) }

func main() {
	s := make([]int, 1, 10)
	tambahSalah(s)
	fmt.Println("salah:", s)
	s = tambahBenar(s)
	fmt.Println("benar:", s)
	tambahPtr(&s)
	fmt.Println("pointer:", s)

	m := map[string]Item{"a": {1}}
	// m["a"].Jumlah = 5  // error: cannot assign to struct field in map
	v := m["a"]
	v.Jumlah = 5
	m["a"] = v
	fmt.Println(m)

	mp := map[string]*Item{"b": {1}}
	mp["b"].Jumlah = 9 // boleh: pointer dapat dideferensikan
	fmt.Println(mp["b"].Jumlah)
}
```

### Mengamati Garbage Collector

```go
// gc.go
package main

import (
	"fmt"
	"runtime"
)

func mb(b uint64) uint64 { return b / 1024 / 1024 }

func main() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	awal := ms.NumGC

	var simpan [][]byte
	for i := 0; i < 50; i++ {
		simpan = append(simpan, make([]byte, 1<<20)) // 1 MiB
	}
	runtime.ReadMemStats(&ms)
	fmt.Println("heap saat data masih dipegang (MiB) >= 50:", mb(ms.HeapAlloc) >= 50)

	simpan = nil
	runtime.GC()
	runtime.ReadMemStats(&ms)
	fmt.Println("heap setelah GC (MiB) < 10:", mb(ms.HeapAlloc) < 10)
	fmt.Println("GC berjalan:", ms.NumGC > awal)
}
```

Jalankan dengan trace GC aktif untuk melihat setiap siklus:

```bash
GODEBUG=gctrace=1 go run gc.go
GOGC=50 go run gc.go          # GC lebih sering, heap lebih kecil
GOMEMLIMIT=64MiB go run gc.go # batas lunak memori total
```

## Latihan

1. Tulis fungsi `Reset(p *int)` yang mengatur nilai menjadi `0`, lalu tulis versi yang menerima `int` biasa. Tunjukkan hanya versi pointer yang mengubah variabel pemanggil, dan jelaskan alasannya dengan istilah *pass-by-value*.
2. Buat struct `Rekening` dengan field `Saldo` dan fungsi `Setor(r *Rekening, jumlah int)` serta `Cek(r Rekening) int`. Jelaskan mengapa `Setor` membutuhkan pointer sedangkan `Cek` tidak.
3. Tulis `NewUser(nama string) *User` yang mengembalikan pointer ke variabel lokal. Mengapa ini aman di Go tetapi berbahaya di C? Sebutkan peran escape analysis.
4. Buat `map[string]Produk` dan coba mengubah `m["kopi"].Stok -= 1`. Catat pesan error compiler, lalu perbaiki dengan dua cara (salin-tulis ulang, dan `map[string]*Produk`). Kapan Anda memilih masing-masing?
5. Tulis program yang menyimpan slice kecil `kecil := besar[:5]` dari slice 100 MB, lalu tunjukkan dengan `runtime.ReadMemStats` bahwa memori belum dibebaskan. Perbaiki dengan `slices.Clone` dan buktikan memori bisa dibebaskan setelah `runtime.GC()`.
6. Jalankan program yang menghasilkan banyak alokasi dengan `GOGC=25`, `GOGC=100`, dan `GOGC=400` bersama `GODEBUG=gctrace=1`. Hitung jumlah siklus GC dan puncak heap pada masing-masing, lalu simpulkan trade-off CPU vs memori.
7. Diskusikan: sebuah layanan Go di Kubernetes dengan limit memori 512 MiB sering di-*OOMKill* walau GC berjalan normal. Jelaskan peran `GOMEMLIMIT`, bagaimana menyetelnya relatif terhadap limit container, dan mengapa batas ini bersifat *lunak*.

---

[← Modul 6: Fungsi](06-fungsi.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 8: Methods & Interfaces →](08-methods-interfaces.md)
