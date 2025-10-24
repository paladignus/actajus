// Package adapter
package adapter

import (
	"testing"
	"time"

	"github.com/paladignus/actajus/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testConfig() config.JWTConfig {
	cfg := config.JWTConfig{
		AccessSecret:  "access-secret",
		RefreshSecret: "refresh-secret",
		AccessExpire:  15 * time.Minute,
		RefreshExpire: 7 * 24 * time.Hour,
		Issuer:        "test-issuer",
	}
	return cfg
}

func TestNewJWTAdapter(t *testing.T) {
	adapter := NewJWTAdapter(testConfig())
	require.NotNil(t, adapter)
	assert.Equal(t, "access-secret", adapter.config.AccessSecret)
	assert.Equal(t, "refresh-secret", adapter.config.RefreshSecret)
	assert.Equal(t, 15*time.Minute, adapter.config.AccessExpire)
	assert.Equal(t, 7*24*time.Hour, adapter.config.RefreshExpire)
	assert.Equal(t, "test-issuer", adapter.config.Issuer)
}

func TestGenerateTokenPair(t *testing.T) {
	t.Run("should return access token and refresh token", func(t *testing.T) {
		adapter := NewJWTAdapter(testConfig())
		userID := "user123"
		tokens, err := adapter.GenerateTokenPair(userID)
		require.NoError(t, err)
		assert.NotEmpty(t, tokens.AccessToken)
		assert.NotEmpty(t, tokens.RefreshToken)
		assert.NotEqual(t, tokens.AccessToken, tokens.RefreshToken)
	})

	t.Run("should validate an access token without error", func(t *testing.T) {
		adapter := NewJWTAdapter(testConfig())
		userID := "user123"
		tokens, err := adapter.GenerateTokenPair(userID)
		require.NoError(t, err)
		claims, err := adapter.ValidateAccessToken(tokens.AccessToken)
		require.NoError(t, err)
		assert.Equal(t, userID, claims.IDUser)
	})

	t.Run("should validate an refresh token without error", func(t *testing.T) {
		adapter := NewJWTAdapter(testConfig())
		userID := "user123"
		tokens, err := adapter.GenerateTokenPair(userID)
		require.NoError(t, err)
		claims, err := adapter.ValidateRefreshToken(tokens.RefreshToken)
		require.NoError(t, err)
		assert.Equal(t, userID, claims.IDUser)
	})

	t.Run("should return access token and refresh token from refresh method", func(t *testing.T) {
		adapter := NewJWTAdapter(testConfig())
		userID := "user123"
		tokens, err := adapter.GenerateTokenPair(userID)
		assert.NoError(t, err)
		newTokens, err := adapter.RefreshAccessToken(tokens.RefreshToken)
		assert.NoError(t, err)
		assert.NotEmpty(t, newTokens.AccessToken)
		assert.NotEmpty(t, newTokens.RefreshToken)
		assert.NotEqual(t, newTokens.AccessToken, newTokens.RefreshToken)
	})

	t.Run("should validate an access token and refresh token with error invalid", func(t *testing.T) {
		adapter := NewJWTAdapter(testConfig())
		_, err := adapter.ValidateAccessToken("invalid-token")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token")
		_, err = adapter.ValidateRefreshToken("invalid-token")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("should validate an access token and refresh token with error expired", func(t *testing.T) {
		cfg := testConfig()
		cfg.AccessExpire = -1 * time.Hour
		cfg.RefreshExpire = -1 * time.Hour
		adapter := NewJWTAdapter(cfg)
		userID := "user123"
		tokens, err := adapter.GenerateTokenPair(userID)
		require.NoError(t, err)
		_, err = adapter.ValidateAccessToken(tokens.AccessToken)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token has expired")
		_, err = adapter.ValidateRefreshToken(tokens.RefreshToken)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token has expired")
	})

	t.Run("should validate an access token and refresh token with error subject", func(t *testing.T) {
		adapter := NewJWTAdapter(testConfig())
		tokens, err := adapter.GenerateTokenPair("")
		require.NoError(t, err)
		_, err = adapter.ValidateAccessToken(tokens.AccessToken)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token")
		_, err = adapter.ValidateRefreshToken(tokens.RefreshToken)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("should return an error if the refresh token is invalid when using the refresh method", func(t *testing.T) {
		cfg := testConfig()
		cfg.RefreshExpire = -1 * time.Hour
		adapter := NewJWTAdapter(cfg)
		tokens, err := adapter.GenerateTokenPair("user-123")
		assert.NoError(t, err)
		_, err = adapter.RefreshAccessToken(tokens.RefreshToken)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token has expired")
		_, err = adapter.RefreshAccessToken(tokens.AccessToken)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token")
	})
}
