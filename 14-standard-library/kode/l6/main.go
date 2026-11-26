// l6: validasi + ekstraksi email dengan named group, dan bukti RE2 linear (anti-ReDoS).
package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Pola praktis (bukan RFC 5322 penuh): bagian lokal, '@', domain berlabel dengan TLD >= 2 huruf.
var reEmail = regexp.MustCompile(`^(?P<nama>[A-Za-z0-9._%+-]+)@(?P<domain>(?:[A-Za-z0-9](?:[A-Za-z0-9-]*[A-Za-z0-9])?\.)+[A-Za-z]{2,})$`)

func Pecah(s string) (nama, domain string, ok bool) {
	m := reEmail.FindStringSubmatch(s)
	if m == nil {
		return "", "", false
	}
	return m[reEmail.SubexpIndex("nama")], m[reEmail.SubexpIndex("domain")], true
}

func main() {
	valid := []string{"andi@contoh.com", "budi.santoso+tag@mail.example.co.id", "x_y-z@sub.domain.org", "A1@b2.io", "siti%x@kampus.ac.id"}
	tidak := []string{"andi@", "@contoh.com", "andi@contoh", "andi contoh@x.com", "andi@@x.com", "andi@-x.com", "andi@x..com", "andi@x.c"}
	fmt.Println("# valid")
	for _, s := range valid {
		n, d, ok := Pecah(s)
		fmt.Printf("%-38s ok=%-5v nama=%q domain=%q\n", s, ok, n, d)
	}
	fmt.Println("# tidak valid")
	for _, s := range tidak {
		_, _, ok := Pecah(s)
		fmt.Printf("%-38s ok=%v\n", s, ok)
	}

	// ReDoS: pola klasik (a+)+$ pada "aaa...ab". Di mesin backtracking ini eksponensial.
	fmt.Println("# waktu (a+)+$ pada n huruf 'a' + 'b' (RE2 linear)")
	re := regexp.MustCompile(`^(a+)+$`)
	for _, n := range []int{1_000, 10_000, 100_000, 1_000_000} {
		in := strings.Repeat("a", n) + "b"
		t := time.Now()
		ok := re.MatchString(in)
		fmt.Printf("n=%-8d cocok=%v waktu=%v\n", n, ok, time.Since(t).Round(time.Microsecond))
	}

	// Fitur yang dikorbankan
	for _, p := range []string{`(a)\1`, `foo(?=bar)`, `(?<!x)y`} {
		_, err := regexp.Compile(p)
		fmt.Printf("%-12s -> %v\n", p, err)
	}
}
