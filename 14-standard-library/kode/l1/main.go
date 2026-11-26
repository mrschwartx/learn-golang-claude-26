// l1: hitung baris, kata, rune, dan 5 kata terbanyak dengan bufio.Scanner.
// Pemakaian: l1 [-buat N] berkas.txt   (-buat N membuat berkas uji N baris)
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Hasil struct {
	Baris, Kata, Rune int
	Frek              map[string]int
}

func Hitung(f *os.File) (Hasil, error) {
	h := Hasil{Frek: map[string]int{}}
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 16*1024*1024) // izinkan baris sampai 16 MB
	for sc.Scan() {
		line := sc.Text()
		h.Baris++
		h.Rune += utf8.RuneCountInString(line) + 1 // +1 untuk '\n'
		for _, k := range strings.FieldsFunc(line, func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) }) {
			h.Kata++
			h.Frek[strings.ToLower(k)]++
		}
	}
	return h, sc.Err()
}

func Teratas(frek map[string]int, n int) []string {
	type kv struct {
		k string
		v int
	}
	var s []kv
	for k, v := range frek {
		s = append(s, kv{k, v})
	}
	sort.Slice(s, func(i, j int) bool {
		if s[i].v != s[j].v {
			return s[i].v > s[j].v
		}
		return s[i].k < s[j].k
	})
	var out []string
	for i := 0; i < n && i < len(s); i++ {
		k := s[i].k
		if len(k) > 20 {
			k = k[:20] + "..."
		}
		out = append(out, fmt.Sprintf("%s=%d", k, s[i].v))
	}
	return out
}

// rssPuncak membaca VmHWM (memori fisik puncak) dari /proc (khusus Linux).
func rssPuncak() string {
	b, err := os.ReadFile("/proc/self/status")
	if err != nil {
		return "n/a"
	}
	for _, l := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(l, "VmHWM:") {
			return strings.TrimSpace(strings.TrimPrefix(l, "VmHWM:"))
		}
	}
	return "n/a"
}

func main() {
	buat := flag.Int("buat", 0, "buat berkas uji dengan N baris lalu keluar")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "butuh 1 argumen: nama berkas")
		os.Exit(2)
	}
	path := flag.Arg(0)
	if *buat > 0 {
		f, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		w := bufio.NewWriter(f)
		for i := 0; i < *buat; i++ {
			fmt.Fprintf(w, "saya belajar go setiap hari baris %d kopi kopi go\n", i%50)
		}
		w.Flush()
		f.Close()
		return
	}
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()
	h, err := Hitung(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	fmt.Printf("baris=%d kata=%d rune=%d\n", h.Baris, h.Kata, h.Rune)
	fmt.Println("top5:", Teratas(h.Frek, 5))
	fmt.Printf("heap yang diminta dari OS: %.1f MiB, RSS puncak: %s\n", float64(ms.HeapSys)/1048576, rssPuncak())
}
