package persistence

import (
	"errors"
	"testing"
	"time"

	"github.com/paladignus/actajus/test/spy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestToken_SaveRevokedToken(t *testing.T) {
	t.Run("success - should saves token with valid expiration", func(t *testing.T) {
		sut := new(spy.Cache)
		tokenRepo := Token{repository: sut}
		token := "valid-jwt-token"
		expiresAt := time.Now().Add(1 * time.Hour)
		expectedKey := "revoked_token:valid-jwt-token"
		sut.On("Set", expectedKey, "1", mock.MatchedBy(func(ttl time.Duration) bool {
			// Aceita TTL com margem de 1 segundo devido ao tempo de execução
			return ttl > 59*time.Minute && ttl <= 1*time.Hour
		})).Return(nil)
		err := tokenRepo.SaveRevokedToken(token, expiresAt)
		assert.NoError(t, err)
		sut.AssertExpectations(t)
	})

	t.Run("success - should does not save token with past expiration", func(t *testing.T) {
		sut := new(spy.Cache)
		tokenRepo := Token{repository: sut}
		token := "expired-token"
		expiresAt := time.Now().Add(-1 * time.Hour) // Já expirado
		// Não deve chamar Set quando TTL <= 0
		err := tokenRepo.SaveRevokedToken(token, expiresAt)
		assert.NoError(t, err)
		sut.AssertNotCalled(t, "Set")
	})

	t.Run("success - should does not save token with zero expiration", func(t *testing.T) {
		sut := new(spy.Cache)
		tokenRepo := Token{repository: sut}
		token := "zero-expiration-token"
		expiresAt := time.Now() // Expira agora (TTL ~0)
		err := tokenRepo.SaveRevokedToken(token, expiresAt)
		assert.NoError(t, err)
		sut.AssertNotCalled(t, "Set")
	})

	t.Run("error - should cache set fails", func(t *testing.T) {
		sut := new(spy.Cache)
		tokenRepo := Token{repository: sut}
		token := "test-token"
		expiresAt := time.Now().Add(30 * time.Minute)
		expectedKey := "revoked_token:test-token"
		expectedErr := errors.New("redis connection error")
		sut.On("Set", expectedKey, "1", mock.AnythingOfType("time.Duration")).
			Return(expectedErr)
		err := tokenRepo.SaveRevokedToken(token, expiresAt)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to save revoked token")
		assert.Contains(t, err.Error(), "redis connection error")
		sut.AssertExpectations(t)
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
				sut := new(spy.Cache)
				tokenRepo := Token{repository: sut}
				token := "test-token-" + tc.name
				expiresAt := time.Now().Add(tc.expiresIn)
				expectedKey := "revoked_token:" + token
				sut.On("Set", expectedKey, "1", mock.AnythingOfType("time.Duration")).
					Return(nil)
				err := tokenRepo.SaveRevokedToken(token, expiresAt)
				assert.NoError(t, err)
				sut.AssertExpectations(t)
			})
		}
	})
}

