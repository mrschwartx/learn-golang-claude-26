# Modul 15: Testing & Benchmarking

[← Modul 14: Standard Library](14-standard-library.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 16: Membangun CLI & Web Development →](16-cli-web-development.md)

**Isi modul ini**

- [Teori](#teori)
  - [15.1 Dasar Paket `testing`](#151-dasar-paket-testing)
  - [15.2 Table-Driven Tests](#152-table-driven-tests)
  - [15.3 Mocks dan Stubs](#153-mocks-dan-stubs)
  - [15.4 `httptest`](#154-httptest)
  - [15.5 Benchmarks](#155-benchmarks)
  - [15.6 Coverage](#156-coverage)
  - [15.7 Fuzzing, Example, dan `synctest`](#157-fuzzing-example-dan-synctest)
  - [15.8 Praktik Baik](#158-praktik-baik)
- [Praktik](#praktik)
  - [Table-Driven Test](#table-driven-test)
  - [Benchmark, Fuzz, dan Example](#benchmark-fuzz-dan-example)
  - [Fake untuk Interface](#fake-untuk-interface)
  - [`httptest` untuk Handler dan Klien](#httptest-untuk-handler-dan-klien)
  - [Menjalankan Seluruhnya](#menjalankan-seluruhnya)
- [Latihan](#latihan)

---


## Teori

### 15.1 Dasar Paket `testing`

Testing adalah **bagian dari toolchain Go**, bukan framework tambahan. Konvensinya:

| Aturan | Penjelasan |
|---|---|
| Berkas bernama **`*_test.go`** | Hanya dikompilasi saat `go test`, tidak masuk binary produksi |
| Fungsi **`TestXxx(t *testing.T)`** | Fungsi uji; `Xxx` diawali huruf kapital |
| Letak berkas | Di **direktori yang sama** dengan kode yang diuji |
| Package | `package harga` (uji *white-box*, bisa mengakses identifier privat) atau `package harga_test` (uji *black-box*, hanya API publik) |

Method penting pada `*testing.T`:

| Method | Efek |
|---|---|
| `t.Error`, `t.Errorf` | Menandai gagal, **lanjut** menjalankan uji |
| `t.Fatal`, `t.Fatalf` | Menandai gagal dan **menghentikan** uji ini seketika |
| `t.Log` | Mencatat pesan (terlihat dengan `-v` atau saat gagal) |
| `t.Helper()` | Menandai fungsi bantu agar baris kegagalan menunjuk pemanggilnya |
| `t.Cleanup(f)` | Mendaftarkan pembersihan yang dijalankan di akhir uji |
| `t.TempDir()` | Direktori sementara yang dihapus otomatis |
| `t.Setenv(k, v)` | Mengatur environment variable selama uji |
| `t.Context()` | `context.Context` yang dibatalkan saat uji selesai (Go 1.24+) |
| `t.Skip(...)` | Melewati uji (misalnya uji integrasi) |

Perintah dasar:

| Perintah | Fungsi |
|---|---|
| `go test ./...` | Jalankan semua uji di seluruh modul |
| `go test -v` | Tampilkan setiap uji |
| `go test -run TestDiskon/negatif` | Jalankan uji yang cocok dengan ekspresi reguler |
| `go test -count=1` | Abaikan cache hasil uji |
| `go test -race` | Aktifkan race detector ([Modul 13](13-sync-context-patterns.md)) |
| `go test -short` | Pertanda untuk melewati uji lambat (`testing.Short()`) |

Perlu pengaturan sekali untuk seluruh package? Definisikan **`TestMain(m *testing.M)`** untuk setup dan teardown global.

### 15.2 Table-Driven Tests

Gaya khas Go: kumpulkan **kasus uji dalam slice struct**, lalu jalankan dalam satu loop. Menambah kasus baru cukup menambah satu baris.

```go
tests := []struct {
    nama      string
    harga, persen int
    want      int
    wantErr   bool
}{
    {"normal", 1000, 10, 900, false},
    {"nol persen", 1000, 0, 1000, false},
    {"persen negatif", 1000, -5, 0, true},
}
for _, tc := range tests {
    t.Run(tc.nama, func(t *testing.T) { /* ... */ })
}
```

`t.Run` membuat **subtest** yang bisa dijalankan terpisah (`-run TestDiskon/normal`). Tambahkan `t.Parallel()` di dalam subtest agar berjalan paralel — aman dengan variabel loop berkat semantik per-iterasi Go 1.22 ([Modul 5](05-control-flow.md)).

Go tidak menyertakan fungsi `assert`; perbandingan ditulis dengan `if got != want`. Untuk struktur kompleks, `reflect.DeepEqual` atau paket `slices`/`maps` (`slices.Equal`) membantu; library populer `testify` dan `go-cmp` menambah kenyamanan jika diperlukan.

### 15.3 Mocks dan Stubs

Karena interface dipenuhi secara implisit dan idealnya didefinisikan di sisi pemakai ([Modul 8](08-methods-interfaces.md)), **mock cukup berupa struct kecil buatan sendiri** — tanpa framework.

| Istilah | Fungsi |
|---|---|
| **Stub** | Mengembalikan jawaban tetap |
| **Fake** | Implementasi sederhana yang bekerja (misalnya penyimpanan in-memory) |
| **Spy/Mock** | Mencatat panggilan agar bisa diperiksa uji |

Alat pendukung bila hasil tulis tangan terlalu banyak: `go.uber.org/mock` (mockgen), `mockery`, dan `testify/mock`. Prinsipnya: **buat interface sekecil mungkin** agar fake mudah ditulis.

### 15.4 `httptest`

Paket `net/http/httptest` membantu menguji HTTP tanpa jaringan sungguhan:

| Alat | Kegunaan |
|---|---|
| `httptest.NewRecorder()` | `ResponseWriter` palsu yang merekam status, header, dan body — untuk menguji **handler** langsung |
| `httptest.NewRequest(method, url, body)` | Membuat `*http.Request` untuk uji |
| `httptest.NewServer(handler)` | Menjalankan **server sungguhan di localhost** pada port acak — untuk menguji **klien HTTP** atau alur end-to-end |

Pola umum: bungkus handler, panggil `handler.ServeHTTP(rec, req)`, lalu periksa `rec.Code` dan `rec.Body` ([Modul 16](16-cli-web-development.md)).

### 15.5 Benchmarks

Benchmark mengukur kinerja kode. Fungsi `BenchmarkXxx(b *testing.B)` berada di berkas `_test.go`:

```go
func BenchmarkBalik(b *testing.B) {
    for b.Loop() {          // Go 1.24+
        Balik("Halo, dunia")
    }
}
```

**`b.Loop()`** (Go 1.24+) menggantikan pola lama `for i := 0; i < b.N; i++`. Ia lebih aman: setup di luar loop tidak ikut terhitung, dan compiler tidak menghapus pemanggilan yang hasilnya tak dipakai.

| Perintah | Fungsi |
|---|---|
| `go test -bench=.` | Jalankan semua benchmark (uji biasa ikut berjalan; tambahkan `-run=^$` untuk melewatinya) |
| `go test -bench=Balik -benchmem` | Tampilkan **alokasi memori** per operasi (`B/op`, `allocs/op`) |
| `-benchtime=3s` atau `-benchtime=1000x` | Atur lama atau jumlah iterasi |
| `-count=10` + `benchstat` | Bandingkan sebelum/sesudah secara statistik |
| `-cpuprofile`, `-memprofile` | Hasilkan profil untuk `pprof` ([Modul 19](19-performance-deployment.md)) |

Hasil dibaca seperti `BenchmarkBalik-8   5000000   230 ns/op   64 B/op   2 allocs/op`: iterasi dijalankan sampai stabil, lalu dilaporkan waktu dan alokasi per operasi. **Ukur dulu, baru optimalkan.**

### 15.6 Coverage

| Perintah | Fungsi |
|---|---|
| `go test -cover ./...` | Persentase kode yang tereksekusi uji |
| `go test -coverprofile=c.out ./...` | Simpan profil cakupan |
| `go tool cover -func=c.out` | Ringkasan per fungsi |
| `go tool cover -html=c.out` | Laporan HTML yang menyorot baris tak tercakup |

Cakupan 100% **bukan jaminan** kebenaran — ia hanya menunjukkan baris mana yang belum pernah dijalankan. Gunakan sebagai pemandu, bukan target tunggal.

### 15.7 Fuzzing, Example, dan `synctest`

| Fitur | Penjelasan |
|---|---|
| **Fuzzing** (Go 1.18+) | `FuzzXxx(f *testing.F)` menghasilkan input acak untuk menemukan bug yang tak terpikir; jalankan `go test -fuzz=FuzzXxx`. Kasus gagal tersimpan di `testdata/fuzz/` sebagai regression test |
| **Example tests** | `ExampleXxx()` dengan komentar `// Output: ...` — diverifikasi saat `go test` **dan** tampil sebagai dokumentasi di pkg.go.dev |
| **`testing/synctest`** (Go 1.25) | Menjalankan kode concurrent dalam "gelembung" dengan **jam virtual**, sehingga uji `time.Sleep` dan timeout berjalan instan dan deterministik |

### 15.8 Praktik Baik

1. Uji **perilaku** (apa yang dijanjikan API), bukan detail implementasi privat.
2. Simpan data uji di direktori **`testdata/`** (diabaikan toolchain); gunakan *golden file* untuk output panjang.
3. Pisahkan **uji unit** (cepat, tanpa I/O) dari **uji integrasi** (database, jaringan) memakai `testing.Short()` atau *build tag* ([Modul 18](18-toolchain-code-quality.md)).
4. Jalankan `go test -race ./...` di CI pada setiap perubahan.
5. Buat uji **deterministik**: hindari `time.Now()` dan `rand` tanpa seed; inject jam atau sumber acak lewat parameter.
6. Pesan kegagalan yang baik memuat **input, hasil aktual, dan hasil yang diharapkan**: `t.Errorf("Diskon(%d,%d) = %d; want %d", ...)`.

## Praktik

Kita bangun satu modul kecil berisi empat package beserta ujinya.

```bash
mkdir uji && cd uji
go mod init belajar/uji
mkdir harga stringutil pendaftaran api
```

### Table-Driven Test

```go
// harga/harga.go
package harga

import (
	"errors"
	"fmt"
)

var ErrPersenTidakValid = errors.New("persen di luar 0-100")

// Diskon mengurangi harga dengan persen tertentu.
func Diskon(harga, persen int) (int, error) {
	if persen < 0 || persen > 100 {
		return 0, fmt.Errorf("diskon %d: %w", persen, ErrPersenTidakValid)
	}
	return harga - harga*persen/100, nil
}
```

```go
// harga/harga_test.go
package harga

import (
	"errors"
	"testing"
)

func TestDiskon(t *testing.T) {
	tests := []struct {
		nama          string
		harga, persen int
		want          int
		wantErr       error
	}{
		{"normal", 1000, 10, 900, nil},
		{"nol persen", 1000, 0, 1000, nil},
		{"seratus persen", 1000, 100, 0, nil},
		{"negatif", 1000, -5, 0, ErrPersenTidakValid},
		{"lebih dari 100", 1000, 150, 0, ErrPersenTidakValid},
	}
	for _, tc := range tests {
		t.Run(tc.nama, func(t *testing.T) {
			t.Parallel()
			got, err := Diskon(tc.harga, tc.persen)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("error = %v; want %v", err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("Diskon(%d,%d) = %d; want %d", tc.harga, tc.persen, got, tc.want)
			}
		})
	}
}
```

```bash
go test ./harga -v
go test ./harga -run TestDiskon/negatif -v
go test ./harga -cover
```

### Benchmark, Fuzz, dan Example

```go
// stringutil/stringutil.go
package stringutil

import (
	"slices"
	"strings"
)

// Balik membalik urutan karakter (rune) sebuah string.
func Balik(s string) string {
	r := []rune(s)
	slices.Reverse(r)
	return string(r)
}

// BalikNaif membalik string dengan penyambungan berulang (sengaja tidak efisien).
func BalikNaif(s string) string {
	hasil := ""
	for _, r := range s {
		hasil = string(r) + hasil
	}
	return hasil
}

// BalikBuilder memakai strings.Builder.
func BalikBuilder(s string) string {
	r := []rune(s)
	var sb strings.Builder
	sb.Grow(len(s))
	for i := len(r) - 1; i >= 0; i-- {
		sb.WriteRune(r[i])
	}
	return sb.String()
}
```

```go
// stringutil/stringutil_test.go
package stringutil

import (
	"fmt"
	"testing"
	"unicode/utf8"
)

func TestBalik(t *testing.T) {
	tests := map[string]string{
		"":      "",
		"abc":   "cba",
		"héllo": "olléh",
		"Go 世界": "界世 oG",
	}
	for in, want := range tests {
		if got := Balik(in); got != want {
			t.Errorf("Balik(%q) = %q; want %q", in, got, want)
		}
	}
}

func FuzzBalik(f *testing.F) {
	f.Add("halo")
	f.Add("世界")
	f.Fuzz(func(t *testing.T, s string) {
		if !utf8.ValidString(s) {
			t.Skip()
		}
		if got := Balik(Balik(s)); got != s {
			t.Errorf("Balik(Balik(%q)) = %q", s, got)
		}
	})
}

func ExampleBalik() {
	fmt.Println(Balik("Golang"))
	// Output: gnaloG
}

const teks = "Pemrograman Go itu menyenangkan dan cepat"

func BenchmarkBalik(b *testing.B) {
	for b.Loop() {
		Balik(teks)
	}
}

func BenchmarkBalikNaif(b *testing.B) {
	for b.Loop() {
		BalikNaif(teks)
	}
}

func BenchmarkBalikBuilder(b *testing.B) {
	for b.Loop() {
		BalikBuilder(teks)
	}
}
```

```bash
go test ./stringutil -v
go test ./stringutil -run '^$' -bench . -benchmem
go test ./stringutil -fuzz FuzzBalik -fuzztime 5s
```

### Fake untuk Interface

```go
// pendaftaran/pendaftaran.go
package pendaftaran

import (
	"errors"
	"fmt"
	"strings"
)

// Pengirim didefinisikan di sisi pemakai: hanya butuh satu method.
type Pengirim interface {
	Kirim(ke, pesan string) error
}

type Layanan struct{ Email Pengirim }

func (l Layanan) Daftar(email string) error {
	if !strings.Contains(email, "@") {
		return errors.New("email tidak valid")
	}
	if err := l.Email.Kirim(email, "Selamat datang!"); err != nil {
		return fmt.Errorf("kirim sambutan: %w", err)
	}
	return nil
}
```

```go
// pendaftaran/pendaftaran_test.go
package pendaftaran

import (
	"errors"
	"testing"
)

type pengirimPalsu struct {
	terkirim []string
	err      error
}

func (p *pengirimPalsu) Kirim(ke, pesan string) error {
	p.terkirim = append(p.terkirim, ke)
	return p.err
}

func TestDaftar(t *testing.T) {
	t.Run("sukses", func(t *testing.T) {
		p := &pengirimPalsu{}
		if err := (Layanan{p}).Daftar("a@b.com"); err != nil {
			t.Fatal(err)
		}
		if len(p.terkirim) != 1 || p.terkirim[0] != "a@b.com" {
			t.Errorf("terkirim = %v", p.terkirim)
		}
	})
	t.Run("email tidak valid tidak mengirim", func(t *testing.T) {
		p := &pengirimPalsu{}
		if err := (Layanan{p}).Daftar("salah"); err == nil {
			t.Error("harus error")
		}
		if len(p.terkirim) != 0 {
			t.Error("tidak boleh mengirim")
		}
	})
	t.Run("pengirim gagal", func(t *testing.T) {
		gagal := errors.New("smtp mati")
		err := (Layanan{&pengirimPalsu{err: gagal}}).Daftar("a@b.com")
		if !errors.Is(err, gagal) {
			t.Errorf("err = %v", err)
		}
	})
}
```

### `httptest` untuk Handler dan Klien

```go
// api/api.go
package api

import (
	"encoding/json"
	"net/http"
)

func Halo(w http.ResponseWriter, r *http.Request) {
	nama := r.URL.Query().Get("nama")
	if nama == "" {
		http.Error(w, "nama wajib diisi", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"pesan": "Halo, " + nama})
}
```

```go
// api/api_test.go
package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHaloHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/halo?nama=Andi", nil)
	rec := httptest.NewRecorder()
	Halo(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	var got map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got["pesan"] != "Halo, Andi" {
		t.Errorf("pesan = %q", got["pesan"])
	}
}

func TestHaloTanpaNama(t *testing.T) {
	rec := httptest.NewRecorder()
	Halo(rec, httptest.NewRequest(http.MethodGet, "/halo", nil))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d; want 400", rec.Code)
	}
}

func TestHaloServerSungguhan(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(Halo))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "?nama=Budi")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if string(b) != "{\"pesan\":\"Halo, Budi\"}\n" {
		t.Errorf("body = %q", b)
	}
}
```

### Menjalankan Seluruhnya

```bash
go vet ./...
go test -race -cover ./...
go test -coverprofile=c.out ./... && go tool cover -func=c.out
go tool cover -html=c.out
```

## Latihan

1. Tambahkan kasus uji ke `TestDiskon` untuk harga `0` dan harga negatif, lalu putuskan apakah `Diskon` perlu memvalidasi harga. Tulis uji lebih dulu, baru ubah kodenya (pendekatan TDD).
2. Hapus `t.Parallel()` lalu tambahkan kembali, dan jalankan dengan `-v`. Jelaskan perbedaan urutan eksekusi dan mengapa semantik variabel loop Go 1.22 membuat pola ini aman.
3. Jalankan benchmark `Balik`, `BalikNaif`, dan `BalikBuilder` dengan `-benchmem`. Urutkan dari tercepat, jelaskan perbedaan `allocs/op`, dan simpan hasil dua kali (`-count=10`) untuk dibandingkan dengan `benchstat` setelah satu optimasi.
4. Sengaja masukkan bug kecil ke `Balik` (misalnya membalik per **byte**, bukan per rune), lalu jalankan `go test -fuzz`. Catat input yang ditemukan fuzzer dan lihat berkas di `testdata/fuzz`.
5. Tulis fake `Pengirim` kedua yang menyimpan pesan di map, lalu uji bahwa `Layanan.Daftar` memanggil `Kirim` tepat sekali. Mengapa interface `Pengirim` lebih baik dideklarasikan di package pemakai daripada di package pengirim email?
6. Tambahkan handler `POST /pesan` ke package `api` yang membaca JSON dari body, lalu uji dengan `httptest.NewRequest` untuk body valid, body JSON rusak, dan field kosong. Periksa status code masing-masing.
7. Diskusikan: sebuah tim memiliki cakupan 95% tetapi bug sering lolos ke produksi. Jelaskan batasan metrik cakupan, jenis uji apa yang hilang (integrasi, fuzz, concurrency dengan `-race`), dan bagaimana Anda menyusun strategi uji bertingkat untuk layanan yang terhubung ke database.

---

[← Modul 14: Standard Library](14-standard-library.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 16: Membangun CLI & Web Development →](16-cli-web-development.md)
