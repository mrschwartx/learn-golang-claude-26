# Modul 2: Instalasi & Setup Environment — Jawaban Latihan

Kode pendukung ada di folder `salam/`:

```bash
cd salam
go run .
go build -o salam . && ls -lh salam
```

Semua output di bawah dijalankan sungguhan di Go **1.24.7** (linux/amd64). Materi course memakai Go 1.27, jadi angka seperti ukuran binary bisa sedikit berbeda di mesin Anda.

---

## Latihan 1 — `go version` dan `go env GOPATH GOROOT`

```
$ go version
go version go1.24.7 linux/amd64

$ go env GOROOT GOPATH GOPROXY GOTOOLCHAIN
/usr/local/go1.24.7
/root/go
https://proxy.golang.org,direct
auto
```

| Variabel | Fungsi |
|---|---|
| `GOROOT` | Lokasi instalasi Go itu sendiri: compiler, linker, standard library, dan tool bawaan. Jarang perlu diubah manual. |
| `GOPATH` | Direktori kerja Go (default `$HOME/go`). Isinya cache modul yang diunduh (`pkg/mod`) dan binary hasil `go install` (`bin`). Kode proyek **tidak** harus berada di sini. |

---

## Latihan 2 — Modul `contoh.com/salam`

```
$ go mod init contoh.com/salam
go: creating new go.mod: module contoh.com/salam

$ cat go.mod
module contoh.com/salam

go 1.24.7
```

`main.go`:

```go
package main

import "fmt"

func main() {
	nama := "Pembelajar Go"
	fmt.Printf("Halo, %s!\n", nama)
}
```

```
$ go run .
Halo, Pembelajar Go!

$ go build -o salam . && ls -lh salam
-rwxr-xr-x 1 root root 2.2M Nov 25 18:41 salam     # 2.225.780 byte

$ ./salam
Halo, Pembelajar Go!
```

**Ukuran binary: sekitar 2,2 MB**, sesuai perkiraan 1,5–2,5 MB di materi. Itu karena runtime Go (scheduler dan garbage collector) serta paket `fmt` tertanam di dalam binary, sehingga bisa dijalankan di server tanpa memasang apa pun.

---

## Latihan 3 — Kurung kurawal pindah ke baris baru

Kode yang dicoba:

```go
func main()
{
	fmt.Println("Halo")
}
```

Pesan error:

```
./main.go:6:1: syntax error: unexpected semicolon or newline before {
```

**Mengapa aturan ini ada.** Lexer Go otomatis menyisipkan titik koma di akhir baris bila token terakhir pada baris itu adalah identifier, literal, kata kunci `return`/`break`/`continue`/`fallthrough`, operator `++`/`--`, atau penutup `)`, `]`, `}`. Baris `func main()` diakhiri `)`, sehingga dibaca sebagai `func main();`. Parser lalu menemukan `{` yang berdiri sendiri dan menolaknya. Kurung kurawal pembuka karena itu **harus** ada di baris yang sama.

Efek sampingnya positif: karena titik koma tidak perlu ditulis manual, sintaks lebih bersih, dan gaya penempatan `{` hanya satu, sehingga debat gaya hilang (sejalan dengan filosofi `gofmt`).

---

## Latihan 4 — Import yang tidak terpakai

```
./main.go:5:2: "os" imported and not used
```

Kode **tidak dikompilasi sama sekali**, jadi ini error, bukan sekadar peringatan.

**Dampak pada proyek besar:**

- Daftar import selalu mencerminkan dependency yang benar-benar dipakai oleh suatu file, sehingga tidak ada import "sisa" yang menumpuk.
- Waktu kompilasi tetap cepat dan graf dependency tidak membengkak oleh paket yang tak terpakai.
- Reviewer tidak perlu memeriksa import yang mubazir.
- Di praktiknya hal ini tidak menyakitkan, karena editor dengan `goimports` atau `gopls` menambah dan membuang import otomatis saat file disimpan.

---

## Latihan 5 — `go run` vs `go build` vs `go install`

| Perintah | Perbedaan | Situasi nyata |
|---|---|---|
| `go run .` | Kompilasi ke file sementara lalu langsung dijalankan; tidak ada binary permanen. | Eksperimen cepat saat belajar, atau mengulang *edit → run* ketika mengembangkan fitur. |
| `go build -o app .` | Menghasilkan binary di direktori saat ini (atau lokasi `-o`). | Membangun artefak untuk deployment, misalnya di tahap build `Dockerfile` atau pipeline CI, dengan `GOOS`/`GOARCH` untuk cross-compile. |
| `go install ./cmd/app` | Kompilasi lalu menaruh binary di `$GOBIN` (default `$GOPATH/bin`), sehingga bisa dipanggil dari terminal bila folder itu ada di `PATH`. | Memasang tool CLI ke komputer sendiri, misalnya `go install honnef.co/go/tools/cmd/staticcheck@latest`. |

