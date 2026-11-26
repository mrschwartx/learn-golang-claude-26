# Modul 8: Methods & Interfaces — Jawaban Latihan

Seluruh kode ada di folder `kode/` (modul Go `contoh.com/modul8`, satu direktori per program). Jalankan dari dalam `kode/`:

```bash
go run ./l1        # latihan 1 sampai ./l7
go test -v ./l7/...    # tes dengan fake untuk latihan 7
go run ./praktik-bentuk   # juga: praktik-receiver, praktik-nilgotcha, praktik-hitungbyte
```

Semua output di bawah dijalankan sungguhan di Go 1.24.7. `go vet ./...`, `gofmt -l`, dan `go test` bersih. Folder `l3-error` berisi kode yang sengaja gagal dikompilasi (diberi build tag `ignore`).

---

## Latihan 1 — Tambah `Segitiga` tanpa mengubah `totalLuas` (`kode/l1`)

```go
type Segitiga struct{ A, B, C float64 }

func (s Segitiga) Keliling() float64 { return s.A + s.B + s.C }
func (s Segitiga) Luas() float64 { // rumus Heron
	p := s.Keliling() / 2
	return math.Sqrt(p * (p - s.A) * (p - s.B) * (p - s.C))
}

var _ Bentuk = Segitiga{}
```

```
main.Persegi             luas=4.0000 keliling=8.0000
main.Lingkaran           luas=3.1416 keliling=6.2832
main.Segitiga            luas=6.0000 keliling=12.0000
total luas: 13.1416 (4 + pi + 6)
```

**Mengapa `totalLuas` tidak perlu diubah?** Fungsi itu hanya bergantung pada **interface `Bentuk`** (`Luas()` dan `Keliling()`), bukan pada tipe konkret mana pun. Di Go, pemenuhan interface bersifat **implisit**: `Segitiga` otomatis menjadi `Bentuk` begitu punya kedua method itu, tanpa kata kunci `implements` dan tanpa menyentuh deklarasi `Bentuk` atau `totalLuas`. Efeknya:

- Menambah tipe baru tidak memicu perubahan di kode yang sudah ada (*open for extension, closed for modification*).
- Hasilnya benar: 4 + π + 6 = 13,1416.
- Baris `var _ Bentuk = Segitiga{}` hanya menambah pengecekan saat kompilasi, sehingga lupa satu method langsung ketahuan.

---

## Latihan 2 — `Celsius` dengan `String()` dan `KeFahrenheit()` (`kode/l2`)

```go
type Celsius float64

func (c Celsius) String() string        { return fmt.Sprintf("%.1f°C", float64(c)) }
func (c Celsius) KeFahrenheit() float64 { return float64(c)*9/5 + 32 }
```

```
36.6°C
36.6°C | 36.6°C | 36.60
100°C = 212.0°F
32 -40
tanpa Stringer: 36.6
nilai  : 300
pointer: KelvinPtr(pointer receiver)
suhu memenuhi fmt.Stringer -> 36.6°C
Polos memenuhi fmt.Stringer? false
```

**Bagaimana `fmt` memutuskan memakai `String()`?** Untuk verba `%v` dan `%s` (dan `Println`), paket `fmt` memeriksa saat runtime apakah nilai yang diberikan memenuhi interface `fmt.Stringer` (`String() string`) lewat *type assertion*. Jika ya, `fmt` memanggil `String()`. Jika tidak (seperti `Polos`), `fmt` mencetak representasi bawaan tipe dasarnya (`36.6`).

Beberapa detail yang terlihat dari percobaan:

- Verba numerik seperti `%.2f` tidak memakai `Stringer`; di contoh itu nilainya dikonversi ke `float64` terlebih dahulu.
- Tipe `KelvinPtr` mendefinisikan `String()` dengan **pointer receiver**. Karena itu **nilainya** (`k`) tidak memenuhi `Stringer` dan dicetak sebagai `300`, sedangkan **pointer**-nya (`&k`) memanggil `String()`. Ini adalah method set (Latihan 3) yang tampil lewat `fmt`.
- Hati-hati agar `String()` tidak memanggil `fmt.Sprint(c)` dengan receiver yang sama, karena akan rekursi tak berujung. Itulah alasan konversi ke `float64(c)` di dalam `String()`.

