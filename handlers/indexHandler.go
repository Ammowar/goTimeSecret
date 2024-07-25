package handlers

import (
	"html/template"
	"net/http"
)

type IndexHandler struct {
	templates *template.Template
}

func NewIndexHandler(templates *template.Template) *IndexHandler {

	return &IndexHandler{
		templates: templates,
	}
}

func (ih *IndexHandler) Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	ih.templates.ExecuteTemplate(w, "index.html", nil)
}
