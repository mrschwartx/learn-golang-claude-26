// kompleks.go
package main

import (
	"fmt"
	"math/cmplx"
)

func main() {
	a := complex(3, 4)        // complex128 (default)
	var b complex128 = 1 + 2i // literal imajiner
	var c complex64 = complex(1, 1)

	fmt.Printf("%T %v\n", a, a)
	fmt.Printf("%T %v\n", c, c)

	fmt.Println(a+b, a*b)           // (4+6i) (-5+10i)
	fmt.Println(real(a), imag(a))   // 3 4
	fmt.Println(cmplx.Abs(a))       // 5  (panjang vektor 3-4-5)
	fmt.Println(cmplx.Sqrt(-1))     // (0+1i)
	fmt.Println(a == complex(3, 4)) // true
	// fmt.Println(a < b)            // error kompilasi: bilangan kompleks tidak terurut
}
