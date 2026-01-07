// Package handler
package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/stretchr/testify/assert"
)

func TestMapDomainErrorToHTTP(t *testing.T) {
	tests := []struct {
		name               string
		inputError         error
		expectedStatusCode int
		expectedResponse   ErrorResponse
	}{
		{
			name:               "should return status code 404 and error response for user not found",
			inputError:         exception.ErrUserNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: ErrorResponse{
				Code:    "USER_NOT_FOUND",
				Message: "user not found",
			},
		},
		{
			name:               "should return status code 401 and error response for invalid credentials",
			inputError:         exception.ErrInvalidCredentials,
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse: ErrorResponse{
				Code:    "INVALID_CREDENTIALS",
				Message: "invalid credentials",
			},
		},
		{
			name:               "should return status code 409 and error response for user already exists",
			inputError:         exception.ErrUserAlreadyExists,
			expectedStatusCode: http.StatusConflict,
			expectedResponse: ErrorResponse{
				Code:    "USER_ALREADY_EXISTS",
				Message: "user already exists",
			},
		},
		{
			name:               "should return status code 401 and error response for unauthorized",
			inputError:         exception.ErrUnauthorized,
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse: ErrorResponse{
				Code:    "UNAUTHORIZED",
				Message: "unauthorized access",
			},
		},
		{
			name:               "should return status code 400 and error response for invalid cpf",
			inputError:         exception.ErrInvalidCPF,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: ErrorResponse{
				Code:    "INVALID_CPF",
				Message: "invalid cpf",
			},
		},
		{
			name:               "should return status code 400 and error response for invalid email",
			inputError:         exception.ErrInvalidEmail,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: ErrorResponse{
				Code:    "INVALID_EMAIL",
				Message: "invalid email",
			},
		},
		{
			name:               "should return status code 500 and error response for internal error",
			inputError:         errors.New("internal server error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: ErrorResponse{
				Code:    "INTERNAL_ERROR",
				Message: "internal server error",
			},
		},
		{
			name:               "should return status code 404 and error response for email not found",
			inputError:         exception.ErrEmailNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: ErrorResponse{
				Code:    "EMAIL_NOT_FOUND",
				Message: "email not found",
			},
		},
		{
			name:               "should return status code 409 and error response for email already exists",
			inputError:         exception.ErrEmailAlreadyExists,
			expectedStatusCode: http.StatusConflict,
			expectedResponse: ErrorResponse{
				Code:    "EMAIL_ALREADY_EXISTS",
				Message: "email already exists",
			},
		},
		{
			name:               "should return status code 400 and error response for invalid password",
			inputError:         exception.ErrInvalidPassword,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: ErrorResponse{
				Code:    "INVALID_PASSWORD",
				Message: "invalid password",
			},
		},
		{
			name:               "should return status code 400 and error response for invalid token",
			inputError:         exception.ErrInvalidToken,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: ErrorResponse{
				Code:    "INVALID_TOKEN",
				Message: "invalid token",
			},
		},
		{
			name:               "should return status code 401 and error response for token expired",
			inputError:         exception.ErrTokenExpired,
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse: ErrorResponse{
				Code:    "TOKEN_EXPIRED",
				Message: "token expired",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, response := MapDomainErrorToHTTP(tt.inputError)
			assert.Equal(t, tt.expectedStatusCode, statusCode)
			assert.Equal(t, tt.expectedResponse, response)
		})
	}
}
