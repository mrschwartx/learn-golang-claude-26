# Modul 18: Go Toolchain, Code Quality & Security

[← Modul 17: gRPC, Database & Logging](17-grpc-database-logging.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 19: Performance, Debugging & Deployment →](19-performance-deployment.md)

**Isi modul ini**

- [Teori](#teori)
  - [18.1 Core Go Commands](#181-core-go-commands)
  - [18.2 Cache dan Build Berulang](#182-cache-dan-build-berulang)
  - [18.3 Code Generation: `go generate`](#183-code-generation-go-generate)
  - [18.4 Build Tags dan Build Constraints](#184-build-tags-dan-build-constraints)
  - [18.5 Kualitas Kode: `gofmt`, `goimports`, `go vet`, `go fix`](#185-kualitas-kode-gofmt-goimports-go-vet-go-fix)
  - [18.6 Linters](#186-linters)
  - [18.7 Keamanan](#187-keamanan)
  - [18.8 Pipeline CI](#188-pipeline-ci)
- [Praktik](#praktik)
  - [`go vet` Menemukan Bug](#go-vet-menemukan-bug)
  - [`go generate` dengan Generator Sendiri](#go-generate-dengan-generator-sendiri)
  - [Build Tags dan Berkas per Platform](#build-tags-dan-berkas-per-platform)
  - [Konfigurasi `golangci-lint` dan `govulncheck`](#konfigurasi-golangci-lint-dan-govulncheck)
- [Latihan](#latihan)

---


## Teori

### 18.1 Core Go Commands

Seluruh siklus pengembangan Go dikendalikan satu perintah: **`go`**. Rangkuman perintah inti:

| Perintah | Fungsi |
|---|---|
| `go run` | Kompilasi sementara dan jalankan |
| `go build` | Kompilasi menjadi binary (`-o nama` menentukan nama keluaran) |
| `go install` | Kompilasi dan pasang binary ke `$GOBIN` |
| `go test` | Menjalankan uji, benchmark, fuzzing ([Modul 15](15-testing-benchmarking.md)) |
| `go fmt` | Memformat kode (membungkus `gofmt`) |
| `go vet` | Analisis statis untuk kesalahan umum |
| `go mod` | Mengelola `go.mod` ([Modul 11](11-modules-packages.md)) |
| `go get` | Menambah/mengubah dependensi dan tool |
| `go list` | Mencetak informasi package/modul (`go list -m all`, `go list ./...`) |
| `go doc` | Dokumentasi di terminal ([Modul 2](02-instalasi-setup.md)) |
| `go env` | Menampilkan/mengatur environment Go |
| `go clean` | Menghapus berkas hasil build; `go clean -cache` / `-modcache` membersihkan cache |
| `go version` | Versi Go; `go version -m ./app` menampilkan **modul dan versi di dalam sebuah binary** |
| `go generate` | Menjalankan generator kode (18.3) |
| `go fix` | Memperbarui kode ke idiom terbaru (18.5) |
| `go tool` | Menjalankan tool bawaan (`pprof`, `cover`, `trace`) dan tool dari direktif `tool` di `go.mod` |
| `go work` | Mengelola workspace multi-modul ([Modul 11](11-modules-packages.md)) |

Wildcard **`./...`** berarti "package di direktori ini dan seluruh subdirektorinya", dipakai di hampir semua perintah (`go test ./...`, `go vet ./...`).

Flag build yang paling sering dipakai:

| Flag | Fungsi |
|---|---|
| `-o nama` | Nama berkas keluaran |
| `-race` | Aktifkan race detector ([Modul 13](13-sync-context-patterns.md)) |
| `-tags a,b` | Aktifkan *build tag* (18.4) |
| `-ldflags "..."` | Flag untuk linker, misalnya menyematkan versi ([Modul 19](19-performance-deployment.md)) |
| `-trimpath` | Menghapus path lokal dari binary (build lebih reproducible) |
| `-v` | Menampilkan package yang dikompilasi |
| `-gcflags "..."` | Flag untuk compiler ([Modul 20](20-advanced-topics.md)) |

### 18.2 Cache dan Build Berulang

Go menyimpan hasil kompilasi dan hasil uji di **build cache**, sehingga build dan test berikutnya sangat cepat bila kode tidak berubah. Hasil uji yang di-cache ditandai `(cached)`; paksa jalan ulang dengan `go test -count=1`. Lokasi cache: `go env GOCACHE`. Di CI, simpan `GOCACHE` dan `GOMODCACHE` antar-run untuk mempercepat pipeline.

### 18.3 Code Generation: `go generate`

Banyak kode Go **dihasilkan oleh program**: stub gRPC dari `.proto`, *mock*, query dari SQL (sqlc), atau `String()` untuk enum (`stringer`). Direktif **`//go:generate`** mendokumentasikan perintah generator langsung di dekat kode yang membutuhkannya:

```go
//go:generate go run ./gen
//go:generate stringer -type=Level
```

| Aspek | Penjelasan |
|---|---|
| Dijalankan dengan `go generate ./...` | **Tidak** otomatis ikut `go build` atau `go test` |
| Hasil di-commit ke repositori | Pengguna paket tidak perlu memasang generator |
| Penanda berkas hasil generate | Baris `// Code generated ... DO NOT EDIT.` membuat linter melewatinya |
| Kunci versi tool | Gunakan direktif `tool` di `go.mod` ([Modul 11](11-modules-packages.md)) agar semua orang memakai generator yang sama |

### 18.4 Build Tags dan Build Constraints

**Build constraint** menentukan berkas mana yang ikut dikompilasi, lewat dua mekanisme:

| Mekanisme | Contoh | Arti |
|---|---|---|
| **Akhiran nama berkas** | `disk_linux.go`, `disk_windows.go`, `x_amd64.go` | Hanya untuk OS/arsitektur tersebut |
| **Direktif `//go:build`** | `//go:build linux && amd64` | Ekspresi boolean (`&&`, `\|\|`, `!`) atas tag; ditulis **di atas** `package`, diikuti baris kosong |

Tag bawaan: nama OS (`linux`, `darwin`, `windows`), arsitektur (`amd64`, `arm64`), `cgo`, dan versi Go (`go1.27`). Anda juga bisa membuat **tag sendiri** dan mengaktifkannya dengan `-tags`:

```go
//go:build integration

package db_test
```

```bash
go test ./...                      # uji integrasi dilewati
go test -tags integration ./...    # uji integrasi ikut berjalan
```

Kegunaan umum: kode khusus platform, memisahkan **uji integrasi** dari uji unit, dan fitur opsional (misalnya build "enterprise" vs "community").

### 18.5 Kualitas Kode: `gofmt`, `goimports`, `go vet`, `go fix`

| Tool | Fungsi |
|---|---|
| **`gofmt` / `go fmt`** | Format standar tunggal; hampir tanpa opsi. `gofmt -l .` mendaftar berkas yang menyimpang, `gofmt -s` menyederhanakan kode |
| **`goimports`** | Seperti `gofmt` plus **menambah/menghapus import otomatis** dan merapikan pengelompokannya |
| **`go vet`** | Analisis statis resmi: format `Printf` yang salah, mutex yang tersalin, `context` yang tidak di-cancel, struct tag keliru, kode tak terjangkau, dan lainnya. **Jalankan di CI** |
| **`go fix`** | Sejak **Go 1.26** diperbarui total menjadi rumah bagi *modernizer*: menulis ulang kode lama ke idiom dan API baru (misalnya memakai `slices`, `min`/`max`, `range` atas integer) tanpa mengubah perilaku. `go fix -diff ./...` menampilkan perubahan tanpa menerapkannya |

`go test` otomatis menjalankan subset pemeriksaan `vet`, sehingga banyak kesalahan terdeteksi sebelum uji dijalankan.

### 18.6 Linters

Linter menangkap masalah yang tidak ditolak compiler:

| Linter | Fokus |
|---|---|
| **staticcheck** | Analisis statis mendalam: bug, API usang, penyederhanaan, kinerja. Banyak dianggap linter paling berharga |
| **revive** | Pengganti `golint` — aturan gaya dan konvensi (nama exported tanpa komentar, nama receiver, dll.) yang bisa dikonfigurasi |
| **errcheck** | Menemukan error yang dikembalikan tetapi **tidak diperiksa** |
| **gosec** | Pola keamanan: kriptografi lemah, SQL string, izin berkas longgar |
| **golangci-lint** | **Pelari (runner)** yang menjalankan puluhan linter secara paralel dengan satu konfigurasi (`.golangci.yml`) dan cache; standar de facto di CI |

Mulai dengan set kecil yang relevan, lalu tambah bertahap — mengaktifkan semua linter sekaligus pada kode lama menghasilkan ribuan temuan yang diabaikan orang. Terapkan linter **hanya pada kode baru** (opsi `new-from-rev`) bila basis kode sudah besar.

### 18.7 Keamanan

| Praktik | Penjelasan |
|---|---|
| **`govulncheck`** | Tool resmi pemindai kerentanan. Berbeda dari pemindai biasa, ia memakai **analisis keterjangkauan**: hanya melaporkan kerentanan pada fungsi yang **benar-benar dipanggil** kode Anda, memakai database `vuln.go.dev`. Bisa memindai source maupun binary |
| **`go mod verify`** | Memastikan isi cache modul sesuai `go.sum` |
| **Checksum database** | `sum.golang.org` mencegah versi modul diubah diam-diam setelah dirilis ([Modul 11](11-modules-packages.md)) |
| **Perbarui Go dan dependensi** | Rilis patch Go sering memuat perbaikan keamanan standard library; hanya dua rilis mayor terakhir yang didukung ([Modul 1](01-pengenalan.md)) |
| **Kriptografi** | Pakai `crypto/rand` (bukan `math/rand`) untuk token dan kunci; jangan menulis algoritma sendiri; hindari `InsecureSkipVerify: true`. Go 1.24 memperkenalkan **Go Cryptographic Module** untuk kebutuhan kepatuhan FIPS 140-3 |
| **Validasi input** | Parameter terikat untuk SQL ([Modul 17](17-grpc-database-logging.md)), `html/template` untuk HTML ([Modul 16](16-cli-web-development.md)), `os.Root` untuk menahan path traversal ([Modul 14](14-standard-library.md)) |
| **Rahasia** | Jangan simpan di kode/repositori; baca dari environment atau secret manager |
| **Race detector dan fuzzing** | Menemukan kelas bug yang sering menjadi celah keamanan ([Modul 13](13-sync-context-patterns.md), [Modul 15](15-testing-benchmarking.md)) |

### 18.8 Pipeline CI

Pipeline minimum untuk proyek Go: **format → vet/lint → test (dengan race) → vulnerability scan → build**. Contohnya dengan GitHub Actions:

```yaml
# .github/workflows/ci.yml
name: ci
on: [push, pull_request]

jobs:
  uji:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v5
      - uses: actions/setup-go@v6
        with:
          go-version-file: go.mod
      - name: Format
        run: test -z "$(gofmt -l .)"
      - name: Vet
        run: go vet ./...
      - name: Lint
        uses: golangci/golangci-lint-action@v8
      - name: Test
        run: go test -race -cover ./...
      - name: Vulnerability scan
        run: go run golang.org/x/vuln/cmd/govulncheck@latest ./...
      - name: Build
        run: go build ./...
```

Versi action di atas hanyalah contoh; gunakan versi mayor terbaru yang tersedia saat Anda menulis pipeline.

## Praktik

### `go vet` Menemukan Bug

```go
// vet_demo.go
package main

import (
	"fmt"
	"sync"
)

type Penghitung struct {
	mu sync.Mutex
	n  int
}

func salin(p Penghitung) int { return p.n } // menyalin mutex

func main() {
	fmt.Printf("jumlah: %d\n", "teks") // format salah
	fmt.Println(salin(Penghitung{}))
}
```

```bash
go vet ./...
```

`go vet` melaporkan dua temuan — argumen `%d` yang diberi string, dan struct berisi mutex yang dikirim sebagai nilai. Perbaiki keduanya, lalu jalankan ulang sampai bersih.

### `go generate` dengan Generator Sendiri

```go
// gen/main.go
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	nama := []string{"Debug", "Info", "Warn", "Error"}
	var sb strings.Builder
	sb.WriteString("// Code generated by gen; DO NOT EDIT.\n\npackage main\n\n")
	sb.WriteString("var namaLevel = map[int]string{\n")
	for i, n := range nama {
		fmt.Fprintf(&sb, "\t%d: %q,\n", i, n)
	}
	sb.WriteString("}\n")
	if err := os.WriteFile("level_gen.go", []byte(sb.String()), 0o644); err != nil {
		panic(err)
	}
}
```

```go
// main.go
package main

import "fmt"

//go:generate go run ./gen

func main() {
	fmt.Println(namaLevel[2])
}
```

```bash
mkdir genny && cd genny
go mod init genny
mkdir gen
# simpan gen/main.go dan main.go di atas, lalu:
go generate ./...
cat level_gen.go
go run .
```

### Build Tags dan Berkas per Platform

```go
// info_linux.go
package main

func platform() string { return "Linux" }
```

```go
// info_other.go
//go:build !linux

package main

func platform() string { return "bukan Linux" }
```

```go
// info.go
package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println(platform(), runtime.GOOS, runtime.GOARCH)
}
```

```bash
go run .
GOOS=windows go build -o app.exe .     # memakai info_other.go
GOOS=darwin  go vet .                  # memakai info_other.go
```

### Konfigurasi `golangci-lint` dan `govulncheck`

```bash
# pasang lewat go install (atau lihat panduan resmi golangci-lint untuk cara lain)
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
```

```yaml
# .golangci.yml
version: "2"
linters:
  default: standard   # errcheck, govet, ineffassign, staticcheck, unused
  enable:
    - revive
    - gosec
formatters:
  enable:
    - gofmt
    - goimports
```

```bash
golangci-lint run ./...
govulncheck ./...
go build -o app . && govulncheck -mode=binary ./app
go version -m ./app
go fix -diff ./...
```

## Latihan

1. Jalankan `go vet` pada `vet_demo.go`, catat dua temuannya, lalu perbaiki keduanya (kirim `*Penghitung` dan ganti `%d` dengan `%s`). Mengapa menyalin `sync.Mutex` berbahaya?
2. Tulis generator `go generate` yang membaca daftar nama dari berkas `.txt` dan menghasilkan konstanta serta fungsi `String()` untuk enum. Mengapa berkas hasil generate sebaiknya di-commit ke repositori?
3. Buat tag `integration` pada sebuah uji, dan buktikan uji itu hanya berjalan dengan `-tags integration`. Kapan pemisahan unit vs integrasi seperti ini penting di CI?
4. Buat dua berkas `os_linux.go` dan `os_windows.go` yang mengembalikan separator path/nama OS, lalu build untuk dua target dengan `GOOS=...`. Apa bedanya akhiran nama berkas dengan direktif `//go:build`?
5. Pasang `golangci-lint`, jalankan pada proyek Anda dengan konfigurasi di atas, dan perbaiki minimal tiga temuan. Aktifkan `gosec`, lalu jelaskan satu temuan yang Anda anggap valid dan satu yang merupakan false positive.
6. Bangun binary dari proyek yang memakai dependensi pihak ketiga, lalu jalankan `govulncheck -mode=binary` dan `go version -m`. Jelaskan perbedaan pemindaian berbasis **keterjangkauan** dengan pemindai yang hanya mencocokkan nama dan versi paket.
7. Diskusikan: tim Anda ingin menerapkan aturan "tidak ada merge tanpa lint dan uji lulus". Rancang pipeline CI lengkap (format, vet, lint, `test -race`, `govulncheck`, build, dan push image Docker), tentukan mana yang memblokir merge dan mana yang hanya peringatan, serta cara menjaga pipeline tetap cepat.

---

[← Modul 17: gRPC, Database & Logging](17-grpc-database-logging.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 19: Performance, Debugging & Deployment →](19-performance-deployment.md)
