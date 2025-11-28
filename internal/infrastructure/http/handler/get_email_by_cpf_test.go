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
	"github.com/paladignus/actajus/internal/domain/repository"
)

// MockGetEmailByCPFService is a mock implementation of service.GetEmailByCPF for testing
type MockGetEmailByCPFService struct {
	expectedOutput dto.GetEmailByCPFOutput
	expectedError  error
}

func (m *MockGetEmailByCPFService) Execute(ctx context.Context, input dto.GetEmailByCPFInput) (dto.GetEmailByCPFOutput, error) {
	return m.expectedOutput, m.expectedError
}

// MockLogger is a mock implementation of repository.Logger for testing
type MockLogger struct{}

func (m *MockLogger) Debug(ctx context.Context, msg string, args ...interface{}) {}
func (m *MockLogger) Info(ctx context.Context, msg string, args ...interface{}) {}
func (m *MockLogger) Warn(ctx context.Context, msg string, args ...interface{}) {}
func (m *MockLogger) Error(ctx context.Context, msg string, args ...interface{}) {}
func (m *MockLogger) With(args ...interface{}) repository.Logger { return m }
func (m *MockLogger) WithError(err error) repository.Logger { return m }

// TestGetEmailByCPFHandler tests the GetEmailByCPF handler functionality
func TestGetEmailByCPFHandler(t *testing.T) {
	
	t.Run("successful request", func(t *testing.T) {
		mockService := &MockGetEmailByCPFService{
			expectedOutput: dto.GetEmailByCPFOutput{
				Email: "test@example.com",
			},
			expectedError: nil,
		}
		
		mockLogger := &MockLogger{}
		
		handler := NewGetEmailByCPF(mockService, mockLogger)
		
		// Create request with valid JSON
		requestBody := `{"cpf": "12345678901"}`
		req := httptest.NewRequest(http.MethodPost, "/get-email", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		
		handler.GetEmailByCPF(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
		
		// Parse response
		var response dto.GetEmailByCPFOutput
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		
		if response.Email != "test@example.com" {
			t.Errorf("Expected email 'test@example.com', got '%s'", response.Email)
		}
	})
	
	t.Run("invalid request body", func(t *testing.T) {
		mockService := &MockGetEmailByCPFService{}
		mockLogger := &MockLogger{}
		
		handler := NewGetEmailByCPF(mockService, mockLogger)
		
		// Create request with invalid JSON
		req := httptest.NewRequest(http.MethodPost, "/get-email", bytes.NewBufferString("{invalid json"))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		
		handler.GetEmailByCPF(w, req)
		
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
	
	t.Run("service returns error", func(t *testing.T) {
		mockService := &MockGetEmailByCPFService{
			expectedOutput: dto.GetEmailByCPFOutput{},
			expectedError:  fmt.Errorf("user not found"),
		}
		
		mockLogger := &MockLogger{}
		
		handler := NewGetEmailByCPF(mockService, mockLogger)
		
		// Create request with valid JSON
		requestBody := `{"cpf": "12345678901"}`
		req := httptest.NewRequest(http.MethodPost, "/get-email", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		
		handler.GetEmailByCPF(w, req)
		
		// The status code depends on the error mapping, but it should be a client error
		if w.Code < 400 || w.Code >= 500 {
			t.Errorf("Expected client error status code (4xx), got %d", w.Code)
		}
	})
}