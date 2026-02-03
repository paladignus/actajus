// Package handler
package handler

import (
	"errors"
	"net/http"

	"github.com/paladignus/actajus/internal/infrastructure/http/handler"
	"github.com/paladignus/actajus/internal/module/company/application/dto"
	"github.com/paladignus/actajus/internal/module/company/application/usecase"
	"github.com/paladignus/actajus/internal/shared/domain"
)

type CompanyHandler struct {
	create     usecase.CreateCompany
	update     usecase.UpdateCompany
	delete     usecase.DeleteCompany
	list       usecase.ListCompanies
	findByCNPJ usecase.FindByCNPJ
	findByID   usecase.FindByID
}

func NewCompanyHandler(
	create usecase.CreateCompany,
	update usecase.UpdateCompany,
	delete usecase.DeleteCompany,
	list usecase.ListCompanies,
	findByCNPJ usecase.FindByCNPJ,
	findByID usecase.FindByID,
) CompanyHandler {
	return CompanyHandler{
		create,
		update,
		delete,
		list,
		findByCNPJ,
		findByID,
	}
}

func (c CompanyHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /companies", c.Create)
	mux.HandleFunc("POST /companies/find/cnpj", c.FindByCNPJ)
	mux.HandleFunc("POST /companies/find/id", c.FindByID)
	mux.HandleFunc("PUT /companies", c.Update)
	mux.HandleFunc("DELETE /companies", c.Delete)
	mux.HandleFunc("GET /companies", c.List)
}

func (c CompanyHandler) Create(w http.ResponseWriter, r *http.Request) {
	req, err := handler.DecodeJSONRequest[dto.CreateCompanyRequest](r)
	if err != nil {
		handler.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	response, err := c.create.Execute(r.Context(), req)
	if err != nil {
		var vErr *domain.ValidationErrors
		if errors.As(err, &vErr) {
			handler.RespondJSON(w, http.StatusBadRequest, map[string]any{
				"errors": vErr.Errors(),
			})
			return
		}
		var dErr *domain.FieldError
		if errors.As(err, &dErr) {
			handler.RespondJSON(w, http.StatusBadRequest, map[string]any{
				"errors": domain.NewValidationErrors([]*domain.FieldError{dErr}).Errors(),
			})
			return
		}
		handler.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	handler.RespondJSON(w, http.StatusCreated, response)
}

func (c CompanyHandler) Update(w http.ResponseWriter, r *http.Request) {
	req, err := handler.DecodeJSONRequest[dto.UpdateCompanyRequest](r)
	if err != nil {
		handler.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	response, err := c.update.Execute(r.Context(), req)
	if err != nil {
		var vErr *domain.ValidationErrors
		if errors.As(err, &vErr) {
			handler.RespondJSON(w, http.StatusBadRequest, map[string]any{
				"errors": vErr.Errors(),
			})
			return
		}
		var dErr *domain.FieldError
		if errors.As(err, &dErr) {
			handler.RespondJSON(w, http.StatusBadRequest, map[string]any{
				"errors": domain.NewValidationErrors([]*domain.FieldError{dErr}).Errors(),
			})
			return
		}
		handler.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	handler.RespondJSON(w, http.StatusCreated, response)
}

func (c CompanyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	req, err := handler.DecodeJSONRequest[dto.DeleteCompanyRequest](r)
	if err != nil {
		handler.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	err = c.delete.Execute(r.Context(), req.IDCompany)
	if err != nil {
		handler.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	handler.RespondJSON(w, http.StatusOK, map[string]string{})
}

func (c CompanyHandler) List(w http.ResponseWriter, r *http.Request) {
	companies, err := c.list.Execute(r.Context(), 1, 10)
	if err != nil {
		handler.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response := map[string]any{
		"companies": companies,
	}
	handler.RespondJSON(w, http.StatusOK, response)

	// page := r.URL.Query().Get("page")
	// fmt.Println(page)
	// return
}

func (c CompanyHandler) FindByCNPJ(w http.ResponseWriter, r *http.Request) {
	req, err := handler.DecodeJSONRequest[dto.FindCompanyByCNPJRequest](r)
	if err != nil {
		handler.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	company, err := c.findByCNPJ.Execute(r.Context(), req.CNPJ)
	if err != nil {
		var dErr *domain.FieldError
		if errors.As(err, &dErr) {
			handler.RespondJSON(w, http.StatusOK, map[string]any{})
			return
		}
		handler.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if company != nil {
		handler.RespondJSON(w, http.StatusOK, company)
		return
	}
	handler.RespondJSON(w, http.StatusOK, map[string]string{})
}

func (c CompanyHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	req, err := handler.DecodeJSONRequest[dto.FindCompanyByIDRequest](r)
	if err != nil {
		handler.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	company, err := c.findByID.Execute(r.Context(), req.ID)
	if err != nil {
		var dErr *domain.FieldError
		if errors.As(err, &dErr) {
			handler.RespondJSON(w, http.StatusOK, map[string]any{})
			return
		}
		handler.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if company != nil {
		handler.RespondJSON(w, http.StatusOK, company)
		return
	}
	handler.RespondJSON(w, http.StatusOK, map[string]string{})
}
