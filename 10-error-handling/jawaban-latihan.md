# Modul 10: Error Handling — Jawaban Latihan

Seluruh kode ada di folder `kode/` (modul Go `contoh.com/modul10`, satu direktori per program). Jalankan dari dalam `kode/`:

```bash
go run ./l1        # latihan 1 sampai ./l7
go run ./praktik-errors

go build -o panicapp ./praktik-panic && ./panicapp; echo "exit code: $?"   # praktik panic
```

Semua output di bawah dijalankan sungguhan di **Go 1.24.7**. `go vet ./...` dan `gofmt -l` bersih.

**Catatan versi.** `errors.AsType` (generik) baru ada di **Go 1.26**. Di Go 1.24 kode asli materi gagal dengan `undefined: errors.AsType`. Pada `praktik-errors` dan Latihan 3 saya memakai helper `asType[T error](err error) (T, bool)` yang dibangun di atas `errors.As` dan berperilaku sama. Di Go 1.26+ Anda bisa langsung memakai `errors.AsType`.

---

## Latihan 1 — `BacaAngka` dan `*strconv.NumError` (`kode/l1`)

```go
func BacaAngka(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("baca angka dari %q: %w", s, err)
	}
	return n, nil
}
```

```
OK: 123
error: baca angka dari "12x": strconv.Atoi: parsing "12x": invalid syntax
  Func="Atoi" Num="12x" Err=invalid syntax
  syntax? true | range? false
error: baca angka dari "": strconv.Atoi: parsing "": invalid syntax
  Func="Atoi" Num="" Err=invalid syntax
  syntax? true | range? false
error: baca angka dari "99999999999999999999": strconv.Atoi: parsing "99999999999999999999": value out of range
  Func="Atoi" Num="99999999999999999999" Err=value out of range
  syntax? false | range? true
```

**Isi field `*strconv.NumError`:**

| Field | Arti | Contoh (`"12x"`) |
|---|---|---|
| `Func` | Nama fungsi `strconv` yang gagal | `"Atoi"` |
| `Num` | Teks masukan yang gagal diurai | `"12x"` |
| `Err` | Penyebab dasar: `strconv.ErrSyntax` (bukan angka) atau `strconv.ErrRange` (di luar rentang) | `invalid syntax` |

`errors.As(err, &ne)` berhasil menembus bungkusan `%w` dan mengisi `ne`. Karena `NumError` mengimplementasikan `Unwrap()`, `errors.Is(err, strconv.ErrSyntax)` dan `errors.Is(err, strconv.ErrRange)` juga bekerja, sehingga pemanggil bisa membedakan "bukan angka" dari "terlalu besar" tanpa membandingkan teks pesan.

---

## Latihan 2 — Sentinel melewati dua lapisan (`kode/l2`)

```go
var ErrTidakDitemukan = errors.New("tidak ditemukan")

func cariUserDB(id int) error { return fmt.Errorf("query user id=%d: %w", id, ErrTidakDitemukan) } // repository
func AmbilProfil(id int) error { /* service */ return fmt.Errorf("ambil profil: %w", cariUserDB(id)) }
```

```
pesan lengkap  : ambil profil: query user id=42: tidak ditemukan
errors.Is      : true
err == sentinel: false
  lapis: "ambil profil: query user id=42: tidak ditemukan" (*fmt.wrapError)
  lapis: "query user id=42: tidak ditemukan"           (*fmt.wrapError)
  lapis: "tidak ditemukan"                             (*errors.errorString)
dibungkus dengan verb v, errors.Is: false
```

**Mengapa `==` gagal tetapi `errors.Is` berhasil?** Setiap lapisan membungkus error dengan `%w` sehingga yang diterima lapisan atas adalah **objek `*fmt.wrapError` baru**, bukan sentinel itu sendiri. `==` membandingkan identitas objek dan jelas berbeda (`false`). `errors.Is` menelusuri rantai dengan memanggil `Unwrap()` berulang kali (tiga lapis di output) sampai menemukan error yang sama dengan target (`true`).

Baris terakhir menegaskan hal sebaliknya: bila lapisan menengah membungkus dengan `%v` (bukan `%w`), rantai terputus dan `errors.Is` bernilai `false`. Karena itu aturan praktisnya: gunakan `%w` bila pemanggil di atas perlu mengenali error aslinya, dan **selalu bandingkan sentinel dengan `errors.Is`, bukan `==`**.

---

