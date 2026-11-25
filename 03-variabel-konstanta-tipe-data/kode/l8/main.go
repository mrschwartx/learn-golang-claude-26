// Latihan 8: akar x^2 + 2x + 5 = 0 dengan complex128.
package main

import (
	"fmt"
	"math/cmplx"
)

func main() {
	a, b, c := complex(1, 0), complex(2, 0), complex(5, 0)

	disk := b*b - 4*a*c // 4 - 20 = -16
	fmt.Println("diskriminan:", disk)

	akar := cmplx.Sqrt(disk) // 4i
	x1 := (-b + akar) / (2 * a)
	x2 := (-b - akar) / (2 * a)

	for i, x := range []complex128{x1, x2} {
		fmt.Printf("x%d = %v  (real=%g, imag=%g)\n", i+1, x, real(x), imag(x))
		hasil := a*x*x + b*x + c
		fmt.Printf("   substitusi: %v  |hasil|=%g  nol? %t\n", hasil, cmplx.Abs(hasil), cmplx.Abs(hasil) < 1e-12)
	}
	// fmt.Println(x1 < x2) // error: invalid operation: operator < not defined on x1 (variable of type complex128)
}
