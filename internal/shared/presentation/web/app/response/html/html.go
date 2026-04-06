package html

import (
	"net/http"

	securityheaders "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/securityheaders"
	pagepresenter "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/page"
	webtemplate "github.com/paladignus/actajus/internal/shared/presentation/web/template"
)

type Responder struct {
	presenter *pagepresenter.Presenter
	headers   *securityheaders.Adapter
}

func New(presenter *pagepresenter.Presenter, headers *securityheaders.Adapter) *Responder {
	return &Responder{presenter: presenter, headers: headers}
}

func (r *Responder) Render(w http.ResponseWriter, status int, templateName string, page webtemplate.Page) error {
	if r.headers != nil {
		r.headers.ApplyDocument(w)
	}
	return r.presenter.Render(w, status, templateName, page)
}
