# Jawaban Latihan — Modul 13: Sync, Context & Concurrency Patterns

Semua kode ada di `kode/` (module `contoh.com/modul13`). Toolchain yang dipakai: **Go 1.24.7**, 2 CPU.
Semua angka di bawah adalah hasil run nyata di mesin ini, bukan salinan dari materi.

> **Catatan versi.** Materi memakai `wg.Go(f)` (Go 1.25+). Toolchain di sini 1.24, jadi di `praktik-perbaikan` dan `praktik-pipeline`
> saya ganti dengan helper `jalankan(wg, f)` (`wg.Add(1)` + `go func(){ defer wg.Done(); f() }()`), perilakunya setara.
> Hal yang bergantung pada fitur Go 1.25+ lainnya tidak saya jalankan.

## Latihan 1 — Race pada counter (`l1`)
1000 goroutine melakukan `hitung++` tanpa sinkronisasi. 20 kali run:
9× hasil 1000 (benar), 11× hasil salah (994–999 dst., terendah 989).
**Pelajaran:** bug race tidak selalu muncul; 9 dari 20 run "kebetulan benar". Karena itu uji harus memakai `-race`, bukan menilai dari hasil.

## Latihan 2 — Cache dengan `sync.RWMutex` (`l2`)
`Cache` dengan `RLock` di `Get`, `Lock` di `Set`. Test 100 goroutine lulus, juga dengan `go test -race`.
Versi buruk (`l2-buruk`, map biasa ditulis 100 goroutine) **selalu** crash: `fatal error: concurrent map writes`, exit code 2.
Berbeda dengan counter, race pada map terdeteksi runtime Go, jadi tidak diam-diam salah.

## Latihan 3 — Mutex vs atomic (`l3`)
| Benchmark | ns/op |
|---|---|
| Mutex serial | 16.0 |
| Atomic serial | 4.7 |
| Mutex parallel | 23.4 |
| Atomic parallel | 30.4 |

Serial: atomic ±3,4× lebih cepat. Parallel di 2 CPU: atomic justru sedikit **lebih lambat** dari mutex (kontensi cache line).
Jadi "atomic selalu lebih cepat" tidak terbukti di sini; ini hanya 2 core, hasil di mesin lain bisa berbeda.
`l3-atomic-tidakcukup`: dua saldo diubah dengan atomic saja, invarian `a+b==1000` dilanggar 325.241 kali; dengan mutex 0 kali.
Atomic menjamin satu operasi utuh, bukan invarian lintas dua variabel.

## Latihan 4 — Context timeout pada HTTP (`l4`)
`AmbilHalaman(ctx, url)` memakai `context.WithTimeout(ctx, 2s)` dan `http.NewRequestWithContext`. Tiga test (httptest) lulus:
cepat → sukses; lambat → `DeadlineExceeded` (2,00 s); parent dibatalkan → `Canceled` (0,10 s). Error dibungkus `%w` sehingga `errors.Is` bekerja.

## Latihan 5 — Pipeline dan goroutine leak (`l5`)
Pipeline generator → kuadrat×3 → gabung, berhenti setelah 3 hasil.
- `gabung` dengan ctx: selisih goroutine 0.
- `gabung` dari materi (tanpa ctx): run pertama selisih 0, tetapi diulang 10×: **6× selisih 0, 4× selisih 2**.

Jadi kebocoran memang terjadi, tetapi tidak deterministik: bergantung pada apakah goroutine pengirim sempat sampai ke send yang tak pernah dibaca
saat hitungan diambil. Satu kali run yang menunjukkan 0 tidak membuktikan tidak ada bug. Saya hanya melaporkan angka terukur; penjelasan mekanismenya
(pengirim tersangkut di send tanpa penerima) konsisten dengan materi, tetapi tidak saya verifikasi dengan dump goroutine.

## Latihan 6 — Semaphore vs worker pool (`l6`)
20 tugas @30 ms, batas 3. Dua pendekatan sama-sama puncak serentak 3 dan waktu 210 ms (dua kali run).
Beda: semaphore membuat 21 goroutine; worker pool hanya 4. Hasil waktu sama, biaya goroutine berbeda.

## Latihan 7 — Graceful shutdown (`l7`)
Server dengan handler 2 detik; klien curl mengirim request, lalu SIGTERM dikirim pada ±0,5 s setelahnya.
- **Tanpa graceful:** proses mati (exit 143), klien putus: `curl exit=52`, tanpa respons.
- **Graceful** (`signal.NotifyContext` + `srv.Shutdown` timeout 5 s): request tetap dilayani, klien menerima HTTP 200 `selesai`, server keluar bersih (exit 0).

Di Kubernetes: SIGTERM dikirim saat pod dihentikan; Pod diberi `terminationGracePeriodSeconds` (default 30 s) sebelum SIGKILL.
Timeout `Shutdown` harus lebih kecil dari nilai itu. Pada rolling update, tanpa graceful shutdown request yang sedang berjalan akan putus seperti kasus pertama.

## Lampiran — Praktik dari materi
`praktik-race`, `praktik-perbaikan`, `praktik-context`, `praktik-pipeline`, `praktik-shutdown` dijalankan dari kode materi
(dengan penyesuaian `wg.Go` di atas). Jalankan semua dengan `cd kode && go vet ./... && go test ./...`.
