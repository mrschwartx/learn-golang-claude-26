# Modul 3: Variabel, Konstanta & Tipe Data — Jawaban Latihan

Seluruh kode ada di folder `kode/` (satu modul Go `contoh.com/modul3`, satu direktori per program). Jalankan dari dalam `kode/`:

```bash
go run ./l1        # latihan 1, dst. sampai ./l8
go run ./praktik-variabel   # juga: praktik-iota, praktik-teks, praktik-kompleks
go vet ./l6        # peringatan untuk latihan 6
```

Semua output di bawah dijalankan sungguhan di Go 1.24.7. `gofmt -l` tidak melaporkan file yang menyimpang.

---

## Latihan 1 — Lima cara deklarasi variabel (`kode/l1`)

```go
var versi = "1.0"           // level paket: wajib var
var umur int = 30           // 1. var dengan tipe
var nama = "Andi"           // 2. var tanpa tipe (inferensi)
var total float64           // 3. var tanpa nilai (zero value)
var ( x = 1; y = 2.5 )      // 4. var berkelompok
kota := "Jakarta"           // 5. short declaration
```

Output:

```
umur  int 30
nama  string Andi
total float64 0
x     int 1
y     float64 2.5
kota  string Jakarta
paket string 1.0 | localhost:8080
```

**Mengapa `:=` tidak bisa di level paket?** Di luar fungsi, setiap pernyataan harus diawali kata kunci (`var`, `const`, `func`, `type`, `import`). `:=` adalah *pernyataan* (statement) yang hanya boleh ada di dalam badan fungsi. Alasan praktisnya: urutan inisialisasi variabel paket ditentukan dari dependensi antar-deklarasi, bukan dari urutan baris, sehingga parser dan compiler memerlukan kata kunci eksplisit agar deklarasi tingkat paket bisa dikenali tanpa menjalankan kode.

---

## Latihan 2 — Shadowing `err` (`kode/l2`)

```go
func salah() error {
	var err error
	if true {
		_, err := operasi() // := membuat err BARU yang menutupi err luar
		_ = err
	}
	return err // tetap nil!
}
```

Output:

```
salah : <nil>
benar : gagal
benar2: gagal
```

**Yang terjadi.** `:=` di dalam blok `if` membuat variabel `err` baru khusus untuk blok itu. Error `"gagal"` disimpan ke variabel baru tersebut, yang lenyap ketika blok selesai, sedangkan `err` luar tidak pernah disentuh sehingga `return err` mengembalikan `nil`. Error hilang tanpa jejak, padahal kodenya valid dan lolos kompilasi.

**Cara memperbaiki:**

1. Pakai `=` (bukan `:=`) bila variabel luar memang yang ingin diisi (`_, err = operasi()`). Kalau variabel kiri lain ikut baru, deklarasikan dulu dengan `var`.
2. Tangani error langsung di scope dalam: `if _, err := operasi(); err != nil { return err }`.
3. Aktifkan analyzer `shadow` (`go vet` dengan analyzer tambahan atau linter) agar kasus seperti ini terdeteksi otomatis.

---

## Latihan 3 — Enum `Level` dan bit flag izin (`kode/l3`)

```go
type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

func (l Level) String() string { /* switch -> "DEBUG", "INFO", "WARN", "ERROR" */ }

type Izin uint8

const (
	Baca  Izin = 1 << iota // 1
	Tulis                  // 2
	Eksek                  // 4
)

func (i Izin) Punya(f Izin) bool { return i&f != 0 }
```

Output:

```
0 -> DEBUG
1 -> INFO
2 -> WARN
3 -> ERROR
9 -> Level(9)
izin=011 baca=true tulis=true eksek=false
izin=101 baca=true tulis=false eksek=true
```

Operator bitwise yang dipakai: `|` menggabungkan izin, `&` mengecek izin, `|=` menambah izin, `&^=` (AND NOT) mencabut izin. Karena `Level` punya method `String()`, `fmt` otomatis memakainya saat mencetak dengan `%s`.

---

## Latihan 4 — `len` vs rune, membalik string (`kode/l4`)

```
len (byte): 17
rune      : 14
dibalik   : 🌅 igap tamaleS
balik byte: "\x85\x8c\x9f\xf0 igap tamaleS" (valid UTF-8: false)
```

