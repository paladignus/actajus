// Package handler
package handler

import (
	"bytes"
	"context"
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

// MockRecoverPasswordService is a mock implementation of service.RecoverPassword for testing
type MockRecoverPasswordService struct {
	expectedError error
}

func (m *MockRecoverPasswordService) Execute(ctx context.Context, input dto.RecoverPasswordInput) error {
	return m.expectedError
}

// TestRecoverPasswordHandler tests the RecoverPassword handler functionality
func TestRecoverPasswordHandler(t *testing.T) {
	mockService := &MockRecoverPasswordService{
		expectedError: nil,
	}
	logger := &spy.Logger{}
	t.Run("should successful request", func(t *testing.T) {
		sut := NewRecoverPassword(mockService, logger)
		requestBody := `{"email": "test@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/recover-password", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		logger.On("Info", req.Context(), "recover password successful", "email", "test@example.com")
		sut.RecoverPassword(w, req)
	})

	t.Run("should invalid request body", func(t *testing.T) {
		mockService := &MockRecoverPasswordService{}
		logger := &spy.Logger{}
		sut := NewRecoverPassword(mockService, logger)
		req := httptest.NewRequest(http.MethodPost, "/recover-password", bytes.NewBufferString("{invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		logger.On("Error", req.Context(), "failed to decode request body", "error", mock.Anything)
		sut.RecoverPassword(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should returns error of internal server error", func(t *testing.T) {
		mockService := &MockRecoverPasswordService{
			expectedError: fmt.Errorf("internal server error"),
		}
		logger := &spy.Logger{}
		sut := NewRecoverPassword(mockService, logger)
		requestBody := `{"email": "test@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/recover-password", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		logger.On("Error", req.Context(), "recover password failed with server error", "error", mock.Anything, "email", "test@example.com")
		sut.RecoverPassword(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should returns error of email not found", func(t *testing.T) {
		mockService := &MockRecoverPasswordService{
			expectedError: exception.ErrEmailNotFound,
		}
		logger := &spy.Logger{}
		sut := NewRecoverPassword(mockService, logger)
		requestBody := `{"email": "nonexistent@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/recover-password", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		logger.On(
			"Warn",
			req.Context(),
			"recover password failed",
			"error", mock.Anything,
			"email", "nonexistent@example.com",
			"status_code",
			http.StatusNotFound,
		)
		sut.RecoverPassword(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
