package template

import (
	"html/template"
	"net/http"
)

type Engine struct {
	templates *template.Template
}

func NewEngine(pattern string) (*Engine, error) {
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		return nil, err
	}

	return &Engine{
		templates: tmpl,
	}, nil
}

func (e *Engine) Render(w http.ResponseWriter, name string, data interface{}) error {
	return e.templates.ExecuteTemplate(w, name, data)
}
