# Modul 9: Generics & Iterators — Jawaban Latihan

Seluruh kode ada di folder `kode/` (modul Go `contoh.com/modul9`, satu direktori per program). Jalankan dari dalam `kode/`:

```bash
go run ./l1        # latihan 1 sampai ./l7
go run ./praktik-generik   # juga: praktik-stack, praktik-iterator
```

Semua output di bawah dijalankan sungguhan di **Go 1.24.7**. `go vet ./...` bersih.

**Catatan versi (penting).** Materi modul ini memuat fitur **Go 1.27** (method generik, pasal 9.5) yang **tidak dapat dikompilasi di Go 1.24**. Dampaknya:

- Pada `praktik-stack`, `func (s *Stack[T]) Transform[U any](...)` pada Go 1.24 gagal dengan `syntax error: method must have no type parameters`. Saya ganti dengan fungsi generik biasa `Transform(s, f)` yang berefek sama.
- Latihan 5 menyentuh perilaku Go 1.27. Bagian yang bisa saya buktikan (pesan error di Go 1.24) dan bagian yang hanya berdasar materi saya pisahkan dengan jelas.
- Folder `l5-error-interface` dan `l5-error-impl` berisi kode yang sengaja ditolak compiler (build tag `ignore`). Karena memang tidak valid di Go 1.24, `gofmt` juga tidak bisa mem-parse-nya; itu bukan kesalahan.

---

## Latihan 1 — `Keys` dan `Values` generik (`kode/l1`)

```go
func Keys[K comparable, V any](m map[K]V) []K {
	out := make([]K, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
// Values analog; keduanya menerima map apa pun.
```

```
Keys (buatan sendiri)     : [d a c b e]
Values (buatan sendiri)   : [4 1 3 2 5]
maps.Keys + Collect       : [d a c b e]
maps.Values + Collect     : [1 3 2 5 4]
terurut: slices.Sorted    : [a b c d e]
terurut: Keys + Sort      : [a b c d e]
urutan berbeda dari 200 kali Keys(m): 5
```

**Perbandingan.** Hasil dan kompleksitasnya setara. Bedanya:

- `maps.Keys(m)` mengembalikan **iterator** (`iter.Seq[K]`) yang *lazy*: tidak membuat slice sampai Anda memintanya (`slices.Collect`), dan bisa langsung dipakai di `for range` tanpa alokasi.
- `Keys` buatan sendiri langsung membuat slice (dengan kapasitas awal `len(m)`), yang lebih langsung bila memang butuh slice.
- Dalam kode baru, versi standard library sebaiknya dipilih agar tidak menulis ulang.

**Mengapa urutan tidak terjamin?** Iterasi map di Go **sengaja diacak** (lihat Modul 5), dan `maps.Keys`/`maps.Values` hanya meneruskan urutan itu. Pada percobaan, 200 kali pemanggilan `Keys(m)` menghasilkan **5 urutan berbeda** dari satu map berisi 5 elemen.

**Membuatnya deterministik:** urutkan hasilnya, baik dengan `slices.Sorted(maps.Keys(m))` (satu langkah) maupun `Keys(m)` lalu `slices.Sort`. Untuk `Values`, urutkan berdasarkan key atau nilainya sesuai kebutuhan.

---

## Latihan 2 — `Bilangan` dan `Rata[T Bilangan]` (`kode/l2`)

```go
type Bilangan interface{ ~int | ~float64 }

func Rata[T Bilangan](xs []T) float64 {
	if len(xs) == 0 { return 0 }
	var total T
	for _, x := range xs { total += x }
	return float64(total) / float64(len(xs))
}

type Skor int
```

```
2.5
2
90
0
3
```

(Berturut-turut: `[]int{1,2,3,4}`, `[]float64{1.5,2.5}`, `[]Skor{80,90,100}`, slice `nil`, dan `RataKetat([]int{2,4})`.)

**Peran `~` (tilde).** `~int` berarti "`int` **dan semua tipe yang underlying type-nya `int`**". `Skor` didefinisikan sebagai `type Skor int`, sehingga underlying type-nya `int` tetapi ia adalah tipe yang **berbeda** dari `int`. Dengan `~int`, `Skor` memenuhi constraint dan `Rata([]Skor{...})` berhasil (hasil 90).

