# Modul 8: Methods & Interfaces

[← Modul 7: Pointer & Memory Management](07-pointer-memory.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 9: Generics & Iterators →](09-generics.md)

**Isi modul ini**

- [Teori](#teori)
  - [8.1 Methods vs Functions](#81-methods-vs-functions)
  - [8.2 Pointer Receivers vs Value Receivers](#82-pointer-receivers-vs-value-receivers)
  - [8.3 Interfaces](#83-interfaces)
  - [8.4 Interface Internals dan Jebakan `nil`](#84-interface-internals-dan-jebakan-nil)
  - [8.5 Empty Interface dan `any`](#85-empty-interface-dan-any)
  - [8.6 Embedding Interfaces](#86-embedding-interfaces)
  - [8.7 Type Assertions](#87-type-assertions)
  - [8.8 Type Switch](#88-type-switch)
  - [8.9 Pengecekan Interface saat Kompilasi](#89-pengecekan-interface-saat-kompilasi)
- [Praktik](#praktik)
  - [Interface `Bentuk` dan Type Switch](#interface-bentuk-dan-type-switch)
  - [Pointer Receiver, Stringer, dan Method Set](#pointer-receiver-stringer-dan-method-set)
  - [Jebakan Interface `nil`](#jebakan-interface-nil)
  - [Embedding Interface: Decorator `io.Writer`](#embedding-interface-decorator-iowriter)
- [Latihan](#latihan)

---


## Teori

### 8.1 Methods vs Functions

**Method** adalah fungsi yang memiliki **receiver** — parameter khusus yang ditulis di antara `func` dan nama method. Method dapat didefinisikan pada **tipe apa pun yang dideklarasikan di paket yang sama** (struct, `type Celsius float64`, slice bernama, dan sebagainya), tetapi tidak pada tipe bawaan seperti `int` secara langsung.

```go
type Persegi struct{ Sisi float64 }

// fungsi biasa
func Luas(p Persegi) float64 { return p.Sisi * p.Sisi }

// method: receiver (p Persegi)
func (p Persegi) Luas() float64 { return p.Sisi * p.Sisi }

p := Persegi{4}
p.Luas()
```

Method bukan sihir — `p.Luas()` hanyalah gula sintaks untuk `Persegi.Luas(p)`. Keuntungan method: nama dikelompokkan bersama tipenya, dan yang lebih penting, method menjadikan sebuah tipe dapat **memenuhi interface** (8.3).

### 8.2 Pointer Receivers vs Value Receivers

| | **Value receiver** `(p Persegi)` | **Pointer receiver** `(p *Persegi)` |
|---|---|---|
| Receiver di dalam method | **Salinan** | **Pointer** ke nilai asli |
| Dapat mengubah state asli | Tidak | **Ya** |
| Biaya pemanggilan | Menyalin struct | Menyalin pointer |
| Boleh dipanggil pada `nil` | Tidak (nilai harus ada) | Ya, jika method menanganinya |

**Method set** menentukan interface apa yang dipenuhi sebuah tipe:

| Tipe | Method yang termasuk method set |
|---|---|
| `T` | Hanya method ber-**value receiver** |
| `*T` | Method ber-**value receiver** dan **pointer receiver** |

Karena itu, jika `Tulis()` didefinisikan dengan pointer receiver pada `Berkas`, maka **hanya `*Berkas`** yang memenuhi interface `Penulis` — mengirim `Berkas{}` (nilai) ke fungsi yang membutuhkan `Penulis` adalah error kompilasi.

Panduan memilih:

1. Jika method **harus mengubah** receiver → pointer receiver.
2. Jika struct **besar** atau mengandung `sync.Mutex` (tidak boleh disalin) → pointer receiver.
3. **Konsisten**: jika satu method memakai pointer receiver, gunakan pointer receiver untuk semua method tipe itu.
4. Tipe kecil dan imutabel (`time.Time`, tipe numerik bernama) → value receiver.

Go memudahkan: memanggil method ber-pointer receiver pada **variabel yang addressable** (`u.Tambah()`) otomatis menjadi `(&u).Tambah()`. Hal ini tidak berlaku pada nilai yang tidak addressable seperti elemen map atau hasil kembalian fungsi.

### 8.3 Interfaces

**Interface** mendefinisikan **sekumpulan method**. Tipe apa pun yang memiliki method tersebut **otomatis memenuhi interface** — tanpa keyword `implements`. Ini disebut **implicit satisfaction** (structural typing), dan menjadi fondasi desain Go.

```go
type Bentuk interface {
    Luas() float64
    Keliling() float64
}
```

Interface kecil bernilai besar di standard library:

| Interface | Method | Dipakai oleh |
|---|---|---|
| `error` | `Error() string` | Semua error ([Modul 10](10-error-handling.md)) |
| `fmt.Stringer` | `String() string` | Kustomisasi hasil `fmt.Println` |
| `io.Reader` | `Read([]byte) (int, error)` | File, koneksi jaringan, body HTTP, buffer |
| `io.Writer` | `Write([]byte) (int, error)` | File, response HTTP, `os.Stdout`, logger |
| `sort.Interface` | `Len`, `Less`, `Swap` | Pengurutan kustom |
| `http.Handler` | `ServeHTTP(w, r)` | Handler HTTP ([Modul 16](16-cli-web-development.md)) |

Prinsip desain yang idiomatik:

- **Interface kecil** (satu sampai tiga method) lebih mudah dipenuhi dan diuji.
- **"Accept interfaces, return structs"** — fungsi menerima interface sebagai parameter agar fleksibel, tetapi mengembalikan tipe konkret.
- **Definisikan interface di sisi pemakai (consumer)**, bukan di sisi penyedia — paket yang membutuhkan perilaku tertentu mendeklarasikan interface minimalnya sendiri. Ini membuat mocking mudah ([Modul 15](15-testing-benchmarking.md)).

### 8.4 Interface Internals dan Jebakan `nil`

Secara internal, nilai interface berisi **dua komponen**: **(tipe konkret, nilai)**. Interface bernilai `nil` **hanya jika kedua komponen kosong**.

Jebakan terkenal: menempatkan **pointer `nil` bertipe konkret** ke dalam interface menghasilkan interface yang **tidak `nil`**, karena komponen tipenya terisi.

```go
func cari() error {
    var e *ErrKustom = nil
    return e          // interface error berisi (*ErrKustom, nil) — BUKAN nil
}
err := cari()
err != nil            // true!
```

Aturan aman: fungsi yang mengembalikan `error` harus mengembalikan literal **`nil`** secara eksplisit saat sukses, bukan variabel pointer bertipe konkret.

### 8.5 Empty Interface dan `any`

Interface tanpa method, `interface{}` (diberi alias **`any`** sejak Go 1.18), dipenuhi oleh **semua tipe**. Ia dipakai oleh fungsi yang menerima nilai apa saja, seperti `fmt.Println(a ...any)` dan `json.Unmarshal` ke `map[string]any`.

Konsekuensinya: **informasi tipe hilang saat kompilasi** sehingga kesalahan baru muncul saat runtime. Gunakan `any` seperlunya — sejak ada generics ([Modul 9](09-generics.md)) banyak kebutuhan yang dulu memakai `interface{}` kini bisa diselesaikan tanpa kehilangan keamanan tipe.

### 8.6 Embedding Interfaces

Interface dapat **menyusun interface lain** lewat embedding:

```go
type ReadWriter interface {
    Reader
    Writer
}
```

Embedding interface **di dalam struct** juga umum untuk membuat **decorator**: struct membungkus sebuah implementasi, meneruskan semua method otomatis, dan hanya menimpa method yang ingin diubah. Pola ini dipakai pada *logging wrapper*, *middleware*, dan *test double* ([Modul 15](15-testing-benchmarking.md)).

### 8.7 Type Assertions

**Type assertion** mengambil nilai konkret dari dalam interface:

| Bentuk | Perilaku |
|---|---|
| `v := i.(T)` | Mengembalikan nilai bertipe `T`; **panic** jika `i` bukan `T` |
| `v, ok := i.(T)` | **Aman** — `ok` bernilai `false` (tanpa panic) jika tipe tidak cocok |
| `i.(OtherInterface)` | Memeriksa apakah nilai juga memenuhi interface lain |

Gunakan selalu bentuk **comma-ok** kecuali Anda yakin sepenuhnya.

### 8.8 Type Switch

**Type switch** adalah `switch` yang mencocokkan **tipe dinamis** sebuah interface:

```go
switch v := x.(type) {
case int:
    fmt.Println("int", v+1)
case string:
    fmt.Println("string", strings.ToUpper(v))
case nil:
    fmt.Println("nil")
default:
    fmt.Printf("tipe lain: %T\n", v)
}
```

Di dalam tiap `case`, `v` otomatis bertipe sesuai `case` tersebut (pada `case` dengan beberapa tipe atau `default`, `v` tetap bertipe interface asal). Pola ini umum pada parser, penanganan hasil `any`, dan pemeriksaan jenis error.

### 8.9 Pengecekan Interface saat Kompilasi

Karena pemenuhan interface bersifat implisit, kesalahan (misalnya lupa satu method) baru ketahuan saat tipe dipakai. Anda bisa **memaksa pengecekan di waktu kompilasi** dengan assertion ke blank identifier:

```go
var _ Bentuk = (*Lingkaran)(nil)
```

Baris ini tidak melakukan apa pun saat runtime, tetapi gagal dikompilasi jika `*Lingkaran` tidak memenuhi `Bentuk`.

## Praktik

### Interface `Bentuk` dan Type Switch

```go
// bentuk.go
package main

import (
	"fmt"
	"math"
)

type Bentuk interface {
	Luas() float64
	Keliling() float64
}

type Persegi struct{ Sisi float64 }
type Lingkaran struct{ R float64 }

func (p Persegi) Luas() float64       { return p.Sisi * p.Sisi }
func (p Persegi) Keliling() float64   { return 4 * p.Sisi }
func (l Lingkaran) Luas() float64     { return math.Pi * l.R * l.R }
func (l Lingkaran) Keliling() float64 { return 2 * math.Pi * l.R }

var _ Bentuk = Persegi{}
var _ Bentuk = (*Lingkaran)(nil) // *T juga memenuhi interface bila T memenuhinya

func totalLuas(daftar []Bentuk) float64 {
	total := 0.0
	for _, b := range daftar {
		total += b.Luas()
	}
	return total
}

func deskripsi(b Bentuk) string {
	switch v := b.(type) {
	case Persegi:
		return fmt.Sprintf("persegi sisi %.1f", v.Sisi)
	case Lingkaran:
		return fmt.Sprintf("lingkaran jari-jari %.1f", v.R)
	default:
		return "bentuk tidak dikenal"
	}
}

func main() {
	daftar := []Bentuk{Persegi{2}, Lingkaran{1}}
	fmt.Printf("total luas: %.2f\n", totalLuas(daftar))
	for _, b := range daftar {
		fmt.Println(deskripsi(b))
	}

	var x any = "go"
	if s, ok := x.(string); ok {
		fmt.Println("string:", s)
	}
	if _, ok := x.(int); !ok {
		fmt.Println("bukan int")
	}
}
```

### Pointer Receiver, Stringer, dan Method Set

```go
// receiver.go
package main

import "fmt"

type Penghitung struct{ n int }

func (p *Penghitung) Tambah()       { p.n++ }
func (p Penghitung) Nilai() int     { return p.n }
func (p Penghitung) String() string { return fmt.Sprintf("Penghitung(%d)", p.n) }

type Penambah interface{ Tambah() }

func main() {
	var c Penghitung
	c.Tambah() // otomatis (&c).Tambah()
	c.Tambah()
	fmt.Println(c.Nilai(), c) // String() dipakai oleh fmt

	var t Penambah = &c // harus pointer: Tambah ada di method set *Penghitung
	t.Tambah()
	fmt.Println(c)

	// var t2 Penambah = c // error: Penghitung does not implement Penambah (method Tambah has pointer receiver)
}
```

### Jebakan Interface `nil`

```go
// nilgotcha.go
package main

import "fmt"

type ErrKustom struct{}

func (*ErrKustom) Error() string { return "error kustom" }

func salah() error {
	var e *ErrKustom // nil
	return e
}

func benar() error {
	return nil
}

func main() {
	fmt.Println("salah() == nil:", salah() == nil)
	fmt.Println("benar() == nil:", benar() == nil)
}
```

### Embedding Interface: Decorator `io.Writer`

```go
// hitung_byte.go
package main

import (
	"fmt"
	"io"
	"os"
)

// PenghitungByte membungkus io.Writer dan menghitung byte yang lewat.
type PenghitungByte struct {
	io.Writer // embedded interface: Write diteruskan otomatis
	Total     int
}

func (p *PenghitungByte) Write(b []byte) (int, error) {
	n, err := p.Writer.Write(b)
	p.Total += n
	return n, err
}

func main() {
	pb := &PenghitungByte{Writer: os.Stdout}
	fmt.Fprintln(pb, "Halo interface!")
	fmt.Fprintf(pb, "angka: %d\n", 42)
	fmt.Println("total byte:", pb.Total)
}
```

## Latihan

1. Tambahkan tipe `Segitiga` ke program `Bentuk` dan hitung total luas dari sebuah `[]Bentuk` campuran. Mengapa Anda tidak perlu mengubah satu baris pun pada `totalLuas`?
2. Buat tipe `Celsius float64` dengan method `String()` dan `KeFahrenheit()`. Cetak dengan `fmt.Println` dan jelaskan bagaimana `fmt` memutuskan memakai `String()`.
3. Definisikan method `Reset()` dengan pointer receiver, lalu coba memenuhi interface `Resetter` memakai nilai biasa dan pointer. Catat pesan error compiler dan jelaskan konsep *method set*.
4. Reproduksi jebakan `nil` pada interface `error`: buat fungsi yang mengembalikan pointer `nil` bertipe konkret sebagai `error`, tunjukkan bahwa `err != nil` bernilai `true`, lalu perbaiki.
5. Tulis fungsi `Ringkas(x any) string` dengan type switch untuk `int`, `float64`, `string`, `[]int`, `error`, `nil`, dan `default`. Uji dengan masing-masing tipe.
6. Buat `io.Reader` kustom `HurufBesar` yang membungkus `io.Reader` lain dan mengubah isi menjadi huruf kapital, lalu salin ke `os.Stdout` memakai `io.Copy`. Hubungkan dengan prinsip interface kecil.
7. Diskusikan: sebuah paket `layananpesanan` membutuhkan fungsi `Simpan(p Pesanan) error` dari database. Mengapa lebih baik mendeklarasikan interface `Penyimpan` di paket `layananpesanan` (sisi pemakai) daripada di paket database? Bagaimana ini mempermudah pengujian dengan fake ([Modul 15](15-testing-benchmarking.md))?

---

[← Modul 7: Pointer & Memory Management](07-pointer-memory.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 9: Generics & Iterators →](09-generics.md)
