// Latihan 6: defer di dalam loop vs fungsi bantu.
// Batas file descriptor diturunkan agar masalahnya tampak jelas.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

func jumlahFDTerbuka() int {
	e, _ := os.ReadDir("/proc/self/fd")
	return len(e)
}

// VERSI 1: defer langsung di loop. Semua file tetap terbuka sampai fungsi selesai.
func BacaSemuaSalah(paths []string, puncak *int) error {
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			return fmt.Errorf("membuka %s: %w", filepath.Base(p), err)
		}
		defer f.Close() // baru jalan saat BacaSemuaSalah return!
		if n := jumlahFDTerbuka(); n > *puncak {
			*puncak = n
		}
	}
	return nil
}

// VERSI 2: isi loop dipindah ke fungsi bantu; defer berjalan di akhir tiap file.
func bacaSatu(p string) error {
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Stat()
	return err
}

func BacaSemua(paths []string, puncak *int) error {
	for _, p := range paths {
		if err := bacaSatu(p); err != nil {
			return fmt.Errorf("membaca %s: %w", filepath.Base(p), err)
		}
		if n := jumlahFDTerbuka(); n > *puncak {
			*puncak = n
		}
	}
	return nil
}

func main() {
	dir, _ := os.MkdirTemp("", "modul6")
	defer os.RemoveAll(dir)

	const jumlah = 1000
	paths := make([]string, jumlah)
	for i := range paths {
		paths[i] = filepath.Join(dir, fmt.Sprintf("f%04d.txt", i))
		os.WriteFile(paths[i], []byte("x"), 0o644)
	}

	// turunkan batas FD ke 256 untuk simulasi (soft limit)
	var rl syscall.Rlimit
	syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rl)
	rl.Cur = 256
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rl); err != nil {
		fmt.Println("setrlimit gagal:", err)
	}
	fmt.Println("batas file descriptor:", rl.Cur, "| berkas:", jumlah)

	var p1 int
	err1 := BacaSemuaSalah(paths, &p1)
	fmt.Printf("\nVERSI 1 (defer di loop): puncak FD terbuka = %d\n  error: %v\n", p1, err1)

	var p2 int
	err2 := BacaSemua(paths, &p2)
	fmt.Printf("\nVERSI 2 (fungsi bantu) : puncak FD terbuka = %d\n  error: %v\n", p2, err2)
}
