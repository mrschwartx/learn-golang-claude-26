# Modul 16: Membangun CLI & Web Development

[← Modul 15: Testing & Benchmarking](15-testing-benchmarking.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 17: gRPC, Database & Logging →](17-grpc-database-logging.md)

**Isi modul ini**

- [Teori](#teori)
  - [16.1 Membangun CLI](#161-membangun-cli)
  - [16.2 `net/http` — Server](#162-nethttp--server)
  - [16.3 `net/http` — Client](#163-nethttp--client)
  - [16.4 Framework Web](#164-framework-web)
  - [16.5 Struktur Aplikasi Web](#165-struktur-aplikasi-web)
  - [16.6 Realtime Communication](#166-realtime-communication)
  - [16.7 Template dan Berkas Statis](#167-template-dan-berkas-statis)
- [Praktik](#praktik)
  - [REST API dengan Standard Library](#rest-api-dengan-standard-library)
  - [Server-Sent Events](#server-sent-events)
  - [Klien HTTP dengan Timeout](#klien-http-dengan-timeout)
  - [Contoh CLI dengan Cobra dan Framework Gin (Pihak Ketiga)](#contoh-cli-dengan-cobra-dan-framework-gin-pihak-ketiga)
- [Latihan](#latihan)

---


## Teori

### 16.1 Membangun CLI

Go sangat cocok untuk **command-line tool**: satu binary, startup instan, dan cross-compilation mudah ([Modul 19](19-performance-deployment.md)). `kubectl`, Docker CLI, Hugo, dan GitHub CLI semuanya dibangun dengan Go.

| Pilihan | Karakteristik |
|---|---|
| **`flag`** (standard library) | Cukup untuk tool sederhana; subcommand lewat `flag.NewFlagSet` ([Modul 14](14-standard-library.md)) |
| **Cobra** (`spf13/cobra`) | Standar de facto untuk CLI bersubcommand: flag POSIX, bantuan otomatis, shell completion; sering dipasangkan dengan **Viper** untuk konfigurasi. Dipakai `kubectl` dan Docker CLI |
| **urfave/cli** | API deklaratif berbasis struct; ringan dan mudah dipelajari |
| **Bubble Tea** (`charmbracelet/bubbletea`) | Framework **TUI** (antarmuka terminal interaktif) dengan arsitektur ala Elm: `Model`, `Update`, `View`; sering dipakai bersama Lip Gloss dan Bubbles |

Tips merancang CLI yang baik: kode keluar (`os.Exit`) **0** untuk sukses dan bukan-nol untuk gagal; tulis hasil ke **stdout** dan pesan error/log ke **stderr**; sediakan `--help` dan `--version`; dukung konfigurasi lewat flag, environment variable, lalu berkas konfigurasi (urutan prioritas jelas).

### 16.2 `net/http` — Server

Paket **`net/http`** cukup matang untuk server produksi, dan inilah yang dipakai Go secara default.

| Konsep | Penjelasan |
|---|---|
| **`http.Handler`** | Interface dengan satu method `ServeHTTP(w http.ResponseWriter, r *http.Request)` ([Modul 8](08-methods-interfaces.md)) |
| **`http.HandlerFunc`** | Adaptor yang membuat fungsi biasa memenuhi `Handler` |
| **`http.ServeMux`** | Router bawaan. Sejak **Go 1.22** mendukung **method dan wildcard** pada pola |
| **Middleware** | Fungsi `func(http.Handler) http.Handler` yang membungkus handler (logging, autentikasi, recover) |
| **`http.Server`** | Struct konfigurasi server dengan timeout dan `Shutdown` |

Pola routing Go 1.22+:

| Pola | Arti |
|---|---|
| `"GET /items"` | Hanya method GET pada path `/items` |
| `"GET /items/{id}"` | Wildcard; baca dengan `r.PathValue("id")` |
| `"GET /files/{path...}"` | Wildcard yang menangkap sisa path |
| `"GET /{$}"` | Cocok **hanya** untuk `/` persis |

Method yang tidak cocok otomatis dijawab **405 Method Not Allowed**.

**Selalu konfigurasi timeout** pada `http.Server`; nilai default-nya tanpa batas sehingga rentan klien yang lambat atau berbahaya (*Slowloris*):

| Field | Fungsi |
|---|---|
| `ReadHeaderTimeout` | Batas waktu membaca header request |
| `ReadTimeout` / `WriteTimeout` | Batas membaca seluruh request / menulis response |
| `IdleTimeout` | Batas koneksi keep-alive menganggur |

Untuk mematikan server dengan rapi gunakan `srv.Shutdown(ctx)` ([Modul 13](13-sync-context-patterns.md)).

### 16.3 `net/http` — Client

```go
client := &http.Client{Timeout: 5 * time.Second}
req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
resp, err := client.Do(req)
if err != nil { return err }
defer resp.Body.Close()
```

| Aturan | Alasan |
|---|---|
| **Jangan pakai `http.DefaultClient`/`http.Get` tanpa timeout** | Default tanpa batas waktu — request bisa menggantung selamanya |
| **Selalu `defer resp.Body.Close()`** | Mencegah kebocoran koneksi |
| **Gunakan `NewRequestWithContext`** | Pembatalan dan deadline lewat `context` ([Modul 13](13-sync-context-patterns.md)) |
| **Buat satu `Client` dan pakai ulang** | `Client` aman dipakai bersamaan dan menjaga *connection pool* |
| **Periksa `resp.StatusCode`** | Status 4xx/5xx **bukan** `err` — hanya kegagalan transport yang menjadi error |

### 16.4 Framework Web

Banyak layanan Go berjalan baik hanya dengan standard library. Framework menambah kenyamanan (binding, validasi, grouping route, middleware siap pakai):

| Framework | Karakteristik |
|---|---|
| **Gin** | Paling populer; router cepat, binding/validasi JSON, middleware kaya |
| **Echo** | Mirip Gin, dengan API handler yang mengembalikan `error` dan dokumentasi rapi |
| **Fiber** | Gaya Express.js; dibangun di atas **fasthttp**, **bukan `net/http`** sehingga middleware `net/http` tidak otomatis kompatibel |
| **Beego** | Framework *full-stack* (MVC, ORM, CLI `bee`) yang lebih "berpendirian" |

Panduan memilih: untuk layanan baru, mulai dari **`net/http` + `ServeMux` Go 1.22** — sering sudah cukup. Pilih Gin atau Echo jika tim menginginkan produktivitas ekstra; pilih Fiber hanya bila Anda memahami konsekuensi `fasthttp`. Apa pun pilihannya, **pisahkan logika bisnis dari lapisan HTTP** agar framework mudah diganti.

### 16.5 Struktur Aplikasi Web

Susunan berlapis yang umum, dengan arah dependensi satu arah ([Modul 11](11-modules-packages.md)):

```
HTTP handler  →  service (logika bisnis)  →  repository (akses data)
 (parse/encode)    (aturan, validasi)         (SQL, cache, API luar)
```

- **Handler** hanya menerjemahkan HTTP ↔ panggilan fungsi (parse input, tulis response, pilih status code).
- **Service** menerima dan mengembalikan tipe domain, tidak tahu soal HTTP.
- **Repository** disembunyikan di balik **interface** sehingga mudah diganti fake saat uji ([Modul 15](15-testing-benchmarking.md)).
- **Dependency injection** lewat konstruktor (`NewService(repo)`) — tanpa framework DI.
- **Konfigurasi** dibaca dari environment variable di `main`, lalu diteruskan sebagai argumen.

### 16.6 Realtime Communication

Untuk komunikasi dua arah atau push data ke klien:

| Teknologi | Penjelasan |
|---|---|
| **Server-Sent Events (SSE)** | Satu arah server → klien lewat HTTP biasa; cukup dengan `http.Flusher` pada standard library |
| **WebSocket** | Dua arah penuh. Pustaka populer: `gorilla/websocket`, `coder/websocket` |
| **Melody** (`olahol/melody`) | Lapisan tipis di atas WebSocket dengan *session* dan *broadcast* yang sederhana |
| **Centrifugo** | **Server realtime mandiri** (ditulis dalam Go) yang menangani koneksi WebSocket/SSE, channel pub/sub, dan presence; aplikasi Anda hanya mempublikasikan pesan lewat HTTP/gRPC API-nya. Cocok jika backend utama bukan Go atau ingin memisahkan beban koneksi |

Setiap koneksi realtime menahan sebuah goroutine (atau lebih), sehingga model goroutine Go ([Modul 12](12-goroutines-channels.md)) cocok untuk puluhan ribu koneksi — tetapi tetap perlu batas, timeout ping/pong, dan penanganan klien lambat.

### 16.7 Template dan Berkas Statis

Untuk halaman server-rendered: **`html/template`** (otomatis meng-escape output sehingga aman dari XSS), dan `http.FileServerFS` untuk berkas statis. Dikombinasikan dengan `go:embed` ([Modul 14](14-standard-library.md)), seluruh template dan aset tertanam di satu binary.

## Praktik

### REST API dengan Standard Library

```go
// api.go
package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Item struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
}

type Toko struct {
	mu    sync.RWMutex
	items map[int]Item
	next  int
}

func NewToko() *Toko { return &Toko{items: make(map[int]Item), next: 1} }

func tulisJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func tulisError(w http.ResponseWriter, status int, pesan string) {
	tulisJSON(w, status, map[string]string{"error": pesan})
}

func (t *Toko) daftar(w http.ResponseWriter, r *http.Request) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	hasil := make([]Item, 0, len(t.items))
	for i := 1; i < t.next; i++ {
		if it, ok := t.items[i]; ok {
			hasil = append(hasil, it)
		}
	}
	tulisJSON(w, http.StatusOK, hasil)
}

func (t *Toko) ambil(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		tulisError(w, http.StatusBadRequest, "id harus angka")
		return
	}
	t.mu.RLock()
	it, ok := t.items[id]
	t.mu.RUnlock()
	if !ok {
		tulisError(w, http.StatusNotFound, "item tidak ditemukan")
		return
	}
	tulisJSON(w, http.StatusOK, it)
}

func (t *Toko) buat(w http.ResponseWriter, r *http.Request) {
	var in Item
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		tulisError(w, http.StatusBadRequest, "JSON tidak valid: "+err.Error())
		return
	}
	if in.Nama == "" || in.Harga <= 0 {
		tulisError(w, http.StatusUnprocessableEntity, "nama wajib diisi dan harga harus positif")
		return
	}
	t.mu.Lock()
	in.ID = t.next
	t.next++
	t.items[in.ID] = in
	t.mu.Unlock()
	tulisJSON(w, http.StatusCreated, in)
}

// Middleware: mencatat setiap request.
func logging(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mulai := time.Now()
		next.ServeHTTP(w, r)
		log.Info("request", "method", r.Method, "path", r.URL.Path, "durasi", time.Since(mulai))
	})
}

func main() {
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	toko := NewToko()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /items", toko.daftar)
	mux.HandleFunc("GET /items/{id}", toko.ambil)
	mux.HandleFunc("POST /items", toko.buat)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           logging(log, mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Info("server berjalan", "addr", srv.Addr)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Error("server berhenti", "err", err)
	}
}
```

Jalankan di satu terminal, lalu uji dengan `curl` di terminal lain:

```bash
go run api.go

curl -i -X POST localhost:8080/items -d '{"nama":"Kopi","harga":25000}'
curl -s localhost:8080/items
curl -s localhost:8080/items/1
curl -s localhost:8080/items/99
curl -s -X POST localhost:8080/items -d '{"nama":"","harga":0}'
curl -i -X DELETE localhost:8080/items/1
```

Perhatikan: item yang tidak ada → `404`, validasi gagal → `422`, dan method yang tidak didaftarkan (`DELETE`) otomatis → `405`.

### Server-Sent Events

```go
// sse.go
package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("GET /waktu", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		fl, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming tidak didukung", http.StatusInternalServerError)
			return
		}
		tick := time.NewTicker(time.Second)
		defer tick.Stop()
		for {
			select {
			case <-r.Context().Done(): // klien terputus
				return
			case t := <-tick.C:
				fmt.Fprintf(w, "data: %s\n\n", t.Format(time.TimeOnly))
				fl.Flush()
			}
		}
	})
	http.ListenAndServe(":8081", nil)
}
```

```bash
go run sse.go
curl -N localhost:8081/waktu
```

### Klien HTTP dengan Timeout

```go
// klien.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

func ambilJSON(ctx context.Context, c *http.Client, url string, v any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status tak terduga: %s", resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(v)
}

func main() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/lambat" {
			time.Sleep(300 * time.Millisecond)
		}
		w.Write([]byte(`{"pesan":"ok"}`))
	}))
	defer srv.Close()

	c := &http.Client{Timeout: 100 * time.Millisecond}
	var hasil map[string]string
	fmt.Println(ambilJSON(context.Background(), c, srv.URL+"/cepat", &hasil), hasil)
	fmt.Println(ambilJSON(context.Background(), c, srv.URL+"/lambat", &hasil) != nil)
}
```

### Contoh CLI dengan Cobra dan Framework Gin (Pihak Ketiga)

Kode berikut memakai library pihak ketiga; pasang dengan `go get` di modul Anda.

```go
// cobra_main.go
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var ulang int
	root := &cobra.Command{Use: "sapa", Short: "Aplikasi penyapa"}
	halo := &cobra.Command{
		Use:   "halo [nama]",
		Short: "Menyapa seseorang",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for range ulang {
				fmt.Fprintf(cmd.OutOrStdout(), "Halo, %s!\n", args[0])
			}
			return nil
		},
	}
	halo.Flags().IntVarP(&ulang, "ulang", "n", 1, "jumlah pengulangan")
	root.AddCommand(halo)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
```

```bash
go mod init sapa && go get github.com/spf13/cobra@latest
go run cobra_main.go halo Andi -n 2
go run cobra_main.go --help
```

```go
// gin_main.go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/halo/:nama", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"pesan": "Halo, " + c.Param("nama")})
	})
	r.Run(":8080")
}
```

Bandingkan jumlah baris, binding, dan perilaku error dengan versi `net/http` di atas.

## Latihan

1. Perluas REST API dengan `PUT /items/{id}` dan `DELETE /items/{id}`. Pastikan `DELETE` pada id yang tidak ada mengembalikan `404`, dan `PUT` memvalidasi input seperti `POST`.
2. Tambahkan middleware `recoverPanic` yang menangkap panic di handler, mencatatnya dengan `slog`, dan mengembalikan `500`. Uji dengan handler yang sengaja `panic`. Hubungkan dengan konsep `recover` di [Modul 10](10-error-handling.md).
3. Tambahkan middleware autentikasi sederhana yang memeriksa header `Authorization: Bearer <token>` dan menolak dengan `401`. Jelaskan mengapa middleware berbentuk `func(http.Handler) http.Handler` mudah dirangkai.
4. Tulis uji `httptest` ([Modul 15](15-testing-benchmarking.md)) untuk seluruh endpoint API di atas, termasuk kasus error dan `405`. Jalankan dengan `go test -race` untuk memastikan `RWMutex` bekerja benar.
5. Tambahkan graceful shutdown ke server ([Modul 13](13-sync-context-patterns.md)) dan buktikan request yang sedang berjalan diselesaikan sebelum proses berhenti.
6. Buat klien yang memanggil API di atas dengan `http.Client` ber-timeout dan **retry** terbatas (maksimal 3 kali, jeda bertambah) untuk status `5xx`. Mengapa tidak semua error aman di-retry?
7. Diskusikan: sebuah tim ingin memakai Gin untuk layanan baru, sementara rekan lain bersikeras memakai `net/http` saja sejak Go 1.22. Susun kriteria keputusan (kebutuhan fitur, ekosistem middleware, keterikatan framework, performa, keterampilan tim) dan tentukan pilihan Anda beserta alasannya.

---

[← Modul 15: Testing & Benchmarking](15-testing-benchmarking.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 17: gRPC, Database & Logging →](17-grpc-database-logging.md)
