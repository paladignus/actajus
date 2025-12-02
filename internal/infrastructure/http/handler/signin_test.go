// Package handler
package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/test/spy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSignInService struct {
	expectedOutput dto.SignInOutput
	expectedError  error
}

func (m *MockSignInService) Execute(ctx context.Context, input dto.SignInInput) (dto.SignInOutput, error) {
	return m.expectedOutput, m.expectedError
}

func TestSignInsut(t *testing.T) {
	mockService := &MockSignInService{
		expectedOutput: dto.SignInOutput{
			IDUser:       "user-123",
			FirstName:    "John",
			LastName:     "Doe",
			Email:        "john.doe@example.com",
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
		},
		expectedError: nil,
	}
	logger := &spy.Logger{}
	t.Run("should successful request", func(t *testing.T) {
		sut := NewSignIn(mockService, logger)
		requestBody := `{"cpf": "12345678901", "password": "password123"}`
		req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		logger.On("Info", req.Context(), "signin successful", "cpf", "12345678901")
		sut.SignIn(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.SignInOutput
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.IDUser, "user-123")
		assert.Equal(t, response.Email, "john.doe@example.com")
	})

	t.Run("should an error for invalid request body", func(t *testing.T) {
		mockService := &MockSignInService{}
		sut := NewSignIn(mockService, logger)
		req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewBufferString("{invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		logger.On("Warn", req.Context(), "invalid request body", "error", mock.Anything)
		sut.SignIn(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("should returns error of user not found", func(t *testing.T) {
		mockService := &MockSignInService{
			expectedOutput: dto.SignInOutput{},
			expectedError:  exception.ErrUserNotFound,
		}
		logger := &spy.Logger{}
		sut := NewSignIn(mockService, logger)
		requestBody := `{"cpf": "12345678901", "password": "password123"}`
		req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		logger.On(
			"Warn",
			req.Context(),
			"signin failed",
			"error",
			mockService.expectedError,
			"cpf",
			"12345678901",
			"status_code",
			http.StatusNotFound,
		)
		sut.SignIn(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("service returns authentication error", func(t *testing.T) {
		mockService := &MockSignInService{
			expectedOutput: dto.SignInOutput{},
			expectedError:  fmt.Errorf("invalid credentials"),
		}
		logger := &spy.Logger{}
		sut := NewSignIn(mockService, logger)
		requestBody := `{"cpf": "12345678901", "password": "wrong-password"}`
		req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		logger.On("Error", req.Context(), "signin failed with server error", "error", mockService.expectedError, "cpf", "12345678901")
		sut.SignIn(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

