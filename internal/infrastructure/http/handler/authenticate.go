// Package handler
package handler

import (
	"net/http"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/application/service"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type Authenticate struct {
	service service.Authenticate
	logger  repository.Logger
}

func NewAuthenticate(service service.Authenticate, logger repository.Logger) Authenticate {
	return Authenticate{service, logger}
}

func (a Authenticate) Authenticate(w http.ResponseWriter, r *http.Request) {
	a.logger.Info(r.Context(), "processing sign_in request")
	req, err := DecodeJSONRequest[dto.AuthenticateInput](r)
	if err != nil {
		a.logger.Warn(r.Context(), "invalid request body", "error", err)
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := a.service.Authenticate(r.Context(), req)
	if err != nil {
		statusCode, errResponse := MapDomainErrorToHTTP(err)
		if statusCode >= 500 {
			a.logger.Error(r.Context(), "authentication failed with server error", "error", err, "cpf", req.CPF)
		} else {
			a.logger.Warn(r.Context(), "authentication failed", "error", err, "cpf", req.CPF, "status_code", statusCode)
		}
		RespondJSON(w, statusCode, errResponse)
		return
	}
	a.logger.Info(r.Context(), "authentication successful", "cpf", req.CPF)
	RespondJSON(w, http.StatusOK, resp)
}
