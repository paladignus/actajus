// Package adapter
package adapter

import (
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/paladignus/actajus/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJWTAdapter(t *testing.T) {
	cfg := config.JWTConfig{
		AccessSecret:  "access-secret",
		RefreshSecret: "refresh-secret",
		AccessExpire:  15 * time.Minute,
		RefreshExpire: 7 * 24 * time.Hour,
		Issuer:        "test-issuer",
	}
	adapter := NewJWTAdapter(cfg)
	require.NotNil(t, adapter)
	assert.Equal(t, "access-secret", adapter.config.AccessSecret)
	assert.Equal(t, "refresh-secret", adapter.config.RefreshSecret)
	assert.Equal(t, 15*time.Minute, adapter.config.AccessExpire)
	assert.Equal(t, 7*24*time.Hour, adapter.config.RefreshExpire)
	assert.Equal(t, "test-issuer", adapter.config.Issuer)
}

func TestGenerateTokenPair(t *testing.T) {
	adapter := newTestAdapter()
	userID := "user123"
	tokens, err := adapter.GenerateTokenPair(userID)
	require.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
	assert.NotEqual(t, tokens.AccessToken, tokens.RefreshToken)
}

func TestValidateAccessToken(t *testing.T) {
	adapter := newTestAdapter()
	userID := "user123"
	tokens, err := adapter.GenerateTokenPair(userID)
	require.NoError(t, err)
	claims, err := adapter.ValidateAccessToken(tokens.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.IDUser)
	// assert.NotEmpty(t, claims.JTI)
	// assert.False(t, claims.Exp.IsZero())
	// assert.False(t, claims.Iat.IsZero())
}

func TestValidateRefreshToken(t *testing.T) {
	adapter := newTestAdapter()
	userID := "user123"
	tokens, err := adapter.GenerateTokenPair(userID)
	require.NoError(t, err)
	claims, err := adapter.ValidateRefreshToken(tokens.RefreshToken)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.IDUser)
	// assert.NotEmpty(t, claims.JTI)
	// assert.False(t, claims.Exp.IsZero())
	// assert.False(t, claims.Iat.IsZero())
}