---

## Latihan 3 — `Reset()` pointer receiver dan interface `Resetter` (`kode/l3`, `kode/l3-error`)

```go
type Resetter interface{ Reset() }
type Counter struct{ n int }
func (c *Counter) Reset() { c.n = 0 }

var r Resetter = Counter{} // nilai biasa -> error
```

Pesan error compiler (dari `l3-error`):

```
./main.go:11:19: cannot use Counter{} (value of struct type Counter) as Resetter value in variable declaration: Counter does not implement Resetter (method Reset has pointer receiver)
```

Versi yang valid (`l3`):

```
setelah pakai(&c): 0
c.Reset() langsung : 0
map berisi pointer : 0
```

Pada `praktik-receiver` pesan yang sama muncul untuk `var t2 Penambah = c`: `Penghitung does not implement Penambah (method Tambah has pointer receiver)`.

**Konsep method set:**

| Tipe | Method yang ada di method set |
|---|---|
| `T` | Hanya method ber-**value receiver** |
| `*T` | Method ber-**value receiver** dan **pointer receiver** |

Karena `Reset` memakai pointer receiver, hanya `*Counter` yang memenuhi `Resetter`, sedangkan `Counter` (nilai) tidak. Alasannya masuk akal: jika nilai diperbolehkan, method akan beroperasi pada **salinan** yang tersimpan di dalam interface, dan perubahannya tidak akan terlihat siapa pun. Untuk memenuhi interface, kirim pointer: `pakai(&c)`.

Perhatikan perbedaan dengan pemanggilan langsung: `c.Reset()` pada variabel *addressable* otomatis menjadi `(&c).Reset()`, sehingga tampak "bekerja pada nilai". Kemudahan ini tidak berlaku pada nilai yang tidak addressable (elemen map, hasil kembalian fungsi) dan tidak berlaku untuk pemenuhan interface. Itulah sebabnya pada contoh map saya menyimpan pointer.

---

## Latihan 4 — Jebakan `nil` pada interface `error` (`kode/l4`)

```go
type ErrValidasi struct{ Field string }
func (e *ErrValidasi) Error() string { return "validasi gagal: " + e.Field }

// SALAH
func validasiSalah(ok bool) error {
	var e *ErrValidasi
	if !ok { e = &ErrValidasi{"nama"} }
	return e // interface berisi (*ErrValidasi, nil)
}

// BENAR
func validasiBenar(ok bool) error {
	if !ok { return &ErrValidasi{"nama"} }
	return nil
}
```

```
validasiSalah(true) != nil : true
  tipe dinamis=*main.ErrValidasi, nilai=true
  memanggil err.Error() -> runtime error: invalid memory address or nil pointer dereference
validasiBenar(true) != nil : false
gagal -> validasi gagal: nama | errors.As: true nama
```

**Penjelasan.** Nilai interface terdiri atas dua komponen: **(tipe dinamis, nilai)**. Interface baru `nil` bila **keduanya** kosong. Pada `validasiSalah(true)` yang sukses, nilainya `nil` tetapi tipenya sudah terisi `*ErrValidasi`, sehingga `err != nil` bernilai **`true`** (dan baris kedua output menunjukkan nilai di dalamnya memang pointer `nil`). Kode pemanggil yang tertib seperti `if err != nil { ... }` akan mengira terjadi kegagalan, bahkan `err.Error()` bisa *panic* karena menyentuh `e.Field` pada pointer `nil`.

**Perbaikan:** fungsi yang mengembalikan `error` harus mengembalikan literal **`nil`** secara eksplisit pada jalur sukses, dan tidak menyimpan hasil dalam variabel bertipe pointer konkret yang lalu dikembalikan begitu saja. Alternatifnya, deklarasikan variabel perantara sebagai `var err error`. Linter seperti `go vet` dan `staticcheck` dapat membantu mendeteksi pola ini.