---

## Latihan 6 — Mengubah versi `go` di `go.mod` ke versi lebih baru

Eksperimen: `go.mod` diubah menjadi `go 1.99.0` (lebih baru dari Go 1.24.7 yang terpasang), lalu `go run .`.

**Dengan `GOTOOLCHAIN=auto` (default):** Go mencoba mengunduh toolchain yang diminta.

```
go: downloading go1.99.0 (linux/amd64)
go: download go1.99.0: golang.org/toolchain@v0.0.1-go1.99.0.linux-amd64: Get "https://proxy.golang.org/...": dial tcp: lookup proxy.golang.org ...: no such host
```

Ini membuktikan perilaku auto-download, tetapi pengunduhan gagal karena lingkungan pengerjaan saya tidak punya internet (dan `1.99.0` memang versi fiktif). Di mesin Anda dengan internet, ganti angkanya dengan versi nyata yang lebih baru dari yang terpasang, misalnya `go 1.27.0`. Toolchain itu akan diunduh sekali ke cache lalu dipakai.

**Dengan `GOTOOLCHAIN=local`:** Go menolak berjalan dan menjelaskan alasannya.

```
go: go.mod requires go >= 1.99.0 (running go 1.24.7; GOTOOLCHAIN=local)
```

**Manfaat toolchain management untuk tim besar:**

- **Versi yang sama untuk semua orang.** Versi Go terkunci di `go.mod`, jadi laptop developer, CI, dan build produksi memakai toolchain yang sama dan masalah "di mesin saya jalan" berkurang.
- **Onboarding mudah.** Developer baru tidak perlu mencari versi Go yang tepat; perintah `go` mengurusnya sendiri.
- **Upgrade terkontrol.** Menaikkan versi cukup satu perubahan di `go.mod` (misalnya `go get go@1.27.4`) yang lewat code review seperti perubahan lain.
- **Beda proyek, beda versi.** Satu mesin bisa mengerjakan beberapa proyek dengan versi Go berbeda tanpa mengganti instalasi.

---

## Latihan 7 — Diskusi: mengapa `gofmt` hanya satu gaya

**Manfaat:**

1. **Tidak ada debat gaya.** Tab atau spasi, posisi kurung, dan lebar baris tidak perlu dirapatkan. Gaya sudah ditetapkan.
2. **Kode dari tim mana pun terasa familiar.** Membaca kode library pihak ketiga atau kode Kubernetes terasa seperti membaca kode sendiri.
3. **Diff dan code review bersih.** Perubahan hanya soal logika, bukan spasi.
4. **Tooling lebih mudah.** Refactoring otomatis dan generator kode bisa menghasilkan kode yang langsung rapi.
5. **Format otomatis.** Tidak ada konfigurasi `.editorconfig` atau linter gaya yang harus dirawat per proyek.

**Kekurangan:**

1. **Tidak ada ruang selera pribadi.** Sebagian developer tidak menyukai tab atau gaya tertentu, dan opsinya tidak ada.
2. **Tidak mengatur semuanya.** `gofmt` tidak membatasi panjang baris atau penamaan; aspek itu tetap perlu aturan tim atau linter.
3. **Sulit menyesuaikan dengan basis kode lama.** Menerapkan `gofmt` pada proyek yang sudah ada menghasilkan commit diff besar sekali.

Kesimpulannya, Go menukar kebebasan estetika dengan konsistensi dan produktivitas tim. Untuk bahasa yang menargetkan tim besar dan proyek berumur panjang (lihat Modul 1, 1.6), pertukaran itu masuk akal.

---

## Praktik tambahan yang sudah dicoba

**Formatting (`gofmt` / `go fmt` / `go vet`).** Dengan `main.go` berindentasi spasi dan tanpa baris kosong:

```
$ gofmt -l .
main.go                       # file ini menyimpang dari format standar

$ gofmt -d main.go
 package main
+
 import "fmt"
+
 func main() {
-    fmt.Println("x")
+	fmt.Println("x")
 }

$ go fmt ./...                # memperbaiki file; setelahnya gofmt -l . kosong
$ go vet ./...                # tidak ada temuan
```

**`go doc`:**

```
$ go doc fmt.Println
func Println(a ...any) (n int, err error)
    Println formats using the default formats for its operands and writes to
    standard output. ...
```

`go doc -src strings.ToUpper` menampilkan kode sumber fungsinya langsung di terminal.
