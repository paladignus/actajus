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
)

// MockSignInService is a mock implementation of service.SignIn for testing
type MockSignInService struct {
	expectedOutput dto.SignInOutput
	expectedError  error
}

func (m *MockSignInService) Execute(ctx context.Context, input dto.SignInInput) (dto.SignInOutput, error) {
	return m.expectedOutput, m.expectedError
}

// TestSignInHandler tests the SignIn handler functionality
func TestSignInHandler(t *testing.T) {
	
	t.Run("successful request", func(t *testing.T) {
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
		
		mockLogger := &MockLogger{}
		
		handler := NewSignIn(mockService, mockLogger)
		
		// Create request with valid JSON
		requestBody := `{"cpf": "12345678901", "password": "password123"}`
		req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		
		handler.SignIn(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
		
		// Parse response
		var response dto.SignInOutput
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		
		if response.IDUser != "user-123" {
			t.Errorf("Expected IDUser 'user-123', got '%s'", response.IDUser)
		}
		
		if response.Email != "john.doe@example.com" {
			t.Errorf("Expected email 'john.doe@example.com', got '%s'", response.Email)
		}
	})
	
	t.Run("invalid request body", func(t *testing.T) {
		mockService := &MockSignInService{}
		mockLogger := &MockLogger{}
		
		handler := NewSignIn(mockService, mockLogger)
		
		// Create request with invalid JSON
		req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewBufferString("{invalid json"))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		
		handler.SignIn(w, req)
		
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
	
	t.Run("service returns error", func(t *testing.T) {
		mockService := &MockSignInService{
			expectedOutput: dto.SignInOutput{},
			expectedError:  fmt.Errorf("user not found"),
		}
		
		mockLogger := &MockLogger{}
		
		handler := NewSignIn(mockService, mockLogger)
		
		// Create request with valid JSON
		requestBody := `{"cpf": "12345678901", "password": "password123"}`
		req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		
		handler.SignIn(w, req)
		
		// The status code depends on the error mapping, but it should be a client error
		if w.Code < 400 || w.Code >= 500 {
			t.Errorf("Expected client error status code (4xx), got %d", w.Code)
		}
	})
	
	t.Run("service returns authentication error", func(t *testing.T) {
		mockService := &MockSignInService{
			expectedOutput: dto.SignInOutput{},
			expectedError:  fmt.Errorf("invalid credentials"),
		}
		
		mockLogger := &MockLogger{}
		
		handler := NewSignIn(mockService, mockLogger)
		
		// Create request with valid JSON
		requestBody := `{"cpf": "12345678901", "password": "wrong-password"}`
		req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		
		handler.SignIn(w, req)
		
		// Should return client error, not server error
		if w.Code >= 500 {
			t.Errorf("Expected client error status code, got %d", w.Code)
		}
	})
}