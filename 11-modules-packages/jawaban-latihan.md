# Modul 11: Modules, Packages & Code Organization — Jawaban Latihan

Modul ini berbasis banyak modul Go kecil, jadi `kode/` **tidak** punya `go.mod` di akarnya; setiap subfolder adalah modul sendiri. Jalankan perintah dari dalam folder yang disebut di tiap latihan, misalnya:

```bash
cd kode/tokoku && go run .
cd kode/l5-mvs/app && go list -m all && go run .
```

Semua output di bawah dijalankan sungguhan di **Go 1.24.7**, **tanpa internet** (`GOPROXY=off`, `GOSUMDB=off`). Karena itu seluruh "dependensi pihak ketiga" dalam latihan ini disimulasikan dengan **modul lokal + direktif `replace`**. Batasannya saya sebut jelas di Latihan 4 dan 6.

Dua penyesuaian terhadap materi:

- Module path memakai `contoh.com/tokoku` (bukan `tokoku` polos) agar menyerupai modul sungguhan.
- Materi memakai `go 1.27.0` pada contoh `go.mod`; di sini `go.mod` berisi versi toolchain yang terpasang (1.24.x).

---

## Latihan 1 — Modul `tokoku` dengan package `laporan` (`kode/tokoku`)

Struktur:

```
tokoku/
  go.mod                      # module contoh.com/tokoku
  main.go
  diskon/diskon.go
  internal/harga/harga.go
  laporan/laporan.go          # package ketiga: mengimpor harga dan diskon
```

```go
// laporan/laporan.go
package laporan

import (
	"fmt"
	"strings"

	"contoh.com/tokoku/diskon"
	"contoh.com/tokoku/internal/harga"
)

func Cetak(barang []int, persen int) (string, error) {
	total := harga.Total(barang)
	akhir, err := diskon.Terapkan(total, persen)
	if err != nil {
		return "", fmt.Errorf("laporan: %w", err)
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Jumlah barang : %d\n", len(barang))
	fmt.Fprintf(&b, "Total         : Rp%d\n", total)
	fmt.Fprintf(&b, "Diskon        : %d%%\n", persen)
	fmt.Fprintf(&b, "Bayar         : Rp%d\n", akhir)
	return b.String(), nil
}
```

Hasil `go run .`:

```
31500 <nil>
diskon 80% di luar batas 0-50%
Jumlah barang : 3
Total         : Rp40000
Diskon        : 20%
Bayar         : Rp32000
err: <nil>
laporan gagal: laporan: diskon 99% di luar batas 0-50%
```

`gofmt -l .` bersih, `go vet ./...` dan `go build ./...` sukses, dan `go list ./...` menampilkan empat package:

```
contoh.com/tokoku
contoh.com/tokoku/diskon
contoh.com/tokoku/internal/harga
contoh.com/tokoku/laporan
```

Catatan desain: `laporan` boleh mengimpor `internal/harga` karena berada di pohon modul yang sama. Error dari `diskon` dibungkus dengan `%w` (Modul 10), jadi `errors.Is/As` tetap bisa menelusurinya.

---

## Latihan 2 — Import melingkar (`kode/l2-melingkar`, `kode/l2-perbaikan`)

Package `a` mengimpor `b` (untuk tipe `Pesanan`), dan `b` mengimpor `a` (untuk tipe `Pengguna`). Hasil `go build ./...`:

```
package contoh.com/melingkar
	imports contoh.com/melingkar/a from main.go
	imports contoh.com/melingkar/b from a.go
	imports contoh.com/melingkar/a from b.go: import cycle not allowed
```

Pesan itu menunjukkan jejak siklusnya: `main → a → b → a`. Go melarang siklus sebagai error kompilasi, sehingga graf dependensi selalu berupa DAG (urutan inisialisasi dan kompilasi jadi jelas).

**Perbaikan 1: pindahkan tipe bersama ke package ketiga** (`l2-perbaikan`):

