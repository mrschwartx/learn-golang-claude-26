// log_regex.go
package main

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"
)

var reTanggal = regexp.MustCompile(`(?P<tahun>\d{4})-(?P<bulan>\d{2})-(?P<hari>\d{2})`)

func main() {
	opsi := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{} // buang waktu agar keluaran deterministik
			}
			return a
		},
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opsi))
	logger.Info("server mulai", "port", 8080)
	req := logger.With("request_id", "abc-123")
	req.Warn("lambat", slog.Int("ms", 1500), slog.Group("user", "id", 7, "role", "admin"))
	req.Error("gagal simpan", "err", fmt.Errorf("koneksi ditolak"))

	teks := "Rilis pada 2026-08-19 dan 2027-02-01."
	fmt.Println(reTanggal.MatchString(teks))
	fmt.Println(reTanggal.FindAllString(teks, -1))
	m := reTanggal.FindStringSubmatch(teks)
	fmt.Println(m[reTanggal.SubexpIndex("tahun")], m[reTanggal.SubexpIndex("bulan")])
	fmt.Println(reTanggal.ReplaceAllString(teks, "${hari}/${bulan}/${tahun}"))
}
