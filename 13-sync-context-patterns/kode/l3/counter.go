package main

import (
	"sync"
	"sync/atomic"
)

type CounterMutex struct {
	mu sync.Mutex
	n  int64
}

func (c *CounterMutex) Add() { c.mu.Lock(); c.n++; c.mu.Unlock() }

type CounterAtomic struct{ n atomic.Int64 }

func (c *CounterAtomic) Add() { c.n.Add(1) }

func main() {}
