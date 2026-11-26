// Package main berisi Cache aman-concurrent dengan sync.RWMutex.
package main

import "sync"

type Cache struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewCache() *Cache { return &Cache{data: make(map[string]string)} }

// Get memakai kunci baca: banyak pembaca boleh bersamaan.
func (c *Cache) Get(k string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.data[k]
	return v, ok
}

// Set memakai kunci tulis: eksklusif.
func (c *Cache) Set(k, v string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[k] = v
}

func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.data)
}

func main() {}
