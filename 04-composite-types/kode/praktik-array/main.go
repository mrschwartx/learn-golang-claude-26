// array.go
package main

import "fmt"

type Titik struct{ X, Y int }

func ubah(a [3]int) { a[0] = 999 } // menerima SALINAN

func main() {
	arr := [3]int{1, 2, 3}
	ubah(arr)
	fmt.Println("array asli:", arr)

	s := arr[:]
	s[0] = 100
	fmt.Println("array berubah lewat slice:", arr)

	salinan := [3]int(s)
	salinan[1] = -1
	fmt.Println("salinan:", salinan, "asli:", arr)

	p1, p2 := Titik{1, 2}, Titik{1, 2}
	fmt.Println("sama:", p1 == p2)

	grid := map[Titik]string{{0, 0}: "asal", {1, 2}: "target"}
	fmt.Println(grid[Titik{1, 2}])
}
