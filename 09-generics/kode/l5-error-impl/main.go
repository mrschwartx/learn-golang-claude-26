//go:build ignore

package main

type Stack[T any] struct{ items []T }

// method generik pada tipe
func (s *Stack[T]) Transform[U any](f func(T) U) []U { return nil }

func main() {}
