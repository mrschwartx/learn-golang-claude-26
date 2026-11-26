// sapa.go
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	nama := flag.String("nama", "dunia", "nama yang disapa")
	ulang := flag.Int("n", 1, "jumlah pengulangan")
	besar := flag.Bool("besar", false, "cetak dengan huruf besar")
	flag.Parse()

	for range *ulang {
		pesan := "Halo, " + *nama + "!"
		if *besar {
			pesan = strings.ToUpper(pesan)
		}
		fmt.Println(pesan)
	}
	if flag.NArg() > 0 {
		fmt.Fprintln(os.Stderr, "argumen sisa:", flag.Args())
	}
}
