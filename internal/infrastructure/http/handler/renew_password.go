// Package handler
package handler

import (
	"net/http"

	"github.com/paladignus/actajus/internal/application/command"
	"github.com/paladignus/actajus/internal/application/service"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type RenewPassword struct {
	service service.RenewPassword
	logger  repository.Logger
}

func NewRenewPassword(
	service service.RenewPassword,
	logger repository.Logger,
) RenewPassword {
	return RenewPassword{
		service,
		logger,
	}
}

func (rp RenewPassword) RenewPassword(w http.ResponseWriter, r *http.Request) {
	req, err := DecodeJSONRequest[command.RenewPasswordCommand](r)
	if err != nil {
		rp.logger.Warn(r.Context(), "invalid request body", "error", err)
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = rp.service.Execute(r.Context(), req)
	if err != nil {
		statusCode, errResponse := MapDomainErrorToHTTP(err)
		if statusCode >= 500 {
			rp.logger.Error(r.Context(), "renew password failed with server error", "error", err, "cpf", req.CPF)
		} else {
			rp.logger.Warn(r.Context(), "renew password failed", "error", err, "cpf", req.CPF)
		}
		RespondJSON(w, statusCode, errResponse)
		return
	}
	rp.logger.Info(r.Context(), "renew password successful", "cpf", req.CPF)
}
