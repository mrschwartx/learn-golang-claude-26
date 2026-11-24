# Modul 1: Pengenalan Go

[Daftar Isi](README.md#daftar-isi) · [Modul 2: Instalasi & Setup Environment →](02-instalasi-setup.md)

**Isi modul ini**

- [Teori](#teori)
  - [1.1 Apa itu Go?](#11-apa-itu-go)
  - [1.2 Sejarah Go](#12-sejarah-go)
  - [1.3 Mengapa Memakai Go?](#13-mengapa-memakai-go)
  - [1.4 Go vs Bahasa Lain](#14-go-vs-bahasa-lain)
  - [1.5 Use Case Umum](#15-use-case-umum)
  - [1.6 Filosofi Go](#16-filosofi-go)
  - [1.7 Go dalam Seri Course Ini dan Roadmap Terkait](#17-go-dalam-seri-course-ini-dan-roadmap-terkait)
  - [1.8 Versi Go dan Kebijakan Rilis](#18-versi-go-dan-kebijakan-rilis)
- [Praktik](#praktik)
- [Latihan](#latihan)

---


## Teori

### 1.1 Apa itu Go?

**Go** (sering disebut **Golang** agar mudah dicari di mesin pencari, namun nama resminya hanya "Go") adalah bahasa pemrograman **open-source, statically typed, dan compiled** yang dirancang di Google. Go dibuat untuk satu tujuan utama: membuat **software berskala besar yang andal dan efisien** menjadi mudah ditulis, dibaca, dan dipelihara oleh tim besar.

Tiga ciri yang paling menentukan karakter Go:

1. **Sederhana secara sengaja** — spesifikasi bahasanya kecil (hanya **25 keyword**), sehingga hampir seluruh bahasa bisa dipelajari dalam hitungan hari.
2. **Concurrency sebagai fitur bahasa** — *goroutine* dan *channel* ada langsung di dalam bahasa ([Modul 12](12-goroutines-channels.md)), bukan library tambahan.
3. **Hasil kompilasi berupa satu file binary** — tanpa runtime eksternal (JVM, interpreter Python, Node.js) yang harus dipasang di server.

### 1.2 Sejarah Go

| Tahun | Peristiwa |
|---|---|
| 2007 | Robert Griesemer, Rob Pike, dan Ken Thompson di Google mulai merancang Go — salah satu pemicunya adalah waktu kompilasi C++ yang sangat lama |
| Nov 2009 | Go dirilis sebagai proyek **open-source** |
| Mar 2012 | **Go 1.0** dirilis bersama **Go 1 compatibility promise** — program Go 1.x yang valid akan tetap dapat dikompilasi di rilis Go 1.x berikutnya |
| 2015 | Go 1.5: compiler ditulis ulang dalam Go sendiri (*self-hosting*) dan GC menjadi concurrent dengan jeda sangat singkat |
| Agu 2018 | Go 1.11: **Go Modules** diperkenalkan sebagai fitur awal, lalu menjadi default di Go 1.16 ([Modul 11](11-modules-packages.md)) |
| Mar 2022 | Go 1.18: **Generics** ([Modul 9](09-generics.md)) dan fuzzing |
| Agu 2023 | Go 1.21: paket `log/slog`, `slices`, `maps`, serta manajemen toolchain otomatis |
| Feb 2024 | Go 1.22: variabel loop per-iterasi dan routing `net/http` yang lebih kaya |
| Agu 2024 | Go 1.23: **iterator** (range-over-func) |
| Feb 2026 | Go 1.26: **Green Tea GC** menjadi default |
| Agu 2026 | **Go 1.27** (versi yang dipakai di course ini): generic methods, `encoding/json/v2`, dan profil goroutine leak |

**Go 1 compatibility promise** adalah salah satu alasan terbesar Go dipercaya di industri: kode yang Anda tulis bertahun-tahun lalu hampir selalu masih bisa dikompilasi hari ini, tanpa drama migrasi besar seperti yang dikenal dari beberapa ekosistem lain.

### 1.3 Mengapa Memakai Go?

| Keunggulan | Penjelasan |
|---|---|
| **Kompilasi sangat cepat** | Proyek besar dikompilasi dalam hitungan detik; siklus *edit → build → run* terasa seperti bahasa scripting |
| **Static binary** | Satu file executable tanpa dependency — cocok sekali untuk Docker image kecil ([Modul 19](19-performance-deployment.md)) |
| **Concurrency bawaan** | Ribuan hingga jutaan goroutine dapat berjalan bersamaan dengan memori sangat kecil per goroutine |
| **Standard library lengkap** | HTTP server/client, JSON, kriptografi, testing, dan profiling tersedia tanpa library pihak ketiga ([Modul 14](14-standard-library.md)) |
| **Tooling terpadu** | `go fmt`, `go test`, `go vet`, `go doc`, dan profiler sudah menjadi bagian dari satu perintah `go` ([Modul 18](18-toolchain-code-quality.md)) |
| **Garbage collector** | Manajemen memori otomatis dengan jeda (pause) sangat singkat, tanpa `malloc/free` manual |
| **Mudah dibaca** | Hanya ada sedikit cara untuk menulis hal yang sama, sehingga kode dari tim lain terasa familiar |

### 1.4 Go vs Bahasa Lain

| Aspek | Go | Python | Java | Rust | Node.js (JavaScript) |
|---|---|---|---|---|---|
| Eksekusi | Compiled ke native | Interpreted | Bytecode di JVM | Compiled ke native | JIT di V8 |
| Typing | Static | Dynamic | Static | Static | Dynamic (TypeScript: static) |
| Manajemen memori | Garbage collector | Garbage collector | Garbage collector | Ownership (tanpa GC) | Garbage collector |
| Concurrency | Goroutine + channel | Thread/asyncio (dibatasi GIL*) | Thread, virtual thread | Thread, async/await | Event loop |
| Kecepatan kompilasi | **Sangat cepat** | — | Sedang | Lambat | — |
| Deployment | **Satu binary** | Interpreter + dependency | JVM + JAR | Satu binary | Node.js + `node_modules` |
| Kurva belajar | **Landai** | Sangat landai | Sedang | Curam | Landai |

\* GIL (Global Interpreter Lock) pada CPython; Python 3.13+ menyediakan build *free-threaded* eksperimental.

**Prinsip penting**: tidak ada bahasa yang "terbaik untuk segalanya". Go sangat kuat untuk **layanan backend, tooling jaringan, dan infrastruktur**; kurang cocok untuk pengembangan UI desktop, scripting cepat sekali pakai, atau komputasi numerik/ML yang didominasi ekosistem Python.

### 1.5 Use Case Umum

| Use Case | Penjelasan |
|---|---|
| **Cloud-native & infrastruktur** | Docker, Kubernetes, Terraform, Prometheus, etcd, dan Consul semuanya ditulis dalam Go |
| **Microservices & REST/gRPC API** | Binary kecil, startup cepat, dan throughput tinggi ([Modul 16](16-cli-web-development.md) dan [Modul 17](17-grpc-database-logging.md)) |
| **CLI tool** | Satu binary mudah didistribusikan ke berbagai OS lewat cross-compilation (kubectl, Hugo, GitHub CLI) |
| **DevOps & automation** | Agen, operator, dan controller Kubernetes |
| **Networking & proxy** | Server HTTP, load balancer, dan reverse proxy (Caddy, Traefik) |
| **Data pipeline & streaming** | Consumer/producer berkinerja tinggi untuk Kafka, Redis, atau Elasticsearch |

### 1.6 Filosofi Go

Beberapa prinsip yang akan terus muncul sepanjang course ini:

1. **Simplicity** — lebih baik sedikit fitur yang konsisten daripada banyak fitur yang saling tumpang tindih. Go sengaja tidak memiliki *class*, *inheritance*, *exception*, atau *operator overloading*.
2. **Composition over inheritance** — perilaku dibangun dengan *struct embedding* dan *interface* ([Modul 8](08-methods-interfaces.md)).
3. **Errors are values** — error adalah nilai biasa yang dikembalikan dan diperiksa secara eksplisit ([Modul 10](10-error-handling.md)).
4. **Share memory by communicating** — *"Don't communicate by sharing memory; share memory by communicating."* ([Modul 12](12-goroutines-channels.md)).
5. **Satu cara format** — `gofmt` memformat seluruh kode Go secara seragam sehingga debat soal gaya penulisan nyaris hilang.
6. **Explicit is better than implicit** — tidak ada konversi tipe otomatis, tidak ada import yang tidak terpakai, tidak ada variabel yang dibiarkan tidak terpakai.

### 1.7 Go dalam Seri Course Ini dan Roadmap Terkait

Course sebelumnya (Docker, Kubernetes, Terraform, EKS, Kafka, Elasticsearch, Keycloak, Redis) membahas **infrastruktur dan platform**. Go menjawab pertanyaan yang berbeda: **"dengan apa kita menulis layanan dan tooling yang berjalan di atas semua itu?"** Memahami Go juga berarti Anda dapat membaca kode sumber Docker, Kubernetes, dan Terraform — dan menulis operator, provider, atau plugin sendiri.

Roadmap Go di roadmap.sh menautkan ke roadmap lain yang saling melengkapi:

| Roadmap Terkait | Relevansi |
|---|---|
| Backend Roadmap | Go sebagai bahasa backend: API, database, caching, autentikasi |
| DevOps Roadmap | Tooling dan otomasi yang ditulis dalam Go |
| Docker Roadmap | Mengemas binary Go dalam image minimal ([Modul 19](19-performance-deployment.md)) |
| Kubernetes Roadmap | Deploy dan operator berbasis Go |
| System Design | Menerapkan konkurensi dan layanan terdistribusi |
| Software Design & Architecture | Struktur proyek, interface, dan dependency injection |

### 1.8 Versi Go dan Kebijakan Rilis

- Go dirilis **dua kali setahun**: Februari dan Agustus (Go 1.26 → Februari 2026, Go 1.27 → Agustus 2026).
- Setiap rilis mayor didukung hingga **dua rilis mayor berikutnya** tersedia — artinya hanya dua versi terakhir yang menerima perbaikan keamanan.
- Nomor `1.x` tetap dipakai karena Go 1 compatibility promise; belum ada rencana Go 2 yang merusak kompatibilitas.
- Dengan **toolchain management** ([Modul 2](02-instalasi-setup.md)), setiap proyek dapat mengunci versi Go yang dibutuhkannya di `go.mod`.

## Praktik

Di modul ini kita belum memasang Go ([Modul 2](02-instalasi-setup.md)). Namun kita bisa langsung mencoba program Go pertama lewat **Go Playground** — compiler Go berbasis web di **https://go.dev/play/** (tanpa instalasi apa pun). Salin program berikut ke Playground lalu klik **Run**.

```go
// hello.go
package main

import "fmt"

func main() {
	fmt.Println("Halo, Go!")

	for i := 1; i <= 3; i++ {
		fmt.Println("Iterasi ke-", i)
	}
}
```

Hasil yang diharapkan:

```
Halo, Go!
Iterasi ke- 1
Iterasi ke- 2
Iterasi ke- 3
```

Beberapa hal yang bisa diamati dari program kecil ini, meski belum dijelaskan mendalam:

- `package main` menandakan ini program **executable** (bukan library).
- `import "fmt"` memuat paket `fmt` dari standard library untuk mencetak ke layar.
- `func main()` adalah titik masuk program.
- Tidak ada titik koma di akhir baris dan tidak ada tanda kurung pada `for`.
- Tidak ada `class` atau `public static void` — hanya paket dan fungsi.

Eksperimen kecil: hapus baris `import "fmt"` lalu jalankan. Compiler Go akan menolak program yang memakai paket tanpa mengimpornya. Lalu tambahkan `import "os"` **tanpa memakainya** — compiler juga menolak, karena Go tidak mengizinkan **import yang tidak terpakai**.

## Latihan

1. Jalankan program Hello di Go Playground, lalu ubah agar mencetak nama Anda dan loop berjalan 5 kali. Amati bagaimana Playground memformat kode Anda (tombol **Format**) dan hubungkan dengan filosofi `gofmt` di 1.6.
2. Sebutkan tiga alasan mengapa **satu file binary tanpa dependency** menguntungkan proses deployment dibanding aplikasi yang membutuhkan interpreter atau runtime terpisah di server.
3. Jelaskan dengan kata-kata Anda sendiri apa isi **Go 1 compatibility promise** dan mengapa janji ini penting bagi perusahaan yang menjalankan kode Go selama bertahun-tahun.
4. Dari tabel 1.5, pilih dua use case yang paling relevan dengan pekerjaan atau minat Anda, lalu jelaskan mengapa Go cocok untuknya dibanding bahasa yang saat ini Anda pakai.
5. Cari di GitHub bahwa Docker (`moby/moby`) dan Kubernetes (`kubernetes/kubernetes`) memang ditulis dalam Go. Dari sudut pandang Anda sebagai pengguna course Docker dan Kubernetes, apa manfaat dari mampu membaca kode sumber tool yang Anda pakai sehari-hari?
6. Diskusikan: sebuah tim sedang memilih bahasa untuk membangun layanan API yang harus menangani puluhan ribu koneksi bersamaan dan di-deploy ke Kubernetes. Sebutkan minimal tiga faktor dari modul ini yang mendukung pilihan Go, serta satu situasi di mana Anda justru akan memilih bahasa lain.

---

[Daftar Isi](README.md#daftar-isi) · [Modul 2: Instalasi & Setup Environment →](02-instalasi-setup.md)