---

## Latihan 5 — `Ringkas(x any) string` dengan type switch (`kode/l5`)

```go
func Ringkas(x any) string {
	switch v := x.(type) {
	case nil:            return "nil"
	case int:            return fmt.Sprintf("int %d (+1 = %d)", v, v+1)
	case float64:        return fmt.Sprintf("float64 %.2f", v)
	case string:         return fmt.Sprintf("string %q (upper %s, %d byte)", v, strings.ToUpper(v), len(v))
	case []int:          /* len dan jumlah */
	case error:          return "error: " + v.Error()
	case int8, int16, int32, int64: // v tetap bertipe any
		return fmt.Sprintf("integer lain %T = %v", v, v)
	default:             return fmt.Sprintf("tipe lain: %T", v)
	}
}
```

```
int 42 (+1 = 43)
float64 3.14
string "halo" (upper HALO, 4 byte)
[]int len=3 jumlah=6
error: gagal
nil
nil
integer lain int64 = 7
tipe lain: bool
tipe lain: map[string]int
tipe lain: *int
```

Hal yang perlu diperhatikan:

- Pada `case` bertipe tunggal, `v` otomatis bertipe sesuai `case`-nya (misalnya `v+1` valid di `case int`). Pada `case` multi-tipe dan `default`, `v` tetap bertipe `any`.
- `kosong` adalah `error` bernilai `nil`. Saat dikirim ke parameter `any`, interface nil itu menjadi `any` nil (tanpa tipe), sehingga cocok dengan `case nil`. Ini berbeda dengan kasus Latihan 4: pointer `nil` bertipe konkret (`ptrNil`) **tidak** cocok dengan `case nil` dan jatuh ke `default` sebagai `*int`.
- **Urutan `case` penting** bila tipenya tumpang tindih: `case error` bisa menangkap tipe konkret apa pun yang punya `Error()`, jadi letakkan `case` yang lebih spesifik lebih dulu.

---

## Latihan 6 — `io.Reader` kustom `HurufBesar` (`kode/l6`)

```go
type HurufBesar struct{ R io.Reader }

func (h HurufBesar) Read(p []byte) (int, error) {
	n, err := h.R.Read(p)
	copy(p[:n], bytes.ToUpper(p[:n]))
	return n, err // io.EOF diteruskan apa adanya
}

io.Copy(os.Stdout, HurufBesar{strings.NewReader("halo, dunia interface go!\n")})
```

```
HALO, DUNIA INTERFACE GO!
disalin: 26 byte, err: <nil>
SELAMAT BE
```

Baris terakhir menunjukkan reader ditumpuk dengan `io.LimitReader(HurufBesar{...}, 10)` sehingga hanya 10 byte pertama yang lewat.

**Hubungan dengan prinsip interface kecil.** `io.Reader` hanya punya **satu method** (`Read([]byte) (int, error)`). Karena sekecil itu:

- Mudah dipenuhi: `HurufBesar` hanya butuh sekitar lima baris, dan otomatis menjadi `io.Reader`.
- Mudah **dikomposisi**: `HurufBesar` membungkus `io.Reader` lain, sehingga sumbernya bisa berupa string, file, koneksi jaringan, atau body HTTP tanpa mengubah satu baris pun. Fungsi `io.Copy`, `io.LimitReader`, `bufio`, `gzip`, dan lain-lain tersambung seperti pipa.
- Mudah diuji: cukup `strings.NewReader` sebagai sumber.

Catatan implementasi: `Read` harus mengubah byte hanya pada `p[:n]` (yang benar-benar terbaca), dan harus meneruskan `err` (termasuk `io.EOF`) agar pemanggil tahu kapan data habis. Konversi `ToUpper` di sini per-chunk, sehingga cocok untuk ASCII; untuk teks multibyte yang terpotong di batas chunk diperlukan penanganan rune yang lebih hati-hati.

---

