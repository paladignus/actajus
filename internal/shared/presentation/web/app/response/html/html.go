package html

import (
	"net/http"

	pagepresenter "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/page"
	webtemplate "github.com/paladignus/actajus/internal/shared/presentation/web/template"
)

type Responder struct {
	presenter *pagepresenter.Presenter
}

func New(presenter *pagepresenter.Presenter) *Responder {
	return &Responder{presenter: presenter}
}

func (r *Responder) Render(w http.ResponseWriter, status int, templateName string, page webtemplate.Page) error {
	return r.presenter.Render(w, status, templateName, page)
}
