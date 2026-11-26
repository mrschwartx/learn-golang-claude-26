# Modul 13: Sync, Context & Concurrency Patterns

[← Modul 12: Goroutines & Channels](12-goroutines-channels.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 14: Standard Library →](14-standard-library.md)

**Isi modul ini**

- [Teori](#teori)
  - [13.1 Paket `sync` dan Data Race](#131-paket-sync-dan-data-race)
  - [13.2 Mutex dan RWMutex](#132-mutex-dan-rwmutex)
  - [13.3 WaitGroup](#133-waitgroup)
  - [13.4 Once, Atomic, dan Pool](#134-once-atomic-dan-pool)
  - [13.5 Paket `context`](#135-paket-context)
  - [13.6 Deadlines dan Cancellations](#136-deadlines-dan-cancellations)
  - [13.7 Use Case Umum Context](#137-use-case-umum-context)
  - [13.8 Concurrency Patterns](#138-concurrency-patterns)
  - [13.9 Race Detection](#139-race-detection)
- [Praktik](#praktik)
  - [Menemukan dan Memperbaiki Data Race](#menemukan-dan-memperbaiki-data-race)
  - [Context: Timeout dan Pembatalan](#context-timeout-dan-pembatalan)
  - [Pipeline Fan-Out/Fan-In dengan Pembatalan](#pipeline-fan-outfan-in-dengan-pembatalan)
  - [Graceful Shutdown (Kerangka)](#graceful-shutdown-kerangka)
- [Latihan](#latihan)

---


## Teori

### 13.1 Paket `sync` dan Data Race

Channel ([Modul 12](12-goroutines-channels.md)) bukan satu-satunya cara mengoordinasikan goroutine. Ketika beberapa goroutine **membaca dan menulis state yang sama**, dibutuhkan sinkronisasi — jika tidak, terjadi **data race**: dua goroutine mengakses variabel yang sama bersamaan dan minimal satu menulis, tanpa urutan yang terjamin. Hasilnya tidak terdefinisi: nilai salah, korupsi memori, atau crash yang sulit direproduksi.

Paket `sync` dan `sync/atomic` menyediakan primitif untuk itu:

| Primitif | Kegunaan |
|---|---|
| `sync.Mutex` / `sync.RWMutex` | Kunci eksklusif untuk melindungi state bersama |
| `sync.WaitGroup` | Menunggu sekumpulan goroutine selesai |
| `sync.Once` | Menjalankan inisialisasi tepat satu kali |
| `sync.Pool` | Menggunakan ulang objek sementara untuk mengurangi alokasi |
| `sync.Cond` | Menunggu/mengumumkan kondisi (jarang dipakai langsung) |
| `sync.Map` | Map aman-concurrent untuk pola akses tertentu |
| `sync/atomic` | Operasi atomik pada angka dan pointer tanpa kunci |

### 13.2 Mutex dan RWMutex

```go
type Penghitung struct {
    mu sync.Mutex
    n  int
}

func (p *Penghitung) Tambah() {
    p.mu.Lock()
    defer p.mu.Unlock()
    p.n++
}
```

| Aturan | Penjelasan |
|---|---|
| **`defer Unlock()`** | Pastikan kunci selalu dilepas, termasuk saat panic atau `return` awal |
| **Jangan salin Mutex** | Menyalin struct yang berisi `sync.Mutex` menyalin state kuncinya; `go vet` (analyzer `copylocks`) memperingatkan. Kirim dan simpan sebagai pointer |
| **Letakkan mutex di samping data yang dilindunginya** | Satu mutex untuk satu kelompok field yang selalu diakses bersama |
| **Hindari memegang kunci saat I/O lambat** | Memperpanjang antrean goroutine yang menunggu |

**`sync.RWMutex`** mengizinkan **banyak pembaca bersamaan** (`RLock`/`RUnlock`) atau **satu penulis** (`Lock`/`Unlock`). Berguna pada data yang jauh lebih sering dibaca daripada ditulis, tetapi lebih berat daripada `Mutex` biasa — ukur sebelum memakainya.

### 13.3 WaitGroup

`sync.WaitGroup` menghitung goroutine yang sedang berjalan. `Add(n)` menambah hitungan, `Done()` menguranginya, `Wait()` memblokir sampai hitungan nol.

```go
var wg sync.WaitGroup
for i := range 3 {
    wg.Add(1)
    go func() {
        defer wg.Done()
        proses(i)
    }()
}
wg.Wait()
```

Sejak **Go 1.25**, tersedia **`wg.Go(f)`** yang menggabungkan `Add(1)`, `go`, dan `Done()` sehingga lebih ringkas dan menghindari lupa memanggil `Add`: `wg.Go(func() { proses(i) })`. Panggil `Add` **sebelum** memulai goroutine, bukan di dalamnya.

### 13.4 Once, Atomic, dan Pool

| Primitif | Contoh dan Catatan |
|---|---|
| **`sync.Once`** | `once.Do(inisialisasi)` — dijamin berjalan sekali, aman dari banyak goroutine. Varian generik: `sync.OnceFunc`, `sync.OnceValue`, `sync.OnceValues` (Go 1.21+) |
| **`sync/atomic`** | Tipe `atomic.Int64`, `atomic.Bool`, `atomic.Pointer[T]` (Go 1.19+) dengan `Add`, `Load`, `Store`, `CompareAndSwap` — cocok untuk counter dan flag sederhana |
| **`sync.Pool`** | `Get`/`Put` untuk objek sementara (misalnya buffer). Isinya bisa dibuang GC kapan saja, jadi **bukan** tempat menyimpan data penting |
| **`sync.Map`** | Hanya unggul untuk dua pola: key ditulis sekali lalu dibaca berkali-kali, atau goroutine mengakses set key yang saling lepas. Selain itu, `map` biasa + `Mutex` lebih jelas dan bertipe aman |

### 13.5 Paket `context`

**`context.Context`** membawa **sinyal pembatalan, deadline, dan nilai terbatas** lintas batas API dan goroutine. Ini cara standar menghentikan pekerjaan yang tidak lagi dibutuhkan — misalnya ketika klien HTTP memutus koneksi.

| Fungsi | Kegunaan |
|---|---|
| `context.Background()` | Konteks akar; dipakai di `main`, init, dan test |
| `context.TODO()` | Penanda sementara saat belum jelas konteks apa yang dipakai |
| `context.WithCancel(parent)` | Membuat konteks yang bisa dibatalkan manual lewat fungsi `cancel` |
| `context.WithTimeout(parent, d)` | Batal otomatis setelah durasi `d` |
| `context.WithDeadline(parent, t)` | Batal otomatis pada waktu `t` |
| `context.WithCancelCause(parent)` | Pembatalan dengan **alasan** yang bisa dibaca lewat `context.Cause` |
| `context.WithValue(parent, key, val)` | Membawa nilai lingkup-request (request ID, user) — **hemat**, bukan untuk parameter fungsi biasa |

Konvensi:

1. `ctx` adalah **parameter pertama** fungsi: `func Ambil(ctx context.Context, id int) (...)`.
2. **Jangan simpan** `Context` di dalam struct; teruskan sebagai argumen.
3. **Selalu panggil `cancel`** (biasanya `defer cancel()`) agar resource dilepas — `go vet` memeriksa ini.
4. Gunakan **tipe key sendiri** (bukan `string`) untuk `WithValue` agar tidak bertabrakan antar paket.

### 13.6 Deadlines dan Cancellations

Context membentuk **pohon**: ketika konteks induk dibatalkan, **semua turunannya ikut batal**. Pekerjaan yang bisa dibatalkan memantau `ctx.Done()`:

```go
select {
case <-ctx.Done():
    return ctx.Err()   // context.Canceled atau context.DeadlineExceeded
case hasil := <-kerja:
    return proses(hasil)
}
```

| `ctx.Err()` | Arti |
|---|---|
| `nil` | Konteks masih aktif |
| `context.Canceled` | Dibatalkan lewat `cancel()` |
| `context.DeadlineExceeded` | Deadline atau timeout terlampaui |

Pembatalan bersifat **kooperatif** — Go tidak bisa menghentikan paksa goroutine; kode harus memeriksa `ctx.Done()` atau memakai API yang menerima `ctx` (`http.NewRequestWithContext`, `db.QueryContext`, [Modul 17](17-grpc-database-logging.md)).

### 13.7 Use Case Umum Context

| Kasus | Penerapan |
|---|---|
| **Timeout panggilan eksternal** | `WithTimeout` di sekitar pemanggilan HTTP/database agar request tidak menggantung |
| **Pembatalan request** | `r.Context()` di handler HTTP batal otomatis saat klien memutus koneksi ([Modul 16](16-cli-web-development.md)) |
| **Graceful shutdown** | `signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)` untuk menghentikan server dengan rapi, termasuk di Kubernetes ([Modul 19](19-performance-deployment.md)) |
| **Menghentikan sekelompok goroutine** | Satu `cancel()` membatalkan seluruh worker yang berbagi konteks |
| **Tracing dan logging** | Menyebarkan request ID lewat `WithValue` |

### 13.8 Concurrency Patterns

| Pola | Penjelasan |
|---|---|
| **Fan-out** | Beberapa worker membaca dari **satu channel** input untuk membagi beban |
| **Fan-in** | Menggabungkan **beberapa channel** menjadi satu |
| **Pipeline** | Rangkaian tahap (`generator → transformasi → konsumen`) yang dihubungkan channel; tiap tahap menutup channel keluarannya |
| **Semaphore** | Channel berbuffer berukuran N sebagai pembatas concurrency |
| **`errgroup`** | `golang.org/x/sync/errgroup` — sekelompok goroutine dengan propagasi error pertama dan pembatalan konteks |
| **Rate limiting** | `time.Ticker` atau `golang.org/x/time/rate` |
| **Generator** | Fungsi yang mengembalikan `<-chan T` atau `iter.Seq[T]` ([Modul 9](09-generics.md)) |

Prinsip pada semua pola: **setiap goroutine harus punya jalur berhenti** (konteks atau channel tertutup), dan **pemilik channel-lah yang menutupnya**.

### 13.9 Race Detection

Go menyertakan **race detector** — instrumentasi runtime yang mendeteksi data race saat program berjalan:

```bash
go run -race main.go
go test -race ./...
go build -race -o app-race .
```

| Aspek | Keterangan |
|---|---|
| Hanya mendeteksi race yang **benar-benar terjadi** saat eksekusi | Beban uji yang bagus menentukan cakupannya |
| **Tidak ada false positive** | Laporan yang muncul adalah race sungguhan |
| Overhead | Kira-kira 5–10× memori dan 2–20× CPU — gunakan untuk test dan CI, bukan produksi |
| Perlu **CGO** | Pada banyak platform, butuh compiler C terpasang |

Laporan memuat stack trace goroutine yang menulis dan yang membaca/menulis bersamaan, termasuk baris kodenya. Jalankan `go test -race ./...` di CI pada setiap perubahan.

Sejak Go 1.25, paket **`testing/synctest`** juga membantu menguji kode concurrent berbasis waktu secara deterministik ([Modul 15](15-testing-benchmarking.md)).

## Praktik

### Menemukan dan Memperbaiki Data Race

```go
// race.go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	hitung := 0
	for range 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			hitung++ // data race!
		}()
	}
	wg.Wait()
	fmt.Println("hasil (seharusnya 1000):", hitung)
}
```

```bash
go run -race race.go
```

Race detector mencetak `WARNING: DATA RACE` disertai stack trace akses baca dan tulis. Tiga cara memperbaikinya:

```go
// perbaikan.go
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Penghitung struct {
	mu sync.Mutex
	n  int
}

func (p *Penghitung) Tambah() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.n++
}

func main() {
	var wg sync.WaitGroup

	// 1. Mutex
	var pc Penghitung
	for range 1000 {
		wg.Go(func() { pc.Tambah() }) // wg.Go: Go 1.25+
	}

	// 2. Atomic
	var at atomic.Int64
	for range 1000 {
		wg.Go(func() { at.Add(1) })
	}

	// 3. Channel: satu goroutine pemilik state
	ch := make(chan struct{})
	selesai := make(chan int)
	go func() {
		n := 0
		for range ch {
			n++
		}
		selesai <- n
	}()
	for range 1000 {
		wg.Go(func() { ch <- struct{}{} })
	}
	wg.Wait()
	close(ch)

	fmt.Println(pc.n, at.Load(), <-selesai)

	var once sync.Once
	for range 3 {
		once.Do(func() { fmt.Println("inisialisasi sekali") })
	}
}
```

### Context: Timeout dan Pembatalan

```go
// context.go
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type kunci string

func kerjaLambat(ctx context.Context, d time.Duration) error {
	select {
	case <-time.After(d):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	fmt.Println("cepat:", kerjaLambat(ctx, 10*time.Millisecond))
	err := kerjaLambat(ctx, 500*time.Millisecond)
	fmt.Println("lambat:", err, errors.Is(err, context.DeadlineExceeded))

	induk, batalkan := context.WithCancel(context.Background())
	anak, cancelAnak := context.WithTimeout(induk, time.Hour)
	defer cancelAnak()
	batalkan() // membatalkan induk -> anak ikut batal
	<-anak.Done()
	fmt.Println("anak:", anak.Err())

	ctxNilai := context.WithValue(context.Background(), kunci("reqID"), "abc-123")
	fmt.Println(ctxNilai.Value(kunci("reqID")))
}
```

### Pipeline Fan-Out/Fan-In dengan Pembatalan

```go
// pipeline.go
package main

import (
	"context"
	"fmt"
	"slices"
	"sync"
)

func generator(ctx context.Context, n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 1; i <= n; i++ {
			select {
			case out <- i:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func kuadrat(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in {
			select {
			case out <- v * v:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func gabung(chans ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	for _, c := range chans {
		wg.Go(func() {
			for v := range c {
				out <- v
			}
		})
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sumber := generator(ctx, 10)
	// fan-out: 3 worker membaca dari channel yang sama
	w1, w2, w3 := kuadrat(ctx, sumber), kuadrat(ctx, sumber), kuadrat(ctx, sumber)

	var hasil []int
	for v := range gabung(w1, w2, w3) { // fan-in
		hasil = append(hasil, v)
	}
	slices.Sort(hasil)
	fmt.Println(hasil)
}
```

### Graceful Shutdown (Kerangka)

```go
// shutdown.go
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{Addr: ":8080"}
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done() // menunggu Ctrl+C atau SIGTERM
	log.Println("mematikan server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("shutdown paksa:", err)
	}
}
```

## Latihan

1. Jalankan `race.go` dengan `-race` beberapa kali, baca laporan race detector, dan tunjukkan baris kode yang menulis serta yang membaca. Jelaskan mengapa hasil tanpa `-race` kadang tetap tampak benar.
2. Implementasikan `Cache` sederhana (`map[string]string`) yang aman-concurrent dengan `sync.RWMutex`. Tulis uji yang memanggilnya dari 100 goroutine dan jalankan dengan `go test -race`.
3. Ganti counter pada latihan sebelumnya dengan `atomic.Int64`, lalu bandingkan kecepatan Mutex vs atomic memakai benchmark ([Modul 15](15-testing-benchmarking.md)). Kapan atomic tidak cukup?
4. Tulis fungsi `AmbilHalaman(ctx, url)` yang memakai `http.NewRequestWithContext` dengan timeout 2 detik. Uji dengan server lambat (`httptest` atau `time.Sleep`) dan periksa `errors.Is(err, context.DeadlineExceeded)`.
5. Ubah pipeline di atas agar berhenti dini: konsumen hanya mengambil 3 hasil lalu memanggil `cancel()`. Buktikan tidak ada goroutine yang bocor, misalnya dengan `runtime.NumGoroutine()` sebelum dan sesudah.
6. Gunakan channel berbuffer sebagai **semaphore** untuk membatasi maksimal 3 goroutine yang berjalan bersamaan dari 20 tugas, lalu bandingkan dengan worker pool ([Modul 12](12-goroutines-channels.md)). Kapan masing-masing lebih tepat?
7. Diskusikan: sebuah layanan di Kubernetes ([Modul 19](19-performance-deployment.md)) menerima `SIGTERM` saat deployment dirolling tetapi sebagian request terputus. Jelaskan bagaimana `signal.NotifyContext` dan `http.Server.Shutdown` memberi *graceful shutdown*, dan hubungannya dengan `terminationGracePeriodSeconds`.

---

[← Modul 12: Goroutines & Channels](12-goroutines-channels.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 14: Standard Library →](14-standard-library.md)
