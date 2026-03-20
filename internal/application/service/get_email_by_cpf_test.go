package service

import (
	"context"
	"errors"
	"testing"

	"github.com/paladignus/actajus/internal/application/command"
	"github.com/paladignus/actajus/internal/application/readmodel"
	"github.com/paladignus/actajus/test/spy"
	"github.com/stretchr/testify/assert"
)

func TestGetEmailByCPFInterface(t *testing.T) {
	ctx := context.Background()
	sut := &spy.GetEmailByCPF{
		ExpectedOutput: readmodel.GetEmailByCPFReadModel{
			Email: "test@example.com",
		},
		ExpectedError: nil,
	}
	input := command.GetEmailByCPFCommand{
		CPF: "12345678901",
	}
	t.Run("should return the same email", func(t *testing.T) {
		output, err := sut.Execute(ctx, input)
		assert.NoError(t, err)
		assert.Equal(t, "test@example.com", output.Email)
	})
	t.Run("should return an empty email", func(t *testing.T) {
		sut = &spy.GetEmailByCPF{
			ExpectedOutput: readmodel.GetEmailByCPFReadModel{},
			ExpectedError:  errors.New("test error"),
		}
		_, err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.Equal(t, "test error", err.Error())
	})
}
