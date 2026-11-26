# Modul 12: Goroutines & Channels — Jawaban Latihan

Seluruh kode ada di folder `kode/` (modul Go `contoh.com/modul12`, satu direktori per program). Jalankan dari dalam `kode/`:

```bash
go run ./l1        # latihan 1 sampai ./l7 (l4 dan l7 lebih baik dibuild dulu: go build -o /tmp/l4 ./l4)
go test -v ./l3    # tes yang membuktikan perbaikan goroutine leak
go run ./praktik-channel   # juga: praktik-select, praktik-workerpool, praktik-deadlock
```

Semua output di bawah dijalankan sungguhan di **Go 1.24.7** pada mesin **2 core** (`GOMAXPROCS=2`). `go vet ./...` dan `gofmt -l` bersih. Angka waktu dan urutan goroutine bergantung mesin; yang penting adalah polanya.

---

## Latihan 1 — 5 goroutine dengan `sync.WaitGroup` (`kode/l1`)

```go
var wg sync.WaitGroup
for id := 1; id <= 5; id++ {
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Print("goroutine-", id, " ")
	}()
}
wg.Wait()
```

Contoh output:

```
goroutine-5 goroutine-1 goroutine-2 goroutine-3 goroutine-4
```

**Mengamati urutan.** Saya menjalankan program ini berulang kali dan menghitung urutan yang muncul:

| Percobaan | Hasil |
|---|---|
| 300 eksekusi, 5 goroutine, `GOMAXPROCS=2` | **Beberapa urutan berbeda** (mis. `5 1 2 3 4` muncul 186 kali, `5 3 4 1 2` 30 kali, `5 2 1 3 4` 22 kali, dan variasi lain, termasuk satu kali `1 5 2 3 4`) |
| 300 eksekusi, `GOMAXPROCS=8` | Variasi lebih banyak |
| 100 eksekusi, `GOMAXPROCS=1` | **Selalu** `5 1 2 3 4` (hanya kebetulan implementasi, bukan jaminan) |
| 200 eksekusi, 20 goroutine | **121 urutan berbeda** |

Jadi pada 5 goroutine yang sangat ringan, satu urutan mendominasi, dan baru terlihat bervariasi bila dijalankan banyak kali atau dengan lebih banyak goroutine. Urutan dominan `5 1 2 3 4` terjadi karena scheduler memberi slot khusus (*runnext*) kepada goroutine yang dibuat paling akhir. Itu optimasi implementasi, bukan perilaku yang boleh diandalkan.

**Mengapa urutannya tidak terjamin?**

1. Goroutine dijadwalkan oleh **scheduler Go** (model G-M-P) ke beberapa thread OS yang berjalan paralel di beberapa core. Tidak ada aturan "yang dibuat lebih dulu berjalan lebih dulu".
2. Kapan sebuah goroutine mulai berjalan bergantung pada kondisi runtime (beban core, pencurian pekerjaan antar-P, preemption), yang berubah antar-eksekusi.
3. Spesifikasi Go hanya menjamin urutan lewat **sinkronisasi** (channel, mutex, `WaitGroup`). `wg.Wait()` menjamin *semua selesai*, bukan *urutan selesainya*.

Bila urutan penting, kumpulkan hasil lalu urutkan (seperti di worker pool), atau gunakan channel untuk menyusun tahapan.

---

## Latihan 2 — Generator `angkaGenap` (`kode/l2`)

```go
func angkaGenap(n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out) // pengirim yang menutup
		for i := 1; i <= n; i++ {
			out <- i * 2
		}
	}()
	return out
}

for v := range angkaGenap(5) { fmt.Print(v, " ") }
```

```
2 4 6 8 10
```

**Siapa yang bertanggung jawab menutup channel?** **Pengirim** (di sini goroutine di dalam generator), bukan penerima. Alasannya:

- Hanya pengirim yang tahu kapan data sudah habis. Bila penerima menutup, pengirim yang masih mengirim akan **panic** (`send on closed channel`).
- Menutup berfungsi sebagai sinyal "tidak ada data lagi", sehingga `for range` di sisi konsumen berhenti sendiri. Tanpa `close`, `range` akan menunggu selamanya (dan program bisa deadlock).
- Penggunaan `defer close(out)` memastikan channel tertutup pada semua jalur keluar goroutine.
- Penulis generator yang mengembalikan `<-chan int` (hanya-terima) juga *memaksa* konsumen tidak bisa menutup atau mengirim, sehingga aturan kepemilikan ditegakkan oleh compiler (lihat Latihan 6).

