package companyhandler

import (
	"net/http"

	companyread "github.com/paladignus/actajus/internal/module/company/application/readmodel"
	companyusecase "github.com/paladignus/actajus/internal/module/company/application/usecase"
	requestmeta "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/requestmeta"
	companybinder "github.com/paladignus/actajus/internal/shared/presentation/web/app/binder/http/company"
	pathbinder "github.com/paladignus/actajus/internal/shared/presentation/web/app/binder/http/path"
	handlerctx "github.com/paladignus/actajus/internal/shared/presentation/web/app/handlerctx"
	companymessage "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/message/company"
	companypage "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/page/company"
	httperror "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/http/error"
)

type Controller struct {
	Context  handlerctx.Context
	Create   companyusecase.CreateCompany
	List     companyusecase.ListCompanies
	Update   companyusecase.UpdateCompany
	Delete   companyusecase.DeleteCompany
	FindByID companyusecase.FindByID
}

func (c Controller) CompaniesPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := c.Context.RequireAuth(w, r); err != nil {
			return
		}
		query := companybinder.BindListQuery(r)
		rm, err := c.List.Execute(r.Context(), query.Filter, query.After, query.Before, query.Limit, requestmeta.BaseURL(r)+"/companies")
		if err != nil {
			c.Context.Logger.Error(r.Context(), "list companies", "error", err)
			httperror.InternalServerError(w)
			return
		}
		page := companypage.CompaniesPage(rm, query.Filter, query.Page, query.Limit, c.Context.Navigator.ReadFlash(w, r, c.Context.FlashCookie), valueOrEmpty(rm.PageInfo.NextURL), valueOrEmpty(rm.PageInfo.PreviousURL))
		if err := c.Context.HTML.Render(w, http.StatusOK, "pages/companies", page); err != nil {
			c.Context.Logger.Error(r.Context(), "render companies page", "error", err)
			httperror.InternalServerError(w)
		}
	})
}

func (c Controller) CompanyNewPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := c.Context.RequireAuth(w, r); err != nil {
			return
		}
		c.renderNewForm(w)
	})
}

func (c Controller) CompanyCreateAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := c.Context.RequireAuth(w, r)
		if err != nil {
			return
		}
		form, cmd, err := companybinder.BindCreateForm(r, state.Claims.IDUser)
		if err != nil {
			httperror.BadRequest(w, "invalid form")
			return
		}
		company, execErr := c.Create.Execute(r.Context(), cmd)
		if execErr != nil {
			form.Error = execErr.Error()
			c.renderForm(w, http.StatusBadRequest, form)
			return
		}
		if c.Context.AuditSuccess != nil {
			c.Context.AuditSuccess(r, "company.create.succeeded", "id_user", state.Claims.IDUser, "id_company", company.ID)
		}
		c.redirectToDetail(w, r, company.ID)
	})
}

func (c Controller) CompanyDetailPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := c.Context.RequireAuth(w, r); err != nil {
			return
		}
		company, ok := c.companyByPath(w, r, "find company by id")
		if !ok {
			return
		}
		if err := c.Context.HTML.Render(w, http.StatusOK, "pages/company-detail", companypage.DetailPage(company)); err != nil {
			c.Context.Logger.Error(r.Context(), "render company detail", "error", err, "id_company", company.ID)
			httperror.InternalServerError(w)
		}
	})
}

func (c Controller) CompanyEditPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := c.Context.RequireAuth(w, r); err != nil {
			return
		}
		company, ok := c.companyByPath(w, r, "find company for edit")
		if !ok {
			return
		}
		c.renderEditForm(w, company)
	})
}

func (c Controller) CompanyUpdateAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := c.Context.RequireAuth(w, r); err != nil {
			return
		}
		company, ok := c.companyByPath(w, r, "find company for edit")
		if !ok {
			return
		}
		form, cmd, err := companybinder.BindUpdateForm(r, company)
		if err != nil {
			httperror.BadRequest(w, "invalid form")
			return
		}
		updated, execErr := c.Update.Execute(r.Context(), cmd)
		if execErr != nil {
			form.Error = execErr.Error()
			c.renderForm(w, http.StatusBadRequest, form)
			return
		}
		if c.Context.AuditSuccess != nil {
			c.Context.AuditSuccess(r, "company.update.succeeded", "id_company", updated.ID)
		}
		c.redirectToDetail(w, r, updated.ID)
	})
}

func (c Controller) CompanyDeleteAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := c.Context.RequireAuth(w, r); err != nil {
			return
		}
		id, err := pathbinder.BindInt64Path(r, "id")
		if err != nil {
			c.Context.RespondNotFound(w, r)
			return
		}
		err = c.Delete.Execute(r.Context(), id)
		if err != nil {
			c.Context.Logger.Error(r.Context(), "delete company", "error", err, "id_company", id)
			if c.Context.AuditFailure != nil {
				c.Context.AuditFailure(r, "company.delete.failed", "id_company", id)
			}
		} else if c.Context.AuditSuccess != nil {
			c.Context.AuditSuccess(r, "company.delete.succeeded", "id_company", id)
		}
		c.redirectDeleteResult(w, r, err)
	})
}

func (c Controller) renderForm(w http.ResponseWriter, status int, data companybinder.FormData) {
	templateName, page := companypage.FormPage(companypage.FormDataFromBinder(data), "pages/company-form")
	if err := c.Context.HTML.Render(w, status, templateName, page); err != nil {
		c.Context.Logger.Error(nil, "render company form", "error", err)
		httperror.InternalServerError(w)
	}
}

func (c Controller) renderNewForm(w http.ResponseWriter) {
	templateName, page := companypage.FormPage(companypage.NewFormData(), "pages/company-form")
	if err := c.Context.HTML.Render(w, http.StatusOK, templateName, page); err != nil {
		httperror.InternalServerError(w)
	}
}

func (c Controller) renderEditForm(w http.ResponseWriter, company *companyread.CompanyReadModel) {
	templateName, page := companypage.FormPage(companypage.FormDataFromReadModel(company), "pages/company-form")
	if err := c.Context.HTML.Render(w, http.StatusOK, templateName, page); err != nil {
		c.Context.Logger.Error(nil, "render company edit form", "error", err, "id_company", company.ID)
		httperror.InternalServerError(w)
	}
}

func (c Controller) companyByPath(w http.ResponseWriter, r *http.Request, action string) (*companyread.CompanyReadModel, bool) {
	id, err := pathbinder.BindInt64Path(r, "id")
	if err != nil {
		c.Context.RespondNotFound(w, r)
		return nil, false
	}
	company, err := c.FindByID.Execute(r.Context(), id)
	if err != nil {
		c.Context.Logger.Error(r.Context(), action, "error", err, "id_company", id)
		httperror.InternalServerError(w)
		return nil, false
	}
	if company == nil {
		c.Context.RespondNotFound(w, r)
		return nil, false
	}
	return company, true
}

func (c Controller) redirectToDetail(w http.ResponseWriter, r *http.Request, id int64) {
	c.Context.Navigator.SeeOther(w, r, companypage.ActionPath(id))
}

func (c Controller) redirectDeleteResult(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		c.Context.Navigator.SeeOtherWithFlash(w, r, "/companies", c.Context.FlashCookie, companymessage.DeleteFailed())
		return
	}
	c.Context.Navigator.SeeOtherWithFlash(w, r, "/companies", c.Context.FlashCookie, companymessage.DeleteSucceeded())
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
