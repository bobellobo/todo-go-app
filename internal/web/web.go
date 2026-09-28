package web

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed templates/*.html
var templateFS embed.FS
var tmpl = template.Must(template.ParseFS(templateFS, "templates/*.html"))

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
