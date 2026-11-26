package main

import (
	"testing"
	"time"
)

func TestHariKe(t *testing.T) {
	d := func(y int, m time.Month, h, jam int) time.Time { return time.Date(y, m, h, jam, 0, 0, 0, time.UTC) }
	if got := HariKe(d(2026, 11, 26, 23), d(2027, 8, 17, 0)); got != 264 {
		t.Errorf("got %d", got) // 26 Nov 2026 -> 17 Agu 2027
	}
	if got := HariKe(d(2027, 8, 17, 23), d(2027, 8, 17, 0)); got != 0 {
		t.Errorf("hari sama: %d", got)
	}
	if got := HariKe(d(2027, 8, 18, 1), d(2027, 8, 17, 0)); got != -1 {
		t.Errorf("lewat: %d", got)
	}
}
