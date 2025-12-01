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
		g.logger.Warn(r.Context(), "invalid request body", "error", err)
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := g.service.Execute(r.Context(), req)
	if err != nil {
		statusCode, errResponse := MapDomainErrorToHTTP(err)
		if statusCode >= 500 {
			g.logger.Error(r.Context(), "get email failed with server error", "error", err, "cpf", req.CPF)
		} else {
			g.logger.Warn(r.Context(), "get email failed", "error", err, "cpf", req.CPF, "status_code", statusCode)
		}
		RespondJSON(w, statusCode, errResponse)
		return
	}
	g.logger.Info(r.Context(), "get email successful", "cpf", req.CPF)
	RespondJSON(w, http.StatusOK, resp)
}
