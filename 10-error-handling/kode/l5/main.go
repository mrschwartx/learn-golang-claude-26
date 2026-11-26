// Latihan 5: Aman(f) mengubah panic menjadi error + stack trace.
package main

import (
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

type ErrPanic struct {
	Nilai any
	Stack string
}

func (e *ErrPanic) Error() string { return fmt.Sprintf("panic: %v", e.Nilai) }

func Aman(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = &ErrPanic{Nilai: r, Stack: string(debug.Stack())}
		}
	}()
	f()
	return nil
}

func bacaDiLuarRentang() {
	s := []int{1, 2, 3}
	i := 5
	_ = s[i]
}

func tulisKeMapNil() {
	var m map[string]int
	m["x"] = 1
}

func main() {
	for nama, f := range map[string]func(){
		"1 slice di luar rentang": bacaDiLuarRentang,
		"2 map nil":               tulisKeMapNil,
		"3 tidak panic":           func() {},
	} {
		err := Aman(f)
		fmt.Printf("%-24s -> err=%v\n", nama, err)

		var ep *ErrPanic
		if errors.As(err, &ep) {
			var re runtime.Error
			isRuntime := false
			if e, ok := ep.Nilai.(error); ok {
				isRuntime = errors.As(e, &re)
			}
			fmt.Printf("   runtime.Error? %t | stack menyebut fungsi sumber? %t | baris stack: %d\n",
				isRuntime,
				strings.Contains(ep.Stack, "bacaDiLuarRentang") || strings.Contains(ep.Stack, "tulisKeMapNil"),
				strings.Count(ep.Stack, "\n"))
		}
	}
	err := Aman(bacaDiLuarRentang)
	fmt.Println("\ncuplikan stack (6 baris pertama):")
	fmt.Println(strings.Join(strings.Split(err.(*ErrPanic).Stack, "\n")[:6], "\n"))
}
