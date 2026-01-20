// Package handler
package handler

import (
	"net/http"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/application/service"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type GetEmailByCPF struct {
	service service.GetEmailByCPF
	logger  repository.Logger
}

func NewGetEmailByCPF(service service.GetEmailByCPF, logger repository.Logger) GetEmailByCPF {
	return GetEmailByCPF{service, logger}
}

func (g GetEmailByCPF) GetEmailByCPF(w http.ResponseWriter, r *http.Request) {
	req, err := DecodeJSONRequest[dto.GetEmailByCPFInput](r)
	if err != nil {
		g.logger.Warn(r.Context(), "failed to decode request body for get email by CPF", "error", err, "method", r.Method, "url", r.URL.Path)
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := g.service.Execute(r.Context(), req)
	if err != nil {
		statusCode, errResponse := MapDomainErrorToHTTP(err)
		if statusCode >= 500 {
			g.logger.Error(r.Context(), "get email by CPF failed with server error", "error", err, "cpf", req.CPF, "method", r.Method, "url", r.URL.Path)
		} else {
			g.logger.Warn(r.Context(), "get email by CPF failed", "error", err, "cpf", req.CPF, "method", r.Method, "url", r.URL.Path)
		}
		RespondJSON(w, statusCode, errResponse)
		return
	}
	g.logger.Info(r.Context(), "get email by CPF successful", "cpf", req.CPF, "email", resp.Email, "method", r.Method, "url", r.URL.Path)
	RespondJSON(w, http.StatusOK, resp)
}