Pembuktian tanpa tilde: constraint `interface{ int | float64 }` menolak `Skor` dengan pesan compiler

```
Skor does not satisfy BilanganKetat (possibly missing ~ for int in BilanganKetat)
```

Pesan itu sendiri sudah menyarankan solusi: tambahkan `~`. Tilde hampir selalu diperlukan bila ingin mendukung tipe bernama buatan pengguna seperti `Rupiah`, `Skor`, atau `Celsius`.

Catatan implementasi: `Rata` mengembalikan `0` untuk slice kosong agar tidak terjadi pembagian nol (`0/0` pada float menghasilkan `NaN`, bukan *panic*, tetapi `NaN` jarang diinginkan).

---

## Latihan 3 — `Set[T comparable]` (`kode/l3`)

```go
type Set[T comparable] struct{ m map[T]struct{} }

func (s *Set[T]) Add(x T)      { s.m[x] = struct{}{} }
func (s *Set[T]) Has(x T) bool { _, ok := s.m[x]; return ok }
func (s *Set[T]) Remove(x T)   { delete(s.m, x) }
func (s *Set[T]) Len() int     { return len(s.m) }
```

```
Len: 3 | Has(a): true | Has(z): false
setelah Remove(a): Len: 2 | Has(a): false
Set[int]: 3
ukuran struct{} : 0 byte
ukuran bool     : 1 byte
ukuran int      : 8 byte
```

**Mengapa `struct{}` sebagai value?**

- Himpunan hanya perlu mencatat **keberadaan** elemen, bukan nilai pendamping. `struct{}` bukan struct kosong yang "menyimpan sesuatu", melainkan tipe dengan **ukuran 0 byte** (terbukti di output), sehingga tidak menambah memori per elemen. Pada `map[T]bool`, setiap entri menyimpan sebuah byte (plus padding) yang tidak berguna.
- Maknanya jelas bagi pembaca: "tidak ada data di sini, hanya key yang penting". `bool` justru menimbulkan pertanyaan apakah `false` berarti "tidak ada" atau "ada dengan nilai false".
- Keduanya sama cepat dalam praktik; keuntungan `struct{}` ada pada memori dan kejelasan niat.

`T` harus `comparable` karena key map harus bisa dibandingkan dengan `==`. Menambah elemen duplikat (`Add("a")` dua kali) tidak mengubah `Len`.

---

## Latihan 4 — `Queue[T any]` (`kode/l4`)

```go
func (q *Queue[T]) Dequeue() (T, bool) {
	var nol T
	if len(q.items) == 0 {
		return nol, false
	}
	x := q.items[0]
	q.items[0] = nol // lepas referensi agar GC bisa membebaskan
	q.items = q.items[1:]
	return x, true
}
```

```
dequeue: 10
dequeue: 0
dequeue: 30
antrean kosong
elemen 0 yang sah -> v: 0 ok: true
antrean kosong     -> v: 0 ok: false
```

**Mengapa mengembalikan `bool` (bukan panic atau zero value saja)?**

- **Zero value ambigu.** Jika hanya mengembalikan `T`, antrean kosong menghasilkan `0` (atau `""`, `nil`, dst.) yang tidak bisa dibedakan dari data sah bernilai nol. Output menunjukkan dua kasus yang sama-sama menghasilkan `v: 0`, hanya `ok` yang membedakan. Antrean yang berisi `0` tidak boleh dianggap kosong.
- **Panic terlalu keras.** Antrean kosong adalah kondisi **normal** (konsumen menunggu kerja), bukan kesalahan pemrograman. Mengharuskan pemanggil memakai `recover` merepotkan. Dengan `bool` ("comma-ok"), pemanggil cukup menulis `if v, ok := q.Dequeue(); ok {...}`.
- Pola ini konsisten dengan idiom Go lain: pembacaan map (`v, ok := m[k]`), type assertion (`v, ok := x.(T)`), dan penerimaan channel (`v, ok := <-ch`).

