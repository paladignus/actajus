package persistence

import (
	"context"
	"testing"

	"github.com/paladignus/actajus/test/utils"
)

func TestAuthenticationRepository(t *testing.T) {
	ctx := context.Background()
	cpf := utils.NewRandomString(10)
	password := utils.NewRandomString(10)
	sut := NewAuthenticateSpy()
	sut.SignIn(ctx, cpf, password)

	t.Run("should call SignIn once", func(t *testing.T) {
		if sut.CallCount != 1 {
			t.Errorf("expected call count to be 1, got %d", sut.CallCount)
		}
	})
}
