package service

import (
	"context"
	"errors"
	"testing"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/test/spy"
	"github.com/stretchr/testify/assert"
)

func TestRecoverPasswordInterface(t *testing.T) {
	ctx := context.Background()
	mockService := &spy.MockRecoverPassword{
		ExpectedError: nil,
	}
	input := dto.RecoverPasswordInput{
		Email: "test@example.com",
	}
	t.Run("should return nil", func(t *testing.T) {
		err := mockService.Execute(ctx, input)
		assert.NoError(t, err)
	})

	t.Run("should return an error", func(t *testing.T) {
		mockService = &spy.MockRecoverPassword{
			ExpectedError: errors.New("password recovery failed"),
		}
		err := mockService.Execute(ctx, input)
		assert.Error(t, err)
		assert.Equal(t, "password recovery failed", err.Error())
	})
}

