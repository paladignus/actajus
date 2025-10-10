// Package handler
package handler

import (
	"errors"
	"net/http"

	"github.com/paladignus/actajus/internal/domain"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func MapDomainErrorToHTTP(err error) (int, ErrorResponse) {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return http.StatusNotFound, ErrorResponse{
			Code:    "USER_NOT_FOUND",
			Message: "user not found",
		}
	case errors.Is(err, domain.ErrInvalidCredentials):
		return http.StatusUnauthorized, ErrorResponse{
			Code:    "INVALID_CREDENTIALS",
			Message: "invalid credentials",
		}
	case errors.Is(err, domain.ErrUserAlreadyExists):
		return http.StatusConflict, ErrorResponse{
			Code:    "USER_ALREADY_EXISTS",
			Message: "user already exists",
		}
	case errors.Is(err, domain.ErrUnauthorized):
		return http.StatusUnauthorized, ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "unauthorized access",
		}
	default:
		return http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "internal server error",
		}
	}
}
