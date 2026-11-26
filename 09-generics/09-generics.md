# Modul 9: Generics & Iterators

[← Modul 8: Methods & Interfaces](08-methods-interfaces.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 10: Error Handling →](10-error-handling.md)

**Isi modul ini**

- [Teori](#teori)
  - [9.1 Mengapa Generics?](#91-mengapa-generics)
  - [9.2 Generic Functions](#92-generic-functions)
  - [9.3 Type Constraints](#93-type-constraints)
  - [9.4 Generic Types](#94-generic-types)
  - [9.5 Generic Methods (Go 1.27)](#95-generic-methods-go-127)
  - [9.6 Generic Interfaces dan Constraint Interfaces](#96-generic-interfaces-dan-constraint-interfaces)
  - [9.7 Type Inference](#97-type-inference)
  - [9.8 Iterators (`iter`, Go 1.23+)](#98-iterators-iter-go-123)
  - [9.9 Kapan (Tidak) Memakai Generics](#99-kapan-tidak-memakai-generics)
- [Praktik](#praktik)
  - [Fungsi Generik dan Constraint](#fungsi-generik-dan-constraint)
  - [Tipe Generik dan Generic Method (Go 1.27)](#tipe-generik-dan-generic-method-go-127)
  - [Iterator](#iterator)
- [Latihan](#latihan)

---


## Teori

### 9.1 Mengapa Generics?

Sebelum Go 1.18 (Maret 2022), menulis fungsi yang bekerja pada banyak tipe hanya punya dua pilihan buruk:

1. **Duplikasi kode** — `SumInt`, `SumFloat64`, `SumInt64`, dst. dengan isi hampir identik.
2. **`interface{}`** — satu fungsi untuk semua tipe, tetapi **keamanan tipe hilang** dan butuh type assertion serta bisa panic saat runtime.

**Generics** (*type parameters*) memecahkan masalah ini: satu fungsi atau tipe ditulis sekali, dan **tipe konkretnya ditentukan saat dipakai**, dengan tetap diperiksa compiler.

### 9.2 Generic Functions

Type parameter ditulis di dalam **kurung siku** setelah nama fungsi, masing-masing dengan **constraint**:

```go
func Cari[T comparable](daftar []T, target T) int {
    for i, v := range daftar {
        if v == target {
            return i
        }
    }
    return -1
}

Cari([]string{"a", "b"}, "b")  // 1
Cari([]int{10, 20}, 20)        // 1
```

Untuk transformasi koleksi, biasanya dibutuhkan lebih dari satu type parameter:

```go
func Map[T, U any](xs []T, f func(T) U) []U { /* ... */ }
```

### 9.3 Type Constraints

**Constraint** menentukan operasi apa yang boleh dilakukan pada type parameter. Constraint adalah **interface** yang mendeskripsikan sekumpulan tipe (*type set*).

| Constraint | Arti |
|---|---|
| `any` | Tipe apa pun (tanpa operasi khusus selain assign dan kirim) |
| `comparable` | Tipe yang mendukung `==` dan `!=` (cocok untuk key map) |
| `cmp.Ordered` | Tipe yang mendukung `<`, `<=`, `>`, `>=` (angka dan string) |
| Union: `~int \| ~float64` | Hanya tipe-tipe yang terdaftar |
| Interface dengan method | Tipe yang memiliki method tersebut |

Tanda **`~`** (*tilde*) berarti "tipe ini **beserta semua tipe turunan** yang underlying type-nya sama". Tanpa `~`, `type Rupiah int` **tidak** memenuhi constraint `int`; dengan `~int`, ia memenuhi.

```go
type Angka interface {
    ~int | ~int64 | ~float64
}

func Jumlah[T Angka](xs []T) T {
    var total T
    for _, x := range xs {
        total += x
    }
    return total
}
```

Paket standar `cmp` menyediakan `cmp.Ordered`, `cmp.Compare`, dan `cmp.Or`; sebelum Go 1.21, constraint semacam ini berasal dari paket eksperimental `golang.org/x/exp/constraints`.

### 9.4 Generic Types

Struct, slice, map, dan tipe lain juga dapat bersifat generik:

```go
type Stack[T any] struct{ items []T }

func (s *Stack[T]) Push(x T) { s.items = append(s.items, x) }

func (s *Stack[T]) Pop() (T, bool) {
    var nol T
    if len(s.items) == 0 {
        return nol, false
    }
    x := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return x, true
}

type Pair[K comparable, V any] struct {
    Key K
    Val V
}
```

Pada method, receiver menyebut type parameter (`Stack[T]`). `var nol T` adalah cara idiomatik mendapatkan **zero value** dari `T`. Sejak **Go 1.24**, *type alias* juga boleh generik (`type Daftar[T any] = []T`).

### 9.5 Generic Methods (Go 1.27)

Sebelum Go 1.27, sebuah method **tidak boleh** mendeklarasikan type parameter sendiri — hanya yang dibawa dari tipenya. **Mulai Go 1.27, method boleh bersifat generik**:

```go
func (s *Stack[T]) Transform[U any](f func(T) U) []U { /* ... */ }
```

Dua batasan penting yang tetap berlaku:

- Sebuah interface **tidak boleh** mendeklarasikan method generik (compiler menolak: *interface method must have no type parameters*).
- Method generik **tidak dapat** dipakai untuk memenuhi method non-generik di interface — pengecekan interface hanya mengenali method dengan signature konkret.

Jadi pakai method generik untuk mengelompokkan operasi pada sebuah tipe, bukan untuk membuat abstraksi lewat interface.

### 9.6 Generic Interfaces dan Constraint Interfaces

Interface yang memuat **union tipe** atau `comparable` hanya dapat dipakai sebagai **constraint**, bukan sebagai tipe variabel biasa. Interface biasa (kumpulan method) tetap bisa dipakai di kedua peran, dan interface juga bisa dibuat generik:

```go
type Repository[T any, ID comparable] interface {
    Get(id ID) (T, error)
    Save(item T) error
}
```

Constraint juga boleh merujuk dirinya sendiri. Mulai **Go 1.26**, restriksi "generic type tidak boleh merujuk dirinya dalam daftar type parameter" dicabut, sehingga pola seperti `type Adder[A Adder[A]] interface { Add(A) A }` menjadi valid.

### 9.7 Type Inference

Compiler biasanya **menyimpulkan type parameter dari argumen**, sehingga Anda tidak perlu menulisnya:

```go
Jumlah([]int{1, 2})            // T disimpulkan = int
Map([]int{1, 2}, strconv.Itoa) // T = int, U = string
```

Instansiasi eksplisit dibutuhkan jika tipe tidak bisa disimpulkan — misalnya fungsi tanpa argumen berparameter tipe: `ptr := Baru[User]()` — atau saat Anda ingin tipe tertentu: `Jumlah[float64](nil)`. Go terus memperluas inferensi pada rilis-rilis terbaru; Go 1.27, misalnya, menggeneralisasi inferensi fungsi generik di berbagai konteks assignment.

### 9.8 Iterators (`iter`, Go 1.23+)

Go 1.23 menambahkan **range-over-func**: sebuah fungsi dapat menjadi sumber iterasi untuk `for range`. Tipe standarnya ada di paket `iter`:

| Tipe | Bentuk | Dipakai untuk |
|---|---|---|
| `iter.Seq[V]` | `func(yield func(V) bool)` | Iterasi satu nilai per langkah |
| `iter.Seq2[K, V]` | `func(yield func(K, V) bool)` | Iterasi pasangan (index/key + value) |

```go
func Hitung(n int) iter.Seq[int] {
    return func(yield func(int) bool) {
        for i := 0; i < n; i++ {
            if !yield(i) {   // false = konsumen berhenti (break)
                return
            }
        }
    }
}

for v := range Hitung(5) { fmt.Println(v) }
```

Fungsi `yield` mengembalikan `false` ketika pemanggil melakukan `break`, sehingga iterator harus **berhenti dan membersihkan resource**. Iterator bersifat **lazy** — nilai dihasilkan satu per satu saat diminta, tanpa membuat slice perantara.

Standard library kini menyediakan banyak fungsi berbasis iterator:

| Fungsi | Fungsi |
|---|---|
| `slices.All(s)` / `slices.Values(s)` | Iterator dari slice |
| `maps.Keys(m)` / `maps.Values(m)` | Iterator key/value map |
| `slices.Collect(seq)` | Iterator → slice |
| `slices.Sorted(seq)` | Iterator → slice terurut |
| `strings.SplitSeq(s, sep)`, `strings.Lines(s)` | Iterasi string tanpa alokasi slice (Go 1.24) |
| `iter.Pull(seq)` | Mengubah iterator "push" menjadi gaya "pull" dengan `next()` dan `stop()` |

### 9.9 Kapan (Tidak) Memakai Generics

| Gunakan generics untuk... | Hindari jika... |
|---|---|
| Fungsi koleksi (`Map`, `Filter`, `Keys`) | Satu interface biasa sudah cukup — misalnya `io.Reader` |
| Struktur data umum (`Stack`, `Set`, `Cache`, `Queue`) | Hanya ada **satu** tipe konkret yang dipakai |
| Menghindari duplikasi kode angka/ordering | Implementasi berbeda **per tipe** (pakai interface/method) |
| Kode yang butuh keamanan tipe tanpa `any` | Demi "fleksibilitas masa depan" yang belum dibutuhkan |

Aturan praktis tim Go: **mulai dari kode konkret**; jadikan generik hanya ketika Anda mulai menulis fungsi yang sama untuk tipe kedua. Banyak kebutuhan sudah dipenuhi paket standar `slices`, `maps`, dan `cmp` ([Modul 4](04-composite-types.md)); `sync/atomic.Pointer[T]` dan `sync.OnceValue` juga contoh API generik di standard library ([Modul 13](13-sync-context-patterns.md)).

## Praktik

### Fungsi Generik dan Constraint

```go
// generik.go
package main

import (
	"cmp"
	"fmt"
	"strconv"
)

type Angka interface {
	~int | ~int64 | ~float64
}

type Rupiah int

func Jumlah[T Angka](xs []T) T {
	var total T
	for _, x := range xs {
		total += x
	}
	return total
}

func Map[T, U any](xs []T, f func(T) U) []U {
	out := make([]U, 0, len(xs))
	for _, x := range xs {
		out = append(out, f(x))
	}
	return out
}

func Filter[T any](xs []T, ok func(T) bool) []T {
	var out []T
	for _, x := range xs {
		if ok(x) {
			out = append(out, x)
		}
	}
	return out
}

func Reduce[T, A any](xs []T, awal A, f func(A, T) A) A {
	acc := awal
	for _, x := range xs {
		acc = f(acc, x)
	}
	return acc
}

func Maks[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Cari[T comparable](xs []T, target T) int {
	for i, x := range xs {
		if x == target {
			return i
		}
	}
	return -1
}

func main() {
	fmt.Println(Jumlah([]int{1, 2, 3}), Jumlah([]float64{1.5, 2.5}), Jumlah([]Rupiah{1000, 500}))
	fmt.Println(Map([]int{1, 2, 3}, strconv.Itoa))
	genap := Filter([]int{1, 2, 3, 4, 5, 6}, func(n int) bool { return n%2 == 0 })
	fmt.Println(genap, Reduce(genap, 0, func(a, n int) int { return a + n*n }))
	fmt.Println(Maks("apel", "jeruk"), Maks(3, 7), Maks(2.5, 1.5))
	fmt.Println(Cari([]string{"a", "b", "c"}, "c"), Cari([]int{1, 2}, 9))
}
```

### Tipe Generik dan Generic Method (Go 1.27)

```go
// stack.go
package main

import "fmt"

type Stack[T any] struct{ items []T }

func (s *Stack[T]) Push(x T) { s.items = append(s.items, x) }

func (s *Stack[T]) Pop() (T, bool) {
	var nol T
	if len(s.items) == 0 {
		return nol, false
	}
	x := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return x, true
}

// Generic method: butuh Go 1.27+
func (s *Stack[T]) Transform[U any](f func(T) U) []U {
	out := make([]U, 0, len(s.items))
	for _, x := range s.items {
		out = append(out, f(x))
	}
	return out
}

type Pair[K comparable, V any] struct {
	Key K
	Val V
}

type Cache[K comparable, V any] struct {
	data map[K]V
}

func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{data: make(map[K]V)}
}

func (c *Cache[K, V]) Set(k K, v V) { c.data[k] = v }

func (c *Cache[K, V]) Get(k K) (V, bool) {
	v, ok := c.data[k]
	return v, ok
}

func main() {
	s := &Stack[string]{}
	s.Push("a")
	s.Push("bb")
	fmt.Println(s.Transform(func(x string) int { return len(x) }))
	v, ok := s.Pop()
	fmt.Println(v, ok)

	c := NewCache[string, Pair[int, bool]]()
	c.Set("x", Pair[int, bool]{1, true})
	fmt.Println(c.Get("x"))
	fmt.Println(c.Get("y"))
}
```

### Iterator

```go
// iterator.go
package main

import (
	"fmt"
	"iter"
	"maps"
	"slices"
)

func Hitung(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; i < n; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

func Filter[T any](seq iter.Seq[T], ok func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range seq {
			if ok(v) && !yield(v) {
				return
			}
		}
	}
}

func Fibonacci() iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		a, b := 0, 1
		for i := 0; ; i++ {
			if !yield(i, a) {
				return
			}
			a, b = b, a+b
		}
	}
}

func main() {
	for v := range Hitung(10) {
		if v == 4 {
			break
		}
		fmt.Print(v, " ")
	}
	fmt.Println()

	genap := Filter(Hitung(10), func(n int) bool { return n%2 == 0 })
	fmt.Println(slices.Collect(genap))

	for i, f := range Fibonacci() {
		if i == 10 {
			break
		}
		fmt.Print(f, " ")
	}
	fmt.Println()

	m := map[string]int{"b": 2, "a": 1, "c": 3}
	fmt.Println(slices.Sorted(maps.Keys(m)))

	next, stop := iter.Pull(Hitung(3))
	defer stop()
	for {
		v, ok := next()
		if !ok {
			break
		}
		fmt.Print("pull:", v, " ")
	}
	fmt.Println()
}
```

## Latihan

1. Tulis `Keys[K comparable, V any](m map[K]V) []K` dan `Values` untuk map apa pun, lalu bandingkan dengan `maps.Keys` + `slices.Collect`. Mengapa urutan hasilnya tidak terjamin, dan bagaimana membuatnya deterministik?
2. Definisikan constraint `Bilangan` (`~int | ~float64`) dan tulis `Rata[T Bilangan](xs []T) float64`. Uji dengan `type Skor int` dan jelaskan peran tanda `~`.
3. Implementasikan `Set[T comparable]` dengan method `Add`, `Has`, `Remove`, dan `Len`, berbasis `map[T]struct{}`. Mengapa `struct{}` dipilih sebagai value?
4. Implementasikan `Queue[T any]` dengan `Enqueue` dan `Dequeue() (T, bool)`. Mengapa `Dequeue` mengembalikan `bool` alih-alih panic atau mengembalikan zero value saja?
5. Coba mendeklarasikan method generik di dalam sebuah interface dan mengisi interface dengan tipe yang memiliki method generik. Catat pesan error compiler, lalu jelaskan batasan Go 1.27 yang Anda amati.
6. Tulis iterator `iter.Seq[int]` untuk deret bilangan prima tak hingga dan ambil 10 bilangan pertama dengan `break`. Jelaskan peran nilai kembalian `yield` dan mengapa iterator harus berhenti saat `yield` mengembalikan `false`.
7. Diskusikan: sebuah tim ingin menulis `func Proses[T any](x T)` yang di dalamnya berisi `switch any(x).(type)` untuk menangani beberapa tipe. Mengapa ini pertanda generics dipakai secara keliru, dan alternatif apa yang lebih idiomatik (interface, atau fungsi terpisah per tipe)?

---

[← Modul 8: Methods & Interfaces](08-methods-interfaces.md) · [Daftar Isi](README.md#daftar-isi) · [Modul 10: Error Handling →](10-error-handling.md)
