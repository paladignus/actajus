package handler

import (
	"net/http"

	requestmeta "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/requestmeta"
	companybinder "github.com/paladignus/actajus/internal/shared/presentation/web/app/binder/http/company"
	pathbinder "github.com/paladignus/actajus/internal/shared/presentation/web/app/binder/http/path"
	companypage "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/page/company"
	httperror "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/http/error"
)

type companyRoutes struct {
	app AppHandler
}

func (h companyRoutes) CompaniesPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.app.requireAuth(w, r); err != nil {
			return
		}
		query := companybinder.BindListQuery(r)
		rm, err := h.app.company.list.Execute(r.Context(), query.Filter, query.After, query.Before, query.Limit, requestmeta.BaseURL(r)+"/companies")
		if err != nil {
			h.app.logger.Error(r.Context(), "list companies", "error", err)
			httperror.InternalServerError(w)
			return
		}
		page := companypage.CompaniesPage(rm, query.Filter, query.Page, query.Limit, h.app.navigator.ReadFlash(w, r, cookieFlash), valueOrEmpty(rm.PageInfo.NextURL), valueOrEmpty(rm.PageInfo.PreviousURL))
		if err := h.app.html.Render(w, http.StatusOK, "pages/companies", page); err != nil {
			h.app.logger.Error(r.Context(), "render companies page", "error", err)
			httperror.InternalServerError(w)
		}
	})
}

func (h companyRoutes) CompanyNewPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.app.requireAuth(w, r); err != nil {
			return
		}
		h.support().renderNewForm(w)
	})
}

func (h companyRoutes) CompanyCreateAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := h.app.requireAuth(w, r)
		if err != nil {
			return
		}
		form, cmd, err := companybinder.BindCreateForm(r, state.Claims.IDUser)
		if err != nil {
			httperror.BadRequest(w, "invalid form")
			return
		}
		company, execErr := h.app.company.create.Execute(r.Context(), cmd)
		if execErr != nil {
			form.Error = execErr.Error()
			h.support().renderForm(w, http.StatusBadRequest, form)
			return
		}
		h.support().redirectToDetail(w, r, company.ID)
	})
}

func (h companyRoutes) CompanyDetailPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.app.requireAuth(w, r); err != nil {
			return
		}
		company, ok := h.support().companyByPath(w, r, "find company by id")
		if !ok {
			return
		}
		if err := h.app.html.Render(w, http.StatusOK, "pages/company-detail", companypage.DetailPage(company)); err != nil {
			h.app.logger.Error(r.Context(), "render company detail", "error", err, "id_company", company.ID)
			httperror.InternalServerError(w)
		}
	})
}

func (h companyRoutes) CompanyEditPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.app.requireAuth(w, r); err != nil {
			return
		}
		company, ok := h.support().companyByPath(w, r, "find company for edit")
		if !ok {
			return
		}
		h.support().renderEditForm(w, company)
	})
}

func (h companyRoutes) CompanyUpdateAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.app.requireAuth(w, r); err != nil {
			return
		}
		company, ok := h.support().companyByPath(w, r, "find company for edit")
		if !ok {
			return
		}
		form, cmd, err := companybinder.BindUpdateForm(r, company)
		if err != nil {
			httperror.BadRequest(w, "invalid form")
			return
		}
		updated, execErr := h.app.company.update.Execute(r.Context(), cmd)
		if execErr != nil {
			form.Error = execErr.Error()
			h.support().renderForm(w, http.StatusBadRequest, form)
			return
		}
		h.support().redirectToDetail(w, r, updated.ID)
	})
}

func (h companyRoutes) CompanyDeleteAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.app.requireAuth(w, r); err != nil {
			return
		}
		id, err := pathbinder.BindInt64Path(r, "id")
		if err != nil {
			h.app.respondNotFound(w, r)
			return
		}
		err = h.app.company.delete.Execute(r.Context(), id)
		if err != nil {
			h.app.logger.Error(r.Context(), "delete company", "error", err, "id_company", id)
		}
		h.support().redirectDeleteResult(w, r, err)
	})
}
