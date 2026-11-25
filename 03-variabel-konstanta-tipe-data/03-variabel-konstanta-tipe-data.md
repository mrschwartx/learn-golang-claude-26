# Modul 3: Variabel, Konstanta & Tipe Data

[← Modul 2: Instalasi & Setup Environment](02-instalasi-setup.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 4: Composite Types (Array, Slice, Map, Struct) →](04-composite-types.md)

**Isi modul ini**

- [Teori](#teori)
  - [3.1 Variabel: `var` vs `:=`](#31-variabel-var-vs-)
  - [3.2 Zero Values](#32-zero-values)
  - [3.3 Konstanta, `const`, dan `iota`](#33-konstanta-const-dan-iota)
  - [3.4 Scope dan Shadowing](#34-scope-dan-shadowing)
  - [3.5 Tipe Data: Boolean dan Numerik](#35-tipe-data-boolean-dan-numerik)
  - [3.6 Runes dan Strings](#36-runes-dan-strings)
  - [3.7 Type Conversion](#37-type-conversion)
  - [3.8 Konvensi Penamaan dan Visibilitas](#38-konvensi-penamaan-dan-visibilitas)
- [Praktik](#praktik)
  - [Deklarasi dan Zero Value](#deklarasi-dan-zero-value)
  - [`const` dan `iota`](#const-dan-iota)
  - [Rune, String, dan Konversi](#rune-string-dan-konversi)
  - [Bilangan Kompleks](#bilangan-kompleks)
- [Latihan](#latihan)

---


## Teori

### 3.1 Variabel: `var` vs `:=`

Go punya dua cara utama mendeklarasikan variabel.

| Bentuk | Contoh | Catatan |
|---|---|---|
| `var` dengan tipe | `var umur int = 30` | Eksplisit; bisa dipakai di level paket maupun fungsi |
| `var` tanpa tipe | `var nama = "Andi"` | Tipe disimpulkan (*type inference*) dari nilainya |
| `var` tanpa nilai | `var total float64` | Variabel mendapat **zero value** (3.2) |
| `var` berkelompok | `var ( x = 1; y = 2 )` | Merapikan banyak deklarasi sekaligus |
| **Short declaration** | `kota := "Jakarta"` | **Hanya di dalam fungsi**; paling umum dipakai |

Aturan `:=` yang sering menjebak pemula:

- Minimal **satu variabel di sisi kiri harus baru**: `a, err := f()` lalu `b, err := g()` valid karena `b` baru (`err` dipakai ulang), tetapi `err := h()` kedua kali dalam scope yang sama adalah error.
- `:=` tidak bisa dipakai di level paket — gunakan `var`.
- **Variabel lokal yang tidak dipakai adalah error kompilasi.** Gunakan *blank identifier* `_` jika nilai memang ingin diabaikan.

### 3.2 Zero Values

Variabel yang dideklarasikan tanpa nilai awal **tidak pernah "kosong/tidak terdefinisi"** — selalu memiliki **zero value** sesuai tipenya. Ini menghilangkan satu kelas bug penuh (variabel tidak diinisialisasi) yang dikenal dari C.

| Tipe | Zero Value |
|---|---|
| `int`, `float64`, dan tipe numerik lain | `0` |
| `bool` | `false` |
| `string` | `""` (string kosong) |
| pointer, slice, map, channel, function, interface | `nil` |
| struct | Struct dengan setiap field bernilai zero value-nya |
| array | Array dengan setiap elemen bernilai zero value-nya |

Desain idiomatik Go berusaha membuat **zero value berguna** — misalnya `sync.Mutex` dan `bytes.Buffer` langsung bisa dipakai tanpa inisialisasi khusus.

### 3.3 Konstanta, `const`, dan `iota`

**Konstanta** dideklarasikan dengan `const`, nilainya ditentukan saat **kompilasi** dan tidak bisa diubah. Hanya tipe dasar (angka, string, boolean, rune) yang bisa menjadi konstanta.

```go
const Pi = 3.14159
const (
    NamaApp = "BelajarGo"
    MaxUser = 100
)
```

**Untyped constant**: konstanta tanpa tipe eksplisit bersifat fleksibel dan memiliki presisi sangat tinggi saat kompilasi, sehingga `const Besar = 1 << 100` valid selama tidak dipakai di tempat yang membutuhkan `int`.

**`iota`** adalah generator nilai berurutan di dalam blok `const` — dimulai dari `0` dan bertambah satu untuk setiap baris. Ini cara Go menirukan **enum**.

| Pola | Hasil |
|---|---|
| `A = iota; B; C` | `A=0, B=1, C=2` (baris tanpa ekspresi mengulang ekspresi sebelumnya) |
| `_ = iota; Satu; Dua` | Melewati nilai `0` dengan `_` |
| `KB = 1 << (10 * (iota + 1))` | Pola ukuran (KB, MB, GB): setiap baris menggeser 10 bit |
| `FlagA = 1 << iota` | Pola **bit flag** yang bisa digabung dengan `\|` |

### 3.4 Scope dan Shadowing

Go memakai **lexical scope** berbasis blok `{ }`:

1. **Universe** — identifier bawaan (`int`, `true`, `nil`, `len`, `append`).
2. **Package** — deklarasi di luar fungsi; terlihat oleh seluruh file dalam paket yang sama.
3. **File** — nama paket yang di-`import`.
4. **Function/Block** — variabel dalam fungsi, `if`, `for`, atau `switch`.

**Shadowing** terjadi saat variabel di scope dalam memakai nama yang sama dengan variabel di scope luar — variabel dalam "menutupi" yang luar di blok tersebut.

```go
x := 1
if true {
    x := 2      // variabel BARU, menutupi x luar
    _ = x
}
fmt.Println(x)  // 1, bukan 2
```

Shadowing paling sering menimbulkan bug saat memakai `:=` di dalam blok untuk `err` — nilai error diisi ke variabel baru yang hilang begitu blok selesai, sementara `err` luar tetap `nil`. Analyzer `shadow` (opsional) pada `go vet` dan beberapa linter ([Modul 18](18-toolchain-code-quality.md)) dapat mendeteksinya. Hindari juga menimpa identifier universe seperti `len` atau `string`; kode tetap valid tetapi membingungkan.

### 3.5 Tipe Data: Boolean dan Numerik

| Kategori | Tipe | Keterangan |
|---|---|---|
| Boolean | `bool` | `true` / `false`; tidak ada konversi implisit dari angka |
| Integer bertanda | `int8` `int16` `int32` `int64` | Rentang dari `-2^(n-1)` sampai `2^(n-1)-1` |
| Integer tak bertanda | `uint8` `uint16` `uint32` `uint64` | Rentang `0` sampai `2^n - 1` |
| Integer platform | `int`, `uint`, `uintptr` | Ukuran mengikuti arsitektur: **32 atau 64 bit** |
| Alias | `byte` = `uint8`, `rune` = `int32` | `byte` untuk data biner, `rune` untuk karakter Unicode |
| Floating point | `float32`, `float64` | `float64` adalah default untuk literal desimal |
| Complex | `complex64`, `complex128` | Bilangan kompleks bawaan (`complex(1, 2)`, `real`, `imag`) |

Hal penting yang perlu diwaspadai:

- **Integer overflow** pada variabel biasa terjadi **diam-diam (wrap-around)**, tidak menimbulkan error maupun panic — misalnya `uint8(255) + 1` menghasilkan `0`. Pada konstanta, overflow adalah error kompilasi.
- **Pembagian integer** membulatkan ke arah nol: `7 / 2 == 3`. Operator `%` mengambil tanda dari bilangan yang dibagi: `-7 % 3 == -1`.
- **Floating point tidak presisi**: `x + y` dengan `x := 0.1` dan `y := 0.2` tidak persis `0.3` (hasilnya `0.30000000000000004`). Perbandingan konstanta literal `0.1+0.2 == 0.3` justru `true` karena konstanta dihitung dengan presisi tak terbatas saat kompilasi. Jangan pakai `float64` untuk uang — gunakan integer satuan terkecil (sen) atau library desimal.
- **Bilangan kompleks** (`complex64` = dua `float32`, `complex128` = dua `float64`) adalah tipe bawaan, bukan library. Dibuat dengan fungsi `complex(real, imag)` atau literal imajiner `3 + 4i`, dan dipecah dengan `real(z)` dan `imag(z)`. Literal kompleks tanpa tipe berjenis `complex128` secara default. Operator `+ - * /` dan `==` berlaku, tetapi `<` dan `>` **tidak** (bilangan kompleks tidak terurut). Fungsi matematika (nilai absolut, akar, eksponen) ada di paket `math/cmplx`. Di praktik sehari-hari tipe ini jarang dipakai; ia relevan untuk pemrosesan sinyal, simulasi, dan komputasi ilmiah.
- Pakai `int` sebagai pilihan default untuk angka bulat, kecuali Anda butuh ukuran tertentu (format file, protokol jaringan, atau memori).

### 3.6 Runes dan Strings

**String** di Go adalah **urutan byte yang tidak bisa diubah (immutable)**, yang secara konvensi berisi teks **UTF-8**. Konsekuensinya:

| Operasi | Hasil untuk `s := "héllo"` |
|---|---|
| `len(s)` | `6` — jumlah **byte** (huruf `é` memakai 2 byte), bukan jumlah karakter |
| `s[1]` | Sebuah **byte** (`195`), bukan karakter `é` |
| `utf8.RuneCountInString(s)` | `5` — jumlah karakter Unicode (rune) |
| `[]rune(s)` | Slice `rune`, tiap elemen satu karakter |
| `for i, r := range s` | Iterasi **per rune**; `i` adalah posisi **byte** awal rune |

**Rune** adalah satu *code point* Unicode (tipe `int32`), ditulis dengan tanda kutip tunggal: `'A'`, `'世'`.

**Dua jenis string literal**:

| Jenis | Sintaks | Perilaku |
|---|---|---|
| **Interpreted** | `"Halo\nDunia"` | Escape sequence (`\n`, `\t`, `\"`, `\u00e9`) diproses |
| **Raw** | `` `C:\data\file.txt` `` | Isi apa adanya, boleh **multi-baris**, tanpa escape — cocok untuk regex, SQL, JSON, path Windows |

Karena string immutable, menyambung string dalam loop dengan `+=` membuat salinan baru setiap iterasi (O(n²)). Gunakan **`strings.Builder`** untuk perakitan string dalam jumlah besar.

### 3.7 Type Conversion

Go **tidak pernah mengonversi tipe secara implisit** — bahkan antara `int` dan `int64` atau `int` dan `float64`. Konversi harus eksplisit dengan sintaks `T(nilai)`.

| Kebutuhan | Cara |
|---|---|
| Antar tipe numerik | `float64(i)`, `int(f)` (**memotong** bagian desimal, bukan membulatkan) |
| Integer ke string desimal | `strconv.Itoa(42)` → `"42"` |
| String ke integer | `strconv.Atoi("42")` → `(42, nil)`, atau `(0, error)` jika gagal |
| String ke float / bool | `strconv.ParseFloat`, `strconv.ParseBool` |
| String ↔ `[]byte` / `[]rune` | `[]byte(s)`, `string(b)`, `[]rune(s)` (menyalin data) |

Jebakan klasik: `string(65)` menghasilkan `"A"` (karakter dengan kode 65), **bukan** `"65"`. `go vet` memperingatkan konversi `int` ke `string` semacam ini. Untuk mengubah angka menjadi teks, gunakan `strconv.Itoa` atau `fmt.Sprint`.

### 3.8 Konvensi Penamaan dan Visibilitas

Go tidak punya keyword `public`/`private`. **Visibilitas ditentukan huruf pertama nama**:

| Awalan | Arti |
|---|---|
| **Huruf kapital** (`Hitung`, `User`) | **Exported** — bisa diakses dari paket lain |
| **Huruf kecil** (`hitung`, `user`) | **Unexported** — hanya terlihat dalam paketnya sendiri |

Konvensi lain: gunakan **MixedCaps/camelCase** (bukan `snake_case`); akronim ditulis konsisten (`userID`, `HTTPServer`, `URL`); nama pendek untuk scope kecil (`i`, `r`, `n`), nama deskriptif untuk scope luas; nama paket huruf kecil, satu kata, tanpa underscore.

## Praktik

### Deklarasi dan Zero Value

```go
// variabel.go
package main

import "fmt"

var versi = "1.0" // level paket: wajib pakai var

func main() {
	var umur int
	var harga float64
	var aktif bool
	var nama string
	var data []int
	var peta map[string]int

	fmt.Printf("%T %v\n", umur, umur)
	fmt.Printf("%T %v\n", harga, harga)
	fmt.Printf("%T %v\n", aktif, aktif)
	fmt.Printf("%T %q\n", nama, nama)
	fmt.Printf("%T %v %t\n", data, data, data == nil)
	fmt.Printf("%T %v %t\n", peta, peta, peta == nil)

	kota := "Jakarta"
	a, b := 10, 3
	fmt.Println(versi, kota, a/b, a%b, float64(a)/float64(b))
}
```

```bash
go run variabel.go
```

### `const` dan `iota`

```go
// iota.go
package main

import "fmt"

type Hari int

const (
	Minggu Hari = iota
	Senin
	Selasa
)

func namaHari(h Hari) string {
	switch h {
	case Minggu:
		return "Minggu"
	case Senin:
		return "Senin"
	default:
		return "Lainnya"
	}
}

const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
	GB
)

const (
	Baca  = 1 << iota // 1
	Tulis             // 2
	Eksek             // 4
)

func main() {
	fmt.Println(Minggu, Senin, Selasa, namaHari(Senin))
	fmt.Println(KB, MB, GB)

	izin := Baca | Tulis
	fmt.Println("izin:", izin, "bisa tulis:", izin&Tulis != 0, "bisa eksek:", izin&Eksek != 0)

	var kecil uint8 = 255
	kecil++
	fmt.Println("overflow uint8:", kecil)

	x, y := 0.1, 0.2
	fmt.Println(x+y, x+y == 0.3) // variabel: dihitung float64 saat runtime
	fmt.Println(0.1+0.2 == 0.3)  // konstanta: presisi tak terbatas saat kompilasi
}
```

### Rune, String, dan Konversi

```go
// teks.go
package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

func main() {
	s := "héllo, 世界"
	fmt.Println("byte:", len(s), "rune:", utf8.RuneCountInString(s))

	for i, r := range "héy" {
		fmt.Printf("posisi %d: %c (%U)\n", i, r, r)
	}

	raw := `baris 1
baris 2 \n tidak diproses`
	fmt.Println(raw)

	angka, err := strconv.Atoi("123abc")
	fmt.Println(angka, err)

	n, _ := strconv.Atoi("42")
	fmt.Println(n+1, strconv.Itoa(n)+"!", string(rune(65)))

	var sb strings.Builder
	for i := 0; i < 3; i++ {
		fmt.Fprintf(&sb, "[%d]", i)
	}
	fmt.Println(sb.String())

	f := 9.99
	fmt.Println(int(f))
}
```

### Bilangan Kompleks

```go
// kompleks.go
package main

import (
	"fmt"
	"math/cmplx"
)

func main() {
	a := complex(3, 4)         // complex128 (default)
	var b complex128 = 1 + 2i  // literal imajiner
	var c complex64 = complex(1, 1)

	fmt.Printf("%T %v\n", a, a)
	fmt.Printf("%T %v\n", c, c)

	fmt.Println(a+b, a*b)            // (4+6i) (-5+10i)
	fmt.Println(real(a), imag(a))    // 3 4
	fmt.Println(cmplx.Abs(a))        // 5  (panjang vektor 3-4-5)
	fmt.Println(cmplx.Sqrt(-1))      // (0+1i)
	fmt.Println(a == complex(3, 4))  // true
	// fmt.Println(a < b)            // error kompilasi: bilangan kompleks tidak terurut
}
```

```bash
go run kompleks.go
```

Perhatikan bahwa `a * b` menghitung `(3+4i)(1+2i) = 3 + 6i + 4i + 8i² = -5 + 10i`, sehingga hasilnya sesuai aturan aljabar kompleks tanpa kode tambahan.

## Latihan

1. Deklarasikan variabel dengan lima cara berbeda (`var` dengan tipe, `var` tanpa tipe, `var` tanpa nilai, `var` berkelompok, `:=`) lalu cetak nilai dan tipenya dengan `%T`. Mengapa `:=` tidak bisa dipakai di level paket?
2. Buat sebuah blok yang memakai `:=` pada `err` di scope dalam sehingga `err` luar tetap `nil`. Jelaskan apa yang terjadi dan bagaimana cara memperbaikinya.
3. Buat enum `Level` (Debug, Info, Warn, Error) dengan `iota`, lalu tulis fungsi yang mengubahnya menjadi string. Tambahkan juga konstanta bit flag untuk tiga izin file dan cek izin dengan operator bitwise.
4. Cetak `len` dan jumlah rune dari `"Selamat pagi 🌅"`. Mengapa hasilnya berbeda? Lalu balik urutan karakter string tersebut dengan benar menggunakan `[]rune`.
5. Tulis program yang membaca string `"12x"` dengan `strconv.Atoi`, menangani error-nya dengan `if err != nil`, dan mencetak pesan yang ramah pengguna.
6. Coba `fmt.Println(string(65))` dan `go vet`. Jelaskan mengapa hasilnya `A` dan apa cara yang benar untuk mendapatkan teks `"65"`.
7. Diskusikan: sebuah toko online menyimpan harga dalam `float64`. Jelaskan mengapa ini berisiko untuk perhitungan uang, dan rancang alternatif yang aman memakai integer.
8. Hitung akar-akar persamaan kuadrat `x² + 2x + 5 = 0` dengan rumus ABC memakai `complex128` dan `cmplx.Sqrt` (diskriminannya negatif). Cetak bagian real dan imajiner tiap akar dengan `real` dan `imag`, lalu verifikasi dengan mensubstitusikan akar ke persamaan. Mengapa `a < b` untuk bilangan kompleks tidak dapat dikompilasi?

---

[← Modul 2: Instalasi & Setup Environment](02-instalasi-setup.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 4: Composite Types (Array, Slice, Map, Struct) →](04-composite-types.md)