func TestToken_IsTokenRevoked(t *testing.T) {
	t.Run("success - esnure token is revoked", func(t *testing.T) {
		sut := new(spy.Cache)
		tokenRepo := Token{repository: sut}
		token := "revoked-token"
		expectedKey := "revoked_token:revoked-token"
		sut.On("Exists", expectedKey).Return(int64(1), nil)
		isRevoked, err := tokenRepo.IsTokenRevoked(token)
		assert.NoError(t, err)
		assert.True(t, isRevoked)
		sut.AssertExpectations(t)
	})

	t.Run("success - should token is not revoked", func(t *testing.T) {
		sut := new(spy.Cache)
		tokenRepo := Token{repository: sut}
		token := "valid-token"
		expectedKey := "revoked_token:valid-token"
		sut.On("Exists", expectedKey).Return(int64(0), nil)
		isRevoked, err := tokenRepo.IsTokenRevoked(token)
		assert.NoError(t, err)
		assert.False(t, isRevoked)
		sut.AssertExpectations(t)
	})

	t.Run("error - should cache exists check fails", func(t *testing.T) {
		sut := new(spy.Cache)
		tokenRepo := Token{repository: sut}
		token := "test-token"
		expectedKey := "revoked_token:test-token"
		expectedErr := errors.New("redis timeout error")
		sut.On("Exists", expectedKey).Return(int64(0), expectedErr)
		isRevoked, err := tokenRepo.IsTokenRevoked(token)
		assert.Error(t, err)
		assert.False(t, isRevoked)
		assert.Contains(t, err.Error(), "failed to check if token is revoked")
		assert.Contains(t, err.Error(), "redis timeout error")
		sut.AssertExpectations(t)
	})

	t.Run("success - should multiple tokens with different states", func(t *testing.T) {
		testCases := []struct {
			name      string
			token     string
			exists    int64
			isRevoked bool
		}{
			{"revoked token 1", "token-1", 1, true},
			{"valid token 1", "token-2", 0, false},
			{"revoked token 2", "token-3", 1, true},
			{"valid token 2", "token-4", 0, false},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				sut := new(spy.Cache)
				tokenRepo := Token{repository: sut}
				expectedKey := "revoked_token:" + tc.token
				sut.On("Exists", expectedKey).Return(tc.exists, nil)
				isRevoked, err := tokenRepo.IsTokenRevoked(tc.token)
				assert.NoError(t, err)
				assert.Equal(t, tc.isRevoked, isRevoked)
				sut.AssertExpectations(t)
			})
		}
	})
}

func TestToken_Integration(t *testing.T) {
	t.Run("full flow - should save and check revoked token", func(t *testing.T) {
		sut := new(spy.Cache)
		tokenRepo := Token{repository: sut}
		token := "integration-test-token"
		expiresAt := time.Now().Add(1 * time.Hour)
		expectedKey := "revoked_token:integration-test-token"
		// Configura mock para salvar
		sut.On("Set", expectedKey, "1", mock.AnythingOfType("time.Duration")).
			Return(nil)
		// Configura mock para verificar
		sut.On("Exists", expectedKey).Return(int64(1), nil)
		// Salva o token
		err := tokenRepo.SaveRevokedToken(token, expiresAt)
		assert.NoError(t, err)
		// Verifica se está revogado
		isRevoked, err := tokenRepo.IsTokenRevoked(token)
		assert.NoError(t, err)
		assert.True(t, isRevoked)
		sut.AssertExpectations(t)
	})

	t.Run("full flow - should token not saved should not be revoked", func(t *testing.T) {
		sut := new(spy.Cache)
		tokenRepo := Token{repository: sut}
		token := "non-revoked-token"
		expectedKey := "revoked_token:non-revoked-token"
		// Apenas verifica se está revogado (sem salvar)
		sut.On("Exists", expectedKey).Return(int64(0), nil)
		isRevoked, err := tokenRepo.IsTokenRevoked(token)
		assert.NoError(t, err)
		assert.False(t, isRevoked)
		sut.AssertExpectations(t)
	})
}

func TestToken_KeyFormat(t *testing.T) {
	t.Run("verify key format consistency", func(t *testing.T) {
		sut := new(spy.Cache)
		tokenRepo := Token{repository: sut}
		tokens := []string{
			"simple-token",
			"token.with.dots",
			"token-with-dashes",
			"token_with_underscores",
			"UPPERCASE-TOKEN",
			"MixedCase-Token",
		}
		for _, token := range tokens {
			expectedKey := "revoked_token:" + token
			expiresAt := time.Now().Add(1 * time.Hour)
			sut.On("Set", expectedKey, "1", mock.AnythingOfType("time.Duration")).
				Return(nil).Once()
			sut.On("Exists", expectedKey).Return(int64(1), nil).Once()
			// Testa save
			err := tokenRepo.SaveRevokedToken(token, expiresAt)
			assert.NoError(t, err)
			// Testa check
			isRevoked, err := tokenRepo.IsTokenRevoked(token)
			assert.NoError(t, err)
			assert.True(t, isRevoked)
		}
		sut.AssertExpectations(t)
	})
}