Detail penting: setelah mengambil elemen, `q.items[0] = nol` dilakukan sebelum memotong slice. Tanpa itu, backing array tetap menahan referensi ke elemen lama sehingga GC tidak bisa membebaskannya (relevan bila `T` berupa pointer atau struct besar; lihat Modul 7).

---

## Latihan 5 — Method generik di interface dan di tipe (`kode/l5`, `kode/l5-error-*`)

**Bagian yang dibuktikan di Go 1.24.** Dua percobaan:

```go
// (a) di dalam interface
type Pengubah interface {
	Transform[U any](f func(int) U) []U
}

// (b) pada tipe
func (s *Stack[T]) Transform[U any](f func(T) U) []U { ... }
```

Pesan compiler:

```
(a) ./main.go:4:11: interface method must have no type parameters
    ./main.go:4:31: undefined: U
(b) ./main.go:6:29: syntax error: method must have no type parameters
```

Pada Go 1.24, **kedua bentuk ditolak**: method tidak boleh punya type parameter sendiri (hanya yang dibawa tipenya, mis. `T` pada `Stack[T]`).

**Bagian yang berdasar materi (tidak bisa saya jalankan di sini, butuh Go 1.27).** Menurut pasal 9.5, mulai Go 1.27:

- Bentuk (b) menjadi **valid**: method pada tipe boleh punya type parameter sendiri.
- Bentuk (a) **tetap ditolak** dengan pesan `interface method must have no type parameters`, yang identik dengan pesan yang saya dapat di Go 1.24 untuk bentuk (a).
- Method generik **tidak bisa dipakai untuk memenuhi** method non-generik di interface, karena pengecekan interface hanya mengenali method dengan signature konkret. Jadi mengisi interface dengan tipe yang punya method generik tidak membuatnya memenuhi interface itu. Bagian ini tidak saya uji sendiri; silakan verifikasi di toolchain Go 1.27.

**Batasan yang bisa disimpulkan:** method generik adalah cara **mengelompokkan operasi pada satu tipe**, bukan alat membuat abstraksi lewat interface. Interface Go bekerja dengan daftar method konkret, sehingga method generik tidak punya satu signature tunggal yang bisa dicocokkan.

**Pola pengganti yang valid di semua versi** (dipakai pada `l5` dan `praktik-stack`): fungsi generik biasa yang menerima tipe sebagai argumen pertama:

```go
func Transform[T, U any](s *Stack[T], f func(T) U) []U
```

---

## Latihan 6 — Iterator bilangan prima tak hingga (`kode/l6`)

```go
func Prima() iter.Seq[int] {
	return func(yield func(int) bool) {
		for n := 2; ; n++ {
			if !laPrima(n) { continue }
			if !yield(n) { // false -> konsumen melakukan break
				return
			}
		}
	}
}
```

```
10 prima pertama:
2 3 5 7 11 13 17 19 23 29   [iterator berhenti & membersihkan diri; yield dipanggil 10 kali]

iterator buruk -> runtime error: range function continued iteration after function for loop body returned false
```

**Peran nilai kembalian `yield`.** `yield` adalah badan loop `for ... range` yang dibungkus menjadi fungsi. `yield(v)` mengirim nilai ke badan loop dan mengembalikan:

- `true` bila loop ingin lanjut ke iterasi berikutnya,
- `false` bila pemanggil keluar dari loop (`break`, `return`, atau `goto` keluar).

**Mengapa iterator harus berhenti saat `yield` mengembalikan `false`?**

1. **Deret ini tak hingga.** Tanpa pengecekan itu, iterator tidak pernah berakhir dan `for` yang sudah di-`break` tetap mencari prima berikutnya selamanya (kebocoran CPU atau goroutine).
2. **Pembersihan resource.** Iterator yang membuka file, koneksi, atau mengunci mutex harus menutupnya ketika konsumen berhenti. `defer` di dalam iterator hanya akan dijalankan bila fungsi iterator benar-benar return (terlihat dari pesan "berhenti & membersihkan diri" setelah 10 kali `yield`).
3. **Kebenaran program.** Go **memeriksa** kontrak ini saat runtime. Iterator "buruk" yang mengabaikan hasil `yield` dan tetap memanggilnya lagi memicu *panic*: `range function continued iteration after function for loop body returned false`.

