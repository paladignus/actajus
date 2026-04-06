// Package adapter centralizes presentation-layer mappings.
package adapter

import (
	"encoding/json"
	"net/http"

	"github.com/paladignus/actajus/internal/shared/domain"
)

func WriteValidationError(w http.ResponseWriter, status int, err domain.ValidationError) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(ValidationPayload(err))
}