## Latihan 3 — `ErrHTTP` dengan `errors.As` dan `AsType` (`kode/l3`)

```go
type ErrHTTP struct {
	Kode  int
	Pesan string
}
func (e *ErrHTTP) Error() string { return fmt.Sprintf("http %d: %s", e.Kode, e.Pesan) }
```

```
/hilang -> As: kode=404 retry=false
/hilang -> asType: tidak ditemukan
/rusak -> As: kode=500 retry=true
/rusak -> asType: kesalahan server
/ok -> sukses
```

Pemanggil membaca `Kode` untuk memutuskan tindakan (mis. coba ulang bila `>= 500`). Dengan `errors.As` Anda perlu variabel perantara (`var eh *ErrHTTP`), sedangkan `AsType` (di sini lewat helper; `errors.AsType` di Go 1.26+) cukup satu baris `v, ok := asType[*ErrHTTP](err)`.

**Kapan custom error lebih baik daripada sentinel?**

| Situasi | Pilihan |
|---|---|
| Pemanggil hanya perlu tahu **jenis** kegagalan ("tidak ditemukan", "stok habis") | Sentinel: sederhana dan tanpa alokasi |
| Pemanggil perlu **data tambahan** (kode status, nama field, durasi *retry-after*, ID yang gagal) | **Custom error type** |
| Satu "jenis" memiliki banyak variasi dengan atribut berbeda | Custom error type, karena sentinel akan berlipat ganda (`ErrHTTP404`, `ErrHTTP500`, ...) |
| Hanya untuk dibaca manusia | `fmt.Errorf` biasa |

Custom error juga bisa mengimplementasikan `Is(target error) bool` atau `Unwrap()` agar tetap kompatibel dengan `errors.Is`. Keduanya (sentinel dan tipe) menjadi bagian dari API publik paket, jadi tambahkan hanya bila pemanggil benar-benar perlu bereaksi berbeda.

---

## Latihan 4 — Kumpulkan semua kegagalan dengan `errors.Join` (`kode/l4`)

```go
func Validasi(f Form) error {
	var errs []error
	if strings.TrimSpace(f.Nama) == "" { errs = append(errs, ErrNamaKosong) }
	if !strings.Contains(f.Email, "@") { errs = append(errs, fmt.Errorf("%w: %q", ErrEmailTidakValid, f.Email)) }
	if f.Umur < 17 { errs = append(errs, fmt.Errorf("%w (diberikan %d)", ErrUmurDiBawah, f.Umur)) }
	return errors.Join(errs...) // nil bila errs kosong
}
```

```
form valid  -> <nil>
form salah  ->
nama wajib diisi
email tidak valid: "bukan-email"
umur minimal 17 (diberikan 10)

Is(NamaKosong)   : true
Is(EmailTidakValid): true
Is(UmurDiBawah)  : true
Is(io lain)      : false
jumlah error di dalamnya: 3
```

**Bagaimana `errors.Is` memeriksa hasil gabungan?** `errors.Join` mengembalikan tipe yang memiliki method `Unwrap() []error` (terbukti dari type assertion di program: 3 error di dalamnya). Sejak Go 1.20, `errors.Is` dan `errors.As` mengenali bentuk `Unwrap() []error` dan menelusuri **semua cabang** secara berurutan (depth-first). Jadi `errors.Is(err, ErrNamaKosong)` berhasil walau error itu hanyalah salah satu dari tiga, bahkan yang sudah dibungkus lagi dengan `%w`.

Catatan:

- `errors.Join` mengabaikan nilai `nil`, dan mengembalikan `nil` bila semua argumen `nil`, sehingga "form valid" menghasilkan `nil` sungguhan (bukan interface berisi error kosong).
- Pesan gabungan dipisahkan baris baru, berguna untuk dikembalikan ke pengguna sekaligus agar mereka tidak memperbaiki satu per satu.
- Pola ini cocok untuk validasi dan pembersihan beberapa resource, tetapi bukan untuk menyembunyikan kegagalan serius di balik daftar panjang.

---

## Latihan 5 — `Aman(f func()) error` (`kode/l5`)

```go
type ErrPanic struct {
	Nilai any
	Stack string
}

func Aman(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = &ErrPanic{Nilai: r, Stack: string(debug.Stack())}
		}
	}()
	f()
	return nil
}
```

