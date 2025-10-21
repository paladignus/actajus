// Package adapter
package adapter

import (
	"testing"
	"time"

	"github.com/paladignus/actajus/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestConfig() config.JWTConfig {
	return config.JWTConfig{
		AccessSecret:  "test-access-secret",
		RefreshSecret: "test-refresh-secret",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
	}
}

func TestGenerateTokenPair(t *testing.T) {
	adapter := NewJWTAdapter(getTestConfig())
	tokenPair, err := adapter.GenerateTokenPair("user-123")
	require.NoError(t, err)
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)
	assert.NotEqual(t, tokenPair.AccessToken, tokenPair.RefreshToken)
}
