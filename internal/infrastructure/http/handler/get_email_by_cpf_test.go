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
	"github.com/paladignus/actajus/internal/application/service"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/internal/domain/repository"
	"github.com/paladignus/actajus/test/infrastructure/adapter/spy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type GetEmailByCPF struct {
	service service.GetEmailByCPF
	logger  repository.Logger
}

func NewGetEmailByCPF(service service.GetEmailByCPF, logger repository.Logger) GetEmailByCPF {
	return GetEmailByCPF{service, logger}
}

func (g *GetEmailByCPF) GetEmailByCPF(w http.ResponseWriter, r *http.Request) {
	g.logger.Info(r.Context(), "processing get_email_by_cpf request")
	req, err := DecodeJSONRequest[dto.GetEmailByCPFInput](r)
	if err != nil {
		g.logger.Warn(r.Context(), "invalid request body", "error", err)
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := g.service.Execute(r.Context(), req)
	if err != nil {
		statusCode, errResponse := MapDomainErrorToHTTP(err)
		if statusCode >= 500 {
			g.logger.Error(r.Context(), "get email failed with server error", "error", err, "cpf", req.CPF)
		} else {
			g.logger.Warn(r.Context(), "get email failed", "error", err, "cpf", req.CPF, "status_code", statusCode)
		}
		RespondJSON(w, statusCode, errResponse)
		return
	}
	g.logger.Info(r.Context(), "get email successful", "cpf", req.CPF)
	RespondJSON(w, http.StatusOK, resp)
}

type SpyGetEmailByCPFService struct {
	mock.Mock
}

func (s *SpyGetEmailByCPFService) Execute(ctx context.Context, input dto.GetEmailByCPFInput) (dto.GetEmailByCPFOutput, error) {
	args := s.Called(ctx, input)
	if args.Get(0) == nil {
		return dto.GetEmailByCPFOutput{}, args.Error(1)
	}
	return args.Get(0).(dto.GetEmailByCPFOutput), args.Error(1)
}

func TestGetEmailByCPF(t *testing.T) {
	sut := SpyGetEmailByCPFService{}
	spyLogger := spy.SpyLogger{}
	handler := NewGetEmailByCPF(&sut, &spyLogger)
	input := dto.GetEmailByCPFInput{
		CPF: "123.456.789-00",
	}
	expectedOutput := dto.GetEmailByCPFOutput{
		Email: "EmPdI@example.com",
	}
	t.Run("should initialize the constructor with its valid parameters", func(t *testing.T) {
		assert.NotNil(t, handler)
		assert.Equal(t, &sut, handler.service)
		assert.Equal(t, &spyLogger, handler.logger)
	})

	t.Run("should get email be successful", func(t *testing.T) {
		spyLogger.On("Info", mock.Anything, "processing get_email_by_cpf request").Once()
		sut.On("Execute", mock.Anything, input).Return(expectedOutput, nil).Once()
		spyLogger.On("Info", mock.Anything, "get email successful", "cpf", input.CPF).Once()
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/authenticate/email", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.GetEmailByCPF(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.GetEmailByCPFOutput
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput.Email, response.Email)
		sut.AssertExpectations(t)
		spyLogger.AssertExpectations(t)
	})

	t.Run("the request parameters should be valid", func(t *testing.T) {
		spyLogger.On("Info", mock.Anything, "processing get_email_by_cpf request", mock.Anything).Once()
		spyLogger.On("Warn", mock.Anything, "invalid request body", "error", mock.Anything).Once()
		req := httptest.NewRequest("POST", "/authenticate/email", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.GetEmailByCPF(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		sut.AssertExpectations(t)
		spyLogger.AssertExpectations(t)
	})

	t.Run("should return an HTTP error greater than or equal to 500", func(t *testing.T) {
		expectedErr := errors.New("database connection failed")
		spyLogger.On("Info", mock.Anything, "processing get_email_by_cpf request").Once()
		sut.On("Execute", mock.Anything, input).Return(nil, expectedErr).Once()
		spyLogger.On("Error", mock.Anything, "get email failed with server error",
			"error", expectedErr, "cpf", input.CPF).Once()
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/authenticate/email", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.GetEmailByCPF(w, req)
		assert.NotEqual(t, http.StatusOK, w.Code)
		sut.AssertExpectations(t)
		spyLogger.AssertExpectations(t)
	})

	t.Run("should return a user not found error with status code 404", func(t *testing.T) {
		expectedErr := exception.ErrEmailNotFound
		spyLogger.On("Info", mock.Anything, "processing get_email_by_cpf request").Once()
		sut.On("Execute", mock.Anything, input).Return(nil, expectedErr).Once()
		spyLogger.On("Warn", mock.Anything, "get email failed",
			"error", expectedErr, "cpf", input.CPF, "status_code", 404).Once()
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/authenticate/email", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.GetEmailByCPF(w, req)
		assert.NotEqual(t, http.StatusOK, w.Code)
		sut.AssertExpectations(t)
		spyLogger.AssertExpectations(t)
	})

	t.Run("should return a invalid cpf error with status code 404", func(t *testing.T) {
		expectedErr := exception.ErrInvalidCPF
		spyLogger.On("Info", mock.Anything, "processing get_email_by_cpf request").Once()
		sut.On("Execute", mock.Anything, input).Return(nil, expectedErr).Once()
		spyLogger.On("Warn", mock.Anything, "get email failed",
			"error", expectedErr, "cpf", input.CPF, "status_code", 400).Once()
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/authenticate/email", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.GetEmailByCPF(w, req)
		assert.NotEqual(t, http.StatusOK, w.Code)
		sut.AssertExpectations(t)
		spyLogger.AssertExpectations(t)
	})
}
