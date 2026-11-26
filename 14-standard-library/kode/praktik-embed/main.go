// main.go
package main

import (
	"embed"
	"fmt"
)

//go:embed salam.txt
var salam string

//go:embed templates/*
var templates embed.FS

func main() {
	fmt.Print(salam)
	b, _ := templates.ReadFile("templates/index.html")
	fmt.Print(string(b))
	daftar, _ := templates.ReadDir("templates")
	fmt.Println(len(daftar), daftar[0].Name())
}
