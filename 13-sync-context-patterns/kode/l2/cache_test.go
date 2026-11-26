package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestCacheConcurrent(t *testing.T) {
	c := NewCache()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			k := fmt.Sprintf("k%d", i%10)
			c.Set(k, fmt.Sprint(i))
			c.Get(k)
			c.Len()
		}()
	}
	wg.Wait()
	if c.Len() != 10 {
		t.Fatalf("Len = %d, seharusnya 10", c.Len())
	}
}
