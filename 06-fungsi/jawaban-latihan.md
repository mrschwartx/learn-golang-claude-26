# Modul 6: Fungsi — Jawaban Latihan

Seluruh kode ada di folder `kode/` (modul Go `contoh.com/modul6`, satu direktori per program). Jalankan dari dalam `kode/`:

```bash
go run ./l1        # latihan 1 sampai ./l7
go run ./praktik-fungsi   # juga: praktik-closure, praktik-defer
```

Semua output di bawah dijalankan sungguhan di Go 1.24.7. `go vet ./...` dan `gofmt -l` bersih. `l6` memakai `/proc` dan `syscall.Setrlimit`, jadi hanya berjalan di Linux.

---

## Latihan 1 — `MinMax(angka ...int) (min, max int, err error)` (`kode/l1`)

```go
var ErrKosong = errors.New("minmax: tidak ada argumen")

func MinMax(angka ...int) (min, max int, err error) {
	if len(angka) == 0 {
		return 0, 0, ErrKosong
	}
	min, max = angka[0], angka[0]
	for _, n := range angka[1:] {
		if n < min { min = n }
		if n > max { max = n }
	}
	return min, max, nil
}
```

```
argumen biasa: -2 9
slice s...   : 3 25 <nil>
satu argumen : 42 42 <nil>
tanpa argumen: minmax: tidak ada argumen | errors.Is: true
```

- Slice diteruskan ke parameter variadic dengan akhiran `...` (`MinMax(s...)`). Tanpa titik tiga itu, kode tidak dikompilasi karena `[]int` bukan `int`.
- Tanpa argumen, `angka` adalah `nil` slice (`len == 0`), sehingga `angka[0]` akan panic; karena itu pengecekan panjang dilakukan lebih dulu dan error dikembalikan.
- Error dideklarasikan sebagai variabel paket (`ErrKosong`) agar pemanggil bisa mengujinya dengan `errors.Is` (dibahas di Modul 10).
- Catatan: nama `min` dan `max` pada named return menimpa fungsi bawaan `min`/`max` (Go 1.21+) di dalam fungsi ini. Itu valid, tetapi sebaiknya dihindari di kode nyata; di sini dipertahankan agar sesuai dengan signature yang diminta soal.

---

## Latihan 2 — `Penjumlah(awal int) func(int) int` (`kode/l2`)

```go
func Penjumlah(awal int) func(int) int {
	total := awal
	return func(n int) int {
		total += n
		return total
	}
}
```

```
a(5)  = 5
a(10) = 15
b(1)  = 101
a(1)  = 16 (a tidak terpengaruh b)
b(1)  = 102 (b tidak terpengaruh a)
```

Setiap pemanggilan `Penjumlah` membuat variabel `total` baru, dan closure yang dikembalikan menangkapnya **secara referensi**. Variabel itu tetap hidup (dipindah ke heap oleh compiler) selama closure masih dipakai. Karena `a` dan `b` berasal dari dua pemanggilan berbeda, masing-masing memegang `total` miliknya sendiri dan state-nya terpisah, seperti terlihat dari output.

---

## Latihan 3 — `Map`, `Filter`, `Reduce` (`kode/l3`)

```go
func Map(data []int, f func(int) int) []int
func Filter(data []int, f func(int) bool) []int
func Reduce(data []int, awal int, f func(akum, v int) int) int
```

Gabungan: kuadratkan semua angka, ambil yang genap, jumlahkan.

```
kuadrat: [1 4 9 16 25 36]
genap  : [4 16 36]
jumlah : 56
digabung: 56
```

Hasil 56 = 4 + 16 + 36. Karena fungsi adalah nilai, perilaku (apa yang dilakukan terhadap tiap elemen) dikirim sebagai argumen dan kerangka perulangannya ditulis sekali. Contoh nyata di standard library: `slices.SortFunc`, `sort.Slice`, dan `http.HandlerFunc`. Sejak Go 1.18, versi generik dapat ditulis untuk tipe elemen apa pun (Modul 9).

---

## Latihan 4 — Output `defer` (`kode/l4`)

Program:

```go
func demo() {
	pesan := "awal"
	defer fmt.Println("defer A: pesan saat didaftarkan =", pesan)
	defer func() {
		fmt.Println("defer B (closure): pesan saat dijalankan =", pesan)
	}()
	for i := range 3 {
		defer fmt.Println("defer loop i =", i)
	}
	pesan = "berubah"
	fmt.Println("badan fungsi selesai, pesan =", pesan)
}
```

**Prediksi sebelum menjalankan:**

