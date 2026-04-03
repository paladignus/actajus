package page

import (
	"net/http"

	webtemplate "github.com/paladignus/actajus/internal/shared/presentation/web/template"
)

type Presenter struct {
	renderer *webtemplate.Renderer
}

func New(renderer *webtemplate.Renderer) *Presenter {
	return &Presenter{renderer: renderer}
}

func (p *Presenter) Render(w http.ResponseWriter, status int, templateName string, page webtemplate.Page) error {
	return p.renderer.RenderHTTP(w, status, templateName, page)
}