## Latihan 7 — Interface `Penyimpan` di sisi pemakai (`kode/l7`)

Struktur:

```
l7/
  db/                  # penyedia: Postgres punya Simpan, Hapus, Cari, Laporan
  layananpesanan/      # pemakai: mendeklarasikan Penyimpan { Simpan(Pesanan) error }
    layanan.go
    layanan_test.go    # memakai fake, tanpa database
  main.go              # merangkai keduanya lewat adapter kecil
```

```go
// layananpesanan/layanan.go (sisi pemakai)
type Penyimpan interface {
	Simpan(p Pesanan) error
}
func Baru(s Penyimpan) *Layanan { ... } // menerima interface, mengembalikan struct
```

```
INSERT ke postgres://localhost/toko: pesanan 1 total 7500
Buat(1, 7500): <nil>
Buat(2, 0)   : total pesanan harus > 0

--- PASS: TestBuatBerhasil
--- PASS: TestBuatValidasi
--- PASS: TestBuatGagalSimpan
ok  	contoh.com/modul8/l7/layananpesanan
```

**Mengapa interface di sisi pemakai (consumer) lebih baik?**

1. **Interface minimal sesuai kebutuhan.** `layananpesanan` hanya butuh `Simpan`. Bila interface dibuat di paket database, ia cenderung memuat semua method (`Hapus`, `Cari`, `Laporan`, ...), sehingga pemakai dipaksa bergantung pada hal yang tidak ia pakai. Interface kecil juga lebih mudah dipenuhi.
2. **Arah dependensi terbalik.** Paket database tidak perlu tahu siapa pemakainya, dan `layananpesanan` tidak perlu mengimpor paket database hanya demi tipe interface. Pemenuhan interface bersifat implisit, jadi `*db.Postgres` (lewat adapter di `main`) langsung cocok tanpa deklarasi apa pun. Hasilnya, kopling antar-paket longgar dan tidak ada siklus impor.
3. **Bebas berubah.** Penyedia boleh menambah method baru tanpa merusak pemakai mana pun. Tiap pemakai bisa mendefinisikan interface yang berbeda dari satu penyedia yang sama.

**Bagaimana ini mempermudah pengujian dengan fake?** Karena `Layanan` hanya meminta `Penyimpan`, tes cukup membuat `fakeStore` kecil (tiga baris method `Simpan` yang menampung data di slice). Tidak perlu database, jaringan, atau framework *mocking*. Tiga tes di atas memverifikasi:

- jalur sukses (data tersimpan),
- validasi (`ErrTotalNol`, dan tidak ada yang tersimpan),
- error dari penyimpan terbungkus dengan benar (`errors.Is` tetap menemukan error aslinya).

Tes berjalan dalam milidetik dan deterministik. Fake juga bisa diprogram untuk mensimulasikan kegagalan yang sulit dibuat pada database asli (koneksi putus, timeout). Konsep ini diperdalam di Modul 15. Satu catatan jujur: pada contoh ini tipe `Pesanan` di `db` dan di `layananpesanan` sengaja dibuat berbeda sehingga perlu adapter kecil di `main`; pada proyek nyata biasanya tipe domain dibagi bersama atau adapter diletakkan di lapisan perakit (*composition root*).

---

## Lampiran: program Praktik dari materi

Keempat program Praktik (`praktik-bentuk`, `praktik-receiver`, `praktik-nilgotcha`, `praktik-hitungbyte`) sudah dijalankan dan sesuai teori:

- `bentuk`: `total luas: 7.14`, type switch memberi deskripsi benar, comma-ok pada `any` aman.
- `receiver`: `2 Penghitung(2)` lalu `Penghitung(3)`; `fmt` memakai `String()`.
- `nilgotcha`: `salah() == nil: false`, `benar() == nil: true`.
- `hitungbyte`: dua kali `Fprint*` ke `PenghitungByte` menghasilkan total 26 byte. Method `Write` yang ditimpa dipanggil, sedangkan sisanya diteruskan ke `os.Stdout` lewat embedding.