---

## Latihan 3 — `ambilDenganTimeout` dan goroutine leak (`kode/l3`)

```go
func ambilDenganTimeout(d time.Duration) (string, error) {
	ch := make(chan string)           // unbuffered
	go lambat(d, ch)                  // tidur d, lalu ch <- "selesai"
	select {
	case s := <-ch:
		return s, nil
	case <-time.After(100 * time.Millisecond):
		return "", fmt.Errorf("timeout setelah 100ms")
	}
}
```

```
goroutine awal: 1
cepat (30ms): selesai <nil>
BOCOR   : setelah 50 timeout & menunggu 2,5 dtk, goroutine = 51
```

**Apa yang terjadi pada goroutine lambat setelah timeout?** `ambilDenganTimeout` sudah kembali dan tidak ada lagi yang membaca `ch`. Ketika goroutine lambat akhirnya selesai tidur, ia melakukan `ch <- "selesai"` pada channel **unbuffered tanpa penerima**, lalu **terblokir selamanya**. Itu adalah **goroutine leak**: goroutine dan stack-nya tidak pernah dibebaskan. Hasilnya: 50 pemanggilan yang timeout meninggalkan **50 goroutine yatim** (1 awal + 50 = 51), sekalipun saya menunggu 2,5 detik hingga semua selesai tidur. Dalam layanan nyata, kebocoran ini menumpuk perlahan hingga memori habis.

**Dua perbaikan** (ada di `l3/main.go`, dibuktikan oleh `go test -v ./l3`):

1. **Channel berbuffer 1**: `make(chan string, 1)`. Goroutine lambat selalu bisa menyimpan hasilnya di buffer lalu selesai, walaupun tidak ada yang membaca.
2. **`context.WithTimeout`**: goroutine ikut memantau `ctx.Done()` dan berhenti lebih awal, sehingga pekerjaan tidak diteruskan sia-sia.

```
fix_test.go:19: tidak bocor: awal=2 akhir=2
--- PASS: TestPerbaikanTidakBocor
```

Jumlah goroutine sebelum dan sesudah 100 pemanggilan yang timeout sama (2 dan 2). Aturan umum: **setiap goroutine harus punya jalur keluar yang pasti**, entah lewat buffer, `done`, atau `context` (Modul 13).

---

## Latihan 4 — 5 worker vs 1 worker, 20 job (`kode/l4`)

Dua jenis beban diuji pada mesin dengan `GOMAXPROCS = 2`:

```
GOMAXPROCS = 2 | NumCPU = 2
BEBAN MENUNGGU (sleep 50ms x 20 job): 1 worker=1.01s | 5 worker=201ms | percepatan=5.0x
BEBAN CPU     (20 job):               1 worker=1.167s | 5 worker=587ms | percepatan=2.0x
```

**Percepatan dan faktor pembatasnya:**

| Jenis beban | Percepatan | Penjelasan |
|---|---|---|
| **Menunggu** (I/O, `sleep`) | **5,0×** (sesuai jumlah worker) | Worker yang sedang menunggu tidak memakai CPU, sehingga 5 worker bisa menunggu bersamaan. 20 job / 5 worker = 4 putaran × 50 ms ≈ 200 ms |
| **CPU murni** | **2,0×** (bukan 5×) | Dibatasi **jumlah core fisik** (di sini 2). Menambah worker melebihi jumlah core tidak mempercepat pekerjaan yang terikat CPU |

Faktor lain yang membatasi percepatan di dunia nyata:

- **Jumlah job** terhadap jumlah worker (dengan 20 job dan 5 worker, tidak mungkin lebih cepat dari 4 putaran).
- **Hukum Amdahl**: bagian program yang berurutan (membuat job, mengumpulkan hasil) tidak bisa diparalelkan.
- **Kontensi**: channel, mutex, atau database bersama bisa menjadi leher botol.
- **Overhead** koordinasi, serta batas resource eksternal (koneksi database, bandwidth).

Pelajaran praktis: untuk kerja I/O, jumlah worker boleh lebih besar dari jumlah core; untuk kerja CPU, ukurannya sekitar `runtime.NumCPU()`.

---

