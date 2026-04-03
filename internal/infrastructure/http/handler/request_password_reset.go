// Package handler
package handler

import (
	"net/http"

	"github.com/paladignus/actajus/internal/application/command"
	"github.com/paladignus/actajus/internal/application/service"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type RequestPasswordReset struct {
	service service.RequestPasswordReset
	logger  repository.Logger
}

func NewRequestPasswordReset(
	service service.RequestPasswordReset,
	logger repository.Logger,
) RequestPasswordReset {
	return RequestPasswordReset{
		service,
		logger,
	}
}

func (rp RequestPasswordReset) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req, err := DecodeJSONRequest[command.RequestPasswordResetCommand](r)
	if err != nil {
		rp.logger.Error(ctx, "failed to decode request body for password reset", "error", err, "method", r.Method, "url", r.URL.Path)
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = rp.service.Execute(ctx, req)
	if err != nil {
		statusCode, errResponse := MapDomainErrorToHTTP(err)
		if statusCode >= 500 {
			rp.logger.Error(ctx, "recover password failed with server error", "error", err, "email", req.Email, "method", r.Method, "url", r.URL.Path)
		} else {
			rp.logger.Warn(ctx, "recover password failed", "error", err, "email", req.Email, "method", r.Method, "url", r.URL.Path)
		}
		RespondJSON(w, statusCode, errResponse)
		return
	}
	rp.logger.Info(ctx, "recover password successful", "email", req.Email, "method", r.Method, "url", r.URL.Path)
}
