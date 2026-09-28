package web

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
)

//go:embed templates/*.html static
var templateFS embed.FS
var tmpl = template.Must(template.ParseFS(templateFS, "templates/*.html"))

var staticFS = func() fs.FS {
	files, err := fs.Sub(templateFS, "static")
	if err != nil {
		panic(err)
	}
	return files
}()

func RegisterStaticRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /manifest.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		http.ServeFileFS(w, r, templateFS, "static/manifest.json")
	})
	mux.HandleFunc("GET /sw.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Service-Worker-Allowed", "/")
		http.ServeFileFS(w, r, templateFS, "static/sw.js")
	})
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))
}

// render the initial shell index.html
func RenderPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

// render a single todo item html fragment
func RenderItem(w http.ResponseWriter, item interface{}) {
	w.Header().Set("Content-Type", "text/html")
	tmpl.ExecuteTemplate(w, "todo_item.html", item)
}

// renders a list of todo items html fragments
func RenderList(w http.ResponseWriter, items []interface{}) {
	w.Header().Set("Content-Type", "text/html")
	for _, item := range items {
		tmpl.ExecuteTemplate(w, "todo_item.html", item)
	}
}
