// Package handler
package handler

import (
	"context"
	"net/http"
	"testing"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/application/service"
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

	t.Run("should initialize the constructor with its valid parameters", func(t *testing.T) {
		assert.NotNil(t, handler)
		assert.Equal(t, &sut, handler.service)
		assert.Equal(t, &spyLogger, handler.logger)
	})
}
