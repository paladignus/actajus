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

type MockGetEmailByCPFService struct {
	expectedOutput dto.GetEmailByCPFOutput
	expectedError  error
}

func (m *MockGetEmailByCPFService) Execute(ctx context.Context, input dto.GetEmailByCPFInput) (dto.GetEmailByCPFOutput, error) {
	return m.expectedOutput, m.expectedError
}

func TestGetEmailByCPFHandler(t *testing.T) {
	mockService := &MockGetEmailByCPFService{
		expectedOutput: dto.GetEmailByCPFOutput{
			Email: "test@example.com",
		},
		expectedError: nil,
	}
	logger := &spy.Logger{}

	t.Run("successful request", func(t *testing.T) {
		sut := NewGetEmailByCPF(mockService, logger)
		requestBody := `{"cpf": "12345678901"}`
		req := httptest.NewRequest(http.MethodPost, "/get-email", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		logger.On("Info", req.Context(), "get email successful", "cpf", "12345678901")
		sut.GetEmailByCPF(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.GetEmailByCPFOutput
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Email, "test@example.com")
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockService := &MockGetEmailByCPFService{}
		sut := NewGetEmailByCPF(mockService, logger)
		req := httptest.NewRequest(http.MethodPost, "/get-email", bytes.NewBufferString("{invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		logger.On("Warn", req.Context(), "invalid request body", "error", mock.Anything)
		sut.GetEmailByCPF(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service return internal server error", func(t *testing.T) {
		mockService := &MockGetEmailByCPFService{
			expectedOutput: dto.GetEmailByCPFOutput{},
			expectedError:  fmt.Errorf("internal server error"),
		}
		sut := NewGetEmailByCPF(mockService, logger)
		requestBody := `{"cpf": "12345678901"}`
		req := httptest.NewRequest(http.MethodPost, "/get-email", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		logger.On("Error", req.Context(), "get email failed with server error", "error", mockService.expectedError, "cpf", "12345678901")
		sut.GetEmailByCPF(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return 404 if email not found", func(t *testing.T) {
		mockService := &MockGetEmailByCPFService{
			expectedOutput: dto.GetEmailByCPFOutput{},
			expectedError:  exception.ErrEmailNotFound,
		}
		sut := NewGetEmailByCPF(mockService, logger)
		requestBody := `{"cpf": "22345678901"}`
		req := httptest.NewRequest(http.MethodPost, "/get-email", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		logger.On("Warn",
			req.Context(),
			"get email failed", "error",
			mockService.expectedError,
			"cpf",
			"22345678901",
			"status_code",
			http.StatusNotFound)
		sut.GetEmailByCPF(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