```
1 slice di luar rentang  -> err=panic: runtime error: index out of range [5] with length 3
   runtime.Error? true | stack menyebut fungsi sumber? true | baris stack: 13
2 map nil                -> err=panic: assignment to entry in nil map
   runtime.Error? true | stack menyebut fungsi sumber? true | baris stack: 13
3 tidak panic            -> err=<nil>

cuplikan stack (6 baris pertama):
goroutine 1 [running]:
runtime/debug.Stack()
	/usr/local/go1.24.7/src/runtime/debug/stack.go:26 +0x5e
main.Aman.func1()
	/home/claude/work/10-error-handling/kode/l5/main.go:22 +0x3d
panic({0x4b3260?, 0xc000020330?})
```

**Hal-hal penting:**

- `recover()` hanya efektif bila dipanggil **langsung di dalam fungsi `defer`**. Named return `err` memungkinkan `defer` mengisi hasil.
- `debug.Stack()` dipanggil **di dalam** `recover` (masih di dalam rantai panic), sehingga trace memuat fungsi asal (`bacaDiLuarRentang` atau `tulisKeMapNil`). Bila dipanggil setelah fungsi selesai, informasi itu hilang.
- Kedua panic berasal dari `runtime.Error` (terdeteksi via `errors.As`), yang membedakannya dari `panic("teks")` buatan program.
- Nilai panic bisa apa saja (`any`), jadi simpan dalam tipe error sendiri (`ErrPanic`) agar pemanggil dapat memeriksanya dengan `errors.As`.
- Pola ini tepat di **batas sistem** (worker pool, handler per request), supaya satu kegagalan tidak menjatuhkan seluruh proses. Jangan dipakai sebagai pengganti pengecekan error biasa.
- Panic di **goroutine lain** tidak tertangkap oleh `recover` di goroutine pemanggil; setiap goroutine perlu `recover` sendiri.

---

## Latihan 6 — `%w` vs `%v` pada paket yang menyembunyikan error driver (`kode/l6`)

Skenario: paket `repo` memakai driver database yang menghasilkan `*ErrDriver` dan `ErrKoneksiDriver`.

```
A pesan: simpan user: driver: kode 23505 (driver: koneksi ditolak)
A: pemanggil bisa melihat error driver? true | Is(ErrKoneksiDriver): true
B pesan: simpan user: repo: data duplikat
B: error driver terlihat? false | Is(ErrDuplikat): true
C (dibungkus dengan v): pesan: simpan user: driver: kode 23505 (driver: koneksi ditolak) | Is(ErrKoneksiDriver): false
```

| | `%w` (versi A) | `%v` / terjemahan ke error domain (B, C) |
|---|---|---|
| Pesan | Lengkap | Lengkap (pesan tetap terbaca) |
| Rantai error | **Tetap terbuka**: pemanggil bisa `errors.As(*ErrDriver)` atau `errors.Is(ErrKoneksiDriver)` | **Terputus**: error driver tak terlihat |
| Konsekuensi untuk API | Error driver menjadi **bagian dari API paket Anda** | Detail driver tetap rahasia internal |

**Kapan memutus rantai justru lebih baik?**

1. **Menjaga abstraksi dan kebebasan berubah.** Bila pemanggil mulai bergantung pada `*ErrDriver` atau kode `23505`, mengganti database atau versi driver menjadi perubahan yang memecahkan kompatibilitas (*breaking change*) bagi semua pemakai paket. Kontrak yang stabil sebaiknya berupa error **milik Anda** (`ErrDuplikat`), seperti pada versi B.
2. **Menghindari kebocoran detail internal** (nama tabel, constraint, kode vendor) ke lapisan atas atau klien.
3. **Kopling antar-paket.** Pemanggil tidak perlu mengimpor paket driver hanya untuk mengenali error.

Praktik yang direkomendasikan: **terjemahkan** error dependensi yang bermakna ke sentinel atau tipe domain (`fmt.Errorf("simpan user: %w", ErrDuplikat)`), dan sisanya dibungkus `%v` agar pesan tetap informatif untuk log. Pakai `%w` untuk error yang memang Anda niatkan sebagai bagian dari API (mis. `fs.ErrNotExist`, `context.Canceled`). Sejak Go 1.20 satu `Errorf` boleh memuat beberapa `%w` bila memang ingin mengekspos dua-duanya.

---

## Latihan 7 — `panic` untuk semua kegagalan konfigurasi (`kode/l7`)

Pendekatan rekan (versi keliru):

