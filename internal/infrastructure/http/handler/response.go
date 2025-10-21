// Package handler
package handler

import (
	"encoding/json"
	"net/http"
)

func DecodeJSONRequest[T any](r *http.Request) (T, error) {
	var req T
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func RespondJSON[T any](w http.ResponseWriter, statusCode int, data T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func RespondError(w http.ResponseWriter, statusCode int, message string) {
	http.Error(w, message, statusCode)
}
