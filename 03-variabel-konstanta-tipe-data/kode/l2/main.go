// Latihan 2: shadowing err dengan := di scope dalam.
package main

import (
	"errors"
	"fmt"
)

func operasi() (int, error) { return 0, errors.New("gagal") }

func salah() error {
	var err error
	if true {
		_, err := operasi() // BARU: menutupi err luar
		_ = err
	}
	return err // tetap nil!
}

func benar() error {
	var err error
	if true {
		_, err = operasi() // "=" memakai err luar
	}
	return err
}

func benar2() error {
	if true {
		if _, err := operasi(); err != nil { // atau tangani langsung
			return err
		}
	}
	return nil
}

func main() {
	fmt.Println("salah :", salah())
	fmt.Println("benar :", benar())
	fmt.Println("benar2:", benar2())
}