```
model/  -> Pengguna, Pesanan        (tidak mengimpor a maupun b)
a/      -> mengimpor model
b/      -> mengimpor model
```

```
$ go run .
pesanan 1 milik Andi
```

Sekarang panah dependensi hanya mengarah ke `model` (`a → model`, `b → model`), tanpa siklus.

**Perbaikan 2: balik arah dependensi dengan interface** (Modul 8). Misalnya `b` tidak mengimpor `a`, tetapi mendeklarasikan interface kecil `type Pemilik interface{ NamaPemilik() string }`; `a` yang memenuhinya secara implisit. Pilih salah satu: tipe data murni cocok dipindah ke package `model`, sedangkan *perilaku* cocok dibalik dengan interface.

Akar masalahnya biasanya desain: dua package yang saling butuh sebenarnya mewakili satu konsep yang terlalu dipecah, atau ada tipe bersama yang belum punya rumah sendiri.

---

## Latihan 3 — Aturan `internal/` (`kode/l3-modul-luar`)

Modul terpisah `contoh.com/pihaklain` mengimpor `contoh.com/tokoku/internal/harga`:

```go
import "contoh.com/tokoku/internal/harga" // internal milik modul lain
```

Hasil `go build ./...`:

```
package contoh.com/pihaklain
	main.go:6:2: use of internal package contoh.com/tokoku/internal/harga not allowed
```

Sebagai pembanding, mengimpor package **publik** (`contoh.com/tokoku/diskon`) dari modul yang sama berhasil (`go run .` mencetak `90000 <nil>`). Kode pelanggarnya disimpan sebagai `main_pelanggar.go.txt` agar folder tetap bisa dibangun.

**Aturan:** package di bawah `.../internal/...` hanya bisa diimpor oleh kode yang berada di **pohon direktori induk `internal`** itu (di sini, di dalam `tokoku/`).

**Manfaat bagi penulis library dibanding sekadar konvensi:**

1. **Ditegakkan compiler, bukan kesepakatan lisan.** Konvensi seperti awalan huruf kecil atau komentar "jangan dipakai" mudah dilanggar, apalagi ketika pengguna sedang terburu-buru. Dengan `internal/`, pelanggaran gagal dibangun.
2. **Ruang gerak untuk berubah.** Apa pun di `internal/` boleh diubah, dipindah, atau dihapus tanpa melanggar janji kompatibilitas semver, karena tidak ada pihak luar yang bisa bergantung padanya. Yang publik harus dijaga stabil (lihat Latihan 6).
3. **Permukaan API kecil dan jelas.** Pengguna hanya melihat package yang memang dimaksudkan sebagai API, sehingga dokumentasi dan dukungan jauh lebih sederhana.
4. Ini juga bekerja untuk *aplikasi*: `internal/` menjamin bahwa kode aplikasi tidak bocor dipakai proyek lain.

---

## Latihan 4 — Dependensi pihak ketiga (`kode/l4-dependensi`)

**Bagian yang tidak bisa saya lakukan.** `go get github.com/google/uuid@latest` gagal karena tidak ada internet:

```
go: github.com/google/uuid@latest: module github.com/google/uuid: Get "https://proxy.golang.org/github.com/google/uuid/@v/list": dial tcp: lookup proxy.golang.org on 8.8.8.8:53: no such host
```

Jadi perubahan `go.sum` berisi hash sungguhan tidak dapat saya tampilkan. Sebagai gantinya saya mensimulasikan alurnya dengan modul lokal `contoh.com/ext/uuidmini` lewat `replace`.

**Hasil simulasi:**

```
--- go.mod sebelum
module contoh.com/app
go 1.24

--- sesudah go mod edit -require ... -replace ... dan go mod tidy
module contoh.com/app
go 1.24
require contoh.com/ext/uuidmini v1.0.0
replace contoh.com/ext/uuidmini => ../ext/uuidmini

--- go mod why contoh.com/ext/uuidmini
# contoh.com/ext/uuidmini
contoh.com/app
contoh.com/ext/uuidmini
```

