// Package handler
package handler

import (
	"net/http"

	"github.com/paladignus/actajus/internal/infrastructure/http/handler"
	"github.com/paladignus/actajus/internal/module/company/application/dto"
	"github.com/paladignus/actajus/internal/module/company/application/usecase"
)

type CompanyHandler struct {
	create usecase.CreateCompany
}

func NewCompanyHandler(
	create usecase.CreateCompany,
) CompanyHandler {
	return CompanyHandler{
		create,
	}
}

func (c CompanyHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /companies", c.Create)
}

func (c CompanyHandler) Create(w http.ResponseWriter, r *http.Request) {
	req, err := handler.DecodeJSONRequest[dto.CreateCompanyRequest](r)
	if err != nil {
		handler.RespondError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}
	response, err := c.create.Execute(r.Context(), req)
	if err != nil {
		handler.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	handler.RespondJSON(w, http.StatusCreated, response)
}
