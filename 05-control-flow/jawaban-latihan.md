# Modul 5: Control Flow (Conditionals & Loops) — Jawaban Latihan

Seluruh kode ada di folder `kode/` (modul Go `contoh.com/modul5`, satu direktori per program). Jalankan dari dalam `kode/`:

```bash
go run ./l1        # latihan 1 sampai ./l6
go run ./praktik-kondisi   # juga: praktik-loop, praktik-closure, praktik-goto

# Latihan 7 memakai dua modul terpisah (go.mod berbeda):
cd l7-go121 && go run .    # go 1.21 -> semantik lama
cd ../l7-go122 && go run . # go 1.22 -> semantik baru
```

Semua output di bawah dijalankan sungguhan di Go 1.24.7. `gofmt -l` bersih, dan `go vet ./...` bersih (kecuali `l7-go121` yang sengaja dipakai untuk memperlihatkan bug lama).

---

## Latihan 1 — `Klasifikasi(suhu float64) string` (`kode/l1`)

```go
func Klasifikasi(suhu float64) string {
	switch {
	case suhu < 15:
		return "dingin"
	case suhu < 25:
		return "sejuk"
	case suhu < 32:
		return "hangat"
	default:
		return "panas"
	}
}
```

```
 -5.00°C -> dingin
 14.90°C -> dingin
 15.00°C -> sejuk
 24.99°C -> sejuk
 25.00°C -> hangat
 31.50°C -> hangat
 32.00°C -> panas
 40.00°C -> panas
```

`case` dievaluasi dari atas ke bawah dan berhenti pada yang pertama benar, sehingga urutan kondisi membuat batas bawah tidak perlu ditulis. Nilai batas (15, 25, 32) sengaja diuji: batas jatuh ke kategori di atasnya karena memakai `<`.

---

## Latihan 2 — FizzBuzz sebagai fungsi (`kode/l2`)

Versi `switch`:

```go
func FizzBuzzSwitch(n int) []string {
	hasil := make([]string, 0, n)
	for i := 1; i <= n; i++ {
		switch {
		case i%15 == 0:
			hasil = append(hasil, "FizzBuzz")
		case i%3 == 0:
			hasil = append(hasil, "Fizz")
		case i%5 == 0:
			hasil = append(hasil, "Buzz")
		default:
			hasil = append(hasil, strconv.Itoa(i))
		}
	}
	return hasil
}
```

Versi `if-else` memakai logika yang sama dengan rantai `if i%15 == 0 {...} else if i%3 == 0 {...} else if ...` (lengkapnya di `kode/l2/main.go`).

```
[1 2 Fizz 4 Buzz Fizz 7 8 Fizz Buzz 11 Fizz 13 14 FizzBuzz]
kedua versi sama: true
```

**Mana yang lebih mudah dibaca?** Versi `switch` tanpa ekspresi, dengan alasan berikut:

- Setiap `case` berdiri sejajar, jadi pembaca langsung melihat daftar aturan sebagai tabel kondisi. Pada `if-else if` yang panjang, indentasi dan kurung kurawal menambah noise visual.
- Menambah aturan baru (misalnya kelipatan 7) cukup menambah satu `case` di posisi yang tepat.
- Tidak ada `break`, tidak ada risiko fall-through, sehingga lebih aman daripada `switch` di C/Java.

Secara perilaku keduanya identik, dan compiler biasanya menghasilkan kode serupa. Rantai `if-else` tetap wajar untuk dua atau tiga cabang. Perhatikan juga bahwa urutan penting: kasus kelipatan 15 harus diperiksa lebih dulu, karena kelipatan 15 juga kelipatan 3 dan 5.

---

## Latihan 3 — Tabel perkalian dan `break` berlabel (`kode/l3`)

```go
tabel:
	for i := 1; i <= 9; i++ {
		for j := 1; j <= 9; j++ {
			hasil := i * j
			if hasil > 50 {
				fmt.Printf("\nberhenti: %d x %d = %d (> 50)\n", i, j, hasil)
				break tabel
			}
			fmt.Printf("%3d", hasil)
		}
		fmt.Println()
	}
```

