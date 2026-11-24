# Modul 17: gRPC, Database & Logging

[← Modul 16: Membangun CLI & Web Development](16-cli-web-development.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 18: Go Toolchain, Code Quality & Security →](18-toolchain-code-quality.md)

**Isi modul ini**

- [Teori](#teori)
  - [17.1 gRPC dan Protocol Buffers](#171-grpc-dan-protocol-buffers)
  - [17.2 Akses Database dengan `database/sql`](#172-akses-database-dengan-databasesql)
  - [17.3 `pgx` — Driver PostgreSQL](#173-pgx--driver-postgresql)
  - [17.4 ORM: GORM](#174-orm-gorm)
  - [17.5 Memilih Pendekatan Akses Data](#175-memilih-pendekatan-akses-data)
  - [17.6 Migrasi Skema](#176-migrasi-skema)
  - [17.7 Logging Berperforma Tinggi: Zerolog dan Zap](#177-logging-berperforma-tinggi-zerolog-dan-zap)
  - [17.8 Menyambung ke Course Lain](#178-menyambung-ke-course-lain)
- [Praktik](#praktik)
  - [Repository dengan `database/sql` (Hanya Standard Library)](#repository-dengan-databasesql-hanya-standard-library)
  - [Menjalankan PostgreSQL dan Menyambungkannya](#menjalankan-postgresql-dan-menyambungkannya)
  - [Contoh gRPC (Kerangka)](#contoh-grpc-kerangka)
  - [Zap dan Zerolog (Perbandingan Singkat)](#zap-dan-zerolog-perbandingan-singkat)
- [Latihan](#latihan)

---


## Teori

> **Catatan:** library pada modul ini adalah **pihak ketiga** dan versinya terus berkembang. Contoh kode memakai API stabil yang umum, tetapi selalu cocokkan dengan dokumentasi versi yang Anda pasang (`go get nama@latest`, lalu baca pkg.go.dev).

### 17.1 gRPC dan Protocol Buffers

**gRPC** adalah framework RPC berperforma tinggi yang berjalan di atas **HTTP/2** dan memakai **Protocol Buffers (protobuf)** sebagai format data biner. Go adalah salah satu bahasa utamanya — banyak sistem cloud-native (etcd, Kubernetes API internal, banyak layanan mikro) memakainya.

| Aspek | REST/JSON | gRPC/protobuf |
|---|---|---|
| Format data | Teks (JSON) | Biner, lebih ringkas dan cepat di-encode |
| Kontrak | Opsional (OpenAPI) | **Wajib**, ditulis di berkas `.proto` |
| Transport | HTTP/1.1 atau HTTP/2 | HTTP/2 (multiplexing, streaming) |
| Streaming | Terbatas (SSE/WebSocket) | **Bawaan**: server, client, dan dua arah |
| Kode klien/server | Ditulis manual | **Digenerate** dari `.proto` |
| Dukungan browser | Langsung | Butuh gRPC-Web atau gateway |

Alur kerja gRPC:

1. Tulis kontrak di **`.proto`** (message dan service).
2. Generate kode Go dengan `protoc` (plugin `protoc-gen-go` dan `protoc-gen-go-grpc`) atau **`buf`** (tool modern yang menyederhanakan proses ini).
3. Implementasikan interface server hasil generate, dan panggil stub klien dari sisi pemanggil.

Empat jenis RPC: **unary** (satu request → satu response), **server streaming**, **client streaming**, dan **bidirectional streaming**.

Hal-hal penting dalam praktik: kirim **deadline lewat `context`** pada setiap panggilan ([Modul 13](13-sync-context-patterns.md)); kembalikan error memakai paket **`status`** dan **`codes`** (`codes.NotFound`, `codes.InvalidArgument`); pasang **interceptor** (setara middleware HTTP) untuk logging, autentikasi, dan tracing; gunakan TLS di produksi. Alternatif yang patut diketahui: **Connect-RPC** (kompatibel dengan gRPC tetapi juga berjalan di HTTP/1.1 dan browser) dan **grpc-gateway** (membuka layanan gRPC sebagai REST).

### 17.2 Akses Database dengan `database/sql`

Standard library menyediakan **`database/sql`** — antarmuka generik ke database relasional. Driver spesifik database didaftarkan lewat **blank import**.

| Konsep | Penjelasan |
|---|---|
| `sql.Open(driver, dsn)` | Membuat `*sql.DB` — **connection pool**, bukan satu koneksi. Tidak langsung menghubungi database |
| `db.PingContext(ctx)` | Memastikan koneksi benar-benar bisa dibuat |
| `QueryContext` / `QueryRowContext` / `ExecContext` | Baca banyak baris / satu baris / eksekusi tanpa hasil baris (INSERT, UPDATE, DELETE) |
| `rows.Scan(&a, &b)` | Menyalin kolom ke variabel; **wajib `defer rows.Close()`** dan cek `rows.Err()` setelah loop |
| `sql.ErrNoRows` | Dikembalikan `QueryRow(...).Scan` bila tidak ada baris |
| `sql.NullString`, `sql.NullInt64` | Menangani kolom yang bisa **NULL** |
| `db.BeginTx(ctx, nil)` | Memulai **transaksi**; selalu `defer tx.Rollback()` lalu `tx.Commit()` |
| `SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime` | Mengatur ukuran dan umur pool |

**Aturan emas keamanan: selalu pakai parameter terikat** (`$1`, `?`) — **jangan pernah** menyambung string SQL dengan input pengguna (SQL injection).

### 17.3 `pgx` — Driver PostgreSQL

**`jackc/pgx`** adalah driver dan toolkit PostgreSQL paling populer di Go. Dua cara pemakaian:

| Mode | Penjelasan |
|---|---|
| **Native** (`pgxpool`) | API khusus PostgreSQL: `pgxpool.New`, `Query`, `CollectRows`, `CopyFrom` (bulk insert), `Batch`, `LISTEN/NOTIFY`, tipe PostgreSQL kaya. Lebih cepat dan lengkap |
| **Lewat `database/sql`** | Import `github.com/jackc/pgx/v5/stdlib` sebagai driver agar kompatibel dengan tool yang mengharapkan `*sql.DB` |

Pilih **native** untuk aplikasi yang hanya memakai PostgreSQL; pilih **mode `database/sql`** bila butuh portabilitas atau library pihak ketiga yang berbasis `*sql.DB`.

### 17.4 ORM: GORM

**GORM** adalah ORM paling populer di Go: struct dipetakan ke tabel, dengan `AutoMigrate`, query berantai, relasi (*has many*, *belongs to*, *many2many*), hook, dan transaksi.

| Kelebihan | Kekurangan |
|---|---|
| Cepat untuk CRUD sederhana dan prototipe | Perilaku "ajaib" — query SQL sebenarnya tersembunyi |
| Migrasi otomatis dan relasi bawaan | Rawan masalah **N+1 query** bila lupa `Preload` |
| Ekosistem plugin dan driver luas | Overhead kinerja dibanding SQL langsung |
| | Query rumit sering lebih sulit dibaca daripada SQL biasa |

### 17.5 Memilih Pendekatan Akses Data

| Pendekatan | Karakter | Cocok untuk |
|---|---|---|
| `database/sql` / `pgx` langsung | Kontrol penuh, SQL eksplisit | Layanan dengan query penting dan performa tinggi |
| **sqlc** | Menulis SQL, tool **menghasilkan kode Go bertipe aman** | Tim yang nyaman dengan SQL dan ingin keamanan tipe saat kompilasi |
| **sqlx** | Ekstensi tipis `database/sql` (scan ke struct) | Mengurangi boilerplate tanpa ORM |
| **GORM / Ent / Bun** | ORM penuh atau query builder | CRUD cepat, model domain dengan banyak relasi |

Tidak ada pilihan tunggal yang benar; yang penting **konsisten** dan paham SQL yang dihasilkan. Apa pun pilihannya, sembunyikan di balik **interface repository** ([Modul 16](16-cli-web-development.md)).

### 17.6 Migrasi Skema

Perubahan skema harus **berversi dan dapat diulang**. Tool yang umum: **golang-migrate**, **goose**, dan **Atlas**. File migrasi (`001_buat_tabel.up.sql` / `.down.sql`) disimpan di repositori, bisa disematkan dengan `go:embed` ([Modul 14](14-standard-library.md)), dan dijalankan otomatis di pipeline deploy atau saat startup (dengan hati-hati pada banyak replika).

### 17.7 Logging Berperforma Tinggi: Zerolog dan Zap

Untuk kebanyakan aplikasi, **`log/slog`** ([Modul 14](14-standard-library.md)) sudah memadai. Zap dan Zerolog muncul lebih dulu dan unggul pada **alokasi minimal** di jalur panas (*hot path*).

| | **Zap** (`go.uber.org/zap`) | **Zerolog** (`rs/zerolog`) |
|---|---|---|
| Gaya API | `logger.Info("pesan", zap.String("k","v"))`; `SugaredLogger` untuk gaya longgar | Chaining: `log.Info().Str("k","v").Msg("pesan")` |
| Fokus | Performa dan konfigurasi fleksibel (encoder, sampling, level dinamis) | Zero-allocation, JSON murni, API ringkas |
| Output | JSON (produksi) atau konsol (dev) | JSON, mode konsol untuk dev |

Panduan: **mulai dari `slog`**. Pindah ke Zap/Zerolog bila profiling ([Modul 19](19-performance-deployment.md)) menunjukkan logging sebagai bottleneck atau tim sudah memakainya. Pola penting apa pun library-nya: **logging terstruktur** (key-value), **level** yang bisa diatur lewat konfigurasi, dan **request ID** di setiap baris log agar jejak antar-layanan bisa digabung — misalnya dikirim ke Elasticsearch (course Elasticsearch).

### 17.8 Menyambung ke Course Lain

| Teknologi | Library Go yang umum |
|---|---|
| **Redis** (course Redis) | `redis/go-redis` (`github.com/redis/go-redis/v9`) |
| **Kafka** (course Kafka) | `segmentio/kafka-go`, `twmb/franz-go`, `confluent-kafka-go` |
| **Elasticsearch** (course Elasticsearch) | `elastic/go-elasticsearch` (cocokkan versi mayor dengan klaster) |
| **Keycloak** (course Keycloak) | Verifikasi token OIDC/JWT dengan `coreos/go-oidc` dan `golang.org/x/oauth2` |
| **PostgreSQL** | `pgx`, `sqlc`, GORM seperti di atas |

Pola yang sama berlaku untuk semuanya: buat **client sekali di `main`**, teruskan lewat konstruktor, gunakan `context` untuk timeout, dan bungkus di balik interface kecil sehingga mudah diuji ([Modul 15](15-testing-benchmarking.md)).

## Praktik

### Repository dengan `database/sql` (Hanya Standard Library)

```go
// pengguna/repo.go
package pengguna

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrTidakDitemukan = errors.New("pengguna tidak ditemukan")

type Pengguna struct {
	ID    int64
	Nama  string
	Email sql.NullString
}

type Repo struct{ db *sql.DB }

func NewRepo(db *sql.DB) *Repo { return &Repo{db: db} }

func (r *Repo) Cari(ctx context.Context, id int64) (Pengguna, error) {
	var p Pengguna
	err := r.db.QueryRowContext(ctx,
		`SELECT id, nama, email FROM pengguna WHERE id = $1`, id,
	).Scan(&p.ID, &p.Nama, &p.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return Pengguna{}, ErrTidakDitemukan
	}
	if err != nil {
		return Pengguna{}, fmt.Errorf("cari pengguna %d: %w", id, err)
	}
	return p, nil
}

func (r *Repo) Daftar(ctx context.Context) ([]Pengguna, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, nama, email FROM pengguna ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hasil []Pengguna
	for rows.Next() {
		var p Pengguna
		if err := rows.Scan(&p.ID, &p.Nama, &p.Email); err != nil {
			return nil, err
		}
		hasil = append(hasil, p)
	}
	return hasil, rows.Err()
}

func (r *Repo) Buat(ctx context.Context, nama string) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO pengguna (nama) VALUES ($1) RETURNING id`, nama,
	).Scan(&id)
	return id, err
}

// Pindah dua operasi dalam satu transaksi: semua berhasil atau semua dibatalkan.
func (r *Repo) GantiNama(ctx context.Context, id int64, nama string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // aman dipanggil setelah Commit (tidak berefek)

	if _, err := tx.ExecContext(ctx, `UPDATE pengguna SET nama = $1 WHERE id = $2`, nama, id); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO audit (pengguna_id, aksi) VALUES ($1, 'ganti_nama')`, id); err != nil {
		return err
	}
	return tx.Commit()
}
```

### Menjalankan PostgreSQL dan Menyambungkannya

```bash
docker run --name pg-latihan -e POSTGRES_PASSWORD=rahasia -e POSTGRES_DB=toko \
  -p 5432:5432 -d postgres:17

docker exec -it pg-latihan psql -U postgres -d toko -c "
CREATE TABLE pengguna (id BIGSERIAL PRIMARY KEY, nama TEXT NOT NULL, email TEXT);
CREATE TABLE audit (id BIGSERIAL PRIMARY KEY, pengguna_id BIGINT, aksi TEXT);"
```

```bash
mkdir dbdemo && cd dbdemo
go mod init contoh.com/dbdemo
go get github.com/jackc/pgx/v5@latest
```

```go
// main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // mendaftarkan driver "pgx"

	"contoh.com/dbdemo/pengguna"
)

func main() {
	dsn := os.Getenv("DATABASE_URL") // postgres://postgres:rahasia@localhost:5432/toko?sslmode=disable
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal("tidak bisa terhubung: ", err)
	}

	repo := pengguna.NewRepo(db)
	id, err := repo.Buat(ctx, "Andi")
	if err != nil {
		log.Fatal(err)
	}
	p, err := repo.Cari(ctx, id)
	fmt.Println(p, err)

	_, err = repo.Cari(ctx, 9999)
	fmt.Println(err == pengguna.ErrTidakDitemukan)
}
```

Pindahkan `repo.go` ke `dbdemo/pengguna/repo.go`, lalu:

```bash
export DATABASE_URL="postgres://postgres:rahasia@localhost:5432/toko?sslmode=disable"
go mod tidy
go run .
```

### Contoh gRPC (Kerangka)

```protobuf
// proto/salam.proto
syntax = "proto3";

package salam.v1;
option go_package = "contoh.com/grpcdemo/proto;salampb";

service SalamService {
  rpc Sapa(SapaRequest) returns (SapaReply);
}

message SapaRequest { string nama = 1; }
message SapaReply   { string pesan = 1; }
```

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/salam.proto
```

```go
// server/main.go
package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "contoh.com/grpcdemo/proto"
)

type server struct {
	pb.UnimplementedSalamServiceServer
}

func (s *server) Sapa(ctx context.Context, req *pb.SapaRequest) (*pb.SapaReply, error) {
	if req.GetNama() == "" {
		return nil, status.Error(codes.InvalidArgument, "nama wajib diisi")
	}
	return &pb.SapaReply{Pesan: "Halo, " + req.GetNama() + "!"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	pb.RegisterSalamServiceServer(s, &server{})
	log.Println("gRPC berjalan di :50051")
	log.Fatal(s.Serve(lis))
}
```

```go
// klien/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "contoh.com/grpcdemo/proto"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewSalamServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Sapa(ctx, &pb.SapaRequest{Nama: "Andi"})
	fmt.Println(r.GetPesan(), err)
}
```

### Zap dan Zerolog (Perbandingan Singkat)

```go
// logging.go
package main

import (
	"os"

	"github.com/rs/zerolog"
	"go.uber.org/zap"
)

func main() {
	z, _ := zap.NewProduction()
	defer z.Sync()
	z.Info("login", zap.String("user", "andi"), zap.Int("percobaan", 2))

	l := zerolog.New(os.Stdout).With().Timestamp().Str("service", "api").Logger()
	l.Info().Str("user", "andi").Int("percobaan", 2).Msg("login")
}
```

## Latihan

1. Jalankan PostgreSQL dengan Docker, buat tabel `pengguna`, lalu implementasikan `Repo.Daftar` dan `Repo.Cari`. Uji kasus pengguna yang tidak ada dan pastikan hasilnya `ErrTidakDitemukan`, bukan error mentah dari driver.
2. Sengaja tulis query rentan SQL injection dengan menyambung string (`"... WHERE nama = '" + input + "'"`), lalu perbaiki dengan parameter `$1`. Coba input `' OR '1'='1` pada kedua versi dan jelaskan perbedaannya.
3. Tulis fungsi transfer saldo antar-rekening dalam satu **transaksi**. Simulasikan kegagalan di tengah proses dan buktikan `Rollback` membatalkan perubahan pertama.
4. Ubah `Repo` menjadi memakai interface kecil (`Querier` dengan `QueryRowContext`/`ExecContext`) agar bisa diuji dengan fake tanpa database ([Modul 15](15-testing-benchmarking.md)), lalu tulis satu uji tabel-driven.
5. Definisikan layanan gRPC dengan satu RPC *server streaming* (misalnya `HitungMundur`) di berkas `.proto`, generate kodenya, dan implementasikan server dan klien. Apa keuntungan streaming dibanding polling REST?
6. Implementasikan model yang sama dengan **GORM** (`AutoMigrate`, `Create`, `First`, `Where`) dan dengan SQL langsung. Aktifkan log SQL GORM dan identifikasi satu skenario yang menghasilkan N+1 query serta perbaikannya dengan `Preload`.
7. Diskusikan: sebuah layanan Go menerima request REST dari klien web, lalu memanggil lima layanan internal lain dan menyimpan hasil di PostgreSQL, Redis, dan Kafka. Pilih protokol komunikasi antar-layanan (REST vs gRPC), strategi akses data (SQL langsung, sqlc, atau ORM), dan pendekatan logging/tracing — beserta alasan dan trade-off-nya.

---

[← Modul 16: Membangun CLI & Web Development](16-cli-web-development.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 18: Go Toolchain, Code Quality & Security →](18-toolchain-code-quality.md)
