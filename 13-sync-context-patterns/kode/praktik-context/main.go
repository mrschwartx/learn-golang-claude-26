// context.go
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type kunci string

func kerjaLambat(ctx context.Context, d time.Duration) error {
	select {
	case <-time.After(d):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	fmt.Println("cepat:", kerjaLambat(ctx, 10*time.Millisecond))
	err := kerjaLambat(ctx, 500*time.Millisecond)
	fmt.Println("lambat:", err, errors.Is(err, context.DeadlineExceeded))

	induk, batalkan := context.WithCancel(context.Background())
	anak, cancelAnak := context.WithTimeout(induk, time.Hour)
	defer cancelAnak()
	batalkan() // membatalkan induk -> anak ikut batal
	<-anak.Done()
	fmt.Println("anak:", anak.Err())

	ctxNilai := context.WithValue(context.Background(), kunci("reqID"), "abc-123")
	fmt.Println(ctxNilai.Value(kunci("reqID")))
}
