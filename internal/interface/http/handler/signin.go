// Package handler
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/interface/controller"
)

type SignIn struct {
	controller controller.SignIn
}

func NewSignIn(controller controller.SignIn) SignIn {
	return SignIn{controller}
}

func (s SignIn) SignIn(w http.ResponseWriter, r *http.Request) {
	var req dto.AuthenticatedInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := s.controller.SignIn(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
