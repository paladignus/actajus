// Package handler
package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLogger é um mock para repository.Logger
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {
	m.Called(ctx, msg, keysAndValues)
}

func (m *MockLogger) Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	m.Called(ctx, msg, keysAndValues)
}

func (m *MockLogger) Warn(ctx context.Context, msg string, keysAndValues ...interface{}) {
	m.Called(ctx, msg, keysAndValues)
}

func (m *MockLogger) Error(ctx context.Context, msg string, keysAndValues ...interface{}) {
	m.Called(ctx, msg, keysAndValues)
}

func (m *MockLogger) With(keysAndValues ...interface{}) repository.Logger {
	return m.With(keysAndValues...)
}

func (m *MockLogger) WithError(err error) repository.Logger {
	return m.WithError(err)
}

// MockAuthenticateService é um mock para usecase.Authenticate
type MockAuthenticateService struct {
	mock.Mock
}

func (m *MockAuthenticateService) Authenticate(ctx context.Context, input *dto.AuthenticateInput) (*dto.AuthenticateOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AuthenticateOutput), args.Error(1)
}

func TestNewAuthenticate(t *testing.T) {
	mockUseCase := new(MockAuthenticateService)
	mockLogger := new(MockLogger)

	handler := NewAuthenticate(mockUseCase, mockLogger)

	assert.NotNil(t, handler)
	assert.Equal(t, mockUseCase, handler.usecase)
	assert.Equal(t, mockLogger, handler.logger)
}

func TestAuthenticate_Success(t *testing.T) {
	mockUseCase := new(MockAuthenticateService)
	mockLogger := new(MockLogger)

	// Setup expectations
	expectedInput := &dto.AuthenticateInput{CPF: "12345678900"}
	expectedOutput := &dto.AuthenticateOutput{
		Token:     "jwt-token-here",
		ExpiresIn: 3600,
	}

	mockLogger.On("Info", mock.Anything, "processing sign_in request", mock.Anything).Once()
	mockLogger.On("Info", mock.Anything, "authentication successful", "cpf", "12345678900").Once()
	mockUseCase.On("Execute", mock.Anything, expectedInput).Return(expectedOutput, nil)

	handler := Authenticate{
		usecase: mockUseCase,
		logger:  mockLogger,
	}

	// Create request body
	requestBody := map[string]string{
		"cpf": "12345678900",
	}
	body, _ := json.Marshal(requestBody)

	// Create HTTP request
	req := httptest.NewRequest("POST", "/authenticate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute handler
	handler.Authenticate(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.AuthenticateOutput
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput.Token, response.Token)
	assert.Equal(t, expectedOutput.ExpiresIn, response.ExpiresIn)

	// Verify mock expectations
	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestAuthenticate_InvalidJSON(t *testing.T) {
	mockUseCase := new(MockAuthenticateService)
	mockLogger := new(MockLogger)

	// Setup expectations
	mockLogger.On("Info", mock.Anything, "processing sign_in request", mock.Anything).Once()
	mockLogger.On("Warn", mock.Anything, "invalid request body", "error", mock.Anything).Once()

	handler := Authenticate{
		usecase: mockUseCase,
		logger:  mockLogger,
	}

	// Create invalid JSON request
	req := httptest.NewRequest("POST", "/authenticate", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute handler
	handler.Authenticate(w, req)

	// Verify response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Verify mock expectations
	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestAuthenticate_UseCaseError(t *testing.T) {
	mockUseCase := new(MockAuthenticateService)
	mockLogger := new(MockLogger)

	// Setup expectations
	expectedInput := &dto.AuthenticateInput{CPF: "12345678900"}
	expectedError := errors.New("invalid credentials")

	mockLogger.On("Info", mock.Anything, "processing sign_in request", mock.Anything).Once()
	mockLogger.On("Warn", mock.Anything, "authentication failed", "error", expectedError, "cpf", "12345678900", "status_code", http.StatusUnauthorized).Once()
	mockUseCase.On("Execute", mock.Anything, expectedInput).Return(nil, expectedError)

	handler := Authenticate{
		usecase: mockUseCase,
		logger:  mockLogger,
	}

	// Create request body
	requestBody := map[string]string{
		"cpf": "12345678900",
	}
	body, _ := json.Marshal(requestBody)

	// Create HTTP request
	req := httptest.NewRequest("POST", "/authenticate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute handler
	handler.Authenticate(w, req)

	// Verify response
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Verify mock expectations
	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestAuthenticate_UseCaseServerError(t *testing.T) {
	mockUseCase := new(MockAuthenticateService)
	mockLogger := new(MockLogger)

	// Setup expectations
	expectedInput := &dto.AuthenticateInput{CPF: "12345678900"}
	expectedError := errors.New("database connection failed")

	mockLogger.On("Info", mock.Anything, "processing sign_in request", mock.Anything).Once()
	mockLogger.On("Error", mock.Anything, "authentication failed with server error", "error", expectedError, "cpf", "12345678900").Once()
	mockUseCase.On("Execute", mock.Anything, expectedInput).Return(nil, expectedError)

	handler := Authenticate{
		usecase: mockUseCase,
		logger:  mockLogger,
	}

	// Create request body
	requestBody := map[string]string{
		"cpf": "12345678900",
	}
	body, _ := json.Marshal(requestBody)

	// Create HTTP request
	req := httptest.NewRequest("POST", "/authenticate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute handler
	handler.Authenticate(w, req)

	// Verify response - assuming MapDomainErrorToHTTP returns 500 for server errors
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Verify mock expectations
	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestAuthenticate_MissingCPF(t *testing.T) {
	mockUseCase := new(MockAuthenticateService)
	mockLogger := new(MockLogger)

	// Setup expectations
	mockLogger.On("Info", mock.Anything, "processing sign_in request", mock.Anything).Once()
	mockLogger.On("Warn", mock.Anything, "invalid request body", "error", mock.Anything).Once()

	handler := Authenticate{
		usecase: mockUseCase,
		logger:  mockLogger,
	}

	// Create request with missing CPF
	requestBody := map[string]string{}
	body, _ := json.Marshal(requestBody)

	// Create HTTP request
	req := httptest.NewRequest("POST", "/authenticate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute handler
	handler.Authenticate(w, req)

	// Verify response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Verify mock expectations
	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
