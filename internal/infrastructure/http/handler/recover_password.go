// Package handler
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/application/service"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type RecoverPassword struct {
	service service.RecoverPassword
	logger  repository.Logger
}

func NewRecoverPassword(service service.RecoverPassword, logger repository.Logger) RecoverPassword {
	return RecoverPassword{service, logger}
}

func (rp RecoverPassword) RecoverPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.RecoverPasswordInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		rp.logger.Error(ctx, "failed to decode request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	output, err := rp.service.Execute(ctx, req)
	// const statusCode = http.StatusOK
	if err != nil {
		statusCode, errResponse := MapDomainErrorToHTTP(err)
		if statusCode >= 500 {
			rp.logger.Error(ctx, "recover password failed with server error", "error", err, "email", req.Email)
		} else {
			rp.logger.Warn(ctx, "recover password failed", "error", err, "email", req.Email, "status_code", statusCode)
		}
		RespondJSON(w, statusCode, errResponse)
		return
	}
	rp.logger.Info(ctx, "recover password successful", "email", req.Email)
	RespondJSON(w, http.StatusOK, output)
}
