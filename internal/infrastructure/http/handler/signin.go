// Package handler
package handler

import (
	"net/http"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/application/service"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type SignIn struct {
	service service.SignIn
	logger  repository.Logger
}

func NewSignIn(service service.SignIn, logger repository.Logger) SignIn {
	return SignIn{service, logger}
}

func (a SignIn) SignIn(w http.ResponseWriter, r *http.Request) {
	a.logger.Info(r.Context(), "processing sign_in request")
	req, err := DecodeJSONRequest[dto.SignInInput](r)
	if err != nil {
		a.logger.Warn(r.Context(), "invalid request body", "error", err)
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := a.service.Execute(r.Context(), req)
	if err != nil {
		statusCode, errResponse := MapDomainErrorToHTTP(err)
		if statusCode >= 500 {
			a.logger.Error(r.Context(), "signin failed with server error", "error", err, "cpf", req.CPF)
		} else {
			a.logger.Warn(r.Context(), "signin failed", "error", err, "cpf", req.CPF, "status_code", statusCode)
		}
		RespondJSON(w, statusCode, errResponse)
		return
	}
	a.logger.Info(r.Context(), "signin successful", "cpf", req.CPF)
	RespondJSON(w, http.StatusOK, resp)
}