## Latihan 5 — Tiga kegagalan: deadlock, send on closed channel, goroutine leak

Setiap kegagalan ada di direktori sendiri (`kode/l5-deadlock`, `l5-sendclosed`, `l5-leak`).

### a. Deadlock (`l5-deadlock`)

```go
ch := make(chan int)
ch <- 1 // tidak ada penerima
```

**Gejala:**

```
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
main.main()
	.../l5-deadlock/main.go:8 +0x36
exit code: 2
```

Runtime mendeteksi bahwa **seluruh** goroutine terblokir dan mematikan program. (Deadlock yang hanya melibatkan sebagian goroutine **tidak terdeteksi**; program hanya macet.)

**Perbaikan:** kirim dari goroutine terpisah sehingga ada penerima yang siap, atau beri buffer jika memang desainnya antrean, misalnya:

```go
go func() { ch <- 1 }()
fmt.Println(<-ch)
```

### b. `send on closed channel` (`l5-sendclosed`)

```
ditutup; sekarang mengirim lagi...
panic: send on closed channel

goroutine 1 [running]:
main.main()
	.../l5-sendclosed/main.go:22 +0x138
exit code: 2
```

**Gejala:** *panic* (exit code 2), bukan error biasa. Pola umum penyebabnya: beberapa goroutine saling menutup channel yang sama, atau mengirim setelah `close`.

**Perbaikan:** tunjuk **satu pemilik** (pengirim) yang menutup, dan hanya setelah semua pengirim selesai. Untuk banyak pengirim, tutup channel sesudah `wg.Wait()` di satu goroutine koordinator:

```go
go func() { wg.Wait(); close(ch) }()
```

Menutup channel yang sudah tertutup juga *panic*, jadi pastikan hanya sekali.

### c. Goroutine leak (`l5-leak`)

```go
// produsen mengirim tanpa batas; konsumen hanya mengambil 1 nilai
go func() { for i := 0; ; i++ { out <- i } }()
```

**Gejala:** program tidak crash, tetapi jumlah goroutine terus naik:

```
BOCOR : goroutine awal = 1 -> sesudah 100 panggilan = 101
AMAN  : 100 panggilan tambahan, goroutine sekarang = 101 (tidak bertambah; selisih 100 berasal dari bagian BOCOR)
```

Setelah 100 panggilan versi bocor, ada 100 goroutine yatim yang terblokir pada `out <- i`. Versi perbaikan (100 panggilan tambahan) **tidak menambah** goroutine sama sekali. Pada sistem nyata gejalanya adalah memori yang merayap naik, dan terlihat dengan `runtime.NumGoroutine()`, profil goroutine, atau metrik.

**Perbaikan:** beri produsen jalur keluar lewat channel `done` (atau `context`):

```go
select {
case out <- i:
case <-done:
	return
}
```

dan konsumen memanggil `close(done)` ketika selesai.

---

## Latihan 6 — Pipeline tiga tahap dengan channel terarah (`kode/l6`)

```go
func generator(n int) <-chan int            // hanya menghasilkan
func kuadrat(in <-chan int) <-chan int      // membaca <-chan, menghasilkan <-chan baru
func cetak(in <-chan int)                   // hanya membaca

cetak(kuadrat(generator(5)))
```

```
hasil: 1
hasil: 4
hasil: 9
hasil: 16
hasil: 25
```

Setiap tahap menutup channel keluarannya dengan `defer close(out)`, sehingga penutupan menjalar dari generator sampai `cetak`, dan `for range` di tiap tahap berhenti dengan sendirinya.

**Keuntungan channel terarah dibanding `chan int`:** *kontrak penggunaan ditegakkan compiler*. Sebuah fungsi yang menerima `<-chan int` tidak bisa mengirim atau menutup. Bukti dari kode yang sengaja salah (`kode/l6-error`):

```
./main.go:4:2: invalid operation: cannot send to receive-only channel in (variable of type <-chan int)
./main.go:5:8: invalid operation: cannot close receive-only channel in (variable of type <-chan int)
./main.go:9:8: invalid operation: cannot receive from send-only channel out (variable of type chan<- int)
```

Manfaatnya:

1. **Niat terdokumentasi** di signature: pembaca langsung tahu fungsi itu produsen atau konsumen.
2. **Mencegah bug kepemilikan**: konsumen tidak bisa menutup channel yang bukan miliknya (sumber utama `panic: close of closed channel` dan `send on closed channel`).
3. Konversi dari `chan T` ke `<-chan T` atau `chan<- T` terjadi **otomatis** saat dipanggil, sehingga tidak merepotkan.

---

## Latihan 7 — Melindungi database dengan batas concurrency (`kode/l7`)

Simulasi: database hanya sanggup **10 query bersamaan**; permintaan di atas itu ditolak (mewakili "tumbang"). Ada 500 request, tiap query 20 ms.

```
tanpa batas:                 berhasil=10   ditolak=490  puncak koneksi=12   waktu=20ms
worker pool (10):            berhasil=500  ditolak=0    puncak koneksi=10   waktu=1.04s
semaphore berbuffer (10):    berhasil=500  ditolak=0    puncak koneksi=10   waktu=1.02s
```

**Cara kerja perlindungannya:**

- **Tanpa batas** (satu goroutine per request): seluruh 500 request menyerbu database sekaligus. Hanya 10 yang lolos, **490 gagal**. Angka "puncak koneksi=12" terlihat kecil hanya karena request yang ditolak segera kembali; pada database sungguhan, banjir koneksi inilah yang membuatnya kehabisan memori atau koneksi dan tumbang.
- **Worker pool 10 worker**: hanya 10 goroutine yang memanggil database, sisanya mengantre di channel `jobs`. **Seluruh 500 berhasil, 0 gagal**, dengan puncak koneksi tepat 10.
- **Channel berbuffer sebagai semaphore**: `sem := make(chan struct{}, 10)`. Goroutine mengambil slot sebelum query (`sem <- struct{}{}`) dan melepasnya sesudahnya (`<-sem`). Hasil sama baik, dengan kode yang lebih sederhana bila goroutine per request sudah ada.

Biayanya terlihat jelas: total waktu naik dari 20 ms (yang gagal) menjadi sekitar 1 detik (yang berhasil), karena pekerjaan yang diterima dilayani bertahap, 10 sekaligus (500 / 10 × 20 ms ≈ 1 detik).

**Kapan membatasi concurrency, bukan menambah kapasitas?**

| Pilih membatasi concurrency bila... | Pilih menambah kapasitas bila... |
|---|---|
| Resource bersama punya batas keras (koneksi DB, file descriptor, API pihak ketiga ber-*rate limit*) | Beban rata-rata memang melebihi kapasitas, bukan sekadar lonjakan |
| Lonjakan bersifat sementara dan antrean bisa menyerapnya | Latensi antrean sudah tak bisa diterima pengguna |
| Biaya menambah kapasitas tinggi atau lambat (mis. replika baru) | Kapasitas bisa ditambah murah dan cepat (autoscaling) |
| Anda ingin kegagalan yang **terkendali** (antre atau tolak dengan jelas) alih-alih kegagalan beruntun | Bottleneck-nya ada di tempat lain |

Praktik tambahan: beri antrean **batas dan timeout** (tolak dengan error bila penuh, jangan antre tanpa batas), dan pantau panjang antrean. Dalam praktik biasanya kedua langkah dipadukan: batasi concurrency agar sistem tidak tumbang, lalu tambah kapasitas bila antrean terus memanjang. Catatan: jumlah worker (10) di sini ditetapkan dari kapasitas database, bukan dari jumlah core CPU.

---

## Lampiran: program Praktik dari materi

Semua dijalankan dan sesuai teori:

- `praktik-channel`: `halo dari goroutine`, `terima: 1..3`, `2 2` (buffer len dan cap), `satu dua`, lalu `9 true` dan `0 false` untuk channel tertutup.
- `praktik-select`: ping-pong `1`, `11`, `21`; `timeout`; `belum ada data` (select non-blocking dengan `default`).
- `praktik-workerpool`: sembilan hasil `job N -> N*N` terurut. Seperti dicatat materi, hasil per job selalu sama sedangkan worker pemrosesnya berbeda tiap eksekusi.
- `praktik-deadlock`: `fatal error: all goroutines are asleep - deadlock!` dengan jejak `goroutine 1 [chan send]`. Dengan `go run` terlihat `exit status 2`.

Hal yang tidak saya jalankan: `go run -race` (race detector) karena tidak ada latihan di modul ini yang memerlukannya; topik data race dibahas di Modul 13.
