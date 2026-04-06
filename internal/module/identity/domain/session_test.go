// Package domain
package domain

import (
	"testing"
	"time"

	"github.com/paladignus/actajus/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

func TestSessionBuilder_Build_Success(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)
	var hash [32]byte
	copy(hash[:], "test_hash_value")

	session, err := NewSessionBuilder().
		WithID(1).
		WithIDUser(123).
		WithRefreshHash(hash).
		WithExpiresAt(expiresAt).
		WithIP("192.168.1.1").
		WithUserAgent("Mozilla/5.0").
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, int64(1), session.ID().Value())
	assert.Equal(t, int64(123), session.IDUser().Value())
	assert.Equal(t, hash, session.RefreshHash())
	assert.Equal(t, expiresAt, session.ExpiresAt())
	assert.Equal(t, "192.168.1.1", session.IP())
	assert.Equal(t, "Mozilla/5.0", session.UserAgent())
	assert.False(t, session.IsRevoked())
	assert.False(t, session.IsExpired(now))
}

func TestSessionBuilder_Build_InvalidUserID(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	session, err := NewSessionBuilder().
		WithID(1).
		WithIDUser(0).
		WithExpiresAt(expiresAt).
		Build()

	assert.Error(t, err)
	assert.Nil(t, session)
	fieldErr, ok := err.(*domain.FieldError)
	assert.True(t, ok)
	assert.Equal(t, "user_id", fieldErr.Field)
}

func TestSessionBuilder_Build_MissingExpiresAt(t *testing.T) {
	t.Parallel()

	session, err := NewSessionBuilder().
		WithID(1).
		WithIDUser(123).
		Build()

	assert.Error(t, err)
	assert.Nil(t, session)
	fieldErr, ok := err.(*domain.FieldError)
	assert.True(t, ok)
	assert.Equal(t, "expires_at", fieldErr.Field)
}

func TestSession_SetID_Success(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	session, err := NewSessionBuilder().
		WithIDUser(123).
		WithExpiresAt(expiresAt).
		Build()
	assert.NoError(t, err)

	err = session.SetID(1)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), session.ID().Value())
}

func TestSession_SetID_AlreadySet(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	session, err := NewSessionBuilder().
		WithID(1).
		WithIDUser(123).
		WithExpiresAt(expiresAt).
		Build()
	assert.NoError(t, err)

	err = session.SetID(2)

	assert.Error(t, err)
	fieldErr, ok := err.(*domain.FieldError)
	assert.True(t, ok)
	assert.Equal(t, "id", fieldErr.Field)
}

func TestSession_SetID_InvalidZero(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	session, err := NewSessionBuilder().
		WithIDUser(123).
		WithExpiresAt(expiresAt).
		Build()
	assert.NoError(t, err)

	err = session.SetID(0)

	assert.Error(t, err)
	fieldErr, ok := err.(*domain.FieldError)
	assert.True(t, ok)
	assert.Equal(t, "id", fieldErr.Field)
}

func TestSession_IsActive_True(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	session, err := NewSessionBuilder().
		WithID(1).
		WithIDUser(123).
		WithExpiresAt(expiresAt).
		Build()
	assert.NoError(t, err)

	assert.True(t, session.IsActive(now))
}

func TestSession_IsActive_False_Revoked(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	session, err := NewSessionBuilder().
		WithID(1).
		WithIDUser(123).
		WithExpiresAt(expiresAt).
		Build()
	assert.NoError(t, err)

	session.Revoke(now)

	assert.False(t, session.IsActive(now))
}

func TestSession_IsActive_False_Expired(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expiresAt := now.Add(-1 * time.Hour)

	session, err := NewSessionBuilder().
		WithID(1).
		WithIDUser(123).
		WithExpiresAt(expiresAt).
		Build()
	assert.NoError(t, err)

	assert.False(t, session.IsActive(now))
}

func TestSession_Revoke_Success(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	session, err := NewSessionBuilder().
		WithID(1).
		WithIDUser(123).
		WithExpiresAt(expiresAt).
		Build()
	assert.NoError(t, err)

	oldUpdatedAt := session.UpdatedAt()
	revokeTime := now.Add(100 * time.Millisecond) // Tempo futuro para garantir a diferença

	err = session.Revoke(revokeTime)

	assert.NoError(t, err)
	assert.True(t, session.IsRevoked())
	assert.NotNil(t, session.RevokedAt())
	assert.True(t, session.UpdatedAt().After(oldUpdatedAt))
	assert.Equal(t, revokeTime, *session.RevokedAt())
}

