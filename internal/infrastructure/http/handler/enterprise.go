// Package handler
package handler

import (
	"net/http"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/application/service"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type Enterprise struct {
	create service.IEnterpriseCreate
	update service.IEnterpriseUpdate
	getAll service.IEnterpriseGetAll
	logger repository.Logger
}

func NewEnterprise(
	create service.IEnterpriseCreate,
	update service.IEnterpriseUpdate,
	getAll service.IEnterpriseGetAll,
	logger repository.Logger,
) Enterprise {
	return Enterprise{
		create,
		update,
		getAll,
		logger,
	}
}

func (e Enterprise) Create(w http.ResponseWriter, r *http.Request) {
	req, err := DecodeJSONRequest[dto.EnterpriseInput](r)
	if err != nil {
		e.logger.Warn(r.Context(), "failed to decode request body for enterprise create", "error", err)
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = e.create.Execute(r.Context(), req)
	if err != nil {
		statusCode, errResponse := MapDomainErrorToHTTP(err)
		if statusCode >= 500 {
			e.logger.Error(r.Context(), "create enterprise failed with server error", "error", err, "cnpj", req.CNPJ, "method", r.Method, "url", r.URL.Path)
		} else {
			e.logger.Warn(r.Context(), "create enterprise failed", "error", err, "cnpj", req.CNPJ, "method", r.Method, "url", r.URL.Path)
		}
		RespondJSON(w, statusCode, errResponse)
		return
	}
	e.logger.Info(r.Context(), "create enterprise successful", "cnpj", req.CNPJ, "registered_by", req.RegisteredBy, "method", r.Method, "url", r.URL.Path)
}

func (e Enterprise) Update(w http.ResponseWriter, r *http.Request) {
	req, err := DecodeJSONRequest[dto.EnterpriseInput](r)
	if err != nil {
		e.logger.Warn(r.Context(), "failed to decode request body for enterprise update", "error", err)
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err = e.update.Execute(r.Context(), req); err != nil {
		statusCode, errResponse := MapDomainErrorToHTTP(err)
		if statusCode >= 500 {
			e.logger.Error(r.Context(), "update enterprise failed with server error", "error", err, "cnpj", req.CNPJ, "method", r.Method, "url", r.URL.Path)
		} else {
			e.logger.Warn(r.Context(), "update enterprise failed", "error", err, "cnpj", req.CNPJ, "method", r.Method, "url", r.URL.Path)
		}
		RespondJSON(w, statusCode, errResponse)
		return
	}
	e.logger.Info(r.Context(), "update enterprise successful", "cnpj", req.CNPJ, "registered_by", req.RegisteredBy, "method", r.Method, "url", r.URL.Path)
}

func (e Enterprise) GetAll(w http.ResponseWriter, r *http.Request) {
	enterprises, err := e.getAll.Execute(r.Context())
	if err != nil {
		e.logger.Warn(r.Context(), "get all enterprise failed", "error", err, "method", r.Method, "url", r.URL.Path)
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	e.logger.Info(r.Context(), "get all enterprise successful", "method", r.Method, "url", r.URL.Path)
	RespondJSON(w, http.StatusOK, enterprises)
}
