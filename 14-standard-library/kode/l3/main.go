// l3: selisih hari ke 17 Agustus tahun depan, tiga layout, dan bukti pentingnya zona waktu.
package main

import (
	"fmt"
	"time"
)

// HariKe menghitung selisih hari kalender antara dua tanggal (mengabaikan jam),
// dengan membandingkan tengah malam UTC dari tanggal kalender masing-masing.
func HariKe(dari, ke time.Time) int {
	a := time.Date(dari.Year(), dari.Month(), dari.Day(), 0, 0, 0, 0, time.UTC)
	b := time.Date(ke.Year(), ke.Month(), ke.Day(), 0, 0, 0, 0, time.UTC)
	return int(b.Sub(a).Hours() / 24)
}

func main() {
	now := time.Now()
	target := time.Date(now.Year()+1, time.August, 17, 0, 0, 0, 0, now.Location())
	fmt.Println("hari ini :", now.Format("2006-01-02"), "zona", now.Location())
	fmt.Println("selisih  :", HariKe(now, target), "hari sampai", target.Format("2006-01-02"))
	fmt.Println("layout 1 :", target.Format("2006-01-02"))
	fmt.Println("layout 2 :", target.Format("02/01/2006"))
	fmt.Println("layout 3 :", target.Format("Monday, 2 January 2006"), "|", target.Format(time.RFC3339))

	// Jebakan: Sub langsung bergantung pada jam saat ini.
	fmt.Printf("Sub langsung: %.2f hari (pecahan karena jam %s)\n", target.Sub(now).Hours()/24, now.Format("15:04"))

	// Jebakan format: "YYYY-MM-DD" bukan layout, tercetak apa adanya.
	fmt.Println("salah    :", target.Format("YYYY-MM-DD"))
	// == vs Equal
	u := time.Date(2027, 8, 17, 7, 0, 0, 0, time.UTC)
	if jkt, err := time.LoadLocation("Asia/Jakarta"); err == nil {
		w := u.In(jkt)
		fmt.Println("instan sama, == :", u == w, "| Equal:", u.Equal(w), "|", w.Format(time.RFC3339))
	} else {
		fmt.Println("tzdata tidak tersedia:", err)
	}
}
