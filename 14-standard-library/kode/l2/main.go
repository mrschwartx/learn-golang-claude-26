// l2: CLI dengan dua subcommand (tambah, hapus), masing-masing FlagSet sendiri.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// Jalankan mengembalikan kode keluar; out dan errOut dapat diinjeksi untuk diuji.
func Jalankan(args []string, out, errOut io.Writer) int {
	if len(args) < 1 {
		fmt.Fprintln(errOut, "pemakaian: app <tambah|hapus> [flag] ...")
		return 2
	}
	switch args[0] {
	case "tambah":
		fs := flag.NewFlagSet("tambah", flag.ContinueOnError)
		fs.SetOutput(errOut)
		nama := fs.String("nama", "", "nama item (wajib)")
		jumlah := fs.Int("jumlah", 1, "jumlah item")
		tag := fs.String("tag", "", "tag dipisah koma")
		if err := fs.Parse(args[1:]); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				return 0 // -h bukan kesalahan
			}
			return 2
		}
		if *nama == "" {
			fmt.Fprintln(errOut, "tambah: -nama wajib diisi")
			return 2
		}
		fmt.Fprintf(out, "tambah %q x%d tag=%v sisa=%v\n", *nama, *jumlah, strings.Split(*tag, ","), fs.Args())
	case "hapus":
		fs := flag.NewFlagSet("hapus", flag.ContinueOnError)
		fs.SetOutput(errOut)
		id := fs.Int("id", 0, "ID item (wajib)")
		paksa := fs.Bool("paksa", false, "hapus tanpa konfirmasi")
		if err := fs.Parse(args[1:]); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				return 0 // -h bukan kesalahan
			}
			return 2
		}
		if *id == 0 {
			fmt.Fprintln(errOut, "hapus: -id wajib diisi")
			return 2
		}
		fmt.Fprintf(out, "hapus id=%d paksa=%v\n", *id, *paksa)
	default:
		fmt.Fprintf(errOut, "subcommand tidak dikenal: %q\n", args[0])
		return 2
	}
	return 0
}

func main() { os.Exit(Jalankan(os.Args[1:], os.Stdout, os.Stderr)) }