func TestValidateAccessToken_WithRefreshToken_ShouldFail(t *testing.T) {
	adapter := newTestAdapter()
	userID := "user123"
	tokens, err := adapter.GenerateTokenPair(userID)
	require.NoError(t, err)
	_, err = adapter.ValidateAccessToken(tokens.RefreshToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")
}

func TestValidateRefreshToken_WithAccessToken_ShouldFail(t *testing.T) {
	adapter := newTestAdapter()
	userID := "user123"
	tokens, err := adapter.GenerateTokenPair(userID)
	require.NoError(t, err)
	_, err = adapter.ValidateRefreshToken(tokens.AccessToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")
}

func TestValidateAccessToken_InvalidToken(t *testing.T) {
	adapter := newTestAdapter()
	_, err := adapter.ValidateAccessToken("invalid-token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")
}

func TestValidateRefreshToken_InvalidToken(t *testing.T) {
	adapter := newTestAdapter()
	_, err := adapter.ValidateRefreshToken("invalid-token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")
}

func TestValidateAccessToken_ExpiredToken(t *testing.T) {
	cfg := config.JWTConfig{
		AccessSecret:  "access-secret",
		RefreshSecret: "refresh-secret",
		AccessExpire:  -1 * time.Hour, // Token já expirado
		RefreshExpire: 7 * 24 * time.Hour,
		Issuer:        "test-issuer",
	}
	adapter := NewJWTAdapter(cfg)

	tokens, err := adapter.GenerateTokenPair("user123")
	require.NoError(t, err)

	_, err = adapter.ValidateAccessToken(tokens.AccessToken)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token has expired")
}

func TestValidateRefreshToken_ExpiredToken(t *testing.T) {
	cfg := config.JWTConfig{
		AccessSecret:  "access-secret",
		RefreshSecret: "refresh-secret",
		AccessExpire:  15 * time.Minute,
		RefreshExpire: -1 * time.Hour, // Token já expirado
		Issuer:        "test-issuer",
	}
	adapter := NewJWTAdapter(cfg)

	tokens, err := adapter.GenerateTokenPair("user123")
	require.NoError(t, err)

	_, err = adapter.ValidateRefreshToken(tokens.RefreshToken)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token has expired")
}

func TestRefreshAccessToken(t *testing.T) {
	adapter := newTestAdapter()
	userID := "user123"

	tokens, err := adapter.GenerateTokenPair(userID)
	require.NoError(t, err)

	newTokens, err := adapter.RefreshAccessToken(tokens.RefreshToken)

	require.NoError(t, err)
	assert.NotEmpty(t, newTokens.AccessToken)
	assert.NotEmpty(t, newTokens.RefreshToken)
	assert.NotEqual(t, tokens.AccessToken, newTokens.AccessToken)
	assert.NotEqual(t, tokens.RefreshToken, newTokens.RefreshToken)

	// Validar que os novos tokens funcionam
	claims, err := adapter.ValidateAccessToken(newTokens.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.IDUser)
}

func TestRefreshAccessToken_InvalidRefreshToken(t *testing.T) {
	adapter := newTestAdapter()

	_, err := adapter.RefreshAccessToken("invalid-token")

	assert.Error(t, err)
	// assert.Contains(t, err.Error(), "invalid refresh token")
}

func TestRefreshAccessToken_WithAccessToken_ShouldFail(t *testing.T) {
	adapter := newTestAdapter()

	tokens, err := adapter.GenerateTokenPair("user123")
	require.NoError(t, err)

	_, err = adapter.RefreshAccessToken(tokens.AccessToken)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")
}

func TestRefreshAccessToken_ExpiredRefreshToken(t *testing.T) {
	cfg := config.JWTConfig{
		AccessSecret:  "access-secret",
		RefreshSecret: "refresh-secret",
		AccessExpire:  15 * time.Minute,
		RefreshExpire: -1 * time.Hour,
		Issuer:        "test-issuer",
	}
	adapter := NewJWTAdapter(cfg)
	tokens, err := adapter.GenerateTokenPair("user123")
	require.NoError(t, err)
	_, err = adapter.RefreshAccessToken(tokens.RefreshToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token has expired")
}

func TestGenerateToken_TokenStructure(t *testing.T) {
	adapter := newTestAdapter()
	userID := "user123"

	tokens, err := adapter.GenerateTokenPair(userID)
	require.NoError(t, err)

	// Parse do access token
	token, err := jwt.ParseString(
		tokens.AccessToken,
		jwt.WithKey(jwa.HS256(), []byte(adapter.config.AccessSecret)),
		jwt.WithValidate(true),
	)
	require.NoError(t, err)

	subject, _ := token.Subject()
	assert.Equal(t, userID, subject)

	// assert.Equal(t, "test-issuer", token.Issuer())
	// assert.NotEmpty(t, token.JwtID())
	// assert.False(t, token.IssuedAt().IsZero())
	// assert.False(t, token.Expiration().IsZero())
}

func TestValidateToken_WrongSecret(t *testing.T) {
	adapter := newTestAdapter()
	userID := "user123"

	tokens, err := adapter.GenerateTokenPair(userID)
	require.NoError(t, err)

	// Criar outro adapter com secret diferente
	wrongAdapter := NewJWTAdapter(config.JWTConfig{
		AccessSecret:  "wrong-secret",
		RefreshSecret: "wrong-refresh-secret",
		AccessExpire:  15 * time.Minute,
		RefreshExpire: 7 * 24 * time.Hour,
		Issuer:        "test-issuer",
	})

	_, err = wrongAdapter.ValidateAccessToken(tokens.AccessToken)
	assert.Error(t, err)

	_, err = wrongAdapter.ValidateRefreshToken(tokens.RefreshToken)
	assert.Error(t, err)
}

// func TestGenerateJTI(t *testing.T) {
// 	jti1, err := generateJTI()
// 	require.NoError(t, err)
// 	assert.NotEmpty(t, jti1)
// 	assert.Len(t, jti1, 32) // 16 bytes em hex = 32 caracteres
//
// 	jti2, err := generateJTI()
// 	require.NoError(t, err)
// 	assert.NotEmpty(t, jti2)
//
// 	// JTIs devem ser diferentes
// 	assert.NotEqual(t, jti1, jti2)
// }

func TestTokenExpiration_AccessToken(t *testing.T) {
	cfg := config.JWTConfig{
		AccessSecret:  "access-secret",
		RefreshSecret: "refresh-secret",
		AccessExpire:  1 * time.Second,
		RefreshExpire: 7 * 24 * time.Hour,
		Issuer:        "test-issuer",
	}
	adapter := NewJWTAdapter(cfg)

	tokens, err := adapter.GenerateTokenPair("user123")
	require.NoError(t, err)

	// Token válido inicialmente
	_, err = adapter.ValidateAccessToken(tokens.AccessToken)
	require.NoError(t, err)

	// Aguardar expiração
	time.Sleep(2 * time.Second)

	// Token deve estar expirado
	_, err = adapter.ValidateAccessToken(tokens.AccessToken)
	assert.Error(t, err)
}

// Helper function
func newTestAdapter() jwtAdapter {
	config := config.JWTConfig{
		AccessSecret:  "access-secret",
		RefreshSecret: "refresh-secret",
		AccessExpire:  15 * time.Minute,
		RefreshExpire: 7 * 24 * time.Hour,
		Issuer:        "test-issuer",
	}
	return NewJWTAdapter(config)
}

func TestAccessTokenSubject(t *testing.T) {
	adapter := newTestAdapter()
	tokens, err := adapter.GenerateTokenPair("")
	require.NoError(t, err)
	_, err = adapter.ValidateAccessToken(tokens.AccessToken)
	assert.Error(t, err)
}