Iterator bersifat *lazy*: hanya 10 prima yang pernah dihitung. Hasil `yield dipanggil 10 kali` membuktikannya.

---

## Latihan 7 — `Proses[T any]` dengan `switch any(x).(type)` (`kode/l7`)

```go
func ProsesKeliru[T any](x T) string {
	switch v := any(x).(type) {
	case int:    return fmt.Sprint("int ", v*2)
	case string: return "string " + strings.ToUpper(v)
	default:     return "tidak didukung"
	}
}
```

```
int 42 | string GO
tidak didukung <- float64 lolos kompilasi, baru gagal saat RUNTIME
int 42 | string GO          (fungsi terpisah per tipe)
int 42 | string GO          (interface Pemroses)
```

Pada alternatif interface, `ProsesInterface(3.14)` ditolak **saat kompilasi**:

```
cannot use 3.14 (constant of type float64) as Pemroses value in argument to ProsesInterface: float64 does not implement Pemroses (missing method Proses)
```

**Mengapa ini pertanda generics dipakai keliru?**

- **Tidak benar-benar generik.** Generics bernilai bila **satu algoritma yang sama** bekerja untuk banyak tipe. Di sini perilaku **berbeda per tipe** (`int` dikali 2, `string` di-uppercase), dan semuanya dipaksa lewat `any(x)` untuk kembali ke type switch. Itu sama saja dengan fungsi `interface{}` lama, hanya dengan tambahan kurung siku.
- **Keamanan tipe hilang.** Constraint `any` menerima semua tipe, jadi `ProsesKeliru(3.14)` lolos kompilasi dan baru ketahuan di `default` saat runtime (terlihat di output). Itu persis masalah yang ingin dipecahkan generics (9.1).
- **Daftar tipe tersembunyi di badan fungsi.** Pemanggil tidak bisa tahu tipe mana yang didukung dari signature, dan menambah tipe baru berarti mengedit fungsi pusat.

**Alternatif yang lebih idiomatik:**

1. **Interface.** Tiap tipe mengimplementasikan perilakunya sendiri (`Proses() string`). Fungsi menerima interface itu. Tipe yang tidak didukung ditolak compiler, dan tipe baru tidak mengubah kode yang ada.
2. **Fungsi terpisah per tipe** (`ProsesInt`, `ProsesString`). Paling sederhana ketika jumlah tipe kecil dan tetap, dan tidak ada kebutuhan polimorfisme.
3. Bila memang ada operasi bersama, pakai **constraint dengan method** atau constraint tipe (`~int | ~float64`) agar badan fungsi hanya memakai operasi yang dijamin constraint, tanpa pengecekan tipe manual.

Sejalan dengan pasal 9.9: mulai dari kode konkret, dan jadikan generik hanya bila logika yang sama benar-benar diulang untuk tipe lain. `switch` pada tipe di dalam fungsi generik hampir selalu pertanda abstraksi yang salah.

---

## Lampiran: program Praktik dari materi

- `praktik-generik` — `6 4 1500` (`Jumlah` untuk `int`, `float64`, dan `Rupiah` yang memenuhi `~int`), `[1 2 3]`, `[2 4 6] 56`, `jeruk 7 2.5`, `2 -1`.
- `praktik-stack` — dijalankan dengan `Transform` sebagai fungsi biasa (catatan versi di atas): `[1 2]`, `bb true`, `{1 true} true`, `{0 false} false`. Cache generik dengan `Pair[int, bool]` bekerja normal.
- `praktik-iterator` — `0 1 2 3`, `[0 2 4 6 8]`, `0 1 1 2 3 5 8 13 21 34`, `[a b c]`, dan `iter.Pull` menghasilkan `pull:0 pull:1 pull:2`. Fibonacci tak hingga aman karena dihentikan dengan `break`.

Fitur Go 1.26/1.27 lain pada materi (constraint yang merujuk dirinya sendiri `Adder[A Adder[A]]`, inferensi generik yang diperluas) tidak saya coba karena toolchain di sini Go 1.24.
