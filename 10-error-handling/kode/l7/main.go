// Latihan 7: panic untuk kegagalan konfigurasi vs desain idiomatik.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
)

type Config struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

// ---------- KELIRU ----------
func muatPanic(path string) Config {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err) // file tidak ada? crash
	}
	var c Config
	if err := json.Unmarshal(b, &c); err != nil {
		panic(err)
	}
	return c
}

func mainKeliru() {
	defer func() {
		if r := recover(); r != nil { // "program tidak pernah crash"
			fmt.Println("[keliru] dipulihkan:", r, "-> lanjut dengan Config kosong?!")
		}
	}()
	c := muatPanic("tidak-ada.json")
	fmt.Println("[keliru] tidak sampai sini", c)
}

// ---------- IDIOMATIK ----------
var ErrConfigTidakAda = errors.New("config tidak ditemukan")

func Muat(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return Config{}, fmt.Errorf("muat %s: %w", path, ErrConfigTidakAda)
	}
	if err != nil {
		return Config{}, fmt.Errorf("muat %s: %w", path, err)
	}
	var c Config
	if err := json.Unmarshal(b, &c); err != nil {
		return Config{}, fmt.Errorf("parse %s: %w", path, err)
	}
	if c.Port == 0 {
		return Config{}, fmt.Errorf("parse %s: port wajib diisi", path)
	}
	return c, nil
}

func run() error {
	c, err := Muat("tidak-ada.json")
	if errors.Is(err, ErrConfigTidakAda) {
		fmt.Println("[idiomatik] config tidak ada, pakai default")
		c = Config{Host: "localhost", Port: 8080}
	} else if err != nil {
		return err // config RUSAK: gagal dengan pesan jelas
	}
	fmt.Println("[idiomatik] config dipakai:", c)

	_, err = Muat("/dev/null")
	fmt.Println("[idiomatik] config rusak ->", err)
	return err
}

func main() {
	mainKeliru()
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		defer os.Exit(1) // keluar dengan kode jelas, bukan crash/stack trace
	}
}
