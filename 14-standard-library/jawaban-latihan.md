# Jawaban Latihan — Modul 14: Standard Library

Kode ada di `kode/` (module `contoh.com/modul14`). Toolchain: **Go 1.24.7**, 2 CPU. Semua angka di bawah adalah hasil run nyata.
Semua Praktik (berkas, flag, time+json, slog+regexp, go:embed) dijalankan dari kode materi tanpa perubahan; hasilnya cocok dengan penjelasan materi.
`go vet ./...`, `gofmt -l`, dan `go test ./...` bersih.

> **Catatan versi.** Bagian materi tentang `encoding/json/v2` (Go 1.27) tidak bisa dijalankan di Go 1.24, jadi hanya diambil dari materi, tidak saya verifikasi.
> `os.OpenRoot` (Go 1.24) ada di toolchain ini, tetapi tidak dipakai latihan mana pun.

## Latihan 1 — Hitung baris/kata/rune dan 5 kata teratas (`l1`)
`bufio.Scanner` + `FieldsFunc` (pemisah: bukan huruf/angka) + map frekuensi; rune dihitung dengan `utf8.RuneCountInString`. Ada test (termasuk karakter non-ASCII `çé`, yang dihitung 1 rune, bukan 2 byte).

Berkas uji 95 MiB (2.000.000 baris, dibuat oleh program itu sendiri):
```
baris=2000000 kata=20000000 rune=99600000
top5: [go=4000000 kopi=4000000 baris=2000000 belajar=2000000 hari=2000000]
```
Perbandingan memori (puncak RSS dari `/proc/self/status`, VmHWM):

| Cara | RSS puncak |
|---|---|
| `bufio.Scanner` (per baris) | ±7,4 MB |
| `os.ReadFile` (seluruh berkas) | ±99,7 MB |

**Mengapa `os.ReadFile` kurang cocok untuk berkas puluhan GB:** ia memuat seluruh isi ke satu `[]byte`, jadi memori naik sebanding ukuran berkas (di sini ≈ 1× ukuran berkas), sedangkan Scanner memori-nya hampir konstan (hanya buffer + map kata). Untuk puluhan GB, `ReadFile` bisa membuat proses kehabisan memori (OOM). Catatan jujur: map frekuensi tetap tumbuh sebanding jumlah kata *unik*, bukan jumlah kata total; pada berkas uji ini hanya segelintir kata unik.

Batas baris: tanpa `sc.Buffer(...)`, baris 100 KB gagal: `baris terbaca: 0 | sc.Err(): bufio.Scanner: token too long` (`l1-default`). Dengan `Buffer` 16 MB, baris yang sama terbaca. Karena itu `sc.Err()` selalu dicek.

## Latihan 2 — CLI dua subcommand (`l2`)
`tambah` (`-nama` wajib, `-jumlah`, `-tag`) dan `hapus` (`-id` wajib, `-paksa`), masing-masing `flag.NewFlagSet` dengan `ContinueOnError`. Logika ada di `Jalankan(args, out, errOut) int` sehingga bisa diuji tanpa `os.Exit`.
```
$ app tambah -nama kopi -jumlah 3 -tag a,b sisa1   → tambah "kopi" x3 tag=[a b] sisa=[sisa1]   (exit 0)
$ app hapus -id 7 -paksa                           → hapus id=7 paksa=true                    (exit 0)
$ app hapus -h                                     → Usage of hapus: ...                      (exit 0)
$ app tambah -jumlah x                             → invalid value "x" for flag -jumlah       (exit 2)
$ app ubah                                         → subcommand tidak dikenal: "ubah"         (exit 2)
```
9 kasus test lulus, termasuk bukti bahwa flag milik subcommand lain (`hapus -nama x`) ditolak. Temuan: dengan `ContinueOnError`, `-h` mengembalikan `flag.ErrHelp`; awalnya saya mengembalikan kode 2 untuk itu, lalu saya ubah menjadi 0 karena meminta bantuan bukan kesalahan.

