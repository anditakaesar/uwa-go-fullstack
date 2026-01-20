package web

import (
	"net/http"
	"text/template"
)

type Renderer struct {
	t *template.Template
}

func NewRenderer() *Renderer {
	t := template.Must(
		template.ParseFS(
			TemplatesFS,
			"templates/*.html",
		),
	)

	return &Renderer{t: t}
}

func (r *Renderer) Render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := r.t.ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}