1. Badan fungsi berjalan lebih dulu dan mencetak `pesan = berubah`.
2. Saat fungsi selesai, `defer` dijalankan **LIFO** (yang terakhir didaftarkan, yang pertama berjalan): tiga `defer` loop dengan `i` = 2, 1, 0 lebih dulu, lalu `defer B`, lalu `defer A`.
3. `defer A` memakai `fmt.Println(..., pesan)` langsung, sehingga **argumennya dievaluasi saat `defer` didaftarkan** → mencetak `awal`.
4. `defer B` adalah closure yang membaca `pesan` **saat dijalankan** → mencetak `berubah`.
5. `defer loop` mengevaluasi `i` saat didaftarkan; di Go 1.22+ tiap iterasi punya `i` sendiri, jadi nilainya 0, 1, 2 dicetak dalam urutan terbalik.

**Verifikasi (hasil nyata):**

```
badan fungsi selesai, pesan = berubah
defer loop i = 2
defer loop i = 1
defer loop i = 0
defer B (closure): pesan saat dijalankan = berubah
defer A: pesan saat didaftarkan = awal
```

Tebakan terbukti benar. Aturan ringkasnya: *argumen* `defer` dievaluasi saat deklarasi, tetapi *isi closure* dievaluasi saat eksekusi.

---

## Latihan 5 — Go adalah pass-by-value (`kode/l5`)

```
struct         : {1}
*struct        : {100}
slice elemen   : [100 2 3]
slice append   : [100 2 3] len = 3
slice diganti  : [100 2 3]
map            : map[baru:1 k:100]
alamat x di main: 0xc000010198
nilai pointer   : 0xc000010198 (sama, karena yang disalin adalah alamatnya)
alamat variabel p: 0xc000056048 (berbeda, p adalah salinan)
```

| Argumen | Yang disalin | Perubahan di dalam fungsi terlihat oleh pemanggil? |
|---|---|---|
| `struct` | Seluruh isi | **Tidak** (`{1}` tetap) |
| `*struct` | Alamat | **Ya** (`{100}`): lewat pointer yang disalin, objek yang ditunjuk sama |
| `[]int`: ubah elemen | Header (pointer, len, cap) | **Ya** (`100` di indeks 0) |
| `[]int`: `append` | Header | **Tidak**: `len` pemanggil tetap 3 |
| `[]int`: ganti dengan slice lain | Header | **Tidak**: hanya header salinan yang menunjuk ke array baru |
| `map[string]int` | Referensi ke struktur internal | **Ya**: baik ubah nilai (`k:100`) maupun tambah key (`baru:1`) |

**Bukti pass-by-value.** Pada contoh pointer terakhir, nilai pointer yang diterima fungsi sama dengan alamat `x`, tetapi alamat *variabel parameter* `p` berbeda. Jadi `p` adalah salinan, dan yang disalin kebetulan sebuah alamat. Inilah sebabnya `*p = ...` mengubah data pemanggil, sedangkan `p = nil` tidak.

**Mengapa `append` di dalam fungsi tidak mengubah panjang slice pemanggil?** Header slice (pointer, `len`, `cap`) disalin saat dipanggil. `append` mengembalikan header baru dengan `len` bertambah, dan header baru itu hanya menggantikan salinan di dalam fungsi. Header milik pemanggil tidak tersentuh, sehingga `len` tetap 3. Karena itu pola yang benar adalah mengembalikan hasilnya: `s = tambah(s)`. Bila kapasitas masih cukup, elemen baru bahkan tertulis ke array bersama, tetapi pemanggil tetap tidak melihatnya karena `len`-nya tidak berubah (dibahas di Modul 4).

---

## Latihan 6 — `defer` di dalam loop (`kode/l6`)

Versi 1 (bermasalah) dan versi 2 (diperbaiki):

```go
// VERSI 1: defer langsung di loop
for _, p := range paths {
	f, err := os.Open(p)
	if err != nil { return err }
	defer f.Close() // baru berjalan saat BacaSemuaSalah return
	...
}

// VERSI 2: isi loop dipindah ke fungsi bantu
func bacaSatu(p string) error {
	f, err := os.Open(p)
	if err != nil { return err }
	defer f.Close() // berjalan di akhir bacaSatu
	_, err = f.Stat()
	return err
}
```

Hasil dengan 1000 berkas dan batas file descriptor diturunkan ke 256 (agar masalahnya tampak jelas):

```
batas file descriptor: 256 | berkas: 1000

VERSI 1 (defer di loop): puncak FD terbuka = 256
  error: membuka f0251.txt: open /tmp/modul6425199429/f0251.txt: too many open files

VERSI 2 (fungsi bantu) : puncak FD terbuka = 6
  error: <nil>
```

