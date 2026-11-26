package main

import (
	"runtime"
	"testing"
	"time"
)

func TestPerbaikanTidakBocor(t *testing.T) {
	awal := runtime.NumGoroutine()
	for i := 0; i < 50; i++ {
		ambilBufer(300 * time.Millisecond)
		ambilContext(300 * time.Millisecond)
	}
	time.Sleep(600 * time.Millisecond)
	if akhir := runtime.NumGoroutine(); akhir > awal {
		t.Fatalf("bocor: awal=%d akhir=%d", awal, akhir)
	}
	t.Logf("tidak bocor: awal=%d akhir=%d", awal, runtime.NumGoroutine())
}
