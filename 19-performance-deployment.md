# Modul 19: Performance, Debugging & Deployment

[← Modul 18: Go Toolchain, Code Quality & Security](18-toolchain-code-quality.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 20: Advanced Topics →](20-advanced-topics.md)

**Isi modul ini**

- [Teori](#teori)
  - [19.1 Prinsip Optimasi](#191-prinsip-optimasi)
  - [19.2 Profiling dengan `pprof`](#192-profiling-dengan-pprof)
  - [19.3 `runtime/trace`](#193-runtimetrace)
  - [19.4 Race Detector](#194-race-detector)
  - [19.5 Teknik Optimasi yang Efektif](#195-teknik-optimasi-yang-efektif)
  - [19.6 Debugging](#196-debugging)
  - [19.7 Cross-Compilation](#197-cross-compilation)
  - [19.8 Membangun Executable untuk Distribusi](#198-membangun-executable-untuk-distribusi)
  - [19.9 Deployment dengan Docker](#199-deployment-dengan-docker)
  - [19.10 Deployment ke Kubernetes](#1910-deployment-ke-kubernetes)
- [Praktik](#praktik)
  - [Profiling CPU dengan `runtime/pprof`](#profiling-cpu-dengan-runtimepprof)
  - [Endpoint pprof pada Layanan](#endpoint-pprof-pada-layanan)
  - [Menyematkan Versi dan Info Build](#menyematkan-versi-dan-info-build)
  - [Cross-Compilation](#cross-compilation)
  - [Image Docker Multi-Stage](#image-docker-multi-stage)
  - [Manifest Kubernetes Minimal](#manifest-kubernetes-minimal)
- [Latihan](#latihan)

---


## Teori

### 19.1 Prinsip Optimasi

Aturan pertama optimasi: **ukur, jangan menebak**. Urutan kerja yang sehat:

1. Tentukan **target** yang jelas (latensi p99, throughput, memori).
2. **Ukur** dengan benchmark ([Modul 15](15-testing-benchmarking.md)) atau profiling di lingkungan realistis.
3. Temukan **bottleneck** sebenarnya dari profil.
4. Ubah **satu hal**, ukur ulang, bandingkan (`benchstat`).
5. Berhenti saat target tercapai — optimasi berlebihan merusak keterbacaan.

Go menyediakan seluruh alat ini di dalam toolchain: **pprof** untuk profiling, **trace** untuk melihat perilaku runtime, dan **race detector** untuk bug konkurensi.

### 19.2 Profiling dengan `pprof`

**pprof** mengambil sampel eksekusi program dan menunjukkan di mana waktu dan memori dihabiskan.

| Profil | Menjawab pertanyaan |
|---|---|
| **CPU** | Fungsi mana yang memakai waktu CPU paling banyak? |
| **Heap** | Siapa yang mengalokasikan memori, dan apa yang masih hidup? |
| **Allocs** | Siapa yang paling sering mengalokasikan (sumber tekanan GC)? |
| **Goroutine** | Berapa banyak goroutine dan sedang menunggu apa? |
| **Block** / **Mutex** | Di mana goroutine terblokir menunggu channel/mutex? |
| **Goroutine leak** | Goroutine mana yang **terblokir permanen** (profil `goroutineleak`; tersedia umum di Go 1.27) |

Cara mendapatkan profil:

| Cara | Penjelasan |
|---|---|
| **Layanan berjalan** | `import _ "net/http/pprof"` mendaftarkan endpoint `/debug/pprof/*` pada `http.DefaultServeMux`. **Jangan buka ke publik** — jalankan di port internal terpisah |
| **Program batch** | `runtime/pprof`: `pprof.StartCPUProfile(f)` lalu `defer pprof.StopCPUProfile()`; `pprof.WriteHeapProfile` |
| **Saat benchmark** | `go test -bench . -cpuprofile cpu.out -memprofile mem.out` |

Analisis dengan:

```bash
go tool pprof -http=:8080 cpu.out     # UI web; tampilan default: flame graph (Go 1.26+)
go tool pprof -top cpu.out            # ringkasan teks
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```

Cara membaca: **flat** = waktu di fungsi itu sendiri; **cum** (kumulatif) = fungsi itu beserta yang dipanggilnya. Fungsi dengan `flat` tinggi adalah kandidat optimasi langsung.

### 19.3 `runtime/trace`

Jika pprof menunjukkan **di mana** waktu habis, **execution trace** menunjukkan **kapan dan mengapa**: goroutine yang berjalan atau menunggu, aktivitas GC, penjadwalan di tiap processor, serta jeda akibat blocking.

```bash
go test -trace trace.out ./...
go tool trace trace.out
```

`runtime/trace` juga menyediakan `trace.NewTask` dan `trace.WithRegion` untuk menandai bagian kode sendiri. Go 1.25 menambahkan **flight recorder** — menyimpan jejak beberapa detik terakhir di memori dan menuliskannya hanya saat terjadi sesuatu yang aneh, sehingga cocok untuk produksi. Trace sangat membantu mendiagnosis latensi yang tinggi tanpa CPU yang tinggi.

### 19.4 Race Detector

Race detector ([Modul 13](13-sync-context-patterns.md)) tetap bagian dari alat debugging utama. Ingat tiga hal: aktifkan dengan `-race` pada `go test` dan `go run`; ia hanya mendeteksi race yang **terjadi** selama eksekusi (kualitas uji menentukan); dan biayanya cukup besar sehingga dipakai di CI, bukan di produksi. Untuk build staging yang sengaja diinstrumentasi, `go build -race` juga tersedia.

### 19.5 Teknik Optimasi yang Efektif

| Teknik | Penjelasan |
|---|---|
| **Kurangi alokasi** | Alokasi memperbanyak kerja GC. Alokasikan kapasitas slice/map di awal (`make([]T, 0, n)`), gunakan `strings.Builder`, hindari konversi `string` ↔ `[]byte` berulang. Ukur dengan `-benchmem` |
| **`sync.Pool`** | Menggunakan ulang buffer sementara pada jalur panas ([Modul 13](13-sync-context-patterns.md)) |
| **Pilih struktur data tepat** | `map[T]struct{}` untuk set, slice terurut + pencarian biner bila memori ketat |
| **Kurangi kontensi lock** | Persempit area kritis, pecah (*shard*) state, atau pakai `atomic` |
| **Batching dan buffering** | `bufio`, penulisan bulk ke database (`CopyFrom`, [Modul 17](17-grpc-database-logging.md)) |
| **Tuning GC** | `GOGC` dan `GOMEMLIMIT` ([Modul 7](07-pointer-memory.md)); Green Tea GC (default Go 1.26) dan alokasi berukuran-khusus (Go 1.27) mengurangi biaya GC tanpa perubahan kode |
| **Profile-Guided Optimization (PGO)** | Sertakan profil CPU produksi sebagai `default.pgo` di direktori package `main`; compiler memakainya untuk inlining dan optimasi lain, dengan peningkatan kisaran beberapa hingga belasan persen pada banyak beban kerja |
| **`GOMAXPROCS` di container** | Sejak Go 1.25 runtime membaca **batas CPU cgroup** otomatis, sehingga tidak perlu lagi library tambahan seperti `automaxprocs` |

### 19.6 Debugging

| Alat | Kegunaan |
|---|---|
| **Delve** (`dlv`) | Debugger khusus Go: breakpoint, step, inspeksi variabel dan goroutine. `dlv debug`, `dlv test`, `dlv attach <pid>`, `dlv connect` untuk debug jarak jauh |
| Stack trace panic | Cetak otomatis; atur rincian dengan `GOTRACEBACK=all` |
| **Sinyal `SIGQUIT`** | Pada program Go yang berjalan, `kill -QUIT <pid>` (atau Ctrl+\) mencetak **dump seluruh goroutine** lalu keluar — penyelamat saat program macet |
| `GODEBUG` | Variabel untuk perilaku runtime: `gctrace=1`, `schedtrace=1000` |
| `runtime.NumGoroutine()` | Indikator kasar goroutine leak |
| Logging terstruktur | `slog` dengan request ID ([Modul 14](14-standard-library.md)) |

### 19.7 Cross-Compilation

Go dapat membangun binary untuk **OS dan arsitektur lain dari satu mesin** hanya dengan mengubah dua environment variable:

```bash
GOOS=linux   GOARCH=amd64 go build -o app-linux-amd64 .
GOOS=linux   GOARCH=arm64 go build -o app-linux-arm64 .
GOOS=windows GOARCH=amd64 go build -o app.exe .
GOOS=darwin  GOARCH=arm64 go build -o app-mac .
```

| Variabel | Arti |
|---|---|
| `GOOS` | Sistem operasi target (`linux`, `windows`, `darwin`, ...) |
| `GOARCH` | Arsitektur (`amd64`, `arm64`, `arm`, `386`, `riscv64`, `wasm`, ...) |
| `CGO_ENABLED=0` | Menonaktifkan cgo agar hasilnya **statis murni** — **wajib** untuk cross-compile sederhana dan image `scratch`/distroless ([Modul 20](20-advanced-topics.md)) |

`go tool dist list` mencantumkan seluruh pasangan target yang didukung. Cross-compile menjadi rumit hanya bila kode memakai **cgo**, karena membutuhkan toolchain C untuk target tersebut.

### 19.8 Membangun Executable untuk Distribusi

| Teknik | Penjelasan |
|---|---|
| **Strip simbol** | `-ldflags="-s -w"` menghapus tabel simbol dan info debug — binary lebih kecil (pada contoh Praktik sekitar sepertiganya), tetapi sulit di-debug |
| **Sematkan versi** | `-ldflags="-X main.version=1.2.3"` mengisi variabel `string` saat build |
| **`-trimpath`** | Menghapus path mesin build dari binary → lebih reproducible dan tidak membocorkan path lokal |
| **Info build bawaan** | `runtime/debug.ReadBuildInfo()` membaca versi Go, modul dependensi, dan info VCS (commit) yang otomatis tertanam; `go version -m app` menampilkannya dari luar |
| **Static binary** | `CGO_ENABLED=0` menghasilkan satu berkas tanpa dependensi library sistem |
| **Rilis banyak target** | Tool seperti **GoReleaser** mengotomatiskan build multi-platform, arsip, checksum, dan publikasi rilis |

### 19.9 Deployment dengan Docker

Karena Go menghasilkan binary statis, image Docker bisa **sangat kecil** dengan **multi-stage build** (course Docker): tahap pertama mengompilasi, tahap kedua hanya membawa binary.

```dockerfile
# syntax=docker/dockerfile:1
FROM golang:1.27 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/app /app
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app"]
```

| Pilihan image akhir | Karakteristik |
|---|---|
| `scratch` | Kosong sama sekali (~ukuran binary); tanpa shell, tanpa sertifikat CA kecuali disalin |
| **distroless/static** | Berisi sertifikat CA, zona waktu, dan user non-root; tanpa shell — pilihan aman dan umum |
| `alpine` | Kecil dan punya shell untuk debugging; hati-hati pada perbedaan libc jika memakai cgo |

Praktik penting: salin `go.mod`/`go.sum` dan jalankan `go mod download` **sebelum** menyalin kode agar layer dependensi ter-cache; jalankan sebagai **non-root**; buat file `.dockerignore`.

### 19.10 Deployment ke Kubernetes

Aplikasi Go cocok dijalankan di Kubernetes (course Kubernetes dan EKS). Hal yang perlu diperhatikan:

| Aspek | Penerapan |
|---|---|
| **Health check** | Sediakan endpoint `/healthz` (liveness) dan `/readyz` (readiness) |
| **Graceful shutdown** | Tangani `SIGTERM` dan panggil `srv.Shutdown(ctx)` ([Modul 13](13-sync-context-patterns.md)); samakan dengan `terminationGracePeriodSeconds` |
| **Resource limit** | Atur `requests`/`limits`; sejak Go 1.25 `GOMAXPROCS` mengikuti limit CPU, dan setel **`GOMEMLIMIT`** sedikit di bawah limit memori (misalnya 80–90%) agar tidak di-*OOMKill* |
| **Konfigurasi** | Lewat environment variable, `ConfigMap`, dan `Secret` |
| **Observabilitas** | Log JSON ke stdout ([Modul 14](14-standard-library.md)), metrik Prometheus (`prometheus/client_golang`), tracing OpenTelemetry |
| **Image** | Tag tetap berversi (jangan `latest`), jalankan non-root, filesystem read-only bila memungkinkan |

## Praktik

### Profiling CPU dengan `runtime/pprof`

```go
// profil.go
package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"strings"
)

func fib(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func bangunTeks(n int) string {
	hasil := ""
	for i := 0; i < n; i++ {
		hasil += "x" // O(n^2): alokasi baru tiap iterasi
	}
	return hasil
}

func bangunTeksBaik(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte('x')
	}
	return sb.String()
}

func main() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	fmt.Println(fib(38))
	fmt.Println(len(bangunTeks(200000)))
	fmt.Println(len(bangunTeksBaik(200000)))
}
```

```bash
go run profil.go
go tool pprof -top cpu.prof
go tool pprof -http=:8080 cpu.prof     # buka flame graph di browser
```

Pada keluaran `-top`, lihat kolom `flat` dan `cum`. Di mesin penulis, baris teratas justru `runtime.memmove` (±37%) disusul fungsi pemindai GC seperti `runtime.scanObject` — jejak penyambungan string `+=` di `bangunTeks` yang menyalin ulang dan menekan GC — sedangkan `main.fib` hanya beberapa persen. Angka persisnya berbeda di tiap mesin, tetapi pelajarannya sama: **profil sering membantah dugaan Anda** tentang bagian yang lambat. Optimasi bagian yang benar-benar teratas (ganti `+=` dengan `strings.Builder`), lalu ukur ulang.

### Endpoint pprof pada Layanan

```go
// pprof_server.go
package main

import (
	"log"
	"net/http"
	_ "net/http/pprof" // mendaftarkan /debug/pprof/* pada DefaultServeMux
)

func main() {
	// Port internal terpisah, JANGAN diekspos ke publik.
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	select {} // ganti dengan server aplikasi Anda
}
```

```bash
go run pprof_server.go &
go tool pprof http://localhost:6060/debug/pprof/heap
curl http://localhost:6060/debug/pprof/goroutine?debug=1
```

### Menyematkan Versi dan Info Build

```go
// versi.go
package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

var version = "dev" // diisi saat build lewat -ldflags

func main() {
	fmt.Println("versi:", version, "| go:", runtime.Version(), "|", runtime.GOOS+"/"+runtime.GOARCH)
	if bi, ok := debug.ReadBuildInfo(); ok {
		fmt.Println("modul:", bi.Main.Path)
	}
}
```

```bash
go run versi.go
go build -o app-biasa versi.go
go build -ldflags="-s -w -X main.version=1.2.3" -trimpath -o app-rilis versi.go
./app-rilis
ls -lh app-biasa app-rilis        # bandingkan ukuran
go version -m app-rilis
```

### Cross-Compilation

```bash
GOOS=linux   GOARCH=amd64 go build -o dist/app-linux-amd64 versi.go
GOOS=linux   GOARCH=arm64 go build -o dist/app-linux-arm64 versi.go
GOOS=windows GOARCH=amd64 go build -o dist/app-windows.exe versi.go
GOOS=darwin  GOARCH=arm64 go build -o dist/app-darwin-arm64 versi.go
file dist/*
go tool dist list | head
```

### Image Docker Multi-Stage

Simpan `Dockerfile` dari 19.9 di akar proyek (dengan `cmd/api/main.go`), lalu:

```bash
docker build -t toko-api:1.0 .
docker images toko-api
docker run --rm -p 8080:8080 toko-api:1.0
```

Bandingkan ukuran image akhir dengan membangun satu stage saja memakai `golang:1.27` — selisihnya biasanya ratusan megabyte.

### Manifest Kubernetes Minimal

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: toko-api
spec:
  replicas: 2
  selector:
    matchLabels: { app: toko-api }
  template:
    metadata:
      labels: { app: toko-api }
    spec:
      terminationGracePeriodSeconds: 30
      containers:
        - name: api
          image: toko-api:1.0
          ports:
            - containerPort: 8080
          env:
            - name: GOMEMLIMIT
              value: "400MiB"          # sekitar 80% dari limit memori
          resources:
            requests: { cpu: 100m, memory: 128Mi }
            limits:   { cpu: 500m, memory: 512Mi }
          readinessProbe:
            httpGet: { path: /readyz, port: 8080 }
          livenessProbe:
            httpGet: { path: /healthz, port: 8080 }
```

## Latihan

1. Jalankan `profil.go`, buka `cpu.prof` dengan `go tool pprof -top` dan `-http`, lalu identifikasi tiga fungsi dengan `flat` tertinggi. Optimasi `bangunTeks` dan buktikan peningkatannya dengan benchmark `-benchmem` ([Modul 15](15-testing-benchmarking.md)).
2. Tambahkan endpoint `net/http/pprof` ke REST API di [Modul 16](16-cli-web-development.md), beri beban dengan `hey` atau loop `curl`, lalu ambil profil heap dan CPU 15 detik. Apa perbedaan profil `heap` dan `allocs`?
3. Tulis program yang sengaja membocorkan goroutine (menunggu channel yang tak pernah terisi), lalu deteksi lewat `runtime.NumGoroutine()` dan profil `goroutine` (`/debug/pprof/goroutine?debug=2`). Perbaiki dengan `context`.
4. Ambil execution trace (`runtime/trace` atau `go test -trace`) dari program dengan worker pool ([Modul 12](12-goroutines-channels.md)), buka dengan `go tool trace`, dan amati penjadwalan goroutine antar-processor serta jeda GC.
5. Bangun aplikasi untuk lima target berbeda lewat `GOOS`/`GOARCH`, periksa dengan `file`, dan jelaskan mengapa `CGO_ENABLED=0` memudahkan cross-compilation dan image `scratch`.
6. Tulis Dockerfile multi-stage untuk aplikasi Go Anda, bandingkan ukuran image antara `golang:1.27`, `alpine`, dan `distroless/static`, lalu jalankan sebagai non-root. Apa trade-off `scratch` dibanding distroless?
7. Diskusikan: layanan Go di Kubernetes dengan limit memori 512 MiB mengalami restart berkala (`OOMKilled`) padahal profil heap terlihat wajar. Rancang langkah diagnosis (profil heap vs memori proses, `GOMEMLIMIT`, goroutine leak, buffer besar) serta perbaikan, dan jelaskan peran `GOMAXPROCS` yang kini sadar-container di Go 1.25.

---

[← Modul 18: Go Toolchain, Code Quality & Security](18-toolchain-code-quality.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 20: Advanced Topics →](20-advanced-topics.md)
