// Package handler
package handler

import (
	"fmt"
	"net/http"

	"github.com/paladignus/actajus/internal/application/command"
	"github.com/paladignus/actajus/internal/application/service"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type Company struct {
	create service.ICreateCompany
	update service.IUpdateCompany
	getAll service.IGetAllCompany
	delete service.IDeleteCompany
	logger repository.Logger
}

func NewCompany(
	create service.ICreateCompany,
	update service.IUpdateCompany,
	getAll service.IGetAllCompany,
	delete service.IDeleteCompany,
	logger repository.Logger,
) Company {
	return Company{
		create,
		update,
		getAll,
		delete,
		logger,
	}
}

func (c Company) Create(w http.ResponseWriter, r *http.Request) {
	req, err := DecodeJSONRequest[command.CreateCompanyCommand](r)
	if err != nil {
		c.logger.Warn(r.Context(), "failed to decode request body for company create", "error", err)
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = c.create.Execute(r.Context(), req)
	if err != nil {
		statusCode, errResponse := MapDomainErrorToHTTP(err)
		if statusCode >= 500 {
			c.logger.Error(r.Context(), "create company failed with server error", "error", err, "cnpj", req.CNPJ, "method", r.Method, "url", r.URL.Path)
		} else {
			c.logger.Warn(r.Context(), "create company failed", "error", err, "cnpj", req.CNPJ, "method", r.Method, "url", r.URL.Path)
		}
		RespondJSON(w, statusCode, errResponse)
		return
	}
	c.logger.Info(r.Context(), "create Company successful", "cnpj", req.CNPJ, "registered_by", req.RegisteredBy, "method", r.Method, "url", r.URL.Path)
}

func (c Company) Update(w http.ResponseWriter, r *http.Request) {
	req, err := DecodeJSONRequest[command.UpdateCompanyCommand](r)
	if err != nil {
		c.logger.Warn(r.Context(), "failed to decode request body for company update", "error", err)
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err = c.update.Execute(r.Context(), req); err != nil {
		statusCode, errResponse := MapDomainErrorToHTTP(err)
		if statusCode >= 500 {
			c.logger.Error(r.Context(), "update company failed with server error", "error", err, "cnpj", req.CNPJ, "method", r.Method, "url", r.URL.Path)
		} else {
			c.logger.Warn(r.Context(), "update company failed", "error", err, "cnpj", req.CNPJ, "method", r.Method, "url", r.URL.Path)
		}
		RespondJSON(w, statusCode, errResponse)
		return
	}
	c.logger.Info(r.Context(), "update company successful", "cnpj", req.CNPJ, "method", r.Method, "url", r.URL.Path)
}

func (c Company) GetAll(w http.ResponseWriter, r *http.Request) {
	companies, err := c.getAll.Execute(r.Context())
	if err != nil {
		c.logger.Warn(r.Context(), "get all companies failed", "error", err, "method", r.Method, "url", r.URL.Path)
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.logger.Info(r.Context(), "get all companies successful", "method", r.Method, "url", r.URL.Path)
	RespondJSON(w, http.StatusOK, companies)
}

func (c Company) Delete(w http.ResponseWriter, r *http.Request) {
	req, err := DecodeJSONRequest[command.DeleteCompanyCommand](r)
	fmt.Println(req)
	if err != nil {
		c.logger.Warn(r.Context(), "failed to decode request body for company delete", "error", err)
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := c.delete.Execute(r.Context(), uint(req.IDCompany)); err != nil {
		c.logger.Warn(r.Context(), "delete company failed", "error", err, "id", req.IDCompany, "method", r.Method, "url", r.URL.Path)
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.logger.Info(r.Context(), "delete company successful", "id", req.IDCompany, "method", r.Method, "url", r.URL.Path)
	// RespondJSON(w, http.StatusOK, fmt.Sprintf("company %s deleted successful",idCompany)
}
