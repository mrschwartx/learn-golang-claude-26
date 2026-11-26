package main

import (
	"sync"
	"testing"
)

func BenchmarkMutexSerial(b *testing.B) {
	var c CounterMutex
	for i := 0; i < b.N; i++ {
		c.Add()
	}
}

func BenchmarkAtomicSerial(b *testing.B) {
	var c CounterAtomic
	for i := 0; i < b.N; i++ {
		c.Add()
	}
}

func BenchmarkMutexParallel(b *testing.B) {
	var c CounterMutex
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Add()
		}
	})
}

func BenchmarkAtomicParallel(b *testing.B) {
	var c CounterAtomic
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Add()
		}
	})
}

func TestKeduanyaBenar(t *testing.T) {
	var cm CounterMutex
	var ca CounterAtomic
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); cm.Add(); ca.Add() }()
	}
	wg.Wait()
	if cm.n != 1000 || ca.n.Load() != 1000 {
		t.Fatalf("mutex=%d atomic=%d", cm.n, ca.n.Load())
	}
}
