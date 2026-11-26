// waktu_json.go
package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Status int

const (
	Aktif Status = iota + 1
	Nonaktif
)

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[Status]string{Aktif: "aktif", Nonaktif: "nonaktif"}[s])
}

type Pengguna struct {
	Nama   string    `json:"nama"`
	Umur   int       `json:"umur,omitempty"`
	Status Status    `json:"status"`
	Lahir  time.Time `json:"lahir"`
}

func main() {
	t := time.Date(2026, time.November, 23, 14, 30, 0, 0, time.UTC)
	fmt.Println(t.Format("2006-01-02 15:04:05"), t.Format("02/01/2006"), t.Format(time.RFC3339))
	fmt.Println(t.Add(36*time.Hour).Weekday(), t.Sub(t.Add(-90*time.Minute)))

	parsed, err := time.Parse("2006-01-02", "2027-01-15")
	fmt.Println(parsed.Sub(t).Round(time.Hour), err)

	p := Pengguna{Nama: "Andi", Status: Aktif, Lahir: time.Date(1995, 5, 17, 0, 0, 0, 0, time.UTC)}
	b, _ := json.MarshalIndent(p, "", "  ")
	fmt.Println(string(b))

	var umum map[string]any
	_ = json.Unmarshal([]byte(`{"id":7,"tag":["a","b"],"profil":{"kota":"Jakarta"}}`), &umum)
	fmt.Printf("%v %T\n", umum["id"], umum["id"])
	fmt.Println(umum["profil"].(map[string]any)["kota"])

	dec := json.NewDecoder(strings.NewReader(`{"nama":"Budi","tidak_dikenal":1}`))
	dec.DisallowUnknownFields()
	var q struct {
		Nama string `json:"nama"`
	}
	fmt.Println(dec.Decode(&q))
}
