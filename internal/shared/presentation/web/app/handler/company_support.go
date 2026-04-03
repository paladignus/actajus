package handler

import (
	"net/http"

	companyread "github.com/paladignus/actajus/internal/module/company/application/readmodel"
	companybinder "github.com/paladignus/actajus/internal/shared/presentation/web/app/binder/http/company"
	pathbinder "github.com/paladignus/actajus/internal/shared/presentation/web/app/binder/http/path"
	companymessage "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/message/company"
	companypage "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/page/company"
	httperror "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/http/error"
)

type companySupport struct {
	routes companyRoutes
}

func (h companyRoutes) support() companySupport {
	return companySupport{routes: h}
}

func (s companySupport) renderForm(w http.ResponseWriter, status int, data companybinder.FormData) {
	templateName, page := companypage.FormPage(companypage.FormDataFromBinder(data), "pages/company-form")
	if err := s.routes.app.html.Render(w, status, templateName, page); err != nil {
		s.routes.app.logger.Error(nil, "render company form", "error", err)
		httperror.InternalServerError(w)
	}
}

func (s companySupport) renderNewForm(w http.ResponseWriter) {
	templateName, page := companypage.FormPage(companypage.NewFormData(), "pages/company-form")
	if err := s.routes.app.html.Render(w, http.StatusOK, templateName, page); err != nil {
		httperror.InternalServerError(w)
	}
}

func (s companySupport) renderEditForm(w http.ResponseWriter, company *companyread.CompanyReadModel) {
	templateName, page := companypage.FormPage(companypage.FormDataFromReadModel(company), "pages/company-form")
	if err := s.routes.app.html.Render(w, http.StatusOK, templateName, page); err != nil {
		s.routes.app.logger.Error(nil, "render company edit form", "error", err, "id_company", company.ID)
		httperror.InternalServerError(w)
	}
}

func (s companySupport) companyByPath(w http.ResponseWriter, r *http.Request, action string) (*companyread.CompanyReadModel, bool) {
	id, err := pathbinder.BindInt64Path(r, "id")
	if err != nil {
		s.routes.app.respondNotFound(w, r)
		return nil, false
	}
	company, err := s.routes.app.company.findByID.Execute(r.Context(), id)
	if err != nil {
		s.routes.app.logger.Error(r.Context(), action, "error", err, "id_company", id)
		httperror.InternalServerError(w)
		return nil, false
	}
	if company == nil {
		s.routes.app.respondNotFound(w, r)
		return nil, false
	}
	return company, true
}

func (s companySupport) redirectToDetail(w http.ResponseWriter, r *http.Request, id int64) {
	s.routes.app.navigator.SeeOther(w, r, companypage.ActionPath(id))
}

func (s companySupport) redirectDeleteResult(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		s.routes.app.navigator.SeeOtherWithFlash(w, r, "/companies", cookieFlash, companymessage.DeleteFailed())
		return
	}
	s.routes.app.navigator.SeeOtherWithFlash(w, r, "/companies", cookieFlash, companymessage.DeleteSucceeded())
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
