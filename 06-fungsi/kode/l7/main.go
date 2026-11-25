// Latihan 7: struct konfigurasi vs functional options.
package main

import (
	"fmt"
	"time"
)

// ---------- Pola A: struct konfigurasi ----------
type ConfigA struct {
	Addr    string
	Timeout time.Duration
	MaxConn int
}

func NewServerA(c ConfigA) *ServerA {
	if c.Addr == "" {
		c.Addr = ":8080"
	}
	if c.Timeout == 0 { // masalah: 0 = "tidak diisi" atau "tanpa timeout"?
		c.Timeout = 30 * time.Second
	}
	if c.MaxConn == 0 {
		c.MaxConn = 100
	}
	return &ServerA{c}
}

type ServerA struct{ cfg ConfigA }

// ---------- Pola B: functional options ----------
type ServerB struct {
	addr    string
	timeout time.Duration
	maxConn int
}

type Option func(*ServerB)

func WithAddr(a string) Option           { return func(s *ServerB) { s.addr = a } }
func WithTimeout(d time.Duration) Option { return func(s *ServerB) { s.timeout = d } }
func WithMaxConn(n int) Option {
	return func(s *ServerB) {
		if n > 0 {
			s.maxConn = n
		}
	}
}

func NewServerB(opts ...Option) *ServerB {
	s := &ServerB{addr: ":8080", timeout: 30 * time.Second, maxConn: 100} // default eksplisit
	for _, o := range opts {
		o(s)
	}
	return s
}

func main() {
	a := NewServerA(ConfigA{Addr: ":9000"})
	fmt.Printf("A default : %+v\n", a.cfg)

	// A: tidak bisa menyatakan "tanpa timeout" (0) karena 0 dianggap belum diisi
	a2 := NewServerA(ConfigA{Timeout: 0})
	fmt.Println("A timeout 0 ->", a2.cfg.Timeout, "(ingin 'tanpa timeout' tapi jadi 30s)")

	b := NewServerB(WithAddr(":9000"), WithMaxConn(500))
	fmt.Printf("B         : %+v\n", *b)

	b2 := NewServerB(WithTimeout(0))
	fmt.Println("B timeout 0 ->", b2.timeout, "(eksplisit: 0 benar-benar dipilih)")
}
