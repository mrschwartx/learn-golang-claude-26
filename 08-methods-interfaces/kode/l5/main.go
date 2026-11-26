// Latihan 5: Ringkas(x any) dengan type switch.
package main

import (
	"errors"
	"fmt"
	"strings"
)

func Ringkas(x any) string {
	switch v := x.(type) {
	case nil:
		return "nil"
	case int:
		return fmt.Sprintf("int %d (+1 = %d)", v, v+1)
	case float64:
		return fmt.Sprintf("float64 %.2f", v)
	case string:
		return fmt.Sprintf("string %q (upper %s, %d byte)", v, strings.ToUpper(v), len(v))
	case []int:
		jumlah := 0
		for _, n := range v {
			jumlah += n
		}
		return fmt.Sprintf("[]int len=%d jumlah=%d", len(v), jumlah)
	case error: // urutan penting: error adalah interface
		return "error: " + v.Error()
	case int8, int16, int32, int64: // case multi-tipe: v tetap bertipe any
		return fmt.Sprintf("integer lain %T = %v", v, v)
	default:
		return fmt.Sprintf("tipe lain: %T", v)
	}
}

func main() {
	var kosong error // interface nil
	var ptrNil *int
	for _, x := range []any{
		42, 3.14159, "halo", []int{1, 2, 3}, errors.New("gagal"),
		nil, kosong, int64(7), true, map[string]int{}, ptrNil,
	} {
		fmt.Println(Ringkas(x))
	}
}
