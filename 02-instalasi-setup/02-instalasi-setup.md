# Modul 2: Instalasi & Setup Environment

[← Modul 1: Pengenalan Go](01-pengenalan.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 3: Variabel, Konstanta & Tipe Data →](03-variabel-konstanta-tipe-data.md)

**Isi modul ini**

- [Teori](#teori)
  - [2.1 Cara Memasang Go](#21-cara-memasang-go)
  - [2.2 Variabel Lingkungan Penting](#22-variabel-lingkungan-penting)
  - [2.3 Manajemen Toolchain (`GOTOOLCHAIN`)](#23-manajemen-toolchain-gotoolchain)
  - [2.4 Anatomi Program Go: Hello World](#24-anatomi-program-go-hello-world)
  - [2.5 Perintah `go` Dasar](#25-perintah-go-dasar)
  - [2.6 `go run` vs `go build` vs `go install`](#26-go-run-vs-go-build-vs-go-install)
  - [2.7 Editor dan Tooling](#27-editor-dan-tooling)
  - [2.8 Commands & Docs — Cara Belajar Mandiri](#28-commands--docs--cara-belajar-mandiri)
- [Praktik](#praktik)
  - [Memverifikasi Instalasi](#memverifikasi-instalasi)
  - [Program Pertama sebagai Modul](#program-pertama-sebagai-modul)
  - [Menjelajahi Dokumentasi dari Terminal](#menjelajahi-dokumentasi-dari-terminal)
  - [Formatting dan Vet](#formatting-dan-vet)
- [Latihan](#latihan)

---


## Teori

### 2.1 Cara Memasang Go

Go didistribusikan sebagai satu paket berisi compiler, linker, standard library, dan seluruh tooling. Ada beberapa cara memasangnya:

| Platform | Cara Instalasi |
|---|---|
| **Semua OS** | Unduh installer resmi dari **https://go.dev/dl/** (`.pkg` macOS, `.msi` Windows, `.tar.gz` Linux) |
| **macOS** | `brew install go` |
| **Windows** | `winget install GoLang.Go` atau installer `.msi` |
| **Linux** | Ekstrak tarball resmi ke `/usr/local/go` lalu tambahkan `/usr/local/go/bin` ke `PATH` |
| **Docker** | Image resmi `golang:1.27` — tanpa memasang apa pun di mesin Anda |

> Paket `golang-go` di repositori `apt` distro Linux biasanya **tertinggal beberapa versi**. Untuk belajar dan pengembangan, lebih baik pakai tarball resmi atau Docker.

Verifikasi instalasi dengan:

```bash
go version
```

### 2.2 Variabel Lingkungan Penting

Go memakai beberapa environment variable. Lihat seluruh nilai aktif dengan `go env`, atau satu nilai dengan `go env NAMA`.

| Variabel | Fungsi |
|---|---|
| `GOROOT` | Lokasi instalasi Go itu sendiri (jarang perlu diubah manual) |
| `GOPATH` | Direktori kerja Go; default `$HOME/go`. Berisi cache modul (`pkg/mod`) dan binary hasil `go install` (`bin`) |
| `GOBIN` | Lokasi tujuan `go install`; default `$GOPATH/bin` |
| `GOPROXY` | Daftar proxy untuk mengunduh modul; default `https://proxy.golang.org,direct` |
| `GOPRIVATE` | Pola modul privat yang tidak boleh lewat proxy/checksum publik ([Modul 11](11-modules-packages.md)) |
| `GOTOOLCHAIN` | Mengatur pemilihan versi toolchain (lihat 2.3) |
| `GOOS` / `GOARCH` | Target OS dan arsitektur saat build, untuk cross-compilation ([Modul 19](19-performance-deployment.md)) |
| `GOFLAGS` | Flag default yang ditambahkan ke setiap perintah `go` |

> Tambahkan `$(go env GOPATH)/bin` ke `PATH` agar tool yang Anda pasang lewat `go install` (misalnya `staticcheck` atau `dlv`) bisa dipanggil langsung dari terminal.

Sejak ada Go Modules, **kode proyek Anda tidak lagi harus berada di dalam `GOPATH`**. Proyek bisa diletakkan di folder mana pun yang memiliki file `go.mod`.

### 2.3 Manajemen Toolchain (`GOTOOLCHAIN`)

Sejak Go 1.21, versi Go bukan lagi sekadar versi yang terpasang di mesin — ia menjadi bagian dari **modul**. File `go.mod` bisa memuat dua baris:

```
go 1.27.0
toolchain go1.27.4
```

- `go` — versi Go **minimum** yang dibutuhkan modul.
- `toolchain` — toolchain yang disarankan (opsional).

Jika Go yang terpasang lebih lama dari yang dibutuhkan `go.mod`, perintah `go` akan **mengunduh toolchain yang sesuai secara otomatis** (perilaku default `GOTOOLCHAIN=auto`) lalu menjalankannya. Anda bisa mengubah versi di `go.mod` dengan `go get go@1.27.4`.

### 2.4 Anatomi Program Go: Hello World

```go
package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
```

| Bagian | Penjelasan |
|---|---|
| `package main` | Setiap file Go dimulai dengan deklarasi paket. Paket bernama `main` menghasilkan program **executable** |
| `import "fmt"` | Memuat paket dari standard library. Import yang tidak dipakai menyebabkan **error kompilasi** |
| `func main()` | Fungsi yang dijalankan pertama kali. Tidak menerima argumen dan tidak mengembalikan nilai |
| `fmt.Println` | Fungsi dari paket `fmt`, diawali **huruf kapital** karena diekspor (*exported*) |

Aturan sintaks yang perlu diingat sejak awal: kurung kurawal pembuka **harus** berada di baris yang sama dengan `func`, `if`, atau `for` — jika dipindah ke baris baru, kode tidak akan dikompilasi.

### 2.5 Perintah `go` Dasar

| Perintah | Fungsi |
|---|---|
| `go run file.go` | Kompilasi sementara dan langsung jalankan |
| `go build` | Kompilasi menjadi file executable di direktori saat ini |
| `go install` | Kompilasi lalu pasang binary ke `$GOPATH/bin` |
| `go fmt ./...` | Format seluruh kode sesuai standar Go |
| `go vet ./...` | Analisis statis untuk menemukan kesalahan umum |
| `go mod init nama` | Membuat modul baru (membuat file `go.mod`) |
| `go doc paket.Simbol` | Menampilkan dokumentasi langsung di terminal |
| `go help` | Daftar seluruh perintah |

### 2.6 `go run` vs `go build` vs `go install`

| Aspek | `go run` | `go build` | `go install` |
|---|---|---|---|
| Menghasilkan binary permanen | Tidak (file sementara) | Ya, di direktori saat ini | Ya, di `$GOPATH/bin` (atau `GOBIN`) |
| Cocok untuk | Eksperimen cepat saat belajar | Distribusi dan deployment | Memasang tool CLI ke komputer Anda |
| Contoh | `go run .` | `go build -o app .` | `go install ./cmd/app` |

Tanda `.` berarti "paket di direktori saat ini". Untuk tool dari repositori publik, Anda bisa memasangnya langsung dengan `go install nama-modul/cmd/tool@versi`.

### 2.7 Editor dan Tooling

Go memiliki **language server resmi bernama `gopls`** yang menyediakan autocomplete, navigasi definisi, refactoring, dan diagnostik ke hampir semua editor.

| Editor | Dukungan Go |
|---|---|
| **VS Code** | Ekstensi resmi **Go** (dari tim Go) — otomatis memasang `gopls` dan `dlv` (debugger) |
| **GoLand** | IDE komersial dari JetBrains dengan fitur refactoring lengkap |
| **Neovim/Vim** | Lewat `gopls` melalui LSP |

Hampir semua editor dikonfigurasi untuk menjalankan `gofmt` (atau `goimports`, yang juga merapikan import) setiap kali Anda menyimpan file.

### 2.8 Commands & Docs — Cara Belajar Mandiri

Dokumentasi Go sangat terpusat dan berkualitas. Sumber yang perlu Anda kenali:

| Sumber | Isi |
|---|---|
| **https://go.dev/doc/** | Dokumentasi resmi, tutorial, dan panduan |
| **https://go.dev/tour/** | *A Tour of Go* — tutorial interaktif di browser |
| **https://go.dev/doc/effective_go** | *Effective Go* — idiom penulisan kode Go |
| **https://pkg.go.dev/** | Dokumentasi seluruh paket standard library dan modul publik |
| `go doc` | Dokumentasi di terminal, tanpa membuka browser |
| `go help <topik>` | Bantuan perintah, misalnya `go help build` atau `go help environment` |

Komentar dokumentasi pada kode Go (komentar tepat di atas deklarasi, diawali nama simbol) otomatis menjadi isi `go doc` dan pkg.go.dev.

## Praktik

### Memverifikasi Instalasi

```bash
go version
go env GOROOT GOPATH GOPROXY GOTOOLCHAIN
```

Alternatif **tanpa memasang Go**, memakai Docker (course Docker):

```bash
mkdir -p ~/belajar-go && cd ~/belajar-go
docker run --rm -it -v "$PWD":/app -w /app golang:1.27 bash
# di dalam container, semua perintah go tersedia
```

### Program Pertama sebagai Modul

```bash
mkdir hello && cd hello
go mod init contoh.com/hello
cat go.mod
```

Buat file `main.go`:

```go
// main.go
package main

import "fmt"

func main() {
	nama := "Go"
	fmt.Printf("Halo dari %s!\n", nama)
}
```

```bash
go run .
go build -o hello .
ls -lh hello
./hello
```

Perhatikan ukuran file `hello` — sekitar **1,5–2,5 MB** meski programnya hanya satu baris. Itu karena **runtime Go (scheduler dan garbage collector) dan paket `fmt` sudah tertanam di dalam binary**, itulah sebabnya binary ini bisa dijalankan di server tanpa memasang apa pun.

### Menjelajahi Dokumentasi dari Terminal

```bash
go doc fmt.Println
go doc strings
go doc -src strings.ToUpper
go help build
```

### Formatting dan Vet

Ubah `main.go` menjadi kode yang formatnya berantakan (misalnya indentasi dengan spasi), lalu:

```bash
gofmt -l .
go fmt ./...
go vet ./...
```

`gofmt -l .` mencantumkan file yang formatnya menyimpang; `go fmt` memperbaikinya. Go memakai **tab** untuk indentasi — ini bagian dari standar format.

## Latihan

1. Pasang Go (atau pakai Docker), lalu jalankan `go version` dan `go env GOPATH GOROOT`. Jelaskan fungsi masing-masing direktori tersebut.
2. Buat modul `contoh.com/salam` dengan `go mod init`, tulis program yang mencetak `Halo, <nama>!` dengan nama sebagai variabel, lalu jalankan dengan `go run .` dan `go build`. Berapa ukuran binary yang dihasilkan?
3. Pindahkan kurung kurawal pembuka `func main()` ke baris berikutnya lalu coba kompilasi. Catat pesan error yang muncul, dan jelaskan mengapa aturan ini ada (petunjuk: Go memasukkan titik koma secara otomatis di akhir baris tertentu).
4. Tambahkan `import "os"` tanpa memakainya. Apa yang terjadi? Jelaskan bagaimana kebijakan ini memengaruhi kebersihan kode di proyek besar.
5. Jelaskan perbedaan `go run`, `go build`, dan `go install`, lalu sebutkan situasi nyata untuk masing-masing.
6. Buka `go.mod` modul Anda, lalu ubah baris `go` ke versi yang lebih baru dari Go yang terpasang. Amati perilaku perintah `go` berikutnya (apakah toolchain baru diunduh otomatis?), lalu jelaskan manfaat fitur toolchain management bagi tim yang beranggotakan banyak developer.
7. Diskusikan: mengapa Go menyediakan `gofmt` sebagai tool tunggal dengan sedikit sekali opsi konfigurasi, alih-alih memberi tim kebebasan menentukan gaya penulisan sendiri? Apa manfaat dan kekurangannya?

---

[← Modul 1: Pengenalan Go](01-pengenalan.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 3: Variabel, Konstanta & Tipe Data →](03-variabel-konstanta-tipe-data.md)
