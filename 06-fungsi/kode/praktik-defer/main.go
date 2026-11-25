// defer.go
package main

import (
	"errors"
	"fmt"
)

type Titik struct{ X int }

func ubahNilai(t Titik)    { t.X = 100 }
func ubahPointer(t *Titik) { t.X = 100 }
func ubahIsi(s []int)      { s[0] = 100 }
func tambahSlice(s []int)  { s = append(s, 999) }

func bungkus() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("bungkus: %w", err)
		}
	}()
	return errors.New("gagal asli")
}

func main() {
	for i := 0; i < 3; i++ {
		defer fmt.Println("defer ke-", i)
	}

	x := 1
	defer fmt.Println("x saat defer didaftarkan:", x)
	x = 2

	t := Titik{1}
	ubahNilai(t)
	fmt.Println("setelah ubahNilai:", t)
	ubahPointer(&t)
	fmt.Println("setelah ubahPointer:", t)

	s := []int{1, 2, 3}
	ubahIsi(s)
	tambahSlice(s)
	fmt.Println("slice:", s, "len:", len(s))

	fmt.Println(bungkus())
}
