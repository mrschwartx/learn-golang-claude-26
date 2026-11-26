// nilgotcha.go
package main

import "fmt"

type ErrKustom struct{}

func (*ErrKustom) Error() string { return "error kustom" }

func salah() error {
	var e *ErrKustom // nil
	return e
}

func benar() error {
	return nil
}

func main() {
	fmt.Println("salah() == nil:", salah() == nil)
	fmt.Println("benar() == nil:", benar() == nil)
}
