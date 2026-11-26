package diskon

import "fmt"

var DiskonMaks = 50

// Terapkan mengurangi total dengan persentase diskon.
func Terapkan(total, persen int) (int, error) {
	if persen < 0 || persen > DiskonMaks {
		return 0, fmt.Errorf("diskon %d%% di luar batas 0-%d%%", persen, DiskonMaks)
	}
	return total - total*persen/100, nil
}
