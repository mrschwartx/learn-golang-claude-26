// closure_loop.go
package main

import "fmt"

type User struct {
	Nama  string
	Aktif bool
}

func main() {
	var fungsi []func() int
	for i := range 3 {
		fungsi = append(fungsi, func() int { return i })
	}
	for _, f := range fungsi {
		fmt.Print(f(), " ") // Go 1.22+: 0 1 2
	}
	fmt.Println()

	users := []User{{"Andi", false}, {"Budi", false}}
	for _, u := range users {
		u.Aktif = true // hanya mengubah salinan
	}
	fmt.Println(users)
	for i := range users {
		users[i].Aktif = true
	}
	fmt.Println(users)
}
