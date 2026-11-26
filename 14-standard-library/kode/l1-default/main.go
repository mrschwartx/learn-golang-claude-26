// l1-default: bukti batas 64 KB default bufio.Scanner untuk satu baris.
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	f, _ := os.Open(os.Args[1])
	defer f.Close()
	sc := bufio.NewScanner(f)
	n := 0
	for sc.Scan() {
		n++
	}
	fmt.Println("baris terbaca:", n, "| sc.Err():", sc.Err())
}