func TestSession_Revoke_InvalidID(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	session, err := NewSessionBuilder().
		WithIDUser(123).
		WithExpiresAt(expiresAt).
		Build()
	assert.NoError(t, err)

	err = session.Revoke(now)

	assert.Error(t, err)
	fieldErr, ok := err.(*domain.FieldError)
	assert.True(t, ok)
	assert.Equal(t, "id", fieldErr.Field)
}

func TestSession_Revoke_AlreadyRevoked(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	session, err := NewSessionBuilder().
		WithID(1).
		WithIDUser(123).
		WithExpiresAt(expiresAt).
		Build()
	assert.NoError(t, err)

	// Primeira revogação
	err = session.Revoke(now)
	assert.NoError(t, err)

	// Segunda revogação (deve falhar)
	err = session.Revoke(now)

	assert.Error(t, err)
	fieldErr, ok := err.(*domain.FieldError)
	assert.True(t, ok)
	assert.Equal(t, "revoked_at", fieldErr.Field)
}

func TestSession_Rotate_Success(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)
	var oldHash [32]byte
	var newHash [32]byte
	copy(oldHash[:], "old_hash_value")
	copy(newHash[:], "new_hash_value")
	newExpiresAt := now.Add(48 * time.Hour)

	session, err := NewSessionBuilder().
		WithID(1).
		WithIDUser(123).
		WithRefreshHash(oldHash).
		WithExpiresAt(expiresAt).
		Build()
	assert.NoError(t, err)

	oldUpdatedAt := session.UpdatedAt()
	rotateTime := now.Add(100 * time.Millisecond) // Tempo futuro para garantir a diferença

	err = session.Rotate(newHash, newExpiresAt, rotateTime)

	assert.NoError(t, err)
	assert.Equal(t, newHash, session.RefreshHash())
	assert.Equal(t, newExpiresAt, session.ExpiresAt())
	assert.NotNil(t, session.RotatedAt())
	assert.True(t, session.UpdatedAt().After(oldUpdatedAt))
	assert.Equal(t, rotateTime, *session.RotatedAt())
}

func TestSession_Rotate_RevokedSession(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)
	var hash [32]byte

	session, err := NewSessionBuilder().
		WithID(1).
		WithIDUser(123).
		WithRefreshHash(hash).
		WithExpiresAt(expiresAt).
		Build()
	assert.NoError(t, err)

	// Revoga a sessão
	session.Revoke(now)

	// Tenta rotacionar (deve falhar)
	err = session.Rotate(hash, expiresAt, now)

	assert.Error(t, err)
	assert.Equal(t, ErrSessionRevoked, err)
}

func TestSession_Rotate_InvalidID(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)
	var hash [32]byte

	session, err := NewSessionBuilder().
		WithIDUser(123).
		WithRefreshHash(hash).
		WithExpiresAt(expiresAt).
		Build()
	assert.NoError(t, err)

	err = session.Rotate(hash, expiresAt, now)

	assert.Error(t, err)
	fieldErr, ok := err.(*domain.FieldError)
	assert.True(t, ok)
	assert.Equal(t, "id", fieldErr.Field)
}

func TestSession_CanRotate_True(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	session, err := NewSessionBuilder().
		WithID(1).
		WithIDUser(123).
		WithExpiresAt(expiresAt).
		Build()
	assert.NoError(t, err)

	assert.True(t, session.CanRotate(now))
}

func TestSession_CanRotate_False_Expired(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expiresAt := now.Add(-1 * time.Hour)

	session, err := NewSessionBuilder().
		WithID(1).
		WithIDUser(123).
		WithExpiresAt(expiresAt).
		Build()
	assert.NoError(t, err)

	assert.False(t, session.CanRotate(now))
}

func TestSession_CanRotate_False_Revoked(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	session, err := NewSessionBuilder().
		WithID(1).
		WithIDUser(123).
		WithExpiresAt(expiresAt).
		Build()
	assert.NoError(t, err)

	session.Revoke(now)

	assert.False(t, session.CanRotate(now))
}

