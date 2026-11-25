// Latihan 3: enum Level dengan iota + bit flag izin file.
package main

import "fmt"

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

func (l Level) String() string {
	switch l {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	default:
		return fmt.Sprintf("Level(%d)", int(l))
	}
}

type Izin uint8

const (
	Baca  Izin = 1 << iota // 1
	Tulis                  // 2
	Eksek                  // 4
)

func (i Izin) Punya(f Izin) bool { return i&f != 0 }

func main() {
	for _, l := range []Level{Debug, Info, Warn, Error, 9} {
		fmt.Printf("%d -> %s\n", int(l), l)
	}

	izin := Baca | Tulis
	fmt.Printf("izin=%03b baca=%t tulis=%t eksek=%t\n",
		izin, izin.Punya(Baca), izin.Punya(Tulis), izin.Punya(Eksek))

	izin |= Eksek  // tambah izin
	izin &^= Tulis // cabut izin tulis (AND NOT)
	fmt.Printf("izin=%03b baca=%t tulis=%t eksek=%t\n",
		izin, izin.Punya(Baca), izin.Punya(Tulis), izin.Punya(Eksek))
}
