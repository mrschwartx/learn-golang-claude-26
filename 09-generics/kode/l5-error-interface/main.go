//go:build ignore

package main

type Pengubah interface {
	Transform[U any](f func(int) U) []U // method generik di dalam interface
}

func main() {}
