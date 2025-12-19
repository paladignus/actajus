package gateway

import (
	"testing"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/test/spy"
	"github.com/stretchr/testify/assert"
)

func TestTokenInterface(t *testing.T) {
	sut := &spy.Token{}
	tokenPair, err := sut.GenerateTokenPair("user-id")
	assert.NoError(t, err)
	// resetToken, err := sut.GenerateResetToken("user-id")
	assert.NoError(t, err)
	tokenClaims, err := sut.ValidateAccessToken("access-token")
	assert.NoError(t, err)
	refreshClaims, err := sut.ValidateRefreshToken("refresh-token")
	assert.NoError(t, err)
	// resetClaims, err := sut.ValidateResetToken("reset-token")
	// assert.NoError(t, err)
	refreshedPair, err := sut.RefreshAccessToken("refresh-token")
	assert.NoError(t, err)
	assert.Equal(t, tokenPair, dto.TokenPair{})
	// assert.Equal(t, resetToken, dto.TokenRecover{})
	assert.Equal(t, tokenClaims, dto.TokenClaims{})
	assert.Equal(t, refreshClaims, dto.TokenClaims{})
	// assert.Equal(t, resetClaims, dto.TokenClaims{})
	assert.Equal(t, refreshedPair, dto.TokenPair{})
}
