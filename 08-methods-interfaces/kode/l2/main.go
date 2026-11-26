// Latihan 2: Celsius dengan String() dan KeFahrenheit().
package main

import "fmt"

type Celsius float64

func (c Celsius) String() string        { return fmt.Sprintf("%.1f°C", float64(c)) }
func (c Celsius) KeFahrenheit() float64 { return float64(c)*9/5 + 32 }

type Polos float64 // tanpa String(), pembanding

type KelvinPtr float64

func (k *KelvinPtr) String() string { return "KelvinPtr(pointer receiver)" }

func main() {
	suhu := Celsius(36.6)
	fmt.Println(suhu)                                         // fmt memanggil String()
	fmt.Printf("%v | %s | %.2f\n", suhu, suhu, float64(suhu)) // %v dan %s -> String(); %f pada float64() -> angka
	fmt.Printf("%d°C = %.1f°F\n", 100, Celsius(100).KeFahrenheit())
	fmt.Println(Celsius(0).KeFahrenheit(), Celsius(-40).KeFahrenheit())

	fmt.Println("tanpa Stringer:", Polos(36.6))

	k := KelvinPtr(300)
	fmt.Println("nilai  :", k)  // method set KelvinPtr (nilai) tidak punya String()
	fmt.Println("pointer:", &k) // method set *KelvinPtr punya String()

	// bagaimana fmt memutuskan: type assertion ke fmt.Stringer
	var x any = suhu
	if s, ok := x.(fmt.Stringer); ok {
		fmt.Println("suhu memenuhi fmt.Stringer ->", s.String())
	}
	var y any = Polos(1)
	_, ok := y.(fmt.Stringer)
	fmt.Println("Polos memenuhi fmt.Stringer?", ok)
}