Setelah pemakaian dihapus dari kode:

```
$ go mod why contoh.com/ext/uuidmini
# contoh.com/ext/uuidmini
(main module does not need package contoh.com/ext/uuidmini)

$ go mod tidy        # baris "require" ikut dihapus dari go.mod
```

**Yang berlaku pada dependensi sungguhan (dari materi, tidak saya verifikasi di sini):**

- `go get` menambah baris `require` di `go.mod` **dan** baris hash di `go.sum` (dua baris per versi: hash isi modul dan hash `go.mod`-nya). Dependensi transitif yang dibawa juga ikut tercatat (kadang sebagai `// indirect`).
- `go.sum` ditulis hanya untuk modul yang diunduh dari proxy. Pada simulasi saya `go.sum` **tidak ada** karena modul `replace` ke direktori lokal tidak diverifikasi dengan hash. Itu pembedanya.
- Commit `go.mod` dan `go.sum`. Hash itu yang memastikan isi sebuah versi tidak berubah diam-diam, ditambah verifikasi terhadap checksum database `sum.golang.org`.
- `go mod why` menjawab "kenapa package ini ada di graf saya"; `go mod tidy` membersihkan `require` dan `go.sum` yang tidak terpakai lagi.

---

## Latihan 5 — Minimal Version Selection (`kode/l5-mvs`)

Setup lokal: modul `lib` punya tiga versi (v1.1.0, v1.2.0, v1.3.0, tiap versi diarahkan ke folder berbeda lewat `replace` per versi). `deta` membutuhkan `lib v1.1.0`, `detb` membutuhkan `lib v1.2.0`, dan aplikasi memakai keduanya.

```
$ go mod graph
contoh.com/app contoh.com/deta@v1.0.0
contoh.com/app contoh.com/detb@v1.0.0
contoh.com/app contoh.com/lib@v1.2.0
contoh.com/deta@v1.0.0 contoh.com/lib@v1.1.0
contoh.com/detb@v1.0.0 contoh.com/lib@v1.2.0
...

$ go run .
deta memakai lib v1.2.0
detb memakai lib v1.2.0
yang benar-benar dibangun: lib v1.2.0
```

**Cara MVS memilih:** untuk setiap modul, Go mengambil **versi tertinggi di antara semua versi yang diminta** oleh modul-modul dalam graf. Di sini yang diminta `v1.1.0` (oleh `deta`) dan `v1.2.0` (oleh `detb`), sehingga terpilih **v1.2.0**. Terlihat bahwa `deta` pun dibangun dengan v1.2.0, bukan v1.1.0 yang dimintanya, karena hanya satu versi `lib` yang dipakai dalam satu build (dalam satu major version).

Sifat penting MVS:

- **Deterministik dan tidak melompat ke yang terbaru.** v1.3.0 ada, tetapi tidak dipilih karena tidak ada yang memintanya. Build yang sama selalu memberi hasil yang sama, tanpa file kunci terpisah seperti `package-lock`.
- **Versi terpilih adalah yang *diminta***, bukan "paling baru yang tersedia".

Menaikkan versi:

```
$ go get contoh.com/lib@v1.3.0
go: upgraded contoh.com/lib v1.2.0 => v1.3.0
$ go run . | tail -1
yang benar-benar dibangun: lib v1.3.0
```

Temuan tak terduga saat mencoba **menurunkan** (`go get contoh.com/lib@v1.1.0`):

```
go: removed contoh.com/detb v1.0.0
go: downgraded contoh.com/lib v1.3.0 => v1.1.0
```

Karena `detb` mensyaratkan `lib v1.2.0`, Go tidak bisa memenuhi permintaan v1.1.0 tanpa **menghapus `detb`** dari `require`. Jadi penurunan bisa berefek ke dependensi lain; periksa output `go get` dengan teliti. (Saya mengembalikan `go.mod` ke keadaan awal sebelum membuat zip.)

