package persistence

import (
	"errors"
	"testing"
	"time"

	"github.com/paladignus/actajus/test/infrastructure/persistence/spy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestToken_SaveRevokedToken(t *testing.T) {
	t.Run("success - should saves token with valid expiration", func(t *testing.T) {
		mockCache := new(spy.MockCache)
		tokenRepo := Token{repository: mockCache}
		token := "valid-jwt-token"
		expiresAt := time.Now().Add(1 * time.Hour)
		expectedKey := "revoked_token:valid-jwt-token"
		mockCache.On("Set", expectedKey, "1", mock.MatchedBy(func(ttl time.Duration) bool {
			// Aceita TTL com margem de 1 segundo devido ao tempo de execução
			return ttl > 59*time.Minute && ttl <= 1*time.Hour
		})).Return(nil)
		err := tokenRepo.SaveRevokedToken(token, expiresAt)
		assert.NoError(t, err)
		mockCache.AssertExpectations(t)
	})

	t.Run("success - should does not save token with past expiration", func(t *testing.T) {
		mockCache := new(spy.MockCache)
		tokenRepo := Token{repository: mockCache}
		token := "expired-token"
		expiresAt := time.Now().Add(-1 * time.Hour) // Já expirado
		// Não deve chamar Set quando TTL <= 0
		err := tokenRepo.SaveRevokedToken(token, expiresAt)
		assert.NoError(t, err)
		mockCache.AssertNotCalled(t, "Set")
	})

	t.Run("success - should does not save token with zero expiration", func(t *testing.T) {
		mockCache := new(spy.MockCache)
		tokenRepo := Token{repository: mockCache}
		token := "zero-expiration-token"
		expiresAt := time.Now() // Expira agora (TTL ~0)
		err := tokenRepo.SaveRevokedToken(token, expiresAt)
		assert.NoError(t, err)
		mockCache.AssertNotCalled(t, "Set")
	})

	t.Run("error - should cache set fails", func(t *testing.T) {
		mockCache := new(spy.MockCache)
		tokenRepo := Token{repository: mockCache}
		token := "test-token"
		expiresAt := time.Now().Add(30 * time.Minute)
		expectedKey := "revoked_token:test-token"
		expectedErr := errors.New("redis connection error")
		mockCache.On("Set", expectedKey, "1", mock.AnythingOfType("time.Duration")).
			Return(expectedErr)
		err := tokenRepo.SaveRevokedToken(token, expiresAt)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to save revoked token")
		assert.Contains(t, err.Error(), "redis connection error")
		mockCache.AssertExpectations(t)
	})

	t.Run("success - should saves token with different expiration times", func(t *testing.T) {
		testCases := []struct {
			name      string
			expiresIn time.Duration
		}{
			{"1 minute", 1 * time.Minute},
			{"15 minutes", 15 * time.Minute},
			{"1 hour", 1 * time.Hour},
			{"24 hours", 24 * time.Hour},
			{"7 days", 7 * 24 * time.Hour},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				mockCache := new(spy.MockCache)
				tokenRepo := Token{repository: mockCache}
				token := "test-token-" + tc.name
				expiresAt := time.Now().Add(tc.expiresIn)
				expectedKey := "revoked_token:" + token
				mockCache.On("Set", expectedKey, "1", mock.AnythingOfType("time.Duration")).
					Return(nil)
				err := tokenRepo.SaveRevokedToken(token, expiresAt)
				assert.NoError(t, err)
				mockCache.AssertExpectations(t)
			})
		}
	})
}
