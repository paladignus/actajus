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

type MockRenewPasswordService struct {
	expectedError error
}

func (m *MockRenewPasswordService) Execute(ctx context.Context, input dto.RenewPasswordInput) error {
	return m.expectedError
}

func TestRecoverPasswordHandler(t *testing.T) {
	mockService := &MockRenewPasswordService{
		expectedError: nil,
	}
	logger := &spy.Logger{}
	t.Run("should successful request", func(t *testing.T) {
		sut := NewRenewPassword(mockService, logger)
		requestBody := `{"email": "test@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/recover-password", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		logger.On("Info", req.Context(), "renew password successful", "cpf", "")
		sut.RenewPassword(w, req)
	})

	t.Run("should invalid request body", func(t *testing.T) {
		mockService := &MockRenewPasswordService{}
		logger := &spy.Logger{}
		sut := NewRenewPassword(mockService, logger)
		req := httptest.NewRequest(http.MethodPost, "/recover-password", bytes.NewBufferString("{invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		logger.On("Warn", req.Context(), "invalid request body", "error", mock.Anything)
		sut.RenewPassword(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should returns error of internal server error", func(t *testing.T) {
		mockService := &MockRenewPasswordService{
			expectedError: fmt.Errorf("internal server error"),
		}
		logger := &spy.Logger{}
		sut := NewRenewPassword(mockService, logger)
		requestBody := `{"email": "test@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/recover-password", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		logger.On("Error", req.Context(), "renew password failed with server error", "error", mock.Anything, "cpf", "")
		sut.RenewPassword(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should returns error of email not found", func(t *testing.T) {
		mockService := &MockRenewPasswordService{
			expectedError: exception.ErrEmailNotFound,
		}
		logger := &spy.Logger{}
		sut := NewRenewPassword(mockService, logger)
		requestBody := `{"email": "nonexistent@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/recover-password", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		logger.On(
			"Warn",
			req.Context(),
			"renew password failed",
			"error", mock.Anything,
			"cpf",
			"",
		)
		sut.RenewPassword(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
