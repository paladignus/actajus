// Package handler
package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/paladignus/actajus/internal/module/person/application/command"
	"github.com/paladignus/actajus/internal/module/person/application/usecase"
	sharedAdapter "github.com/paladignus/actajus/internal/shared/application/adapter"
	"github.com/paladignus/actajus/internal/shared/domain"
	"github.com/paladignus/actajus/internal/shared/presentation/validation"
)

type PersonHTTPHandler struct {
	person usecase.CreatePerson
}

func NewPersonHTTPHandler(usecase usecase.CreatePerson) *PersonHTTPHandler {
	return &PersonHTTPHandler{usecase}
}

func (h *PersonHTTPHandler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var input command.CreatePersonCommand
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vs := validation.New().ValidateStruct(input)
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		var ve domain.ValidationError
		if errors.As(err, &ve) {
			_ = sharedAdapter.WriteValidationError(w, http.StatusBadRequest, ve)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	person, err := h.person.Execute(r.Context(), input)
	if err != nil {
		var ve domain.ValidationError
		if errors.As(err, &ve) {
			_ = sharedAdapter.WriteValidationError(w, http.StatusBadRequest, ve)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"id": person.ID})
}