```go
func muatPanic(path string) Config {
	b, err := os.ReadFile(path)
	if err != nil { panic(err) }
	...
}
// di main: defer func() { recover() ... }()  // "tidak pernah crash"
```

```
[keliru] dipulihkan: open tidak-ada.json: no such file or directory -> lanjut dengan Config kosong?!
```

**Mengapa keliru:**

1. **File konfigurasi hilang atau rusak bukan bug programmer**, melainkan kondisi lingkungan yang wajar dan dapat diperkirakan. Tabel 10.7 menyatakan input salah dan file tidak ada harus berupa `error`, bukan `panic`.
2. **`recover` di `main` membuat program berjalan dalam keadaan tidak valid.** Panic memutus fungsi di tengah; program yang "lanjut" biasanya memakai `Config` kosong atau setengah terisi dan gagal di tempat lain dengan penyebab yang makin sulit ditelusuri. Crash yang cepat dan jelas lebih baik daripada keadaan korup yang diam-diam.
3. **Alur kontrol tersembunyi.** Pemanggil tidak melihat dari signature bahwa fungsi bisa panic, sehingga memaksa semua orang memasang `recover`. Ini kebalikan dari filosofi *errors are values*.
4. **Kehilangan informasi dan kesempatan menangani.** Dengan `error`, pemanggil dapat memutuskan: file tidak ada, pakai nilai default; file rusak, gagal dengan pesan jelas. Dengan `panic`, semuanya diperlakukan sama.
5. Panic juga tidak menjalankan pembersihan secara terkontrol bila `recover` ditaruh di tempat yang salah, dan stack trace panic membingungkan pengguna akhir.

**Alternatif idiomatik** (`Muat` di `l7`):

```go
var ErrConfigTidakAda = errors.New("config tidak ditemukan")

func Muat(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return Config{}, fmt.Errorf("muat %s: %w", path, ErrConfigTidakAda)
	}
	if err != nil { return Config{}, fmt.Errorf("muat %s: %w", path, err) }
	// parse dan validasi, setiap kegagalan dibungkus dengan konteks
}
```

```
[idiomatik] config tidak ada, pakai default
[idiomatik] config dipakai: {8080 localhost}
[idiomatik] config rusak -> parse /dev/null: unexpected end of JSON input
fatal: parse /dev/null: unexpected end of JSON input
exit code: 1
```

Rancangannya: (a) `Muat` mengembalikan `error` dengan konteks; (b) pemanggil membedakan kasus **"tidak ada" → pakai default** dari **"rusak" → gagal**; (c) `main` menangani kegagalan fatal di satu tempat, mencetak pesan jelas ke `stderr`, dan keluar dengan kode non-nol (`exit: 1`), tanpa stack trace yang menakutkan. Panic tetap pantas untuk pelanggaran invariant atau kegagalan inisialisasi yang mustahil dilanjutkan (`regexp.MustCompile`), dan `recover` untuk batas sistem seperti Latihan 5.

---

## Lampiran: program Praktik dari materi

**`praktik-errors`** (dengan helper `asType`, lihat catatan versi):

```
buat pesanan: stok habis
stok habis? true false
field bermasalah: email
AsType: wajib diisi
open tidak-ada.txt: no such file or directory
true
langkah 1 gagal
langkah 2 gagal
```

`errors.Is` bernilai `true` sedangkan `==` bernilai `false` karena error terbungkus; `errors.Is(err, fs.ErrNotExist)` mengenali kegagalan sistem berkas lintas-platform.

**`praktik-panic`**: dijalankan sebagai binary hasil build:

```
2 <nil>
stack berisi 'ambil': true
0 pulih dari panic: runtime error: index out of range [5] with length 3
defer tetap dijalankan saat panic
panic: assignment to entry in nil map

goroutine 1 [running]:
main.main()
	.../praktik-panic/main.go:32 +0x1d7
exit code: 2
```

Panic yang tidak dipulihkan menjalankan `defer` yang sudah terdaftar, mencetak stack trace, lalu keluar dengan **kode 2**. Dengan `go run`, kode keluar proses `go` itu sendiri adalah **1** (hanya mencetak `exit status 2`), sesuai peringatan di materi, sehingga pengujian kode keluar harus memakai binary hasil build.

**Delve** (bagian opsional) tidak saya jalankan: instalasi memerlukan akses jaringan yang tidak tersedia di lingkungan ini, dan debugger interaktif tidak cocok untuk eksekusi tanpa pengawasan.