func TestSession_TimeToLive_Active(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	session, err := NewSessionBuilder().
		WithID(1).
		WithIDUser(123).
		WithExpiresAt(expiresAt).
		Build()
	assert.NoError(t, err)

	ttl := session.TimeToLive(now)

	assert.Greater(t, ttl, 23*time.Hour)
	assert.Less(t, ttl, 24*time.Hour+1*time.Second)
}

func TestSession_TimeToLive_Expired(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expiresAt := now.Add(-1 * time.Hour)

	session, err := NewSessionBuilder().
		WithID(1).
		WithIDUser(123).
		WithExpiresAt(expiresAt).
		Build()
	assert.NoError(t, err)

	ttl := session.TimeToLive(now)

	assert.Equal(t, time.Duration(0), ttl)
}

func TestSession_TimeToLive_Revoked(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	session, err := NewSessionBuilder().
		WithID(1).
		WithIDUser(123).
		WithExpiresAt(expiresAt).
		Build()
	assert.NoError(t, err)

	session.Revoke(now)

	ttl := session.TimeToLive(now)

	assert.Equal(t, time.Duration(0), ttl)
}

func TestSession_IsAboutToExpire_True(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expiresAt := now.Add(2 * time.Minute)
	threshold := 5 * time.Minute

	session, err := NewSessionBuilder().
		WithID(1).
		WithIDUser(123).
		WithExpiresAt(expiresAt).
		Build()
	assert.NoError(t, err)

	assert.True(t, session.IsAboutToExpire(now, threshold))
}

func TestSession_IsAboutToExpire_False(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)
	threshold := 5 * time.Minute

	session, err := NewSessionBuilder().
		WithID(1).
		WithIDUser(123).
		WithExpiresAt(expiresAt).
		Build()
	assert.NoError(t, err)

	assert.False(t, session.IsAboutToExpire(now, threshold))
}

func TestSession_IsAboutToExpire_Expired(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expiresAt := now.Add(-1 * time.Hour)
	threshold := 5 * time.Minute

	session, err := NewSessionBuilder().
		WithID(1).
		WithIDUser(123).
		WithExpiresAt(expiresAt).
		Build()
	assert.NoError(t, err)

	assert.False(t, session.IsAboutToExpire(now, threshold))
}

func TestSession_IsAboutToExpire_Revoked(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expiresAt := now.Add(2 * time.Minute)
	threshold := 5 * time.Minute

	session, err := NewSessionBuilder().
		WithID(1).
		WithIDUser(123).
		WithExpiresAt(expiresAt).
		Build()
	assert.NoError(t, err)

	session.Revoke(now)

	assert.False(t, session.IsAboutToExpire(now, threshold))
}

func TestSession_Timestamps_AreSet(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	session, err := NewSessionBuilder().
		WithID(1).
		WithIDUser(123).
		WithExpiresAt(expiresAt).
		WithCreatedAt(now).
		WithUpdatedAt(now).
		Build()
	assert.NoError(t, err)

	assert.Equal(t, now, session.CreatedAt())
	assert.Equal(t, now, session.UpdatedAt())
}

func TestSession_IsExpired(t *testing.T) {
	t.Parallel()

	now := time.Now()

	// Sessão não expirada
	expiresAt := now.Add(24 * time.Hour)
	session1, err := NewSessionBuilder().
		WithID(1).
		WithIDUser(123).
		WithExpiresAt(expiresAt).
		Build()
	assert.NoError(t, err)
	assert.False(t, session1.IsExpired(now))

	// Sessão expirada
	expiresAt2 := now.Add(-1 * time.Hour)
	session2, err := NewSessionBuilder().
		WithID(1).
		WithIDUser(123).
		WithExpiresAt(expiresAt2).
		Build()
	assert.NoError(t, err)
	assert.True(t, session2.IsExpired(now))

	// Sessão expirando exatamente agora
	session3, err := NewSessionBuilder().
		WithID(1).
		WithIDUser(123).
		WithExpiresAt(now).
		Build()
	assert.NoError(t, err)
	assert.True(t, session3.IsExpired(now))
}