```
  1  2  3  4  5  6  7  8  9
  2  4  6  8 10 12 14 16 18
  3  6  9 12 15 18 21 24 27
  4  8 12 16 20 24 28 32 36
  5 10 15 20 25 30 35 40 45
  6 12 18 24 30 36 42 48
berhenti: 6 x 9 = 54 (> 50)
```

Hasil pertama yang melewati 50 adalah 6 × 9 = 54, dan seluruh loop berhenti saat itu juga. Baris 7, 8, dan 9 tidak pernah dicetak.

Sebagai pembanding, `break` tanpa label di loop dalam hanya menghentikan loop dalam itu: baris 7, 8, 9 tetap dicetak, hanya terpotong di kolom yang hasilnya melewati 50 (lihat bagian kedua output program).

---

## Latihan 4 — `v` pada `range` adalah salinan (`kode/l4`)

```
1. salah (salinan)       : [{a 1} {b 2} {c 3}]
   &v==&items[0]? false
2. benar (lewat indeks)   : [{a 10} {b 20} {c 30}]
3. benar (slice pointer)  : {x 10} {y 20}
```

**Bukti:** setelah `for _, v := range items { v.Stok *= 10 }`, isi `items` tidak berubah. Alamat `&v` juga berbeda dari alamat elemennya (`&items[0]`), sehingga `v` memang variabel terpisah berisi salinan struct.

**Dua cara memperbaiki:**

1. **Lewat indeks:** `for i := range items { items[i].Stok *= 10 }` mengubah elemen aslinya.
2. **Slice berisi pointer:** untuk `[]*Item`, `v` tetap salinan, tetapi yang disalin adalah *pointer*, yang menunjuk ke objek yang sama, sehingga `p.Stok *= 10` mengubah objek aslinya. Cara ini cocok bila struct besar (menghindari salinan) atau sudah umum diakses lewat pointer; harganya adalah alokasi per objek dan kemungkinan `nil`.

---

## Latihan 5 — Iterasi `"Go 世界 🌍"` (`kode/l5`)

```
len (byte): 14

loop indeks (per byte):
0:47 1:6f 2:20 3:e4 4:b8 5:96 6:e7 7:95 8:8c 9:20 10:f0 11:9f 12:8c 13:8d
iterasi: 14

range (per rune):
0:G(1 byte) 1:o(1 byte) 2: (1 byte) 3:世(3 byte) 6:界(3 byte) 9: (1 byte) 10:🌍(4 byte)
iterasi: 7
```

| Cara | Iterasi | Nilai |
|---|---|---|
| Loop indeks | **14** | Satu **byte** per iterasi (nilai 0–255) |
| `range` | **7** | Satu **rune** (karakter Unicode) per iterasi; indeks adalah posisi byte awal rune |

**Penjelasan:** string Go berisi byte UTF-8. Karakter ASCII memakai 1 byte (`G`, `o`, dua spasi), karakter CJK memakai 3 byte (`世`, `界`), dan emoji 🌍 memakai 4 byte. Totalnya 1+1+1+3+3+1+4 = 14 byte untuk 7 karakter. Loop indeks melihat tiap byte, termasuk byte-byte lanjutan yang bukan karakter utuh (misalnya `e4 b8 96` adalah satu karakter `世`). `range` mendekode UTF-8 otomatis, sehingga indeksnya melompat (0, 1, 2, 3, **6**, 9, 10). Gunakan `range` untuk teks manusia, dan loop indeks hanya bila memang memproses byte.

---

## Latihan 6 — Map terurut berdasarkan nama (`kode/l6`)

```go
for _, k := range slices.Sorted(maps.Keys(nilai)) {
	fmt.Printf("  %-6s %d\n", k, nilai[k])
}
```

```
terurut berdasarkan nama:
  Andi   90
  Budi   80
  Citra  85
  Dewi   70
  Eko    95
```

Urutan iterasi mentah berubah dari satu eksekusi ke eksekusi lain. Saya menjalankan program enam kali dan urutannya antara lain `Budi Andi Citra Dewi Eko`, `Dewi Eko Budi Andi Citra`, dan `Andi Citra Dewi Eko Budi`.

**Mengapa Go sengaja mengacak urutan iterasi map?**