## Latihan 3 — Selisih hari dan layout (`l3`)
Hari ini (26 Nov 2026, zona Asia/Jakarta) → 17 Agustus 2027: **264 hari**.
Tiga layout: `2027-08-17`, `17/08/2027`, `Tuesday, 17 August 2027` (RFC3339: `2027-08-17T00:00:00+07:00`).
Jebakan yang terukur:
- `target.Sub(now)` langsung memberi **263,30 hari** (pecahan karena jam saat ini 16:41); untuk selisih hari kalender saya bandingkan tanggal pada tengah malam UTC (`HariKe`, dengan test). Hasilnya sama untuk `TZ=UTC` dan `TZ=Asia/Jakarta` (264).
- `Format("YYYY-MM-DD")` mencetak `YYYY-MM-DD` apa adanya (bukan layout).
- Layout dengan angka tertukar (`2006-02-01`) ditangkap oleh `go vet` ("should be 2006-01-02"), sehingga baris itu tidak saya masukkan ke program; ini alasan bagus menjalankan `go vet`.
- Instan sama di zona berbeda: `u == w` → `false`, `u.Equal(w)` → `true`.

**Mengapa tanggal referensi `2006-01-02 15:04:05`?** Layout ditulis sebagai contoh bentuk keluaran: tiap bagian adalah angka unik yang dikenali parser (bulan=1, hari=2, jam=3 (15=24 jam), menit=4, detik=5, tahun=6 (2006), zona=7 (MST/-0700)), sehingga `Mon Jan 2 15:04:05 MST 2006` = `1 2 3 4 5 6 7`. Keuntungannya: format terbaca sebagai contoh nyata, dan tidak ada ambiguitas huruf (`MM` bulan vs menit, yang di bahasa lain harus dibedakan `mm`/`MM`). Kelemahannya: sulit dihafal. Ini penjelasan desain dari materi; alasan historis persisnya tidak saya verifikasi.

## Latihan 4 — `Rupiah` JSON (`l4`)
`type Rupiah int64`; `MarshalJSON` (receiver nilai) menghasilkan `"Rp 15.000"`; `UnmarshalJSON` (receiver pointer) mem-parse balik, memvalidasi awalan `Rp ` dan kelompok ribuan (1–3 digit pertama, sisanya tepat 3), mendukung negatif.
```
{"nama":"Kopi","harga":"Rp 15.000"}
<nil> true {Kopi Rp 15.000}          ← round trip identik
{"harga":15000}      -> rupiah harus string: json: cannot unmarshal number ...
{"harga":"15000"}    -> awalan "Rp " tidak ada
{"harga":"Rp 15.00"} -> format ribuan salah
{"harga":"Rp 1.5000"}-> format ribuan salah
{"harga":"Rp abc"}   -> strconv.ParseInt: ... invalid syntax
```
Test: round trip untuk 0, 5, 999, 1000, 15000, 1.234.567, 1 miliar, −2500; 7 input tidak valid ditolak (termasuk angka di luar rentang int64); juga berfungsi di `map[string]Rupiah`.

## Latihan 5 — `slog` TextHandler + `LOG_LEVEL` (`l5`)
`NewTextHandler`, `logger.With("request_id","abc-123")`, level dari `LOG_LEVEL` lewat `Level.UnmarshalText`. Keluaran (waktu dibuang agar deterministik):
```
LOG_LEVEL kosong  → INFO, WARN, ERROR (DEBUG tersaring)
LOG_LEVEL=debug   → DEBUG, INFO, WARN, ERROR
LOG_LEVEL=warn    → WARN, ERROR
LOG_LEVEL=error   → ERROR
LOG_LEVEL=ngawur  → WARN "memakai level default" err="...unknown name", lalu level INFO
level=INFO msg=login request_id=abc-123 user=andi ip=10.0.0.7
```
Dengan `LOG_FORMAT=json` hasilnya `{"level":"INFO","msg":"login","request_id":"abc-123","user":"andi","ip":"10.0.0.7"}`.