`go list -m -u all` (daftar pembaruan) pada dependensi sungguhan butuh jaringan ke proxy untuk melihat versi terbaru; pada simulasi lokal ia hanya mengenali versi yang ada di `replace`.

---

## Latihan 6 — Menerbitkan modul dan `/v2` (`kode/l6-v2`)

**Bagian yang tidak bisa saya lakukan.** Menerbitkan ke repositori Git publik dengan tag `v0.1.0` lalu `go get` dari proyek lain memerlukan akses jaringan dan akun GitHub Anda. Silakan jalankan langkah di 11.7:

```bash
git tag v0.1.0 && git push origin v0.1.0
go get github.com/<anda>/<modul>@v0.1.0   # dari proyek lain
```

**Bagian yang saya buktikan: perubahan yang merusak membutuhkan `/v2`.** Saya membuat `pustaka` v1 (`Sapa(nama)`) dan v2 (`Sapa(nama, sapaan)`, signature berbeda = *breaking change*):

```
pustaka-v1/go.mod  ->  module contoh.com/pustaka
pustaka-v2/go.mod  ->  module contoh.com/pustaka/v2      # akhiran /v2 WAJIB
```

Aplikasi memakai **keduanya berdampingan**:

```go
import (
	pustaka "contoh.com/pustaka"
	pustakav2 "contoh.com/pustaka/v2"
)
```

```
Halo, Andi
Selamat pagi, Andi
contoh.com/pustaka v1.5.0 => ../pustaka-v1
contoh.com/pustaka/v2 v2.0.0 => ../pustaka-v2
```

Bila v2 *tidak* memakai akhiran `/v2` pada module path, Go menolaknya:

```
require contoh.com/pustaka: version "v2.0.0" invalid: should be v0 or v1, not v2
```

**Yang harus diubah saat merilis perubahan API yang merusak:**

1. Ubah baris `module` di `go.mod` menjadi `.../v2` dan perbarui semua import internal modul itu.
2. Tag rilis `v2.0.0` (bukan `v1.x`). Dua pendekatan repositori: folder `v2/` terpisah atau cabang mayor; keduanya valid.
3. Pengguna lama tetap aman di v1 sampai mereka memutuskan pindah (karena path berbeda, tidak ada konflik); pengguna baru mengimpor `.../v2`.
4. Dokumentasikan migrasi, dan tandai API lama dengan `// Deprecated: ...`.

Aturan versi (11.7): `v0.x` belum stabil; sejak `v1` Anda menjanjikan kompatibilitas; rilis yang salah tidak bisa dihapus dari proxy, jadi gunakan direktif `retract`.

---

## Latihan 7 — Struktur `models/`, `controllers/`, `utils/` (diskusi)

**Akar masalah.** Struktur itu memecah kode berdasarkan **jenis teknis** (apa bentuknya), bukan **tanggung jawab domain** (apa urusannya). Akibatnya:

- Fitur yang sama (misalnya pesanan) tersebar di tiga package yang saling bergantung: `controllers` memakai `models`, `models` butuh helper dari `utils`, `utils` butuh tipe dari `models`, dan `models` memanggil sesuatu di `controllers`. Itu resep siklus seperti di Latihan 2.
- `utils` menjadi tempat sampah dengan nama tanpa makna; semua package menariknya masuk sehingga ia menjadi pusat ketergantungan.
- Mengubah satu fitur memerlukan edit di banyak package, dan tidak ada batas enkapsulasi yang jelas.

**Usulan struktur berbasis domain:**