- Spesifikasi Go tidak menjanjikan urutan apa pun. Dulu banyak program tanpa sadar bergantung pada urutan "kebetulan" dari implementasi tertentu. Ketika implementasi map berubah (dan memang berubah antar versi), program mereka rusak.
- Dengan mengacak urutan secara sengaja, ketergantungan semacam itu langsung ketahuan sejak awal, bukan setelah upgrade versi atau pada produksi.
- Selain itu, implementasi map bebas memilih struktur data yang paling cepat tanpa harus menjaga urutan penyisipan.

Jika butuh urutan yang stabil, urutkan key terlebih dahulu seperti di atas, atau gunakan slice terurut, atau struktur data lain.

---

## Latihan 7 — Bug variabel loop sebelum Go 1.22 (`kode/l7-go121`, `kode/l7-go122`)

**Kode yang sama** dijalankan pada dua modul yang hanya berbeda di baris `go` pada `go.mod`:

```go
for _, u := range users { // users = andi, budi, citra
	wg.Add(1)
	go func() {
		defer wg.Done()
		/* simpan u */
	}()
}
```

Hasil dengan `go 1.21` di `go.mod` (lima kali dijalankan, hasil sama):

```
citra citra citra <- closure
[citra citra citra] <- goroutine
```

Hasil dengan `go 1.22` di `go.mod`:

```
andi budi citra <- closure
[andi budi citra] <- goroutine
```

`go vet` pada modul `go 1.21` juga langsung menandai kodenya: `loop variable u captured by func literal`.

**Penyebab.** Sebelum Go 1.22, `for _, u := range users` membuat **satu variabel `u` untuk seluruh loop**, yang ditimpa nilainya di setiap iterasi. Closure (atau goroutine) hanya menangkap *referensi* ke variabel itu, bukan salinan nilainya. Saat goroutine akhirnya berjalan, loop sudah selesai dan `u` berisi nilai terakhir (`citra`), sehingga semua goroutine memproses user yang sama. Waktu jalannya goroutine tidak deterministik, jadi bug ini kadang muncul dengan nilai campuran dan sulit direproduksi. (Pada percobaan di sini selalu `citra` karena loop selesai lebih dulu.)

**Perbaikan di Go 1.22.** Setiap iterasi sekarang memiliki **salinan variabel sendiri**, sehingga tiap closure menangkap `u` miliknya. Pada kode loop tiga komponen (`for i := 0; ...`) nilai variabel disalin ke iterasi berikutnya, jadi perilaku per-iterasi juga berlaku di sana.

**Yang harus dicek di `go.mod`.** Perbaikan ini hanya aktif bila direktif `go` bernilai **`go 1.22` atau lebih baru**. Tim itu harus memeriksa baris tersebut (misalnya `go 1.21`) dan menaikkannya (`go get go@1.22` atau lebih baru), dengan catatan:

- Toolchain yang terpasang harus Go 1.22+ (lihat toolchain management di Modul 2).
- Semantik berlaku per modul dan per file; modul dengan direktif lama tetap memakai perilaku lama walau toolchain baru.
- Setelah menaikkan versi, uji ulang kode yang bergantung pada perilaku lama (jarang, tetapi ada, misalnya kode yang sengaja menimpa variabel loop). Sampai versi dinaikkan, solusi sementaranya adalah menyalin variabel di dalam loop: `u := u`.

---

## Lampiran: program Praktik dari materi

Keempat program Praktik (`praktik-kondisi`, `praktik-loop`, `praktik-closure`, `praktik-goto`) sudah dijalankan dan hasilnya sesuai teori. Beberapa hal yang bisa diamati dari outputnya:

- `fallthrough` mencetak "tiga (fallthrough, tanpa cek kondisi)" tanpa memeriksa kondisi `case 3`.
- Pada `"hé"`, loop indeks menghasilkan tiga byte (`104 195 169`), sedangkan `range` menghasilkan dua rune (`0:h 1:é`).
- `continue luar` dan `break luar` menghasilkan `(0,0) (0,1) (1,0) (1,1)` saja.
- Mengubah `u` pada `for _, u := range users` tidak mengubah slice (`false false`), sedangkan lewat indeks berhasil (`true true`).
