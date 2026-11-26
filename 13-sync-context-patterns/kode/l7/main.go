// l7: server dengan handler lambat (2 detik) untuk membandingkan shutdown
// tanpa graceful vs graceful saat menerima SIGTERM.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	graceful := flag.Bool("graceful", false, "tangani SIGTERM dengan Shutdown")
	addr := flag.String("addr", "127.0.0.1:8077", "alamat")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/lambat", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(os.Stderr, "[server] request diterima, memproses 2 detik")
		time.Sleep(2 * time.Second)
		fmt.Fprintln(w, "selesai")
		fmt.Fprintln(os.Stderr, "[server] request selesai dilayani")
	})
	srv := &http.Server{Addr: *addr, Handler: mux}

	if !*graceful {
		// tanpa penanganan sinyal: SIGTERM mematikan proses seketika
		fmt.Fprintln(os.Stderr, "[server] mode TIDAK graceful")
		_ = srv.ListenAndServe()
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintln(os.Stderr, "[server] error:", err)
		}
	}()
	fmt.Fprintln(os.Stderr, "[server] mode graceful")
	<-ctx.Done()
	fmt.Fprintln(os.Stderr, "[server] sinyal diterima, Shutdown (maks 5 detik)")
	sctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(sctx); err != nil {
		fmt.Fprintln(os.Stderr, "[server] shutdown gagal:", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "[server] keluar bersih")
}