**Masalah pada versi pertama.** `defer` baru dijalankan ketika **fungsi pembungkus** selesai, bukan di akhir iterasi. Seluruh file yang sudah dibuka tetap terbuka sampai loop berakhir. Dengan ribuan file akibatnya:

- File descriptor menumpuk sampai menembus batas proses (`ulimit -n`, sering 1024 di Linux), lalu `open` gagal dengan **`too many open files`** (terlihat di output pada berkas ke-252).
- Memori dan sumber daya OS ikut tertahan, dan pada Windows file yang masih terbuka bisa menghalangi penghapusan atau perubahan file lain.
- Daftar `defer` itu sendiri juga membesar sesuai jumlah iterasi.
- Bug ini tidak tampak pada pengujian dengan sedikit file dan baru muncul di produksi dengan data besar.

**Perbaikan.** Pindahkan isi loop ke fungsi bantu: `defer` di dalamnya berjalan setiap kali fungsi bantu selesai, sehingga hanya satu file yang terbuka pada satu waktu (puncak tetap rendah di versi 2). Alternatif lain: panggil `f.Close()` secara eksplisit di akhir iterasi (lebih rawan bocor bila ada `return` atau `panic` di tengah) atau gunakan closure yang langsung dipanggil di dalam loop.

---

## Latihan 7 — Struct konfigurasi vs functional options (`kode/l7`)

Contoh singkat kedua pola:

```go
// A: struct konfigurasi
s := NewServerA(ConfigA{Addr: ":9000"})

// B: functional options
s := NewServerB(WithAddr(":9000"), WithMaxConn(500))
```

Dari percobaan:

```
A timeout 0 -> 30s (ingin 'tanpa timeout' tapi jadi 30s)
B timeout 0 -> 0s (eksplisit: 0 benar-benar dipilih)
```

Pada pola A, nilai nol (`0`, `""`) tidak bisa dibedakan dari "tidak diisi", sehingga pengguna tidak dapat meminta timeout nol. Pola B bisa membedakannya karena default ditetapkan di konstruktor dan opsi hanya mengubah yang dipilih.

| Aspek | Struct konfigurasi | Functional options |
|---|---|---|
| Kesederhanaan | **Sangat sederhana**, mudah dibaca dan dibuat (mis. dari JSON/YAML) | Perlu fungsi `With...` untuk tiap opsi |
| Zero value | Ambigu: "belum diisi" vs "diisi nol" | Eksplisit; default ditetapkan di konstruktor |
| Menambah opsi baru | Menambah field (kompatibel mundur untuk struct literal bernama) | Menambah fungsi opsi baru (sangat kompatibel) |
| Validasi | Di konstruktor, terhadap seluruh struct | Per opsi (bisa mengembalikan error) maupun di konstruktor |
| Penemuan API | Autocomplete menampilkan semua field di satu tempat | Opsi tersebar sebagai fungsi `With...` |
| Cocok dipakai | Konfigurasi dibaca dari file, banyak field, tanpa default rumit | API library publik yang akan berkembang, banyak default, opsi opsional |

**Kapan memilih masing-masing:**

- **Struct konfigurasi** untuk kode internal atau aplikasi, jumlah field wajar, dan konfigurasi berasal dari file atau environment. Contoh: `http.Server{...}` di standard library.
- **Functional options** untuk library yang dipakai banyak orang dan perlu menambah opsi tanpa merusak pemanggil, atau yang punya default non-trivial dan nilai nol bermakna. Contoh: `grpc.NewServer(opts ...ServerOption)`.
- Jika hanya satu atau dua parameter wajib, pakai parameter biasa. Kedua pola ini baru berguna ketika jumlah opsi bertambah.

---

## Lampiran: program Praktik dari materi

Ketiga program Praktik (`praktik-fungsi`, `praktik-closure`, `praktik-defer`) sudah dijalankan dan sesuai teori. Beberapa hal yang terlihat dari outputnya:

- `a(), a(), a(), b()` menghasilkan `1 2 3 1`: dua closure punya state terpisah.
- Pada `praktik-defer`, `defer` berjalan LIFO dan nilai `x` yang dicetak adalah `1` (bukan `2`).
- `bungkus()` mengembalikan `bungkus: gagal asli`, karena `defer` memodifikasi named return `err`.
- `ubahNilai` tidak mengubah struct, sedangkan `ubahPointer` mengubahnya. Slice berubah elemennya (`100`), tetapi `len` tetap 3 setelah `append`.
