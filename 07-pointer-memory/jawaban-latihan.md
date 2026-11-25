# Modul 7: Pointer & Memory Management — Jawaban Latihan

Seluruh kode ada di folder `kode/` (modul Go `contoh.com/modul7`, satu direktori per program). Jalankan dari dalam `kode/`:

```bash
go run ./l1          # latihan 1 sampai ./l7 (l6 dan l7 lebih baik dibuild dulu, lihat di bawah)
go run ./praktik-pointer   # juga: praktik-jebakan, praktik-gc

go build -o /tmp/l6 ./l6 && GOGC=25 /tmp/l6      # latihan 6
go build -o /tmp/l7 ./l7 && GOMEMLIMIT=64MiB /tmp/l7   # latihan 7
```

Semua output di bawah dijalankan sungguhan di Go 1.24.7 (mesin dengan 2 core). `go vet ./...` dan `gofmt -l` bersih. Angka memori dan waktu bisa berbeda di mesin Anda; yang penting adalah **polanya**.

**Catatan versi:** Materi memakai `new(30)` (ekspresi di `new`, fitur **Go 1.26**). Di Go 1.24 kode asli itu gagal dikompilasi dengan `30 is not a type`. Pada `praktik-pointer` saya ganti dengan helper generik `ptr(30)` yang berefek sama. Di Go 1.26+ Anda bisa kembali ke `new(30)`.

---

## Latihan 1 — `Reset(*int)` vs `ResetNilai(int)` (`kode/l1`)

```go
func Reset(p *int)     { *p = 0 }
func ResetNilai(n int) { n = 0 }
```

```
setelah ResetNilai(a): 42
setelah Reset(&b)    : 0
alamat a=0xc000082040 b=0xc000082048
alamat parameter n (salinan)  : 0xc000082050
nilai pointer p (alamat milik b): 0xc000082048
```

Hanya versi pointer yang mengubah variabel pemanggil. **Alasannya:** Go selalu *pass-by-value*, sehingga setiap argumen disalin ke parameter.

- `ResetNilai(a)` menerima salinan `int` yang tinggal di alamat berbeda (`...050`, bukan `...040`). Mengubah salinan itu tidak menyentuh `a`.
- `Reset(&b)` juga menyalin argumennya, tetapi yang disalin adalah **alamat** `b`. Salinan alamat itu menunjuk ke variabel `b` yang sama (`...048`), sehingga `*p = 0` menulis ke `b` milik pemanggil.

---

## Latihan 2 — `Setor(*Rekening, int)` dan `Cek(Rekening) int` (`kode/l2`)

```go
func Setor(r *Rekening, jumlah int) { r.Saldo += jumlah }
func Cek(r Rekening) int            { return r.Saldo }
```

```
SetorSalah (nilai) : 100
Setor (pointer)    : 150
```

- **`Setor` butuh pointer** karena tugasnya *mengubah* saldo milik pemanggil. Bila parameternya nilai (`SetorSalah`), perubahan hanya terjadi pada salinan dan saldo tetap 100.
- **`Cek` cukup nilai** karena hanya *membaca*; salinan struct kecil sudah memadai dan fungsinya tidak punya efek samping (aman dipakai bersamaan).
- Untuk struct besar, `Cek` sebaiknya juga menerima pointer demi menghindari salinan. Namun itu soal performa, bukan kebenaran, dan sebaiknya diputuskan setelah pengukuran (lihat 7.8).

---

## Latihan 3 — `NewUser(nama string) *User` (`kode/l3`)

```go
func NewUser(nama string) *User {
	u := User{Nama: nama, Umur: 20}
	return &u // alamat variabel lokal dikembalikan
}
```

```
{Andi 20} {Budi 20} true 1
```

Hasil `go build -gcflags='all=-m' -o /dev/null ./l3`:

```
l3/main.go:13:2: moved to heap: u
```

Variabel `v` di `hitungTanpaLolos` tidak mendapat baris seperti itu, artinya tetap di stack.

**Mengapa aman di Go tetapi berbahaya di C?**

- Di C, variabel lokal hidup di *stack frame* yang **dibuang** ketika fungsi selesai. Pointer yang dikembalikan menunjuk ke memori yang sudah tidak valid (*dangling pointer*), dan pemakaiannya menghasilkan *undefined behavior* (data korup atau crash).
- Di Go, **escape analysis** pada waktu kompilasi mendeteksi bahwa `u` masih dirujuk setelah fungsi selesai, lalu memindahkannya ke **heap** (`moved to heap: u`). Garbage collector menjaga objek itu tetap hidup selama masih ada yang menunjuknya, dan membebaskannya setelah tidak terjangkau.
- Dua pemanggilan `NewUser` menghasilkan dua objek terpisah (`a != b` bernilai `true`), jadi tidak ada keadaan yang dibagi secara tidak sengaja.

Harganya: alokasi heap lebih mahal daripada stack dan menambah beban GC. Karena itu jangan mengubah desain hanya demi menghindarinya sebelum ada bukti dari profiling.

