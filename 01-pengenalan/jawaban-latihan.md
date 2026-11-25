# Modul 1: Pengenalan Go — Jawaban Latihan

Kode pendukung ada di folder `hello/`. Jalankan dengan:

```bash
cd hello
go run .
```

Semua output di bawah ini dijalankan sungguhan (Go 1.24 di lingkungan pengerjaan; kode ini tidak memakai fitur yang khusus Go 1.27).

---

## Latihan 1 — Hello di Go Playground (nama + loop 5 kali)

Kode (`hello/main.go`):

```go
package main

import "fmt"

func main() {
	nama := "Pembelajar Go"
	fmt.Printf("Halo, Go! Saya %s.\n", nama)

	for i := 1; i <= 5; i++ {
		fmt.Println("Iterasi ke-", i)
	}
}
```

Output:

```
Halo, Go! Saya Pembelajar Go.
Iterasi ke- 1
Iterasi ke- 2
Iterasi ke- 3
Iterasi ke- 4
Iterasi ke- 5
```

**Eksperimen import** (hasil nyata dari compiler):

| Eksperimen | Pesan compiler |
|---|---|
| Hapus `import "fmt"` tapi tetap memakai `fmt.Println` | `./main.go:4:2: undefined: fmt` |
| Tambah `import "os"` tanpa memakainya | `./main.go:5:2: "os" imported and not used` |

**Tombol Format / `gofmt`.** Kode berantakan seperti ini:

```go
package main
import "fmt"
func main(){
fmt.Println( "x" )
   for i:=0;i<2;i++{fmt.Println(i)}
}
```

diubah otomatis menjadi:

```go
package main

import "fmt"

func main() {
	fmt.Println("x")
	for i := 0; i < 2; i++ {
		fmt.Println(i)
	}
}
```

Spasi, indentasi (tab), dan pemecahan baris ditentukan oleh tool, bukan selera programmer. Inilah prinsip 1.6 poin 5 (*satu cara format*): karena semua kode Go berformat sama, tim tidak perlu berdebat soal gaya, dan diff di code review hanya berisi perubahan yang bermakna.

---

## Latihan 2 — Tiga alasan satu binary tanpa dependency menguntungkan deployment

1. **Server/image tidak perlu runtime.** Tidak ada JVM, interpreter Python, atau Node.js + `node_modules` yang harus dipasang dan dijaga versinya. Docker image bisa sangat kecil (misalnya `scratch` atau distroless), sehingga *pull* lebih cepat dan permukaan serangan lebih kecil.
2. **Bebas dependency hell di server.** Tidak ada kasus "jalan di laptop, gagal di server" karena versi library atau interpreter berbeda. Artefak yang diuji sama persis dengan yang di-deploy.
3. **Distribusi dan rollback sederhana.** Deploy cukup menyalin satu file lalu menjalankannya; rollback cukup mengganti file dengan versi sebelumnya. Lewat cross-compilation, satu mesin build bisa menghasilkan binary untuk Linux, macOS, dan Windows.

---

## Latihan 3 — Go 1 compatibility promise

**Isi:** sejak Go 1.0 (Maret 2012), program Go 1.x yang valid dijanjikan tetap bisa dikompilasi dan berjalan dengan rilis Go 1.x berikutnya. Karena itu nomor versi tetap `1.x` dan belum ada rencana Go 2 yang merusak kompatibilitas. (Janji ini mencakup bahasa dan standard library; ada pengecualian kecil seperti perbaikan bug keamanan atau paket di bawah `x/`.)

**Mengapa penting bagi perusahaan:**

- Biaya pemeliharaan rendah: upgrade versi Go biasanya cukup mengganti toolchain, tanpa proyek migrasi besar.
- Investasi kode aman: kode yang ditulis bertahun-tahun lalu tetap bernilai.
- Perusahaan lebih mudah mengikuti patch keamanan, karena hanya dua rilis terakhir yang didukung (lihat 1.8) dan upgrade-nya berisiko rendah.
- Perencanaan jangka panjang lebih mudah dan pelatihan ulang tim jarang diperlukan, karena bahasa tidak berubah drastis antar rilis.

---

## Latihan 4 — Dua use case yang relevan (dari tabel 1.5)

