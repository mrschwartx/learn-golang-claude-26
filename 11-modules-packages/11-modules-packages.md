# Modul 11: Modules, Packages & Code Organization

[← Modul 10: Error Handling](10-error-handling.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 12: Goroutines & Channels →](12-goroutines-channels.md)

**Isi modul ini**

- [Teori](#teori)
  - [11.1 Package — Unit Organisasi Kode](#111-package--unit-organisasi-kode)
  - [11.2 Go Modules](#112-go-modules)
  - [11.3 `go mod init`, `tidy`, dan `vendor`](#113-go-mod-init-tidy-dan-vendor)
  - [11.4 Package Import Rules](#114-package-import-rules)
  - [11.5 Menggunakan Package Pihak Ketiga](#115-menggunakan-package-pihak-ketiga)
  - [11.6 Struktur Proyek](#116-struktur-proyek)
  - [11.7 Publishing Modules](#117-publishing-modules)
  - [11.8 Dependensi Tool, Workspace, dan `replace`](#118-dependensi-tool-workspace-dan-replace)
- [Praktik](#praktik)
  - [Membuat Modul dengan Dua Package](#membuat-modul-dengan-dua-package)
  - [Menguji Aturan `internal/`](#menguji-aturan-internal)
  - [`replace` Lokal dan Workspace](#replace-lokal-dan-workspace)
  - [Vendoring dan Inspeksi Dependensi](#vendoring-dan-inspeksi-dependensi)
- [Latihan](#latihan)

---


## Teori

### 11.1 Package — Unit Organisasi Kode

**Package** adalah unit dasar organisasi kode di Go. Aturannya sederhana:

- **Satu direktori = satu package.** Semua file `.go` dalam sebuah direktori harus mendeklarasikan nama package yang sama (kecuali file `_test`, [Modul 15](15-testing-benchmarking.md)).
- Package bernama **`main`** menghasilkan program executable dan wajib memiliki `func main()`. Package lain adalah **library**.
- Identifier berhuruf kapital **diekspor** (terlihat di luar package); huruf kecil bersifat privat bagi package ([Modul 3](03-variabel-konstanta-tipe-data.md)).
- Nama package: huruf kecil, singkat, satu kata, tanpa underscore atau `camelCase` (`harga`, `http`, `diskon`). Hindari nama generik seperti `util` atau `common`.
- Fungsi **`init()`** di dalam package dijalankan otomatis **sebelum `main`**, setelah semua variabel level paket diinisialisasi. Gunakan seperlunya — efek samping saat import mempersulit pengujian.

Urutan inisialisasi: package yang diimpor diinisialisasi lebih dulu (mengikuti graf dependensi), lalu variabel level paket, lalu `init()`, lalu `main()`.

### 11.2 Go Modules

**Module** adalah kumpulan package yang dirilis dan diberi versi bersama, didefinisikan oleh file **`go.mod`** di akarnya. Module menggantikan model lama `GOPATH`.

| Direktif `go.mod` | Fungsi |
|---|---|
| `module jalur/modul` | **Module path** — identitas modul sekaligus prefiks import-nya (biasanya URL repositori) |
| `go 1.27.0` | Versi Go minimum yang dibutuhkan ([Modul 2](02-instalasi-setup.md)) |
| `toolchain go1.27.4` | Toolchain yang disarankan (opsional) |
| `require pkg v1.2.3` | Dependensi beserta versi minimumnya |
| `replace a => ../b` | Mengganti sumber dependensi (path lokal atau fork) |
| `exclude pkg v1.2.3` | Melarang versi tertentu |
| `retract v1.0.1` | Penulis modul menandai versinya sendiri sebagai **bermasalah** |
| `tool pkg/cmd/x` | Mendeklarasikan **dependensi tool** (Go 1.24+, lihat 11.8) |

Berkas pendamping **`go.sum`** menyimpan **hash kriptografis** tiap versi dependensi untuk memastikan isinya tidak berubah. Commit `go.mod` dan `go.sum` ke repositori.

**Versioning** memakai *semantic versioning* (`vMAJOR.MINOR.PATCH`). Untuk versi mayor **v2 ke atas**, module path harus diakhiri `/v2` (`example.com/lib/v2`) sehingga v1 dan v2 dapat dipakai berdampingan. Go memilih versi dengan **Minimal Version Selection (MVS)**: versi tertinggi yang **diminta** oleh salah satu modul dalam graf dependensi — hasilnya deterministik dan tidak otomatis melompat ke rilis terbaru.

### 11.3 `go mod init`, `tidy`, dan `vendor`

| Perintah | Fungsi |
|---|---|
| `go mod init jalur/modul` | Membuat `go.mod` untuk modul baru |
| `go mod tidy` | Menambah dependensi yang dipakai kode dan **menghapus yang tidak terpakai**; jalankan sebelum commit |
| `go mod vendor` | Menyalin seluruh dependensi ke direktori `vendor/` |
| `go mod download` | Mengunduh dependensi ke cache modul |
| `go mod verify` | Memverifikasi hash cache modul terhadap `go.sum` |
| `go mod why pkg` | Menjelaskan mengapa sebuah package dibutuhkan |
| `go mod graph` | Mencetak graf dependensi |
| `go list -m -u all` | Daftar modul beserta pembaruan yang tersedia |
| `go get pkg@v1.2.3` | Menambah atau mengubah versi dependensi (`@latest`, `@none` untuk menghapus) |

**Vendoring** menyimpan salinan dependensi di dalam repositori. Manfaatnya: build **tanpa jaringan**, kontrol penuh atas kode yang dibangun, dan audit lebih mudah — dengan harga ukuran repositori yang lebih besar. Jika ada direktori `vendor/`, `go build` memakainya secara default.

### 11.4 Package Import Rules

| Aturan | Penjelasan |
|---|---|
| **Import path** | `module path` + jalur direktori, misalnya `github.com/andi/tokoku/diskon` |
| **Tidak boleh ada import melingkar** | Jika A mengimpor B, maka B tidak boleh (langsung maupun tidak langsung) mengimpor A — ini error kompilasi |
| **Import tak terpakai = error** | Hapus, atau pakai blank import `_` jika memang hanya butuh efek samping `init()` (misalnya driver database) |
| **Alias** | `import str "strings"` mengatasi konflik nama |
| **Dot import** | `import . "pkg"` — **hindari**, merusak keterbacaan |
| **Direktori `internal/`** | Package di bawah `internal/` hanya bisa diimpor oleh kode **dalam pohon direktori induk `internal`** — penegakan enkapsulasi oleh compiler |
| **Pengelompokan** | Konvensi: standard library, lalu pihak ketiga, lalu package internal proyek, dipisah baris kosong (dirapikan `goimports`) |

Import melingkar biasanya pertanda desain: pindahkan tipe bersama ke package ketiga, atau gunakan interface ([Modul 8](08-methods-interfaces.md)) agar arah dependensi menjadi satu arah.

### 11.5 Menggunakan Package Pihak Ketiga

```bash
go get github.com/google/uuid@latest
```

Lalu `import "github.com/google/uuid"`. Dokumentasi seluruh modul publik ada di **https://pkg.go.dev**. Modul diunduh melalui **module proxy** (`GOPROXY`, default `proxy.golang.org`) dan diverifikasi terhadap **checksum database** (`GOSUMDB`, `sum.golang.org`), sehingga isi sebuah versi tidak bisa diam-diam diubah.

Untuk repositori privat, atur `GOPRIVATE` agar tidak melewati proxy dan checksum publik:

```bash
go env -w GOPRIVATE=github.com/perusahaan/*
```

Pertimbangkan sebelum menambah dependensi: apakah **standard library** sudah cukup (Go 1.27, misalnya, menambah dukungan UUID ke standard library — lihat catatan rilis), seberapa aktif pemeliharaannya, lisensinya, dan berapa banyak dependensi transitif yang ikut masuk. Gunakan `govulncheck` untuk memeriksa kerentanan ([Modul 18](18-toolchain-code-quality.md)).

### 11.6 Struktur Proyek

Go tidak memaksa struktur tertentu. Untuk proyek kecil, **datar** sudah cukup (beberapa file `.go` dalam satu package `main`). Untuk aplikasi yang lebih besar, tata letak umum:

```
tokoku/
├── go.mod
├── go.sum
├── cmd/
│   ├── api/main.go          # entry point server
│   └── migrasi/main.go      # entry point tool lain
├── internal/                # kode privat aplikasi
│   ├── pesanan/
│   ├── pembayaran/
│   └── config/
├── pkg/                     # (opsional) library yang boleh dipakai proyek lain
├── migrations/
└── README.md
```

| Direktori | Fungsi |
|---|---|
| `cmd/<nama>/` | Satu subdirektori per executable; `main.go` dibuat **setipis mungkin** |
| `internal/` | Logika aplikasi yang **tidak boleh** diimpor proyek lain — pilihan default yang aman |
| `pkg/` | Penggunaannya diperdebatkan; hanya pakai jika memang menyediakan library publik |

Kelompokkan package **berdasarkan tanggung jawab domain** (`pesanan`, `pembayaran`), bukan berdasarkan jenis teknis (`models`, `controllers`, `utils`) yang cenderung memicu import melingkar.

### 11.7 Publishing Modules

Menerbitkan modul di Go cukup dengan **repositori Git yang publik** — tidak ada registry terpisah yang harus didaftarkan.

1. Pastikan `module` di `go.mod` sama dengan URL repositori (`github.com/andi/tokoku`).
2. Sertakan `README` dan **`LICENSE`** (pkg.go.dev hanya menampilkan dokumentasi modul berlisensi yang dapat didistribusikan).
3. Beri **tag semver**: `git tag v1.0.0 && git push origin v1.0.0`.
4. Pengguna cukup `go get github.com/andi/tokoku@v1.0.0`. Proxy dan pkg.go.dev akan mengambilnya otomatis saat pertama kali diminta.

Aturan versi: `v0.x` = belum stabil; mulai `v1` Anda berjanji menjaga **kompatibilitas** API; perubahan yang merusak membutuhkan modul `v2` dengan path `/v2`. Versi yang salah rilis tidak bisa dihapus dari proxy — gunakan direktif **`retract`**; API usang ditandai dengan paragraf komentar `// Deprecated: ...`.

### 11.8 Dependensi Tool, Workspace, dan `replace`

| Fitur | Penjelasan |
|---|---|
| **`tool` directive** (Go 1.24+) | `go get -tool golang.org/x/tools/cmd/stringer` mencatat tool di `go.mod`; jalankan dengan `go tool stringer`. Versi tool terkunci bersama proyek |
| **`replace`** | `replace example.com/lib => ../lib` — mengembangkan dua modul sekaligus tanpa publikasi |
| **Workspace** (`go.work`, Go 1.18+) | `go work init ./a ./b` — beberapa modul lokal dikerjakan bersama tanpa mengubah `go.mod`; **jangan commit** `go.work` kecuali sengaja |

## Praktik

### Membuat Modul dengan Dua Package

```bash
mkdir tokoku && cd tokoku
go mod init tokoku
mkdir -p internal/harga diskon
```

```go
// internal/harga/harga.go
package harga

// Total menjumlahkan harga barang (dalam rupiah).
func Total(barang []int) int {
	total := 0
	for _, h := range barang {
		total += h
	}
	return total
}
```

```go
// diskon/diskon.go
package diskon

import "fmt"

var DiskonMaks = 50

// Terapkan mengurangi total dengan persentase diskon.
func Terapkan(total, persen int) (int, error) {
	if persen < 0 || persen > DiskonMaks {
		return 0, fmt.Errorf("diskon %d%% di luar batas 0-%d%%", persen, DiskonMaks)
	}
	return total - total*persen/100, nil
}
```

```go
// main.go
package main

import (
	"fmt"

	"tokoku/diskon"
	"tokoku/internal/harga"
)

func main() {
	total := harga.Total([]int{10000, 25000})
	akhir, err := diskon.Terapkan(total, 10)
	fmt.Println(akhir, err)

	_, err = diskon.Terapkan(total, 80)
	fmt.Println(err)
}
```

```bash
go run .
go vet ./...
go build ./...
```

### Menguji Aturan `internal/`

Buat modul lain di luar `tokoku` yang mencoba `import "tokoku/internal/harga"` — compiler akan menolak dengan pesan *use of internal package ... not allowed*. Inilah penegakan enkapsulasi oleh toolchain.

### `replace` Lokal dan Workspace

```bash
cd ..
mkdir salam && cd salam
go mod init contoh.com/salam
cat > salam.go <<'EOF'
package salam

func Halo(nama string) string { return "Halo, " + nama + "!" }
EOF

cd ../tokoku
go mod edit -require=contoh.com/salam@v0.0.0 -replace=contoh.com/salam=../salam
go mod tidy
cat go.mod
```

Tambahkan `import "contoh.com/salam"` di `main.go` lalu panggil `salam.Halo("Andi")`. Perhatikan baris `replace` di `go.mod`. Alternatifnya, pakai **workspace** agar `go.mod` tidak berubah:

```bash
cd ..
go work init ./tokoku ./salam
cat go.work
```

### Vendoring dan Inspeksi Dependensi

```bash
cd tokoku
go list -m all
go mod graph
go mod verify

go mod vendor
ls vendor
go build -mod=vendor ./...
```

Setelah `vendor/` ada, `go list -m all` menolak berjalan (menggunakan direktori vendor) kecuali ditambah `-mod=mod`, jadi lakukan inspeksi dependensi **sebelum** vendoring.

## Latihan

1. Buat modul `tokoku` seperti di atas, tambahkan package ketiga `laporan` yang mengimpor `harga` dan `diskon`, lalu cetak laporan dari `main`. Jalankan `go vet ./...` dan `go build ./...`.
2. Sengaja buat **import melingkar** antara dua package (A mengimpor B dan B mengimpor A) dan catat pesan error-nya. Lalu perbaiki dengan memindahkan tipe bersama ke package ketiga atau memakai interface.
3. Buktikan aturan `internal/`: dari modul terpisah, impor package di dalam `internal/` modul `tokoku`. Jelaskan manfaat aturan ini bagi penulis library dibanding sekadar mengandalkan konvensi.
4. Tambahkan dependensi pihak ketiga (misalnya `github.com/google/uuid`) dengan `go get`, lihat perubahan `go.mod` dan `go.sum`, lalu jalankan `go mod why github.com/google/uuid` dan `go mod tidy` setelah menghapus pemakaiannya dari kode.
5. Jalankan `go list -m -u all` pada sebuah proyek, naikkan satu dependensi dengan `go get nama@versi`, lalu jelaskan bagaimana **Minimal Version Selection** memilih versi akhir bila dua dependensi meminta versi berbeda.
6. Terbitkan modul percobaan ke repositori Git publik dengan tag `v0.1.0`, lalu `go get` dari proyek lain. Jelaskan apa yang harus diubah jika Anda merilis perubahan API yang merusak (petunjuk: `/v2`).
7. Diskusikan: sebuah tim menaruh seluruh kode di direktori `models/`, `controllers/`, dan `utils/` lalu mulai mengalami import melingkar. Jelaskan akar masalahnya dan usulkan struktur berbasis domain beserta peran `cmd/` dan `internal/` ([Modul 16](16-cli-web-development.md) membahas aplikasi web).

---

[← Modul 10: Error Handling](10-error-handling.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 12: Goroutines & Channels →](12-goroutines-channels.md)
