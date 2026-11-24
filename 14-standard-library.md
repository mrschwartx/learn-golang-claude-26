# Modul 14: Standard Library

[← Modul 13: Sync, Context & Concurrency Patterns](13-sync-context-patterns.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 15: Testing & Benchmarking →](15-testing-benchmarking.md)

**Isi modul ini**

- [Teori](#teori)
  - [14.1 Peta Standard Library](#141-peta-standard-library)
  - [14.2 I/O dan File Handling](#142-io-dan-file-handling)
  - [14.3 Paket `bufio`](#143-paket-bufio)
  - [14.4 Paket `flag`](#144-paket-flag)
  - [14.5 Paket `time`](#145-paket-time)
  - [14.6 Paket `encoding/json`](#146-paket-encodingjson)
  - [14.7 Paket `log/slog`](#147-paket-logslog)
  - [14.8 Paket `regexp`](#148-paket-regexp)
  - [14.9 `go:embed`](#149-goembed)
- [Praktik](#praktik)
  - [File I/O dan `bufio`](#file-io-dan-bufio)
  - [CLI Sederhana dengan `flag`](#cli-sederhana-dengan-flag)
  - [`time` dan `encoding/json`](#time-dan-encodingjson)
  - [`slog` dan `regexp`](#slog-dan-regexp)
  - [Menyematkan Berkas dengan `go:embed`](#menyematkan-berkas-dengan-goembed)
- [Latihan](#latihan)

---


## Teori

### 14.1 Peta Standard Library

Standard library Go sering disebut **"batteries included"** — banyak kebutuhan sehari-hari bisa dipenuhi tanpa library pihak ketiga. Dokumentasi lengkap ada di **https://pkg.go.dev/std**.

| Kelompok | Paket Penting |
|---|---|
| Teks dan string | `fmt`, `strings`, `strconv`, `unicode`, `regexp`, `text/template`, `html/template` |
| I/O dan berkas | `io`, `os`, `bufio`, `path/filepath`, `embed`, `io/fs` |
| Data | `encoding/json`, `encoding/csv`, `encoding/xml`, `encoding/base64`, `encoding/binary` |
| Waktu dan sistem | `time`, `flag`, `os/signal`, `os/exec`, `runtime` |
| Koleksi | `slices`, `maps`, `sort`, `container/heap`, `container/list` |
| Jaringan | `net`, `net/http`, `net/url`, `net/http/httptest` |
| Konkurensi | `sync`, `sync/atomic`, `context` |
| Logging dan observabilitas | `log`, `log/slog`, `runtime/pprof`, `expvar` |
| Keamanan | `crypto/*` (`sha256`, `rand`, `tls`, `bcrypt` ada di `x/crypto`) |
| Testing | `testing`, `testing/fstest`, `testing/iotest` |

Standard library dirawat tim Go dan tunduk pada **Go 1 compatibility promise** ([Modul 1](01-pengenalan.md)), sehingga umumnya lebih stabil daripada dependensi pihak ketiga yang setara.

### 14.2 I/O dan File Handling

Seluruh I/O di Go dibangun di atas dua interface kecil: **`io.Reader`** dan **`io.Writer`** ([Modul 8](08-methods-interfaces.md)). Berkas, koneksi jaringan, body HTTP, buffer memori, dan kompresi semuanya memenuhinya, sehingga dapat dirangkai bebas.

| Kebutuhan | Fungsi |
|---|---|
| Baca seluruh berkas kecil | `os.ReadFile(path)` → `[]byte` |
| Tulis seluruh berkas | `os.WriteFile(path, data, 0o644)` |
| Buka untuk baca / buat baru | `os.Open(path)` / `os.Create(path)` — **selalu `defer f.Close()`** |
| Salin dari reader ke writer | `io.Copy(dst, src)` |
| Baca semua dari reader | `io.ReadAll(r)` |
| Daftar isi direktori | `os.ReadDir(dir)` |
| Gabung path lintas OS | `filepath.Join("data", "file.txt")` |
| Direktori sementara | `os.MkdirTemp`, `os.CreateTemp` |
| Argumen dan environment | `os.Args`, `os.Getenv`, `os.LookupEnv` |
| Akses terbatas pada satu direktori | `os.OpenRoot` / `os.Root` (Go 1.24+) — mencegah path traversal keluar dari root |

Perintah `os.ReadFile` memuat seluruh isi ke memori; untuk berkas besar, baca bertahap dengan `bufio` atau `io.Copy`.

### 14.3 Paket `bufio`

**`bufio`** menambahkan **buffer** di atas reader/writer sehingga operasi kecil yang berulang tidak memicu banyak system call.

| Tipe | Kegunaan |
|---|---|
| `bufio.Scanner` | Membaca **per baris** (atau per kata) — pola utama membaca berkas teks atau `os.Stdin` |
| `bufio.Reader` | Membaca fleksibel: `ReadString('\n')`, `ReadRune`, `Peek` |
| `bufio.Writer` | Menulis berbuffer; **wajib memanggil `Flush()`** di akhir atau data bisa hilang |

Catatan: `Scanner` membatasi panjang satu baris sampai **64 KB** secara default; untuk baris lebih panjang gunakan `scanner.Buffer(buf, maks)`. Selalu cek `scanner.Err()` setelah loop.

### 14.4 Paket `flag`

`flag` mem-parsing argumen baris perintah:

```go
nama := flag.String("nama", "dunia", "nama yang disapa")
ulang := flag.Int("n", 1, "jumlah pengulangan")
verbose := flag.Bool("v", false, "mode verbose")
flag.Parse()
args := flag.Args() // argumen non-flag sisanya
```

Bentuk pemakaian: `./app -nama Andi -n 3 -v file1 file2`. Untuk **subcommand** (`app serve`, `app migrate`) gunakan `flag.NewFlagSet` per subcommand. Untuk CLI yang lebih kaya, lihat Cobra dan urfave/cli di [Modul 16](16-cli-web-development.md).

### 14.5 Paket `time`

| Konsep | Penjelasan |
|---|---|
| `time.Time` | Satu titik waktu (dengan zona waktu dan *monotonic clock*) |
| `time.Duration` | Selisih waktu dalam nanodetik; tulis dengan konstanta: `150 * time.Millisecond` |
| `time.Now()`, `time.Since(t)` | Waktu sekarang dan durasi sejak `t` |
| `time.Sleep(d)` | Menunda goroutine saat ini |
| `time.NewTimer`, `time.NewTicker` | Timer satu kali dan pewaktu berulang (`Stop()` setelah selesai) |
| `t.Add(d)`, `t.Sub(u)`, `t.Before(u)`, `t.Equal(u)` | Aritmetika dan perbandingan waktu |

**Format waktu** Go unik: bukan `YYYY-MM-DD`, tetapi **tanggal referensi** `Mon Jan 2 15:04:05 MST 2006` (angka `1 2 3 4 5 6 7`). Anda menulis layout dengan menunjukkan bagaimana tanggal referensi itu akan terlihat:

| Layout | Contoh Hasil |
|---|---|
| `2006-01-02` | `2026-11-23` |
| `2006-01-02 15:04:05` | `2026-11-23 14:30:00` |
| `02/01/2006` | `23/11/2026` |
| `time.RFC3339` | `2026-11-23T14:30:00Z` |

Jangan membandingkan `time.Time` dengan `==` — gunakan `t.Equal(u)`. Simpan waktu dalam **UTC** dan konversi ke zona lokal hanya saat ditampilkan.

### 14.6 Paket `encoding/json`

| Fungsi | Kegunaan |
|---|---|
| `json.Marshal(v)` / `json.MarshalIndent(v, "", "  ")` | Struct/map/slice → JSON |
| `json.Unmarshal(data, &v)` | JSON → struct/map (argumen kedua **harus pointer**) |
| `json.NewEncoder(w).Encode(v)` / `json.NewDecoder(r).Decode(&v)` | Versi **streaming** (misalnya body HTTP) |
| `json.RawMessage` | Menunda decode sebagian JSON |
| `MarshalJSON` / `UnmarshalJSON` | Method untuk **kustomisasi** encode/decode tipe |
| `decoder.DisallowUnknownFields()` | Menolak field tak dikenal (validasi lebih ketat) |

Catatan penting:

- Hanya **field exported** yang ikut; atur nama dengan tag ([Modul 4](04-composite-types.md)).
- Decode ke `map[string]any` memberi angka bertipe **`float64`**, dan objek bersarang menjadi `map[string]any`.
- **Go 1.27 menambahkan `encoding/json/v2`** (dengan `encoding/json/jsontext` untuk pemrosesan level rendah): opsi lebih kaya dan default lebih ketat (menolak UTF-8 tidak valid dan nama duplikat). Paket `encoding/json` yang ada tetap kompatibel dan kini didukung implementasi v2 sehingga decode lebih cepat.

### 14.7 Paket `log/slog`

**`log/slog`** (Go 1.21+) adalah paket **structured logging** resmi: log berupa pesan plus pasangan key-value yang mudah diproses mesin (Elasticsearch, Loki, dsb.).

| Komponen | Keterangan |
|---|---|
| Level | `Debug`, `Info`, `Warn`, `Error` |
| Handler | `slog.NewTextHandler` (mudah dibaca) atau `slog.NewJSONHandler` (untuk produksi) |
| Atribut | `slog.Info("login", "user", "andi", "ip", ip)` atau tipe aman `slog.String`, `slog.Int` |
| `logger.With(...)` | Logger turunan dengan atribut tetap (misalnya `request_id`) |
| `slog.Group` | Mengelompokkan atribut |
| `slog.SetDefault(l)` | Menjadikan logger default untuk seluruh aplikasi |
| `...Context` (`InfoContext`) | Menyertakan `context.Context` agar handler bisa mengambil metadata |

Pustaka berperforma ekstrem seperti Zap dan Zerolog dibahas di [Modul 17](17-grpc-database-logging.md).

### 14.8 Paket `regexp`

Regex Go memakai mesin **RE2** — dijamin berjalan dalam waktu **linear** terhadap input, sehingga kebal terhadap serangan *ReDoS*. Konsekuensinya: tidak mendukung *backreference* dan *lookahead*.

| Fungsi | Kegunaan |
|---|---|
| `regexp.MustCompile(pola)` | Kompilasi (panic jika pola salah); simpan di variabel level paket agar tidak dikompilasi berulang |
| `re.MatchString(s)` | Cocok atau tidak |
| `re.FindString`, `re.FindAllString(s, -1)` | Ambil kecocokan |
| `re.FindStringSubmatch(s)` | Ambil grup tangkap; mendukung **named group** `(?P<nama>...)` |
| `re.ReplaceAllString(s, "$1")` | Ganti dengan referensi grup |

Pakai **raw string** (`` `\d+` ``) untuk pola agar backslash tidak perlu digandakan ([Modul 3](03-variabel-konstanta-tipe-data.md)).

### 14.9 `go:embed`

Direktif **`//go:embed`** menyematkan berkas ke dalam binary saat kompilasi, sehingga aplikasi tetap **satu file** meski membutuhkan template, aset statis, atau migrasi SQL.

| Tipe variabel | Hasil |
|---|---|
| `string` | Isi berkas sebagai teks |
| `[]byte` | Isi berkas sebagai byte |
| `embed.FS` | Sistem berkas read-only untuk banyak berkas/direktori (pola seperti `static/*` atau `templates`) |

Syarat: variabel dideklarasikan di level paket, `import "embed"` (boleh blank import untuk `string`/`[]byte`), dan berkas berada di dalam direktori modul/paket (tidak boleh naik dengan `..`). `embed.FS` memenuhi `fs.FS`, sehingga bisa langsung dipakai `http.FileServerFS` dan `template.ParseFS`.

## Praktik

### File I/O dan `bufio`

```go
// berkas.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dir, err := os.MkdirTemp("", "belajar-go-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "catatan.txt")
	isi := "baris satu\nbaris dua\nbaris tiga\n"
	if err := os.WriteFile(path, []byte(isi), 0o644); err != nil {
		panic(err)
	}

	data, _ := os.ReadFile(path)
	fmt.Print(string(data))

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	n := 0
	for sc.Scan() {
		n++
		fmt.Printf("%d: %s\n", n, strings.ToUpper(sc.Text()))
	}
	if err := sc.Err(); err != nil {
		panic(err)
	}

	// bufio.Writer: jangan lupa Flush
	out, _ := os.Create(filepath.Join(dir, "keluar.txt"))
	w := bufio.NewWriter(out)
	fmt.Fprintln(w, "ditulis lewat buffer")
	w.Flush()
	out.Close()

	entri, _ := os.ReadDir(dir)
	for _, e := range entri {
		fmt.Println("berkas:", e.Name())
	}

	if _, err := os.Open("tidak-ada.txt"); err != nil {
		fmt.Println(os.IsNotExist(err))
	}
}
```

### CLI Sederhana dengan `flag`

```go
// sapa.go
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	nama := flag.String("nama", "dunia", "nama yang disapa")
	ulang := flag.Int("n", 1, "jumlah pengulangan")
	besar := flag.Bool("besar", false, "cetak dengan huruf besar")
	flag.Parse()

	for range *ulang {
		pesan := "Halo, " + *nama + "!"
		if *besar {
			pesan = strings.ToUpper(pesan)
		}
		fmt.Println(pesan)
	}
	if flag.NArg() > 0 {
		fmt.Fprintln(os.Stderr, "argumen sisa:", flag.Args())
	}
}
```

```bash
go run sapa.go -nama Andi -n 2 -besar extra1 extra2
go run sapa.go -h
```

### `time` dan `encoding/json`

```go
// waktu_json.go
package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Status int

const (
	Aktif Status = iota + 1
	Nonaktif
)

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[Status]string{Aktif: "aktif", Nonaktif: "nonaktif"}[s])
}

type Pengguna struct {
	Nama   string    `json:"nama"`
	Umur   int       `json:"umur,omitempty"`
	Status Status    `json:"status"`
	Lahir  time.Time `json:"lahir"`
}

func main() {
	t := time.Date(2026, time.November, 23, 14, 30, 0, 0, time.UTC)
	fmt.Println(t.Format("2006-01-02 15:04:05"), t.Format("02/01/2006"), t.Format(time.RFC3339))
	fmt.Println(t.Add(36*time.Hour).Weekday(), t.Sub(t.Add(-90*time.Minute)))

	parsed, err := time.Parse("2006-01-02", "2027-01-15")
	fmt.Println(parsed.Sub(t).Round(time.Hour), err)

	p := Pengguna{Nama: "Andi", Status: Aktif, Lahir: time.Date(1995, 5, 17, 0, 0, 0, 0, time.UTC)}
	b, _ := json.MarshalIndent(p, "", "  ")
	fmt.Println(string(b))

	var umum map[string]any
	_ = json.Unmarshal([]byte(`{"id":7,"tag":["a","b"],"profil":{"kota":"Jakarta"}}`), &umum)
	fmt.Printf("%v %T\n", umum["id"], umum["id"])
	fmt.Println(umum["profil"].(map[string]any)["kota"])

	dec := json.NewDecoder(strings.NewReader(`{"nama":"Budi","tidak_dikenal":1}`))
	dec.DisallowUnknownFields()
	var q struct {
		Nama string `json:"nama"`
	}
	fmt.Println(dec.Decode(&q))
}
```

### `slog` dan `regexp`

```go
// log_regex.go
package main

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"
)

var reTanggal = regexp.MustCompile(`(?P<tahun>\d{4})-(?P<bulan>\d{2})-(?P<hari>\d{2})`)

func main() {
	opsi := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{} // buang waktu agar keluaran deterministik
			}
			return a
		},
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opsi))
	logger.Info("server mulai", "port", 8080)
	req := logger.With("request_id", "abc-123")
	req.Warn("lambat", slog.Int("ms", 1500), slog.Group("user", "id", 7, "role", "admin"))
	req.Error("gagal simpan", "err", fmt.Errorf("koneksi ditolak"))

	teks := "Rilis pada 2026-08-19 dan 2027-02-01."
	fmt.Println(reTanggal.MatchString(teks))
	fmt.Println(reTanggal.FindAllString(teks, -1))
	m := reTanggal.FindStringSubmatch(teks)
	fmt.Println(m[reTanggal.SubexpIndex("tahun")], m[reTanggal.SubexpIndex("bulan")])
	fmt.Println(reTanggal.ReplaceAllString(teks, "${hari}/${bulan}/${tahun}"))
}
```

### Menyematkan Berkas dengan `go:embed`

```bash
mkdir embed-demo && cd embed-demo
go mod init embeddemo
mkdir templates
echo "Halo dari berkas yang tertanam!" > salam.txt
echo "<h1>Beranda</h1>" > templates/index.html
```

```go
// main.go
package main

import (
	"embed"
	"fmt"
)

//go:embed salam.txt
var salam string

//go:embed templates/*
var templates embed.FS

func main() {
	fmt.Print(salam)
	b, _ := templates.ReadFile("templates/index.html")
	fmt.Print(string(b))
	daftar, _ := templates.ReadDir("templates")
	fmt.Println(len(daftar), daftar[0].Name())
}
```

```bash
go build -o app . && rm salam.txt && rm -r templates && ./app
```

Binary tetap berjalan meski berkas aslinya sudah dihapus — isinya sudah ada di dalam executable.

## Latihan

1. Tulis program yang membaca berkas teks dengan `bufio.Scanner`, menghitung jumlah baris, kata, dan karakter (rune), lalu mencetak 5 kata paling sering muncul. Mengapa `os.ReadFile` kurang cocok untuk berkas puluhan gigabyte?
2. Perluas program `flag` menjadi CLI dengan dua subcommand (`tambah` dan `hapus`) memakai `flag.NewFlagSet`, masing-masing dengan flag sendiri.
3. Hitung selisih hari antara hari ini dan tanggal 17 Agustus tahun depan, lalu cetak tanggal itu dengan tiga layout berbeda. Jelaskan mengapa Go memakai tanggal referensi `2006-01-02 15:04:05` alih-alih `YYYY-MM-DD`.
4. Buat tipe `Rupiah int` dengan `MarshalJSON` yang menghasilkan string `"Rp 15.000"` dan `UnmarshalJSON` yang membacanya kembali. Uji bolak-balik (round trip).
5. Ganti `slog.NewJSONHandler` dengan `NewTextHandler`, tambahkan `logger.With("request_id", ...)`, dan atur level lewat variabel lingkungan `LOG_LEVEL`. Mengapa log terstruktur lebih mudah dicari di Elasticsearch (course Elasticsearch) dibanding log teks bebas?
6. Tulis regex untuk memvalidasi dan mengekstrak bagian email (`nama` dan `domain`) dengan named group, lalu uji dengan 5 input valid dan tidak valid. Mengapa regex Go aman dari ReDoS, dan fitur apa yang harus dikorbankan?
7. Diskusikan: sebuah aplikasi web membutuhkan template HTML, file CSS, dan migrasi SQL saat berjalan di container. Jelaskan manfaat `go:embed` dibanding menyalin folder aset ke image Docker (course Docker), serta kapan aset sebaiknya tetap berada di luar binary.

---

[← Modul 13: Sync, Context & Concurrency Patterns](13-sync-context-patterns.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 15: Testing & Benchmarking →](15-testing-benchmarking.md)
