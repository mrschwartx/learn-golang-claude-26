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
  - [CLI dengan urfave/cli](#cli-dengan-urfavecli)
  - [TUI dengan Bubble Tea](#tui-dengan-bubble-tea)
  - [Echo, Fiber, dan Beego: Endpoint yang Sama](#echo-fiber-dan-beego-endpoint-yang-sama)
  - [WebSocket dengan Melody dan Publikasi ke Centrifugo](#websocket-dengan-melody-dan-publikasi-ke-centrifugo)
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

Cara memilih di antara keempatnya:

- **`flag`** bila tool hanya punya satu perintah dan beberapa opsi, tanpa dependensi tambahan.
- **Cobra** bila tool punya banyak subcommand bersarang (`app server start`, `app db migrate`), butuh shell completion, atau akan diintegrasikan dengan Viper.
- **urfave/cli** bila ingin definisi perintah yang ringkas dalam satu struktur data dan tidak butuh ekosistem sebesar Cobra. Contohnya ada di bagian Praktik.
- **Bubble Tea** bukan pengganti Cobra atau urfave/cli: ia dipakai untuk antarmuka **interaktif** (menu, daftar pilihan, progress bar, formulir) di dalam terminal. Keduanya sering dipakai bersamaan: Cobra untuk parsing perintah, Bubble Tea untuk tampilan interaktif pada salah satu subcommand.

Bubble Tea memakai pola **Model-Update-View** (arsitektur Elm). `Model` menyimpan seluruh state, `Update(msg)` menerima peristiwa (tombol ditekan, timer, hasil I/O) dan mengembalikan state baru beserta perintah berikutnya, dan `View()` merender state menjadi string. Karena state hanya berubah di `Update`, alur programnya mudah diuji dan diikuti.

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

Perbedaan yang paling sering menentukan pilihan:

| Aspek | Gin | Echo | Fiber | Beego |
|---|---|---|---|---|
| Signature handler | `func(*gin.Context)` (tanpa return) | `func(echo.Context) error` | `func(*fiber.Ctx) error` (v2) | method pada struct `Controller` |
| Penanganan error | Manual (`c.JSON(...)` / `c.AbortWithError`) | Cukup `return err`, ditangani error handler terpusat | Cukup `return err`, ditangani error handler terpusat | Manual di dalam method controller |
| Berbasis `net/http` | Ya | Ya | **Tidak** (`fasthttp`) | Ya |
| Gaya | Minimalis, router + middleware | Minimalis, router + middleware | Mirip Express.js | *Full-stack* MVC dengan ORM dan CLI |

Contoh endpoint yang sama di ketiga framework tambahan ada di bagian Praktik.

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

Memilih di antara opsi realtime:

- **SSE** cukup bila hanya server yang mengirim data (notifikasi, progress, feed). Ia memakai HTTP biasa sehingga ramah proxy dan load balancer, dan klien (`EventSource` di browser) menyambung ulang otomatis.
- **WebSocket langsung** (`coder/websocket`, `gorilla/websocket`) bila Anda ingin kendali penuh atas protokol dan manajemen koneksi.
- **Melody** bila Anda ingin WebSocket dengan *broadcast* dan *session* tanpa menulis sendiri manajemen kumpulan koneksi. Cocok untuk aplikasi kecil sampai menengah dalam satu proses.
- **Centrifugo** bila koneksi realtime perlu dipisahkan dari backend utama: jumlah koneksi besar, banyak instance backend, atau backend bukan Go. Server Centrifugo menahan semua koneksi klien; backend Anda hanya memanggil API HTTP/gRPC untuk mempublikasikan pesan ke sebuah *channel*, dan Centrifugo yang menyebarkannya ke semua pelanggan.

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

### CLI dengan urfave/cli

Versi yang sama dengan contoh Cobra di atas, ditulis dengan `urfave/cli` v3. Seluruh definisi perintah berupa satu struktur data.

```go
// urfave_main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "sapa",
		Usage: "Aplikasi penyapa",
		Commands: []*cli.Command{
			{
				Name:      "halo",
				Usage:     "Menyapa seseorang",
				ArgsUsage: "[nama]",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "ulang",
						Aliases: []string{"n"},
						Value:   1,
						Usage:   "jumlah pengulangan",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					nama := c.Args().First()
					if nama == "" {
						return fmt.Errorf("nama wajib diisi")
					}
					for i := 0; i < int(c.Int("ulang")); i++ {
						fmt.Printf("Halo, %s!\n", nama)
					}
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
```

```bash
go mod init sapa && go get github.com/urfave/cli/v3@latest
go run urfave_main.go halo Andi -n 2
go run urfave_main.go --help
```

Perhatikan perbedaannya dengan Cobra: perintah dan flag dideklarasikan sebagai nilai di dalam struct (bukan dirangkai dengan pemanggilan method), dan `Action` menerima `context.Context` sehingga pembatalan dapat diteruskan ke pekerjaan di dalamnya. API urfave/cli berubah antara versi mayor (v2 dan v3); pastikan contoh sesuai dengan versi yang Anda pasang.

### TUI dengan Bubble Tea

Aplikasi terminal interaktif sederhana: sebuah penghitung yang diubah dengan tombol panah.

```go
// tui_main.go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	hitung int
}

// Init dipanggil sekali saat program mulai; nil berarti tidak ada pekerjaan awal.
func (m model) Init() tea.Cmd { return nil }

// Update menerima peristiwa dan mengembalikan state baru.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			m.hitung++
		case "down", "j":
			m.hitung--
		}
	}
	return m, nil
}

// View merender state menjadi string.
func (m model) View() string {
	return fmt.Sprintf("\nHitungan: %d\n\n↑/↓ (atau k/j) mengubah nilai · q keluar\n", m.hitung)
}

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
```

```bash
go mod init tui && go get github.com/charmbracelet/bubbletea@latest
go run tui_main.go
```

Contoh memakai jalur impor `github.com/charmbracelet/bubbletea` (v1). Bila Anda memasang versi yang lebih baru dengan jalur impor berbeda, ikuti panduan migrasi pada repositori Bubble Tea. Karena `Model` hanyalah tipe Go biasa, `Update` dapat diuji tanpa terminal: panggil `Update(tea.KeyMsg{...})` lalu periksa `View()`.

### Echo, Fiber, dan Beego: Endpoint yang Sama

Ketiga program berikut membuat endpoint `GET /halo/{nama}` yang sama dengan contoh Gin di atas, sehingga perbedaan gaya mudah dibandingkan. Jalankan masing-masing di direktori dan modul terpisah.

```go
// echo_main.go
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Recover()) // panic di handler menjadi 500, bukan proses mati

	e.GET("/halo/:nama", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"pesan": "Halo, " + c.Param("nama")})
	})

	e.Logger.Fatal(e.Start(":8081"))
}
```

```go
// fiber_main.go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/halo/:nama", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"pesan": "Halo, " + c.Params("nama")})
	})

	log.Fatal(app.Listen(":8082"))
}
```

```go
// beego_main.go
package main

import (
	"github.com/beego/beego/v2/server/web"
)

type HaloController struct {
	web.Controller
}

// Get menangani GET pada route yang dipetakan ke controller ini.
func (c *HaloController) Get() {
	nama := c.Ctx.Input.Param(":nama")
	c.Data["json"] = map[string]string{"pesan": "Halo, " + nama}
	_ = c.ServeJSON()
}

func main() {
	web.Router("/halo/:nama", &HaloController{})
	web.Run(":8083")
}
```

```bash
go mod init contoh && go get github.com/labstack/echo/v4 github.com/gofiber/fiber/v2 github.com/beego/beego/v2@latest
go run echo_main.go &
curl http://localhost:8081/halo/Andi
```

Hal yang perlu dicermati dari ketiganya:

- **Echo**: handler mengembalikan `error`, sehingga `return err` cukup dan penanganan dilakukan satu tempat oleh error handler Echo.
- **Fiber**: tampilannya hampir sama dengan Echo, tetapi `*fiber.Ctx` milik `fasthttp`. Middleware dan handler `net/http` standar tidak bisa dipasang langsung; Anda butuh adaptor. Signature handler berbeda antar versi mayor Fiber, jadi periksa dokumentasi versi yang Anda pasang.
- **Beego**: berorientasi controller (MVC). Satu struct dapat memiliki method `Get`, `Post`, `Put`, `Delete`, dan data dikirim lewat `c.Data`. Cocok bila tim terbiasa dengan gaya MVC dan menginginkan ORM serta CLI scaffolding (`bee`) dalam satu paket.

### WebSocket dengan Melody dan Publikasi ke Centrifugo

**Melody**: sebuah ruang obrolan sederhana. Setiap pesan yang diterima disiarkan ke seluruh koneksi.

```go
// melody_main.go
package main

import (
	"log"
	"net/http"

	"github.com/olahol/melody"
)

func main() {
	m := melody.New()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if err := m.HandleRequest(w, r); err != nil {
			log.Println("ws:", err)
		}
	})

	m.HandleConnect(func(s *melody.Session) {
		_ = m.Broadcast([]byte("seseorang bergabung"))
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		_ = m.Broadcast(msg) // siarkan ke semua koneksi
	})

	log.Fatal(http.ListenAndServe(":8084", nil))
}
```

```bash
go mod init obrolan && go get github.com/olahol/melody@latest
go run melody_main.go
# di dua terminal berbeda (butuh websocat atau wscat):
websocat ws://localhost:8084/ws
```

Ketik pesan di satu terminal dan pesan itu muncul di terminal lainnya. Seluruh koneksi hidup di dalam proses ini, sehingga bila Anda menjalankan dua instance server, pengguna di instance A tidak melihat pesan pengguna di instance B. Itulah masalah yang diselesaikan Centrifugo.

**Centrifugo**: server realtime terpisah. Backend Go tidak memegang koneksi klien; ia hanya mempublikasikan pesan lewat HTTP API.

```bash
# 1. Buat konfigurasi awal (berisi kunci API dan rahasia token), lalu jalankan server
docker run --rm -v "$PWD":/centrifugo centrifugo/centrifugo:v6 centrifugo genconfig
docker run --rm -v "$PWD":/centrifugo -p 8000:8000 centrifugo/centrifugo:v6 centrifugo --config=config.json
```

```go
// publish_main.go: backend memublikasikan pesan ke channel "obrolan"
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := strings.NewReader(`{"channel":"obrolan","data":{"teks":"halo dari Go"}}`)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:8000/api/publish", body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", os.Getenv("CENTRIFUGO_API_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	hasil, _ := io.ReadAll(resp.Body)
	fmt.Println(resp.Status, string(hasil))
}
```

```bash
export CENTRIFUGO_API_KEY="<nilai http_api key dari config.json>"
go run publish_main.go
```

Pesan yang dipublikasikan akan diterima oleh semua klien yang berlangganan channel `obrolan` lewat WebSocket atau SSE Centrifugo. Periksa isi balasan, bukan hanya status HTTP: Centrifugo dapat membalas `200` dengan objek `error` di dalam JSON (misalnya channel tidak dikenal). Format berkas konfigurasi dan opsi `genconfig` dapat berubah antar versi mayor, jadi ikuti dokumentasi Centrifugo untuk versi yang Anda pakai.

## Latihan

1. Perluas REST API dengan `PUT /items/{id}` dan `DELETE /items/{id}`. Pastikan `DELETE` pada id yang tidak ada mengembalikan `404`, dan `PUT` memvalidasi input seperti `POST`.
2. Tambahkan middleware `recoverPanic` yang menangkap panic di handler, mencatatnya dengan `slog`, dan mengembalikan `500`. Uji dengan handler yang sengaja `panic`. Hubungkan dengan konsep `recover` di [Modul 10](10-error-handling.md).
3. Tambahkan middleware autentikasi sederhana yang memeriksa header `Authorization: Bearer <token>` dan menolak dengan `401`. Jelaskan mengapa middleware berbentuk `func(http.Handler) http.Handler` mudah dirangkai.
4. Tulis uji `httptest` ([Modul 15](15-testing-benchmarking.md)) untuk seluruh endpoint API di atas, termasuk kasus error dan `405`. Jalankan dengan `go test -race` untuk memastikan `RWMutex` bekerja benar.
5. Tambahkan graceful shutdown ke server ([Modul 13](13-sync-context-patterns.md)) dan buktikan request yang sedang berjalan diselesaikan sebelum proses berhenti.
6. Buat klien yang memanggil API di atas dengan `http.Client` ber-timeout dan **retry** terbatas (maksimal 3 kali, jeda bertambah) untuk status `5xx`. Mengapa tidak semua error aman di-retry?
7. Diskusikan: sebuah tim ingin memakai Gin untuk layanan baru, sementara rekan lain bersikeras memakai `net/http` saja sejak Go 1.22. Susun kriteria keputusan (kebutuhan fitur, ekosistem middleware, keterikatan framework, performa, keterampilan tim) dan tentukan pilihan Anda beserta alasannya.
8. Tulis ulang CLI `sapa` dengan satu subcommand tambahan (`selamat-pagi`) memakai Cobra, lalu memakai urfave/cli. Bandingkan jumlah baris, cara mendeklarasikan flag, dan tampilan `--help`. Kapan Anda memilih masing-masing?
9. Kembangkan aplikasi Bubble Tea menjadi daftar tugas: tampilkan beberapa item, pindahkan kursor dengan panah, tandai selesai dengan spasi. Tulis uji yang memanggil `Update` dengan pesan tombol dan memeriksa `View()` tanpa membuka terminal.
10. Implementasikan endpoint `GET /halo/{nama}` dengan `net/http`, Gin, Echo, Fiber, dan Beego. Buat tabel perbandingan: jumlah baris, cara menangani error, dan apakah middleware `net/http` standar dapat dipasang. Mengapa Fiber memerlukan adaptor?
11. Jalankan dua instance server Melody di port berbeda dan buktikan pesan tidak lintas instance. Lalu pasang Centrifugo, publikasikan pesan dari Go, dan jelaskan bagaimana arsitektur ini menyelesaikan masalah tersebut. Apa yang tetap menjadi tanggung jawab backend Anda (autentikasi, otorisasi channel)?

---

[← Modul 15: Testing & Benchmarking](15-testing-benchmarking.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 17: gRPC, Database & Logging →](17-grpc-database-logging.md)
