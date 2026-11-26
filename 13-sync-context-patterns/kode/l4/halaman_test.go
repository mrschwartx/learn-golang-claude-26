package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAmbilHalamanCepat(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("halo"))
	}))
	defer srv.Close()
	s, err := AmbilHalaman(context.Background(), srv.URL)
	if err != nil || s != "halo" {
		t.Fatalf("s=%q err=%v", s, err)
	}
}

func TestAmbilHalamanLambat(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(5 * time.Second): // server lambat
		case <-r.Context().Done(): // klien memutus -> handler berhenti
		}
	}))
	defer srv.Close()

	mulai := time.Now()
	_, err := AmbilHalaman(context.Background(), srv.URL)
	lama := time.Since(mulai)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("seharusnya DeadlineExceeded, dapat: %v", err)
	}
	if lama < 1900*time.Millisecond || lama > 2500*time.Millisecond {
		t.Fatalf("lama %v di luar ~2 detik", lama)
	}
	t.Logf("err = %v", err)
	t.Logf("lama = %v", lama.Round(10*time.Millisecond))
}

func TestAmbilHalamanDibatalkanInduk(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer srv.Close()
	ctx, cancel := context.WithCancel(context.Background())
	go func() { time.Sleep(100 * time.Millisecond); cancel() }()
	_, err := AmbilHalaman(ctx, srv.URL)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("seharusnya Canceled, dapat: %v", err)
	}
	t.Logf("err = %v", err)
}
