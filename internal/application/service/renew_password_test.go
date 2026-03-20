// Package service
package service

import (
	"context"
	"errors"
	"testing"

	"github.com/paladignus/actajus/internal/application/command"
	"github.com/paladignus/actajus/test/spy"
	"github.com/stretchr/testify/assert"
)

func TestRenewPasswordInterface(t *testing.T) {
	ctx := context.Background()
	sut := &spy.RenewPassword{
		ExpectedError: nil,
	}
	input := command.RenewPasswordCommand{
		CPF:      "111.444.777-35",
		Password: "password",
		Token:    "token",
	}

	t.Run("should return nil", func(t *testing.T) {
		err := sut.Execute(ctx, input)
		assert.NoError(t, err)
	})

	t.Run("should return an error", func(t *testing.T) {
		sut = &spy.RenewPassword{
			ExpectedError: errors.New("renew password failed"),
		}
		err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.Equal(t, "renew password failed", err.Error())
	})
}
