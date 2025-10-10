// Package handler
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/application/usecase"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type SignIn struct {
	usecase usecase.Authenticate
	logger  repository.Logger
}

func NewSignIn(usecase usecase.Authenticate, logger repository.Logger) SignIn {
	return SignIn{usecase, logger}
}

func (s SignIn) SignIn(w http.ResponseWriter, r *http.Request) {
	var req dto.AuthenticatedInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Warn(r.Context(), "invalid request body",
			"error", err,
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := s.usecase.Execute(r.Context(), req.CPF, req.Password)
	if err != nil {
		statusCode, errResponse := MapDomainErrorToHTTP(err)
		if statusCode >= 500 {
			s.logger.Error(r.Context(), "authentication failed with server error",
				"error", err,
				"cpf", req.CPF,
			)
		} else {
			s.logger.Warn(r.Context(), "authentication failed",
				"error", err,
				"cpf", req.CPF,
				"status_code", statusCode,
			)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(errResponse)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
