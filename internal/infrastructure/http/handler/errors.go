// Package handler
package handler

import (
	"errors"
	"net/http"

	"github.com/paladignus/actajus/internal/domain/exception"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func MapDomainErrorToHTTP(err error) (int, ErrorResponse) {
	switch {
	case errors.Is(err, exception.ErrUserNotFound):
		return http.StatusNotFound, ErrorResponse{
			Code:    "USER_NOT_FOUND",
			Message: "user not found",
		}
	case errors.Is(err, exception.ErrInvalidCredentials):
		return http.StatusUnauthorized, ErrorResponse{
			Code:    "INVALID_CREDENTIALS",
			Message: "invalid credentials",
		}
	case errors.Is(err, exception.ErrUserAlreadyExists):
		return http.StatusConflict, ErrorResponse{
			Code:    "USER_ALREADY_EXISTS",
			Message: "user already exists",
		}
	case errors.Is(err, exception.ErrUnauthorized):
		return http.StatusUnauthorized, ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "unauthorized access",
		}
	case errors.Is(err, exception.ErrInvalidCPF):
		return http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_CPF",
			Message: "invalid cpf",
		}
	case errors.Is(err, exception.ErrInvalidEmail):
		return http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_EMAIL",
			Message: "invalid email",
		}
	case errors.Is(err, exception.ErrEmailNotFound):
		return http.StatusNotFound, ErrorResponse{
			Code:    "EMAIL_NOT_FOUND",
			Message: "email not found",
		}
	case errors.Is(err, exception.ErrEmailAlreadyExists):
		return http.StatusConflict, ErrorResponse{
			Code:    "EMAIL_ALREADY_EXISTS",
			Message: "email already exists",
		}
	case errors.Is(err, exception.ErrInvalidPassword):
		return http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_PASSWORD",
			Message: "invalid password",
		}
	case errors.Is(err, exception.ErrInvalidToken):
		return http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_TOKEN",
			Message: "invalid token",
		}
	default:
		return http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "internal server error",
		}
	}
}
