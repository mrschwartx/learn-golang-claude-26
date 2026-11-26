# Modul 10: Error Handling

[← Modul 9: Generics & Iterators](09-generics.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 11: Modules, Packages & Code Organization →](11-modules-packages.md)

**Isi modul ini**

- [Teori](#teori)
  - [10.1 Filosofi: Errors Are Values](#101-filosofi-errors-are-values)
  - [10.2 Interface `error`](#102-interface-error)
  - [10.3 `errors.New` dan `fmt.Errorf`](#103-errorsnew-dan-fmterrorf)
  - [10.4 Wrapping dan Unwrapping Errors](#104-wrapping-dan-unwrapping-errors)
  - [10.5 Sentinel Errors](#105-sentinel-errors)
  - [10.6 Custom Error Types](#106-custom-error-types)
  - [10.7 `panic` dan `recover`](#107-panic-dan-recover)
  - [10.8 Stack Traces dan Debugging](#108-stack-traces-dan-debugging)
  - [10.9 Praktik Terbaik](#109-praktik-terbaik)
- [Praktik](#praktik)
  - [Membungkus dan Memeriksa Error](#membungkus-dan-memeriksa-error)
  - [`panic`, `recover`, dan Stack Trace](#panic-recover-dan-stack-trace)
  - [Debugging dengan Delve (Opsional)](#debugging-dengan-delve-opsional)
- [Latihan](#latihan)

---


## Teori

### 10.1 Filosofi: Errors Are Values

Go **tidak memiliki exception** (`try/catch`). Error diperlakukan sebagai **nilai biasa** yang dikembalikan fungsi dan diperiksa secara eksplisit oleh pemanggil. Alur kegagalan menjadi terlihat jelas di kode — tidak ada lompatan tersembunyi dari dalam fungsi yang jauh.

```go
data, err := os.ReadFile("config.json")
if err != nil {
    return fmt.Errorf("baca config: %w", err)
}
```

Pola `if err != nil` memang berulang, dan sering dikritik panjang — namun itulah harga dari eksplisitnya: setiap titik kegagalan wajib diputuskan secara sadar (dikembalikan, dibungkus, dicatat, atau ditangani).

### 10.2 Interface `error`

`error` adalah interface bawaan yang sangat sederhana:

```go
type error interface {
    Error() string
}
```

Tipe apa pun yang memiliki method `Error() string` adalah error ([Modul 8](08-methods-interfaces.md)). Konvensi pentingnya: **`nil` berarti sukses**, dan error adalah **nilai kembalian terakhir** sebuah fungsi. Perhatikan juga jebakan interface `nil` di [Modul 8](08-methods-interfaces.md) — jangan mengembalikan pointer `nil` bertipe konkret sebagai `error`.

### 10.3 `errors.New` dan `fmt.Errorf`

| Fungsi | Kegunaan |
|---|---|
| `errors.New("pesan")` | Membuat error sederhana dengan pesan tetap |
| `fmt.Errorf("format %d", x)` | Membuat error dengan pesan terformat |
| `fmt.Errorf("konteks: %w", err)` | Membuat error yang **membungkus** (*wrap*) error lain |

Konvensi pesan error: **huruf kecil, tanpa tanda baca di akhir**, dan menambah konteks dari lapisan masing-masing sehingga hasil akhirnya terbaca seperti jejak: `simpan pesanan: query database: connection refused`.

### 10.4 Wrapping dan Unwrapping Errors

Dengan verb **`%w`**, `fmt.Errorf` membungkus error asli sehingga **tetap dapat diperiksa** oleh pemanggil di lapisan atas. Paket `errors` menyediakan alat untuk menelusuri rantai bungkusan:

| Fungsi | Kegunaan |
|---|---|
| `errors.Is(err, target)` | Apakah **ada error dalam rantai** yang sama dengan `target`? (pengganti `==`) |
| `errors.As(err, &target)` | Apakah ada error dalam rantai yang **bertipe** tertentu? Jika ya, isi `target` |
| `errors.AsType[T](err)` | Versi generik dari `As` (Go 1.26+): `v, ok := errors.AsType[*MyErr](err)` — tanpa variabel perantara |
| `errors.Unwrap(err)` | Mengambil error satu lapis di dalam (jarang dipakai langsung) |
| `errors.Join(errs...)` | Menggabungkan beberapa error menjadi satu (Go 1.20+) |

Sejak Go 1.20, satu `fmt.Errorf` boleh memakai **beberapa `%w`**, dan `errors.Is`/`As` menelusuri semuanya.

> Membungkus dengan `%w` berarti **membuka detail error di dalamnya sebagai bagian dari API Anda**. Jika Anda tidak ingin pemanggil bergantung pada error internal (misalnya error driver database), bungkus dengan `%v` agar pesannya tetap ada tetapi rantainya terputus.

### 10.5 Sentinel Errors

**Sentinel error** adalah error yang dideklarasikan sebagai variabel paket, dipakai untuk menandai kondisi tertentu yang dapat dikenali pemanggil. Contoh di standard library: `io.EOF`, `sql.ErrNoRows`, `os.ErrNotExist`.

```go
var ErrStokHabis = errors.New("stok habis")

if errors.Is(err, ErrStokHabis) { /* tangani khusus */ }
```

Konvensi: nama diawali **`Err`**. **Selalu bandingkan dengan `errors.Is`**, bukan `==`, karena error biasanya sudah terbungkus. Sentinel menjadi bagian dari API publik paket — tambahkan hanya jika pemanggil memang perlu bereaksi berbeda.

### 10.6 Custom Error Types

Saat pemanggil membutuhkan **data tambahan** dari error (kode status, nama field), buat tipe error sendiri:

```go
type ErrValidasi struct {
    Field string
    Pesan string
}

func (e *ErrValidasi) Error() string {
    return fmt.Sprintf("validasi %s: %s", e.Field, e.Pesan)
}

var ev *ErrValidasi
if errors.As(err, &ev) {
    fmt.Println(ev.Field)
}
```

| Kebutuhan | Pilihan |
|---|---|
| Hanya perlu **tahu jenis** kegagalannya | Sentinel error |
| Perlu **membaca detail** kegagalan | Custom error type |
| Cukup pesan untuk manusia | `fmt.Errorf` biasa |

### 10.7 `panic` dan `recover`

`panic` menghentikan eksekusi normal fungsi, menjalankan semua `defer` yang sudah terdaftar, lalu menjalar ke pemanggil hingga program **crash** (exit code 2) bila tidak ada yang menanganinya. **`recover`** mengembalikan nilai panic dan menghentikan penjalaran — **hanya efektif bila dipanggil langsung di dalam fungsi `defer`**.

| Situasi | Apakah memakai `panic`? |
|---|---|
| Input salah, file tidak ada, jaringan putus | **Tidak** — kembalikan `error` |
| Pelanggaran yang menandakan **bug programmer** (invariant rusak) | Boleh |
| Kegagalan inisialisasi yang membuat program mustahil berjalan | Boleh (`regexp.MustCompile`, `template.Must`) |
| Panic dari runtime (index di luar rentang, nil map, nil pointer) | Terjadi otomatis |

`recover` dipakai di **batas** sistem untuk mencegah satu kegagalan menjatuhkan seluruh proses — misalnya `net/http` menangkap panic per request, atau worker pool yang mengubah panic menjadi error. **Jangan memakai panic/recover sebagai pengganti alur kontrol biasa.**

### 10.8 Stack Traces dan Debugging

Error biasa di Go **tidak membawa stack trace** — hanya pesan. Alat untuk melacak asal masalah:

| Alat | Kegunaan |
|---|---|
| Output **panic** | Otomatis mencetak pesan dan stack trace seluruh goroutine yang relevan |
| `runtime/debug.Stack()` / `debug.PrintStack()` | Mengambil stack trace saat ini, misalnya di dalam `recover` |
| `GOTRACEBACK=all` (atau `system`, `crash`) | Mengatur seberapa rinci trace saat crash |
| **Delve** (`dlv debug`, `dlv test`) | Debugger interaktif: breakpoint, inspeksi variabel, daftar goroutine |
| `log` / `slog` ([Modul 14](14-standard-library.md)) | Mencatat error beserta konteks terstruktur |
| `go vet`, `staticcheck` ([Modul 18](18-toolchain-code-quality.md)) | Menemukan error yang tidak diperiksa atau dibungkus dengan salah |

Pustaka pihak ketiga ada yang menambahkan stack trace ke error, tetapi praktik umum di Go adalah memperkaya error dengan **konteks pada setiap lapisan** (`%w`) dan mencatat stack trace hanya saat terjadi panic.

### 10.9 Praktik Terbaik

1. **Jangan abaikan error** (`_ = f()`), kecuali memang aman dan Anda menuliskan alasannya.
2. **Tangani error satu kali**: *catat atau kembalikan*, jangan keduanya — mencatat di setiap lapisan menghasilkan log duplikat.
3. **Tambah konteks** saat membungkus: apa yang sedang dilakukan, bukan sekadar meneruskan error mentah.
4. Pakai **`errors.Is`/`errors.As`**, bukan perbandingan `==` atau membandingkan teks `err.Error()`.
5. **Periksa error dari `Close()` pada resource yang ditulis** (file, koneksi) — di sinilah kegagalan flush sering muncul.
6. Dalam goroutine, **teruskan error lewat channel atau `errgroup`** ([Modul 13](13-sync-context-patterns.md)); error tidak bisa "dikembalikan" ke pemanggil `go func()`.

## Praktik

### Membungkus dan Memeriksa Error

```go
// errors.go
package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

var ErrStokHabis = errors.New("stok habis")

type ErrValidasi struct {
	Field string
	Pesan string
}

func (e *ErrValidasi) Error() string {
	return fmt.Sprintf("validasi %s: %s", e.Field, e.Pesan)
}

func cekStok(n int) error {
	if n <= 0 {
		return ErrStokHabis
	}
	return nil
}

func buatPesanan(jumlah int, email string) error {
	if email == "" {
		return &ErrValidasi{Field: "email", Pesan: "wajib diisi"}
	}
	if err := cekStok(jumlah); err != nil {
		return fmt.Errorf("buat pesanan: %w", err)
	}
	return nil
}

func main() {
	err := buatPesanan(0, "a@b.com")
	fmt.Println(err)
	fmt.Println("stok habis?", errors.Is(err, ErrStokHabis), err == ErrStokHabis)

	err = buatPesanan(1, "")
	var ev *ErrValidasi
	if errors.As(err, &ev) {
		fmt.Println("field bermasalah:", ev.Field)
	}
	if v, ok := errors.AsType[*ErrValidasi](err); ok {
		fmt.Println("AsType:", v.Pesan)
	}

	_, err = os.Open("tidak-ada.txt")
	fmt.Println(err)
	fmt.Println(errors.Is(err, fs.ErrNotExist))

	gabungan := errors.Join(errors.New("langkah 1 gagal"), errors.New("langkah 2 gagal"))
	fmt.Println(gabungan)
}
```

### `panic`, `recover`, dan Stack Trace

```go
// panic.go
package main

import (
	"fmt"
	"runtime/debug"
	"strings"
)

func ambil(data []int, i int) int {
	return data[i] // panic jika i di luar rentang
}

// Aman mengubah panic menjadi error di batas sistem.
func ambilAman(data []int, i int) (hasil int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("pulih dari panic: %v", r)
			stack := string(debug.Stack())
			fmt.Println("stack berisi 'ambil':", strings.Contains(stack, "ambil"))
		}
	}()
	return ambil(data, i), nil
}

func main() {
	fmt.Println(ambilAman([]int{1, 2, 3}, 1))
	fmt.Println(ambilAman([]int{1, 2, 3}, 5))

	defer fmt.Println("defer tetap dijalankan saat panic")
	var m map[string]int
	m["x"] = 1 // panic: assignment to entry in nil map (tidak dipulihkan)
}
```

Jalankan dan perhatikan pesan panic, stack trace, serta **exit code 2**. Gunakan binary hasil build untuk melihat kode keluar yang sebenarnya — `go run` sendiri hanya mencetak `exit status 2` lalu keluar dengan kode 1.

```bash
go build -o panicapp panic.go
./panicapp
echo "exit code: $?"
```

### Debugging dengan Delve (Opsional)

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug panic.go
# di prompt dlv: break main.ambil, continue, print i, stack, goroutines, exit
```

## Latihan

1. Tulis fungsi `BacaAngka(s string) (int, error)` yang membungkus error `strconv.Atoi` dengan konteks memakai `%w`, lalu periksa jenis aslinya (`*strconv.NumError`) memakai `errors.As`. Apa isi field `Func`, `Num`, dan `Err`?
2. Definisikan sentinel `ErrTidakDitemukan` dan tulis dua lapisan fungsi (repository → service) yang masing-masing membungkus error. Buktikan `errors.Is` tetap mengenali sentinel di lapisan paling atas, sedangkan `==` gagal.
3. Buat custom error `ErrHTTP{Kode int, Pesan string}`, kembalikan dari fungsi, lalu tangani di pemanggil dengan `errors.As` dan `errors.AsType`. Kapan custom error lebih baik daripada sentinel?
4. Tulis fungsi yang menjalankan tiga langkah validasi dan mengumpulkan **semua** kegagalan dengan `errors.Join` alih-alih berhenti di kegagalan pertama. Bagaimana `errors.Is` memeriksa error hasil gabungan?
5. Buat fungsi `Aman(f func()) error` yang menjalankan `f` dan mengubah panic menjadi error beserta stack trace (`debug.Stack()`). Uji dengan fungsi yang membaca slice di luar rentang dan fungsi yang menulis ke map `nil`.
6. Bandingkan membungkus error dengan `%w` dan `%v` pada sebuah paket yang menyembunyikan error driver database. Kapan memutus rantai error justru lebih baik untuk desain API?
7. Diskusikan: seorang rekan menaruh `panic` untuk setiap kegagalan pembacaan file konfigurasi dan memakai `recover` di `main` agar program "tidak pernah crash". Jelaskan mengapa pendekatan ini keliru dari sisi desain, dan rancang alternatif yang lebih idiomatik.

---

[← Modul 9: Generics & Iterators](09-generics.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 11: Modules, Packages & Code Organization →](11-modules-packages.md)
