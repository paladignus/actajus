// Package template
package template

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"strings"
)

//go:embed templates/*.html
var templatesFS embed.FS

type HTMLRenderer struct {
	tpl      *template.Template
	subjects map[string]string
}

func NewHTMLRenderer() (*HTMLRenderer, error) {
	t, err := template.New("emails").ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("parse templates: %w", err)
	}
	return &HTMLRenderer{
		tpl: t,
		subjects: map[string]string{
			"password-reset": "Recuperação de senha",
		},
	}, nil
}

func (r *HTMLRenderer) RenderHTML(ctx context.Context, templateName string, data any) (string, string, error) {
	subject, ok := r.subjects[templateName]
	if !ok {
		return "", "", fmt.Errorf("unknown template: %s", templateName)
	}
	var b strings.Builder
	if err := r.tpl.ExecuteTemplate(&b, templateName+".html", data); err != nil {
		return "", "", fmt.Errorf("execute template %s: %w", templateName, err)
	}
	return subject, b.String(), nil
}
