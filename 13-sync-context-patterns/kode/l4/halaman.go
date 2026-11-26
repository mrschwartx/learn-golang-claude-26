package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AmbilHalaman mengunduh url dengan batas waktu 2 detik.
func AmbilHalaman(ctx context.Context, url string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("buat request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ambil %s: %w", url, err) // %w: errors.Is tetap bekerja
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("baca body: %w", err)
	}
	return string(b), nil
}

func main() {}