**Mengapa berbeda?** `len` menghitung **byte**, bukan karakter. String Go berisi UTF-8: huruf ASCII memakai 1 byte, sedangkan emoji 🌅 memakai **4 byte**. Rinciannya: 13 karakter ASCII (termasuk spasi) = 13 byte, ditambah 4 byte emoji = **17 byte**, tetapi hanya **14 karakter (rune)**. `utf8.RuneCountInString` menghitung rune.

**Membalik dengan benar:** ubah ke `[]rune`, tukar elemen dari ujung ke ujung, lalu ubah kembali ke `string`:

```go
func balik(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
```

Membalik per-byte (`[]byte`) merusak karakter multi-byte: urutan 4 byte emoji ikut terbalik dan hasilnya bukan UTF-8 valid, seperti terlihat pada baris terakhir output. Catatan: membalik per-rune pun masih belum sempurna untuk karakter gabungan (misalnya huruf + aksen terpisah), karena itu perlu paket tambahan yang memahami *grapheme cluster*.

---

## Latihan 5 — `strconv.Atoi("12x")` (`kode/l5`)

```go
n, err := strconv.Atoi(input)
if err != nil {
	var numErr *strconv.NumError
	if errors.As(err, &numErr) && numErr.Err == strconv.ErrSyntax {
		fmt.Printf("%q bukan angka yang valid. Masukkan angka bulat, misalnya 12.\n", input)
	} else if errors.Is(err, strconv.ErrRange) {
		fmt.Printf("%q terlalu besar atau terlalu kecil.\n", input)
	} else {
		fmt.Println("terjadi kesalahan:", err)
	}
	return
}
```

Output:

```
"12x" bukan angka yang valid. Masukkan angka bulat, misalnya 12.
  (error asli: strconv.Atoi: parsing "12x": invalid syntax )
Berhasil: 12 (n+1 = 13)
"" bukan angka yang valid. Masukkan angka bulat, misalnya 12.
  (error asli: strconv.Atoi: parsing "": invalid syntax )
"99999999999999999999" terlalu besar atau terlalu kecil.
  (error asli: strconv.Atoi: parsing "99999999999999999999": value out of range )
```

Pesan asli dari `strconv` ("invalid syntax") bagus untuk log, tetapi kurang ramah bagi pengguna akhir. Karena itu kodenya membedakan penyebab error (`ErrSyntax` vs `ErrRange`) lalu menampilkan pesan yang jelas dan memberi contoh masukan yang benar. Teknik `errors.As`/`errors.Is` dibahas lebih dalam di Modul 10.

---

## Latihan 6 — `string(65)` (`kode/l6`)

```go
n := 65
fmt.Println(string(n))       // A
fmt.Println(string(rune(n))) // A
fmt.Println(strconv.Itoa(n)) // 65
```

Output program: `A`, `A`, `65`.

Output `go vet ./l6`:

```
l6/main.go:11:14: conversion from int to string yields a string of one rune, not a string of digits
```

**Mengapa hasilnya `A`?** Konversi `string(x)` pada nilai bilangan bulat tidak memformat angka menjadi teks, melainkan menafsirkannya sebagai **code point Unicode**. Kode 65 adalah huruf `A`, sehingga `string(65)` menghasilkan `"A"`. Itu sebabnya `go vet` memperingatkannya: hampir selalu programmer bermaksud mendapatkan `"65"`.

**Cara yang benar:**

- `strconv.Itoa(65)` → `"65"` (cara paling umum dan cepat untuk `int`).
- `fmt.Sprint(65)` → `"65"` (untuk tipe apa pun).
- Jika memang butuh karakter dari kode, tulis eksplisit: `string(rune(65))`. Maksudnya jadi jelas dan `vet` tidak protes.

---

## Latihan 7 — Harga dengan `float64` vs integer (`kode/l7`)

Hasil percobaan:

```
float64 10x0.1  : 0.9999999999999999 | == 1.0? false
float64 0.1+0.2 : 0.30000000000000004 | == 0.3? false
float64 1000x19.99: 19990.0000000001 (seharusnya 19990)
int64 sen 1000x1999: Rp19990,00 | tepat 19990,00? true
diskon 15% dari Rp19,99 = Rp3,00
PPN 11% dari Rp19,99 = Rp2,20
```

