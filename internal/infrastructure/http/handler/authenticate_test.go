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
	"github.com/paladignus/actajus/test/infrastructure/adapter/spy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type SpyAuthenticateService struct {
	mock.Mock
}

func (s *SpyAuthenticateService) Authenticate(ctx context.Context, input dto.AuthenticateInput) (dto.AuthenticateOutput, error) {
	args := s.Called(ctx, input)
	if args.Get(0) == nil {
		return dto.AuthenticateOutput{}, args.Error(1)
	}
	return args.Get(0).(dto.AuthenticateOutput), args.Error(1)
}

func TestAuthenticateHandler(t *testing.T) {
	spyUseCase := SpyAuthenticateService{}
	spyLogger := spy.SpyLogger{}
	handler := NewAuthenticate(&spyUseCase, &spyLogger)
	input := dto.AuthenticateInput{
		CPF: "123.456.789-00",
	}
	expectedOutput := dto.AuthenticateOutput{
		AccessToken:  "AccessToken",
		RefreshToken: "RefreshToken",
		IDUser:       "user-123",
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "EmPdI@example.com",
		Roles:        []string{"manager"},
	}
	t.Run("should initialize the constructor with its valid parameters", func(t *testing.T) {
		assert.NotNil(t, handler)
		assert.Equal(t, &spyUseCase, handler.service)
		assert.Equal(t, &spyLogger, handler.logger)
	})

	t.Run("should authentication be successful", func(t *testing.T) {
		spyLogger.On("Info", mock.Anything, "processing sign_in request").Once()
		spyUseCase.On("Authenticate", mock.Anything, input).Return(expectedOutput, nil).Once()
		spyLogger.On("Info", mock.Anything, "authentication successful", "cpf", input.CPF).Once()
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/authenticate", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.Authenticate(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.AuthenticateOutput
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput.AccessToken, response.AccessToken)
		assert.Equal(t, expectedOutput.RefreshToken, response.RefreshToken)
		spyUseCase.AssertExpectations(t)
		spyLogger.AssertExpectations(t)
	})

	t.Run("the request parameters should be valid", func(t *testing.T) {
		spyLogger.On("Info", mock.Anything, "processing sign_in request", mock.Anything).Once()
		spyLogger.On("Warn", mock.Anything, "invalid request body", "error", mock.Anything).Once()
		req := httptest.NewRequest("POST", "/authenticate", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.Authenticate(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		spyUseCase.AssertExpectations(t)
		spyLogger.AssertExpectations(t)
	})

	t.Run("should return an HTTP error greater than or equal to 500", func(t *testing.T) {
		expectedErr := errors.New("database connection failed")
		spyLogger.On("Info", mock.Anything, "processing sign_in request").Once()
		spyUseCase.On("Authenticate", mock.Anything, input).Return(nil, expectedErr).Once()
		spyLogger.On("Error", mock.Anything, "authentication failed with server error",
			"error", expectedErr, "cpf", input.CPF).Once()
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/authenticate", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.Authenticate(w, req)
		assert.NotEqual(t, http.StatusOK, w.Code)
		spyUseCase.AssertExpectations(t)
		spyLogger.AssertExpectations(t)
	})
}
