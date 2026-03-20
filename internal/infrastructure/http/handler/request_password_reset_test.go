// Package handler
package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/paladignus/actajus/internal/application/command"
	"github.com/paladignus/actajus/internal/application/readmodel"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/test/spy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRequestPasswordResetService struct {
	expectedError error
}

func (m *MockRequestPasswordResetService) Execute(ctx context.Context, input command.RequestPasswordResetCommand) error {
	return m.expectedError
}

func TestRequestPasswordResetHandler(t *testing.T) {
	mockService := &MockRequestPasswordResetService{
		expectedError: nil,
	}
	logger := &spy.Logger{}

	t.Run("should successful request", func(t *testing.T) {
		sut := NewRequestPasswordReset(mockService, logger)
		requestBody := `{}`
		req := httptest.NewRequest(http.MethodPost, "/request-password", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		logger.On("Info", req.Context(), "recover password successful", "email", "", "method", req.Method, "url", req.URL.Path)
		sut.RequestPasswordReset(w, req)
	})

	t.Run("should invalid request body", func(t *testing.T) {
		mockService := &MockRequestPasswordResetService{}
		sut := NewRequestPasswordReset(mockService, logger)
		req := httptest.NewRequest(http.MethodPost, "/request-password", bytes.NewBufferString("{invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		logger.On("Error", req.Context(), "failed to decode request body for password reset", "error", mock.Anything, "method", req.Method, "url", req.URL.Path)
		sut.RequestPasswordReset(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should returns error of internal server error", func(t *testing.T) {
		mockService := &MockRequestPasswordResetService{
			expectedError: fmt.Errorf("internal server error"),
		}
		logger := &spy.Logger{}
		sut := NewRequestPasswordReset(mockService, logger)
		requestBody := `{"email": "test@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/request-password", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		logger.On("Error", req.Context(), "recover password failed with server error", "error", mock.Anything, "email", "test@example.com", "method", req.Method, "url", req.URL.Path)
		sut.RequestPasswordReset(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should returns error of email not found", func(t *testing.T) {
		mockService := &MockRequestPasswordResetService{
			expectedError: exception.ErrEmailNotFound,
		}
		sut := NewRequestPasswordReset(mockService, logger)
		requestBody := `{"email": "nonexistent@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/request-password", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		logger.On(
			"Warn",
			req.Context(),
			"recover password failed",
			"error", mock.Anything,
			"email",
			"nonexistent@example.com",
			"method", req.Method,
			"url", req.URL.Path,
		)
		sut.RequestPasswordReset(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
