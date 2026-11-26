// l5: slog dengan TextHandler, logger.With, dan level dari LOG_LEVEL.
package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// LevelDariEnv membaca LOG_LEVEL (debug|info|warn|error); default info; nilai salah -> info + peringatan.
func LevelDariEnv() (slog.Level, error) {
	v, ok := os.LookupEnv("LOG_LEVEL")
	if !ok || v == "" {
		return slog.LevelInfo, nil
	}
	var l slog.Level
	if err := l.UnmarshalText([]byte(strings.TrimSpace(v))); err != nil {
		return slog.LevelInfo, fmt.Errorf("LOG_LEVEL=%q tidak valid: %w", v, err)
	}
	return l, nil
}

func main() {
	level, lerr := LevelDariEnv()
	opsi := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{} // deterministik untuk demo
			}
			return a
		},
	}
	var h slog.Handler = slog.NewTextHandler(os.Stdout, opsi)
	if os.Getenv("LOG_FORMAT") == "json" {
		h = slog.NewJSONHandler(os.Stdout, opsi)
	}
	logger := slog.New(h)
	if lerr != nil {
		logger.Warn("memakai level default", "err", lerr)
	}
	req := logger.With("request_id", "abc-123")
	req.Debug("detail debug", "langkah", 1)
	req.Info("login", "user", "andi", "ip", "10.0.0.7")
	req.Warn("lambat", "ms", 1500)
	req.Error("gagal simpan", "err", fmt.Errorf("koneksi ditolak"))
}
