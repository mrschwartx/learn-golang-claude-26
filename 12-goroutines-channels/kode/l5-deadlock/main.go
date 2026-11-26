// Kegagalan 1: deadlock (semua goroutine tertidur).
package main

import "fmt"

func main() {
	ch := make(chan int)
	ch <- 1 // tidak ada penerima
	fmt.Println(<-ch)
}