---

## Latihan 4 — Mengubah field struct di dalam map (`kode/l4`, `kode/l4-error`)

Kode yang gagal (`l4-error`, diberi build tag `ignore` agar tidak merusak `go build ./...`):

```go
m := map[string]Produk{"kopi": {10}}
m["kopi"].Stok -= 1
```

Pesan error compiler:

```
./main.go:8:2: cannot assign to struct field m["kopi"].Stok in map
```

Penyebabnya: elemen map **tidak addressable**. Map dapat memindahkan elemen ketika tumbuh (rehash), sehingga Go tidak mengizinkan mengambil alamat atau menugaskan ke field elemen secara langsung.

**Dua perbaikan:**

```go
// 1. Salin, ubah, tulis ulang
p := m["kopi"]
p.Stok--
m["kopi"] = p

// 2. Map berisi pointer
mp := map[string]*Produk{"kopi": {"Kopi", 10}}
mp["kopi"].Stok--
```

```
salin-tulis ulang : map[kopi:{Kopi 9}]
map pointer       : {Kopi 9}
recover: runtime error: invalid memory address or nil pointer dereference
```

Baris terakhir menunjukkan jebakan versi pointer: `mp["teh"]` (key tidak ada) bernilai `nil`, sehingga `mp["teh"].Stok--` *panic*. Versi pointer wajib mengecek keberadaan key dengan comma-ok.

**Kapan memilih yang mana:**

| Pilih... | Jika... |
|---|---|
| **Salin-tulis ulang** (`map[string]T`) | Struct kecil, jarang diubah, ingin semantik nilai dan GC ringan (tidak ada alokasi per elemen), dan map tidak diakses banyak goroutine tanpa sinkronisasi |
| **`map[string]*T`** | Struct besar atau sering diubah di tempat, objek perlu dibagi lewat pointer, atau struct mengandung state yang tidak boleh disalin (mis. `sync.Mutex`). Waspadai `nil`, alokasi per elemen, dan aliasing (dua key menunjuk objek yang sama) |

---

## Latihan 5 — Slice kecil menahan array besar (`kode/l5`)

Program membuat array 100 MiB lalu menyimpan potongan 5 byte. Heap diukur dengan `runtime.GC()` + `runtime.ReadMemStats` (`HeapAlloc`).

```
heap awal                       : 0 MiB
setelah ambilKecilBocor (len=5 cap=100 MiB): 100 MiB <- 100 MiB masih tertahan
setelah kecil=nil + GC          : 0 MiB
setelah ambilKecilAman (len=5 cap=8): 0 MiB <- array besar bebas
```

- **Versi bocor:** `besar[:5]` hanya berisi 5 elemen (`len=5`), tetapi `cap=100 MiB` dan pointer-nya menunjuk ke awal array 100 MiB. Selama potongan itu masih dirujuk, seluruh backing array tidak bisa dibebaskan GC: heap tetap **100 MiB**.
- Setelah referensi dilepas (`kecil = nil`) dan GC berjalan, heap turun lagi ke ~0 MiB, membuktikan itu bukan kebocoran permanen, melainkan memori yang ditahan oleh referensi.
- **Versi aman:** `slices.Clone(besar[:5])` membuat array baru kecil (`cap=8`). Array 100 MiB tidak lagi direferensikan, sehingga bisa dibebaskan segera setelah GC: heap ~**0 MiB** saat fungsi mengembalikan hasilnya.

Fungsi pembantu diberi `//go:noinline` agar array besar tidak dioptimalkan lain oleh compiler pada percobaan ini. Dalam praktik, pola ini sering muncul saat membaca file/respons besar lalu hanya menyimpan sebagian kecil dalam cache jangka panjang.

---

## Latihan 6 — `GOGC` 25, 100, dan 400 (`kode/l6`)

Beban kerja: 32 MiB data hidup tetap, ditambah sekitar 2 GiB alokasi sampah jangka pendek (32.768 potongan × 64 KiB).

```
GOGC=25   siklus GC=287  puncak heap=  41 MiB  pause total=18.371ms  waktu=255ms
GOGC=100  siklus GC=65   puncak heap=  64 MiB  pause total=4.08ms   waktu=180ms
GOGC=400  siklus GC=16   puncak heap= 168 MiB  pause total=848µs    waktu=318ms
```

Dengan `GODEBUG=gctrace=1` pada `GOGC=100`, terhitung 68 baris `gc N @...` (bandingkan 65 dari `NumGC` ditambah siklus awal saat data hidup dialokasikan). Contoh dua baris pertama:

```
gc 1 @0.000s 16%: ... 4->4->4 MB, 4 MB goal, ...
gc 2 @0.000s 20%: ... 9->9->9 MB, 9 MB goal, ...
```

Format `A->B->C MB` pada trace: heap saat GC mulai, saat GC selesai, dan ukuran data hidup.

**Kesimpulan trade-off:**