Karena course ini bagian dari seri Docker/Kubernetes/Kafka/dst., saya memilih dua use case yang paling nyambung dengan seri tersebut. Ganti sesuai konteks Anda sendiri bila berbeda.

**a. CLI tool.**
Go menghasilkan satu binary yang bisa di-*cross-compile* untuk banyak OS, startup-nya instan, dan tidak butuh runtime di mesin pengguna. Dibanding skrip Python/Node, yang perlu interpreter dan dependency terpasang, tool Go cukup diunduh lalu dijalankan. Itu sebabnya `kubectl`, Hugo, dan GitHub CLI memakai Go.

**b. Cloud-native & infrastruktur (operator/controller Kubernetes).**
Library resmi Kubernetes (`client-go`) dan framework operator (Kubebuilder, controller-runtime) berbasis Go, sehingga kode kita memakai tipe dan API yang sama dengan Kubernetes sendiri. Concurrency bawaan (goroutine) cocok untuk controller yang memantau banyak resource sekaligus, dan binary kecil mudah dikemas sebagai container.

---

## Latihan 5 — Docker dan Kubernetes ditulis dalam Go

Dari pengetahuan umum: Docker Engine ada di repo `moby/moby` dan Kubernetes di `kubernetes/kubernetes`, dan keduanya memang ditulis dalam Go (README modul juga menyebut Terraform dan Prometheus). Lingkungan pengerjaan ini tidak punya akses internet, jadi **silakan Anda verifikasi sendiri** dengan membuka kedua repo di GitHub dan melihat bar bahasa (Go mendominasi).

**Manfaat bisa membaca kode sumber tool yang dipakai sehari-hari:**

- Saat perilaku tool tidak jelas atau dokumentasinya kurang, kita bisa membaca kodenya untuk tahu persisnya apa yang terjadi (misalnya bagaimana `kubectl` membentuk request ke API server).
- Pesan error yang membingungkan bisa dilacak sampai asalnya.
- Kita bisa menulis ekstensi: operator, plugin `kubectl`, provider Terraform, atau exporter Prometheus, dengan bahasa dan tipe data yang sama.
- Kita bisa berkontribusi perbaikan bug atau fitur ke proyek upstream.
- Pola idiomatis Go (interface kecil, error sebagai nilai, goroutine) bisa dipelajari dari kode produksi berkualitas tinggi.

---

## Latihan 6 — Diskusi: layanan API puluhan ribu koneksi bersamaan di Kubernetes

**Faktor dari modul ini yang mendukung Go:**

1. **Concurrency bawaan.** Goroutine sangat ringan, sehingga puluhan ribu koneksi bisa ditangani dengan model "satu goroutine per koneksi" yang sederhana, tanpa callback bertingkat atau thread pool rumit (1.1, 1.3).
2. **Static binary.** Image container kecil, *pull* dan startup cepat, jadi scale-out Pod dan rolling update di Kubernetes cepat (1.3, 1.5).
3. **Standard library lengkap.** `net/http` siap produksi tanpa framework tambahan, sehingga lebih sedikit dependency yang harus dirawat (1.3).
4. **Garbage collector dengan jeda sangat singkat.** Latensi API relatif stabil (1.3).
5. **Ekosistem cloud-native.** Client Kubernetes, gRPC, dan Prometheus punya dukungan kelas satu di Go (1.5, 1.7).
6. **Kemudahan dibaca dan Go 1 compatibility.** Tim besar bisa saling membaca kode, dan upgrade jangka panjang berisiko rendah (1.3, 1.2).

**Situasi yang lebih cocok memakai bahasa lain:**

- Jika inti layanan adalah **ML/komputasi numerik atau data science**, ekosistem Python (NumPy, PyTorch, dst.) jauh lebih matang.
- Atau bila tim seluruhnya berpengalaman Java dengan ekosistem dan library enterprise yang sudah ada, biaya pindah bahasa mungkin lebih besar daripada keuntungannya.
- Atau bila butuh kontrol memori tingkat rendah tanpa GC (misalnya komponen yang sangat sensitif terhadap latensi), Rust bisa lebih tepat.

Keputusan akhir tetap bergantung pada keahlian tim dan kebutuhan spesifik, sesuai prinsip di 1.4 bahwa tidak ada bahasa yang terbaik untuk segalanya.
