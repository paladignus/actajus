// Package handler
package handler

import (
	"encoding/json"
	"net/http"

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
	var input dto.CreateCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		c.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	response, err := c.create.Execute(r.Context(), input)
	if err != nil {
		c.respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	c.respondJSON(w, http.StatusCreated, response)
}

func (c CompanyHandler) respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (c CompanyHandler) respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
