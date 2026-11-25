// Latihan 2: Setor(*Rekening) vs Cek(Rekening).
package main

import "fmt"

type Rekening struct{ Saldo int }

func Setor(r *Rekening, jumlah int) { r.Saldo += jumlah }
func Cek(r Rekening) int            { return r.Saldo }

// Pembanding: Setor dengan nilai tidak berguna
func SetorSalah(r Rekening, jumlah int) { r.Saldo += jumlah }

func main() {
	r := Rekening{Saldo: 100}
	SetorSalah(r, 50)
	fmt.Println("SetorSalah (nilai) :", Cek(r))
	Setor(&r, 50)
	fmt.Println("Setor (pointer)    :", Cek(r))
}
