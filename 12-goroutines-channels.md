# Modul 12: Goroutines & Channels

[← Modul 11: Modules, Packages & Code Organization](11-modules-packages.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 13: Sync, Context & Concurrency Patterns →](13-sync-context-patterns.md)

**Isi modul ini**

- [Teori](#teori)
  - [12.1 Concurrency vs Parallelism](#121-concurrency-vs-parallelism)
  - [12.2 Goroutine](#122-goroutine)
  - [12.3 Channel](#123-channel)
  - [12.4 Buffered vs Unbuffered](#124-buffered-vs-unbuffered)
  - [12.5 Select Statement](#125-select-statement)
  - [12.6 Worker Pools](#126-worker-pools)
  - [12.7 Jebakan Umum](#127-jebakan-umum)
  - [12.8 Channel atau Mutex?](#128-channel-atau-mutex)
- [Praktik](#praktik)
  - [Goroutine dan Channel Dasar](#goroutine-dan-channel-dasar)
  - [Ping-Pong dan Select](#ping-pong-dan-select)
  - [Worker Pool](#worker-pool)
  - [Mengamati Deadlock](#mengamati-deadlock)
- [Latihan](#latihan)

---


## Teori

### 12.1 Concurrency vs Parallelism

| Konsep | Arti |
|---|---|
| **Concurrency** | **Menyusun** program sebagai beberapa tugas independen yang bisa berjalan tumpang tindih — soal *struktur* |
| **Parallelism** | **Menjalankan** beberapa tugas benar-benar bersamaan di beberapa core CPU — soal *eksekusi* |

Program concurrent bisa berjalan paralel bila ada lebih dari satu core, tetapi tetap benar di satu core. Go dirancang agar menulis program concurrent semudah menulis program biasa; runtime-lah yang memetakan pekerjaan ke core yang tersedia.

### 12.2 Goroutine

**Goroutine** adalah fungsi yang berjalan secara concurrent, dimulai dengan menambahkan keyword **`go`** di depan pemanggilan fungsi:

```go
go kirimEmail(user)   // kembali seketika; kirimEmail berjalan di goroutine lain
```

Mengapa goroutine disebut *lightweight*:

| Aspek | Goroutine | Thread OS |
|---|---|---|
| Ukuran stack awal | Sekitar **2 KB**, tumbuh dan menyusut dinamis | Sekitar 1–8 MB tetap |
| Biaya pembuatan | Sangat murah — ratusan ribu goroutine wajar | Mahal |
| Penjadwalan | Oleh **Go runtime** di ruang pengguna | Oleh kernel OS |

Runtime memakai **scheduler M:N** (model G-M-P): ribuan goroutine (G) dijadwalkan ke sejumlah thread OS (M) lewat processor logis (P). Jumlah P diatur `GOMAXPROCS` (default = jumlah core CPU). Detail scheduler dibahas di [Modul 20](20-advanced-topics.md).

Aturan penting:

- **Ketika `main()` selesai, seluruh program berhenti** — goroutine lain ikut dimatikan tanpa menunggu.
- Anda harus **menyinkronkan** goroutine (channel atau `sync.WaitGroup`, [Modul 13](13-sync-context-patterns.md)), bukan menebak dengan `time.Sleep`.
- Goroutine tidak mengembalikan nilai; hasilnya dikirim lewat channel atau variabel yang dilindungi sinkronisasi.

### 12.3 Channel

> *"Don't communicate by sharing memory; share memory by communicating."*

**Channel** adalah saluran bertipe yang menghubungkan goroutine — satu goroutine **mengirim** nilai, goroutine lain **menerima**. Channel sekaligus mengirim data dan **menyinkronkan** pengirim dengan penerima.

| Operasi | Sintaks |
|---|---|
| Membuat | `ch := make(chan int)` |
| Mengirim | `ch <- 42` |
| Menerima | `v := <-ch` |
| Menerima + cek status | `v, ok := <-ch` (`ok == false` jika channel tertutup dan kosong) |
| Menutup | `close(ch)` |
| Iterasi sampai tertutup | `for v := range ch { ... }` |

Aturan yang perlu dihafal:

1. **Hanya pengirim yang boleh menutup** channel, dan hanya jika penerima perlu tahu bahwa data sudah habis. Menutup tidak wajib untuk membebaskan memori.
2. **Mengirim ke channel tertutup menyebabkan panic.** Menutup channel yang sudah tertutup juga panic.
3. Menerima dari channel tertutup tidak memblokir — menghasilkan zero value dengan `ok == false`.
4. Channel `nil` memblokir selamanya pada kirim maupun terima.

**Channel terarah** membatasi penggunaan di tanda tangan fungsi, sehingga niat lebih jelas dan compiler membantu: `chan<- int` (hanya kirim) dan `<-chan int` (hanya terima).

### 12.4 Buffered vs Unbuffered

| | **Unbuffered** `make(chan T)` | **Buffered** `make(chan T, n)` |
|---|---|---|
| Kirim memblokir sampai... | Ada penerima yang siap (**rendezvous**) | Buffer penuh |
| Terima memblokir sampai... | Ada pengirim yang siap | Buffer kosong |
| Sifat | Sinkronisasi **ketat** | Antrean dengan kapasitas `n` |
| Cocok untuk | Handshake, sinyal selesai | Memperhalus lonjakan, membatasi concurrency |

`len(ch)` mengembalikan jumlah elemen di buffer dan `cap(ch)` kapasitasnya. Jangan memakai buffer besar hanya untuk "menghilangkan deadlock" — itu biasanya menyembunyikan masalah desain.

**Deadlock** terjadi saat semua goroutine saling menunggu. Runtime Go mendeteksi kasus ketika *seluruh* goroutine terblokir dan mematikan program dengan `fatal error: all goroutines are asleep - deadlock!`. Deadlock yang hanya melibatkan sebagian goroutine **tidak terdeteksi** dan hanya terlihat sebagai program yang macet.

### 12.5 Select Statement

**`select`** menunggu **beberapa operasi channel sekaligus** dan menjalankan salah satu yang siap. Jika beberapa siap bersamaan, dipilih **secara acak**.

```go
select {
case pesan := <-ch1:
    fmt.Println("dari ch1:", pesan)
case ch2 <- 7:
    fmt.Println("terkirim ke ch2")
case <-time.After(2 * time.Second):
    fmt.Println("timeout")
default:
    fmt.Println("tidak ada yang siap")   // membuat select non-blocking
}
```

| Pola | Penjelasan |
|---|---|
| **Timeout** | `case <-time.After(d)` membatasi lama menunggu |
| **Non-blocking** | `default` dijalankan seketika bila tidak ada case yang siap |
| **Pembatalan** | `case <-ctx.Done()` — fondasi pembatalan dengan `context` ([Modul 13](13-sync-context-patterns.md)) |
| **Menonaktifkan case** | Mengatur variabel channel menjadi `nil` membuat case itu tidak pernah terpilih |

### 12.6 Worker Pools

**Worker pool** membatasi jumlah goroutine yang memproses antrean tugas: sejumlah *worker* tetap membaca dari channel `jobs` dan menulis ke channel `results`. Ini mengendalikan konsumsi resource (koneksi database, memori) dibanding membuat satu goroutine per tugas tanpa batas.

```
          ┌─► worker 1 ─┐
jobs ─────┼─► worker 2 ─┼─────► results
          └─► worker 3 ─┘
```

Kerangka umum: buat channel `jobs` dan `results`, jalankan N worker yang `for j := range jobs`, kirim semua tugas lalu `close(jobs)`, dan tunggu seluruh worker selesai sebelum menutup `results` (atau kumpulkan sejumlah hasil yang diketahui).

### 12.7 Jebakan Umum

| Masalah | Penyebab dan penanganan |
|---|---|
| **Goroutine leak** | Goroutine menunggu channel yang tidak akan pernah terisi/tertutup. Selalu pastikan ada jalur keluar (channel `done`, `context`) |
| **Deadlock** | Kirim tanpa penerima, atau saling menunggu. Periksa urutan kirim/terima |
| **Panic: send on closed channel** | Beberapa goroutine menutup channel yang sama. Tunjuk **satu** pemilik yang menutup |
| **Lupa menunggu** | `main` selesai sebelum goroutine selesai — gunakan `WaitGroup` atau channel |
| **Menangkap variabel loop** | Teratasi sejak Go 1.22 ([Modul 5](05-control-flow.md)); di versi lama salin variabel ke dalam loop |
| **Data race** | Beberapa goroutine menulis variabel yang sama tanpa sinkronisasi ([Modul 13](13-sync-context-patterns.md)) |

### 12.8 Channel atau Mutex?

Channel unggul untuk **memindahkan kepemilikan data** dan mengoordinasikan alur kerja (pipeline, worker pool, sinyal). **Mutex** ([Modul 13](13-sync-context-patterns.md)) lebih sederhana dan cepat untuk **melindungi state bersama** seperti cache atau counter. Memilih alat yang tepat jauh lebih penting daripada memaksakan satu gaya.

## Praktik

### Goroutine dan Channel Dasar

```go
// channel.go
package main

import (
	"fmt"
	"sync"
)

func kirimAngka(n int, out chan<- int) {
	for i := 1; i <= n; i++ {
		out <- i
	}
	close(out) // pengirim yang menutup
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("halo dari goroutine")
	}()
	wg.Wait() // tanpa ini main bisa selesai lebih dulu

	// unbuffered: kirim dan terima saling menyinkronkan
	ch := make(chan int)
	go kirimAngka(3, ch)
	for v := range ch {
		fmt.Println("terima:", v)
	}

	// buffered: kirim tidak memblokir selama buffer belum penuh
	buf := make(chan string, 2)
	buf <- "satu"
	buf <- "dua"
	fmt.Println(len(buf), cap(buf))
	fmt.Println(<-buf, <-buf)

	// channel tertutup: ok == false
	tutup := make(chan int, 1)
	tutup <- 9
	close(tutup)
	v, ok := <-tutup
	fmt.Println(v, ok)
	v, ok = <-tutup
	fmt.Println(v, ok)
}
```

### Ping-Pong dan Select

```go
// select.go
package main

import (
	"fmt"
	"time"
)

func main() {
	ping, pong := make(chan int), make(chan int)
	go func() {
		for n := range ping {
			pong <- n + 1
		}
		close(pong)
	}()
	for i := 0; i < 3; i++ {
		ping <- i * 10
		fmt.Println("balasan:", <-pong)
	}
	close(ping)

	// timeout
	lambat := make(chan string)
	go func() {
		time.Sleep(200 * time.Millisecond)
		lambat <- "selesai"
	}()
	select {
	case s := <-lambat:
		fmt.Println(s)
	case <-time.After(50 * time.Millisecond):
		fmt.Println("timeout")
	}

	// non-blocking
	kosong := make(chan int)
	select {
	case v := <-kosong:
		fmt.Println(v)
	default:
		fmt.Println("belum ada data")
	}
}
```

### Worker Pool

```go
// workerpool.go
package main

import (
	"fmt"
	"slices"
	"sync"
	"time"
)

type Hasil struct {
	Job, Worker, Kuadrat int
}

func worker(id int, jobs <-chan int, results chan<- Hasil, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		time.Sleep(10 * time.Millisecond) // simulasi kerja
		results <- Hasil{Job: j, Worker: id, Kuadrat: j * j}
	}
}

func main() {
	jobs := make(chan int, 9)
	results := make(chan Hasil, 9)

	var wg sync.WaitGroup
	for id := 1; id <= 3; id++ {
		wg.Add(1)
		go worker(id, jobs, results, &wg)
	}

	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	close(jobs)

	wg.Wait()
	close(results)

	var semua []Hasil
	for h := range results {
		semua = append(semua, h)
	}
	slices.SortFunc(semua, func(a, b Hasil) int { return a.Job - b.Job })
	for _, h := range semua {
		fmt.Printf("job %d -> %d\n", h.Job, h.Kuadrat)
	}
}
```

Hasil per job selalu sama, tetapi **worker mana** yang memprosesnya berbeda setiap eksekusi. Itu sebabnya hasil diurutkan sebelum dicetak.

### Mengamati Deadlock

```bash
cat > deadlock.go <<'EOF'
package main

func main() {
	ch := make(chan int)
	ch <- 1 // tidak ada penerima
}
EOF
go run deadlock.go
```

Hasilnya `fatal error: all goroutines are asleep - deadlock!`. Perbaiki dengan menjalankan pengirim di goroutine, atau memberi buffer dan memahami konsekuensinya.

## Latihan

1. Tulis program yang menjalankan 5 goroutine, masing-masing mencetak ID-nya, lalu tunggu semuanya dengan `sync.WaitGroup`. Jalankan beberapa kali dan amati urutan output yang berbeda. Mengapa urutannya tidak terjamin?
2. Buat generator `func angkaGenap(n int) <-chan int` yang mengirim `n` bilangan genap pertama lalu menutup channel. Konsumsi dengan `for range` dan jelaskan siapa yang bertanggung jawab menutup channel.
3. Implementasikan fungsi `ambilDenganTimeout` yang menunggu hasil dari goroutine lambat maksimal 100ms dengan `select` dan `time.After`. Apa yang terjadi pada goroutine lambat itu setelah timeout (kaitkan dengan goroutine leak)?
4. Modifikasi worker pool menjadi 5 worker dengan 20 job dan ukur total waktu dibanding 1 worker. Berapa percepatan yang Anda lihat, dan faktor apa yang membatasinya?
5. Sengaja reproduksi tiga kegagalan: deadlock, `send on closed channel`, dan goroutine leak. Untuk masing-masing, catat gejalanya dan tulis perbaikannya.
6. Tulis pipeline tiga tahap — generator angka → kuadrat → cetak — dengan channel terarah (`chan<-` dan `<-chan`) di tanda tangan fungsi. Apa keuntungan channel terarah dibanding `chan int` biasa?
7. Diskusikan: sebuah layanan membuat satu goroutine baru untuk setiap request masuk yang memanggil database, dan database tumbang saat traffic melonjak. Jelaskan bagaimana worker pool atau channel berbuffer sebagai semaphore dapat melindungi database, dan kapan Anda memilih membatasi concurrency daripada menambah kapasitas.

---

[← Modul 11: Modules, Packages & Code Organization](11-modules-packages.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 13: Sync, Context & Concurrency Patterns →](13-sync-context-patterns.md)
