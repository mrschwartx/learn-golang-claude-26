# Course Go (Golang): Basic to Advance

> **Catatan sumber:** Kurikulum ini disusun mengikuti **roadmap Go dari roadmap.sh** (file `golang.pdf` yang Anda unggah): Introduction to Go, Language Basics, Functions, Pointers, Methods & Interfaces, Generics, Error Handling, Code Organization, Goroutines & Concurrency, Standard Library, Testing & Benchmarking, Ecosystem & Popular Libraries, Go Toolchain and Tools, Deployment, dan Advanced Topics. Seluruh node roadmap tersebut tercakup di modul-modul di bawah. Materi dilengkapi **dokumentasi resmi Go** (go.dev/doc dan pkg.go.dev) untuk **Go 1.27** (rilis Agustus 2026). Seluruh contoh kode yang hanya memakai standard library telah dikompilasi dan dijalankan dengan Go 1.27; contoh yang memakai library pihak ketiga (Cobra, Gin, gRPC, pgx, GORM, Zap, Zerolog) ditulis berdasarkan API stabil masing-masing dan perlu disesuaikan dengan versi yang Anda pasang.

Go melengkapi seri Anda dari sudut **bahasa pemrograman**: Docker/Kubernetes/EKS menjawab "bagaimana menjalankan aplikasi", Terraform "bagaimana menyediakan infrastrukturnya", Kafka "bagaimana data mengalir", Elasticsearch "bagaimana data dicari", Keycloak "siapa yang boleh akses", Redis "bagaimana membuat semuanya cepat" — Go menjawab **"dengan apa kita menulis layanan dan tooling yang berjalan di atas semua itu"**. Docker, Kubernetes, Terraform, dan Prometheus sendiri ditulis dalam Go, sehingga memahami Go berarti mampu membaca, memperluas, dan membangun tool di ekosistem cloud-native yang sudah Anda pelajari.

## Daftar Isi

| # | Modul |
|---|-------|
| 1 | [Pengenalan Go](01-pengenalan.md) |
| 2 | [Instalasi & Setup Environment](02-instalasi-setup.md) |
| 3 | [Variabel, Konstanta & Tipe Data](03-variabel-konstanta-tipe-data.md) |
| 4 | [Composite Types (Array, Slice, Map, Struct)](04-composite-types.md) |
| 5 | [Control Flow (Conditionals & Loops)](05-control-flow.md) |
| 6 | [Fungsi](06-fungsi.md) |
| 7 | [Pointer & Memory Management](07-pointer-memory.md) |
| 8 | [Methods & Interfaces](08-methods-interfaces.md) |
| 9 | [Generics & Iterators](09-generics.md) |
| 10 | [Error Handling](10-error-handling.md) |
| 11 | [Modules, Packages & Code Organization](11-modules-packages.md) |
| 12 | [Goroutines & Channels](12-goroutines-channels.md) |
| 13 | [Sync, Context & Concurrency Patterns](13-sync-context-patterns.md) |
| 14 | [Standard Library](14-standard-library.md) |
| 15 | [Testing & Benchmarking](15-testing-benchmarking.md) |
| 16 | [Membangun CLI & Web Development](16-cli-web-development.md) |
| 17 | [gRPC, Database & Logging](17-grpc-database-logging.md) |
| 18 | [Go Toolchain, Code Quality & Security](18-toolchain-code-quality.md) |
| 19 | [Performance, Debugging & Deployment](19-performance-deployment.md) |
| 20 | [Advanced Topics](20-advanced-topics.md) |

## Cara Menggunakan Materi Ini

1. Ikuti urutan modul 1-20. Modul 3-8 (dasar bahasa, fungsi, pointer, method & interface) menjadi fondasi untuk seluruh modul berikutnya, dan Modul 12-13 (concurrency) adalah kekuatan utama Go.
2. Setiap modul berisi **Teori**, **Praktik** (program Go dan perintah terminal siap coba), dan **Latihan**.
3. Praktik memakai **Go 1.27** yang dipasang lokal (go.dev/dl) atau lewat **Docker** (`golang:1.27`, course Docker Anda). Modul 1 bisa dicoba langsung di Go Playground (go.dev/play) tanpa instalasi.
4. Ketik sendiri kode di bagian Praktik alih-alih menyalin; jalankan juga dengan `go vet` dan `go test -race` bila relevan.
5. Modul 16-17 menyinggung PostgreSQL, Redis, Kafka, Elasticsearch, dan Keycloak, dan Modul 19 menyinggung Docker, Kubernetes, dan EKS dari course Anda.

## Prasyarat

- Familiar dengan command line dan dasar-dasar Git.
- Pernah menulis program dalam bahasa apa pun (variabel, kondisi, loop, fungsi) akan sangat membantu.
- Untuk Modul 17 dan 19: familiar dengan Docker (course Docker) dan konsep dasar Kubernetes (course Kubernetes/EKS).
- Tidak perlu pengalaman Go sebelumnya.

## Sumber Referensi

- Dokumentasi resmi: https://go.dev/doc/
- A Tour of Go: https://go.dev/tour/
- Effective Go: https://go.dev/doc/effective_go
- Dokumentasi paket dan modul: https://pkg.go.dev/
- Catatan rilis Go 1.27: https://go.dev/doc/go1.27
- Roadmap Go: https://roadmap.sh/golang

## Navigasi

Setiap modul punya tautan *Sebelumnya · Daftar Isi · Berikutnya* di atas dan di bawah, plus daftar isi modul itu sendiri.
