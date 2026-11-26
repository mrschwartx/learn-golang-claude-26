// l7: aplikasi web satu-binary: template, CSS, dan migrasi SQL disematkan dengan go:embed.
package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
)

//go:embed templates/*.html
var tplFS embed.FS

//go:embed static
var staticFS embed.FS

//go:embed migrasi/*.sql
var migrasiFS embed.FS

func daftarMigrasi() ([]string, error) {
	e, err := fs.ReadDir(migrasiFS, "migrasi")
	if err != nil {
		return nil, err
	}
	var n []string
	for _, x := range e {
		n = append(n, x.Name())
	}
	sort.Strings(n) // ReadDir sudah terurut; eksplisit agar niatnya jelas
	return n, nil
}

func Handler() (http.Handler, error) {
	tpl, err := template.ParseFS(tplFS, "templates/*.html")
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.FileServerFS(staticFS))
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		m, _ := daftarMigrasi()
		tpl.ExecuteTemplate(w, "index.html", map[string]any{"Judul": "Beranda", "Migrasi": m})
	})
	return mux, nil
}

func main() {
	h, err := Handler()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, path := range []string{"/", "/static/app.css", "/static/tidak-ada.css"} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
		fmt.Printf("%-22s -> %d %q\n", path, rec.Code, rec.Body.String())
	}
	sql, _ := migrasiFS.ReadFile("migrasi/001_buat_pengguna.sql")
	fmt.Printf("isi migrasi 001: %s", sql)
}