**Mengapa log terstruktur lebih mudah dicari di Elasticsearch:** tiap atribut menjadi field terindeks sendiri (`request_id`, `user`, `ms`), sehingga query berupa filter/agregasi (`request_id:abc-123 AND level:ERROR`, rata-rata `ms`) tanpa parsing teks bebas dengan regex yang rapuh terhadap perubahan kalimat. Tipe data juga terjaga (angka tetap angka). Itu alasan konseptual dari materi; saya tidak menjalankan Elasticsearch di sini.

## Latihan 6 — Regex email dengan named group (`l6`)
Pola: `^(?P<nama>[A-Za-z0-9._%+-]+)@(?P<domain>(?:[A-Za-z0-9](?:[A-Za-z0-9-]*[A-Za-z0-9])?\.)+[A-Za-z]{2,})$` (praktis, bukan RFC 5322 penuh).
- 5 valid, semua cocok, mis. `budi.santoso+tag@mail.example.co.id` → nama `budi.santoso+tag`, domain `mail.example.co.id`.
- 8 tidak valid, semua ditolak: `andi@`, `@contoh.com`, `andi@contoh`, spasi, `@@`, `andi@-x.com`, `andi@x..com`, `andi@x.c`.

**Mengapa aman dari ReDoS:** RE2 tidak melakukan backtracking, waktu proporsional dengan panjang input. Pola klasik `^(a+)+$` pada `aaa…ab`:

| n | Go (RE2) | Python `re` (backtracking, pembanding) |
|---|---|---|
| 1.000 | 41 µs | |
| 10.000 | 529 µs | |
| 100.000 | 5,2 ms | |
| 1.000.000 | 36,8 ms | |
| 16 / 20 / 24 / 26 | | 0,003 s / 0,038 s / 0,60 s / 2,42 s |

Go tumbuh kira-kira linear; Python menjadi ≈4× lebih lambat tiap n bertambah 2 (eksponensial). Angka waktu bergantung mesin, tren-nya yang penting.
**Yang dikorbankan:** backreference (`(a)\1` → error `invalid escape sequence`), lookahead (`(?=`), lookbehind (`(?<!`): ketiganya ditolak saat kompilasi.

## Latihan 7 — `go:embed` vs menyalin aset ke image Docker (`l7`)
Aplikasi web kecil: `templates/*.html` (embed.FS + `template.ParseFS`), `static/` (`http.FileServerFS`), `migrasi/*.sql` (dibaca dengan `fs.ReadDir`). Diuji lewat `httptest`:
```
/                     -> 200  (HTML ter-render, menyebut 2 migrasi tertanam)
/static/app.css       -> 200
/static/tidak-ada.css -> 404
```
Binary (±10 MB) disalin sendirian ke direktori kosong dan tetap menyajikan semuanya: aset memang ada di dalam executable.
Aturan yang dibuktikan: `//go:embed ../rahasia.txt` → error kompilasi `invalid pattern syntax`; dengan pola direktori (`//go:embed d`), berkas berawalan `.` (`.tersembunyi`) **tidak** ikut tertanam (hanya `a.txt`).

**Manfaat dibanding menyalin folder aset ke image Docker:** (1) satu artefak yang tak bisa "tidak sinkron" dengan kodenya: versi binary = versi templat/migrasi; (2) image bisa `scratch`/distroless tanpa `COPY` aset dan tanpa soal path kerja (`WORKDIR`); (3) tidak ada kegagalan runtime karena berkas hilang/salah path, error muncul saat kompilasi; (4) migrasi SQL ikut versi rilis; (5) tidak mungkin dimodifikasi di server. Harga: binary lebih besar, dan mengubah satu CSS berarti build ulang.
**Kapan sebaiknya tetap di luar binary:** aset besar (video, dataset, model; memperbesar binary dan memori), aset yang perlu diganti tanpa rilis ulang (tema/branding per pelanggan, konfigurasi, templat email yang diedit tim non-engineer), data yang berubah saat runtime, dan rahasia (kunci, sertifikat; yang tertanam bisa diekstrak dari binary). Aset statis besar yang lebih baik dilayani CDN juga sebaiknya tidak ikut. Ini pertimbangan desain; di sini saya tidak membangun image Docker (tidak ada Docker).