| GOGC | Siklus GC | Puncak heap | Arti |
|---|---|---|---|
| 25 | **paling banyak** (287) | **paling kecil** (41 MiB) | GC sangat sering: hemat memori, tetapi lebih banyak CPU untuk GC dan jeda total lebih besar (18 ms) |
| 100 (default) | sedang (65) | sedang (64 MiB) | Heap sekitar dua kali data hidup |
| 400 | **paling sedikit** (16) | **paling besar** (168 MiB) | GC jarang dan murah, tetapi memori boros |

Jumlah siklus berkurang kira-kira berbanding terbalik dengan GOGC, sedangkan puncak heap naik mengikuti. **GOGC rendah menukar CPU dengan memori; GOGC tinggi menukar memori dengan CPU.** Kolom "waktu" di sini **tidak** bisa dipakai sebagai kesimpulan: mesin ini kecil dan bising (GOGC=400 justru lebih lambat dari 100 pada satu kali percobaan). Untuk membandingkan waktu secara layak perlu benchmark berulang (Modul 15).

---

## Latihan 7 — OOMKill di Kubernetes dan `GOMEMLIMIT` (`kode/l7`)

Simulasi: "layanan" dengan 40 MiB data hidup dan `GOGC` dilonggarkan ke 1000 (GC jarang). Hasil:

```
GOMEMLIMIT=""      siklus GC=2   puncak HeapSys=443 MiB
GOMEMLIMIT="64MiB" siklus GC=30  puncak HeapSys=59 MiB
GOMEMLIMIT="48MiB" siklus GC=179 puncak HeapSys=47 MiB
```

**Mengapa OOMKill terjadi walau GC normal?** Dengan `GOGC` saja, GC dipicu berdasarkan **pertumbuhan heap relatif** terhadap data hidup, tanpa tahu batas memori container. Heap bisa membesar hingga melewati limit 512 MiB sebelum GC berikutnya, dan kernel membunuh proses (OOMKill) tanpa peringatan. Pada simulasi, tanpa `GOMEMLIMIT` heap mencapai 443 MiB hanya dengan 2 siklus GC.

**Peran `GOMEMLIMIT`.** Ini adalah batas memori **total** yang diketahui runtime Go (heap + stack + struktur runtime lain). Ketika memori mendekati batas, GC dijalankan **lebih agresif** walau `GOGC` belum memicunya. Hasilnya, puncak turun dari 443 MiB menjadi 59 MiB (64MiB) dan 47 MiB (48MiB), dengan harga lebih banyak siklus GC (CPU).

**Cara menyetel relatif terhadap limit container (rekomendasi umum):**

- Set `GOMEMLIMIT` di kisaran **90% dari limit container** (misalnya 450 MiB untuk limit 512 MiB, sesuai contoh di materi). Sisakan ruang untuk memori di luar kendali runtime Go: stack, cgo, mmap, buffer kernel, dan sebagainya.
- Di Kubernetes, nilai ini sering diisi lewat environment variable dari `resources.limits.memory` pada Pod spec (misalnya dengan `resourceFieldRef`), supaya tetap sinkron saat limit diubah. Bila *requests* dan *limits* berbeda, limit yang menjadi acuan OOM.
- Biarkan `GOGC` pada nilai default atau kombinasikan (mis. `GOGC=off` bersama `GOMEMLIMIT` hanya untuk beban yang dipahami betul).
- Pantau dengan `GODEBUG=gctrace=1` atau metrik runtime untuk memastikan waktu GC tidak berlebihan.

**Mengapa bersifat *lunak*?**

- Runtime Go tidak bisa menolak alokasi bila data hidup memang lebih besar dari batas. Jika data hidup sudah di atas `GOMEMLIMIT`, program tetap menggunakan memori sebanyak yang dibutuhkan dan tidak crash karena batas ini.
- Untuk mencegah GC yang terus-menerus (*death spiral*) saat melewati batas, runtime membatasi porsi CPU yang boleh dipakai GC (sekitar 50%) sehingga program tetap berjalan, meski memori bisa melebihi batas.
- Karena itu `GOMEMLIMIT` bukan pengganti perencanaan kapasitas. Bila data hidup memang mendekati limit container, solusinya adalah menaikkan limit atau mengurangi pemakaian memori (mengatasi kebocoran, membatasi cache, dan sebagainya).

---

## Lampiran: program Praktik dari materi

- `praktik-pointer` — dijalankan dengan perubahan `new(30)` → `ptr(30)` (lihat catatan versi di atas). Output: `20 20`, `2 1`, `true`, `7`, `30`, `{Andi 31}`.
- `praktik-jebakan` — `tambahSalah` tidak mengubah `len` (`[0]`), `tambahBenar` dan `tambahPtr` berhasil (`[0 99]`, `[0 99 99]`). Mengubah field struct dalam map lewat salin-tulis ulang berhasil (`map[a:{5}]`), dan lewat pointer juga (`9`).
- `praktik-gc` — ketiga pernyataan bernilai `true`: heap ≥ 50 MiB saat data dipegang, turun di bawah 10 MiB setelah `runtime.GC()`, dan GC memang berjalan.
