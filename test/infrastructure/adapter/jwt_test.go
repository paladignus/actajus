package adapter

import (
	"errors"
	"testing"
	"time"

	"github.com/paladignus/actajus/internal/infrastructure/adapter/spy"
	"github.com/paladignus/actajus/internal/infrastructure/config"
)

func TestJWT(t *testing.T) {
	config := config.Load()
	sut := spy.NewJWTAdapter(config.JWT.AccessSecret, config.JWT.AccessExpiry)
	t.Run("should create a new JWT adapter instance", func(t *testing.T) {
		if sut.SecretKey != config.JWT.AccessSecret || sut.ExpiresIn != config.JWT.AccessExpiry {
			t.Errorf("expected SecretKey %s and ExpiresIn %v, got SecretKey %s and ExpiresIn %v", sut.SecretKey, sut.ExpiresIn, config.JWT.AccessSecret, config.JWT.AccessExpiry)
		}
	})
	t.Run("should have correct expiration time", func(t *testing.T) {
		expectedExpiry := 15 * time.Minute
		if sut.ExpiresIn != expectedExpiry {
			t.Errorf("expected ExpiresIn %v, got %v", expectedExpiry, sut.ExpiresIn)
		}
	})
	t.Run("should call NewToken once", func(t *testing.T) {
		sut.NewToken()
		if sut.CallsCount != 1 {
			t.Errorf("expected calCount to be 1, got %d", sut.CallsCount)
		}
	})
	t.Run("should return error when ErrGenerateToken is set", func(t *testing.T) {
		sut.ErrGenerateToken = errors.New("error generating token")
		_, err := sut.NewToken()
		if err == nil {
			t.Error("expected error to be not nil")
		}
	})
	t.Run("should return valid token when no error", func(t *testing.T) {
		sut.ErrGenerateToken = nil
		token, err := sut.NewToken()
		if err != nil {
			t.Error("expected error to be nil, got", err)
		}
		if token != "valid_token" {
			t.Errorf("expected token to be 'valid_token', got %s", token)
		}
	})
}
