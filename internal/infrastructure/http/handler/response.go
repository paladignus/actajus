// Package handler
package handler

import (
	"net/http"

	sharedhttp "github.com/paladignus/actajus/internal/shared/infrastructure/http/handler"
)

func DecodeJSONRequest[T any](r *http.Request) (T, error) {
	return sharedhttp.DecodeJSONRequest[T](r)
}

func RespondJSON[T any](w http.ResponseWriter, statusCode int, data T) error {
	return sharedhttp.RespondJSON(w, statusCode, data)
}

func RespondError(w http.ResponseWriter, statusCode int, message string) {
	sharedhttp.RespondError(w, statusCode, message)
}
