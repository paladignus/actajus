// Package handler
package handler

import (
	"context"
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
	t.Run("should initialize the constructor with its valid parameters", func(t *testing.T) {
		assert.NotNil(t, handler)
		assert.Equal(t, &spyUseCase, handler.service)
		assert.Equal(t, &spyLogger, handler.logger)
	})
}