```
tokoku/
├── go.mod
├── cmd/
│   ├── api/main.go           # entry point server HTTP (tipis: rangkai dan jalankan)
│   └── migrasi/main.go       # entry point tool lain
├── internal/
│   ├── pesanan/              # tipe, aturan bisnis, repository untuk pesanan
│   ├── pembayaran/           # domain pembayaran
│   ├── katalog/              # domain produk
│   ├── transport/http/       # handler HTTP: menerjemahkan request ke pemanggilan domain
│   └── config/
└── migrations/
```

**Peran tiap bagian:**

- **`cmd/<nama>/main.go`** hanya merangkai dependensi (membaca konfigurasi, membuat repository dan service, menjalankan server). Logikanya sesedikit mungkin, sehingga mudah diuji dan satu modul bisa punya beberapa executable.
- **`internal/`** mengunci kode aplikasi dari impor proyek lain (Latihan 3), sehingga struktur internal boleh bebas berubah.
- **Package per domain** (`pesanan`, `pembayaran`): tipe, logika, dan akses data yang berubah bersama berada berdekatan (kohesif).

**Mencegah siklus:** domain tidak boleh mengimpor lapisan transport; arah dependensi satu arah: `cmd → transport → domain`. Bila `pesanan` butuh sesuatu dari `pembayaran`, definisikan **interface kecil di sisi pemakai** (`pesanan` mendeklarasikan interface `Pembayar`) dan biarkan `cmd` yang menyambungkan implementasinya (Modul 8, Latihan 7). Tipe yang benar-benar dipakai bersama dipindah ke package kecil bernama jelas, bukan `utils`.

---

## Lampiran: bagian Praktik dari materi

**`replace` lokal** (`kode/praktik-replace`, `kode/salam`). Perintah yang dijalankan:

```bash
go mod edit -require=contoh.com/salam@v0.0.0 -replace=contoh.com/salam=../salam
go mod tidy
```

Temuan: urutan di materi (`edit` → `tidy` → baru menambah `import`) membuat `go mod tidy` **menghapus baris `require`** karena belum ada kode yang memakai `salam`. Setelah `import "contoh.com/salam"` ditambahkan, build gagal dengan:

```
module contoh.com/salam provides package contoh.com/salam and is replaced but not required; to add it:
	go get contoh.com/salam
```

Jalankan `go mod tidy` (atau `go get`) **sekali lagi setelah import ditambahkan**, dan hasilnya:

```
module contoh.com/tokoku
go 1.24.7
replace contoh.com/salam => ../salam
require contoh.com/salam v0.0.0-00010101000000-000000000000
```

`go run .` kemudian mencetak `Halo, Andi!` diikuti keluaran `tokoku`. (Versi `v0.0.0-0001...` adalah *pseudo-version* otomatis untuk modul lokal.)

**Workspace** (`kode/praktik-workspace`):

```
$ go work init ./tokoku ./salam
$ cat go.work
go 1.24.7

use (
	./salam
	./tokoku
)
```

`go.mod` milik `tokoku` **tetap tanpa `require`/`replace` untuk `salam`**, namun `go run .` berhasil (`Halo, Budi!`), sebab workspace menyatukan dua modul lokal. Karena itu `go.work` biasanya tidak di-commit (11.8).

**Vendoring dan inspeksi** (`kode/praktik-vendor`):

```
$ go list -m all
contoh.com/tokoku
contoh.com/salam v0.0.0-00010101000000-000000000000 => ../salam

$ go mod verify
all modules verified

$ go mod vendor && ls vendor
contoh.com  modules.txt
$ go build -mod=vendor ./...        # sukses

$ go list -m all
go: can't compute 'all' using the vendor directory
	(Use -mod=mod or -mod=readonly to bypass.)
```

Perilaku terakhir persis seperti peringatan di materi: lakukan inspeksi dependensi **sebelum** vendoring, atau tambahkan `-mod=mod`.

Hal yang tidak saya coba di sini: `go get -tool` (direktif `tool`) dan `GOPRIVATE`, karena keduanya bergantung pada unduhan dari proxy atau repositori privat.