**Mengapa `float64` berisiko untuk uang:**

- Pecahan desimal seperti 0,1 tidak dapat direpresentasikan persis dalam biner, sehingga setiap operasi membawa galat kecil yang menumpuk pada penjumlahan berulang (lihat 1000 × 19,99).
- Perbandingan `==` jadi tidak andal. Total pesanan yang "seharusnya sama" bisa dianggap berbeda, yang berujung pada rekonsiliasi yang gagal atau pembayaran yang ditolak.
- Galat tidak muncul sebagai error atau panic, jadi bug ini lolos diam-diam dan baru terlihat pada audit.

**Alternatif aman:** simpan uang sebagai **integer satuan terkecil** (sen atau rupiah) dengan tipe tersendiri:

```go
type Uang int64 // jumlah sen

func (u Uang) KaliPersen(p int64) Uang { return Uang((int64(u)*p + 50) / 100) } // pembulatan eksplisit
```

Penjumlahan dan perbandingan integer selalu tepat. Pembulatan (misalnya untuk diskon atau pajak) dilakukan satu kali, secara eksplisit, dengan aturan yang dipilih sendiri (di sini half-up) sehingga hasilnya konsisten dan dapat diaudit. Format ke teks dilakukan hanya saat ditampilkan (`String()`). Untuk kebutuhan lanjutan (banyak mata uang, pecahan sangat kecil), pakai library desimal presisi tetap; tipe `Uang` di atas hanyalah ilustrasi. Hindari juga `int` bila aplikasi dapat berjalan di platform 32-bit; `int64` lebih aman.

---

## Latihan 8 — Persamaan kuadrat dengan `complex128` (`kode/l8`)

Persamaan: `x² + 2x + 5 = 0`, sehingga `a=1, b=2, c=5`.

```
diskriminan: (-16+0i)
x1 = (-1+2i)  (real=-1, imag=2)
   substitusi: (0+0i)  |hasil|=0  nol? true
x2 = (-1-2i)  (real=-1, imag=-2)
   substitusi: (0+0i)  |hasil|=0  nol? true
```

**Perhitungan:** D = b² − 4ac = 4 − 20 = −16, `√D = cmplx.Sqrt(-16) = 4i`, sehingga `x = (−2 ± 4i) / 2 = −1 ± 2i`.

**Verifikasi substitusi** untuk `x = −1 + 2i`: x² = 1 − 4i + 4i² = −3 − 4i, 2x = −2 + 4i, maka x² + 2x + 5 = (−3 − 2 + 5) + (−4 + 4)i = 0. Hasil program juga `0+0i`. Untuk akar bilangan kompleks yang lebih rumit, galat pembulatan mungkin membuat hasil bukan persis nol, sehingga pengecekan memakai batas toleransi (`cmplx.Abs(hasil) < 1e-12`) lebih aman daripada `== 0`.

**Mengapa `a < b` tidak bisa dikompilasi?** Bilangan kompleks **tidak punya urutan alami**. Tidak ada cara baku menyatakan apakah `1+2i` lebih kecil dari `2+1i`, karena bidang kompleks dua dimensi, bukan garis bilangan. Karena itu spesifikasi Go hanya mendefinisikan `==` dan `!=` untuk tipe kompleks, dan operator `<`, `<=`, `>`, `>=` ditolak compiler dengan pesan `operator < not defined`. Bila perlu membandingkan, bandingkan sesuatu yang berurutan, misalnya besaran `cmplx.Abs(a) < cmplx.Abs(b)`.

---

## Lampiran: program Praktik dari materi

Keempat program di bagian Praktik materi (`praktik-variabel`, `praktik-iota`, `praktik-teks`, `praktik-kompleks`) juga sudah saya jalankan dan hasilnya cocok dengan penjelasan teori, misalnya `overflow uint8: 0`, `0.30000000000000004 false` untuk variabel dibanding `true` untuk konstanta `0.1+0.2 == 0.3`, serta `byte: 14 rune: 9` untuk `"héllo, 世界"`.
