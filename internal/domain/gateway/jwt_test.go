package gateway

import (
	"testing"

	"github.com/paladignus/actajus/internal/application/readmodel"
	"github.com/paladignus/actajus/test/spy"
	"github.com/stretchr/testify/assert"
)

func TestTokenInterface(t *testing.T) {
	sut := &spy.JWT{}
	tokenPair, err := sut.GenerateTokenPair("user-id")
	assert.NoError(t, err)
	tokenClaims, err := sut.ValidateAccessToken("access-token")
	assert.NoError(t, err)
	refreshClaims, err := sut.ValidateRefreshToken("refresh-token")
	assert.NoError(t, err)
	refreshedPair, err := sut.RefreshAccessToken("refresh-token")
	assert.NoError(t, err)
	assert.Equal(t, tokenPair, readmodel.TokenPairReadModel{})
	assert.Equal(t, tokenClaims, readmodel.TokenClaimsReadModel{})
	assert.Equal(t, refreshClaims, readmodel.TokenClaimsReadModel{})
	assert.Equal(t, refreshedPair, readmodel.TokenPairReadModel{})
}
