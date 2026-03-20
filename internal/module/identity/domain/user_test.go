// Package domain
package domain

import (
	"testing"
	"time"

	"github.com/paladignus/actajus/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// MockPasswordHasher é um implementacao mock para testes
type MockPasswordHasher struct {
	hashFunc    func(string) (string, error)
	compareFunc func(string, string) error
	needsRehash func(string) bool
}

func (m *MockPasswordHasher) Hash(password string) (string, error) {
	if m.hashFunc != nil {
		return m.hashFunc(password)
	}
	return "hashed_" + password, nil
}

func (m *MockPasswordHasher) Compare(hash, password string) error {
	if m.compareFunc != nil {
		return m.compareFunc(hash, password)
	}
	if hash == "hashed_"+password {
		return nil
	}
	return ErrInvalidCredentials
}

func (m *MockPasswordHasher) NeedsRehash(hash string) bool {
	if m.needsRehash != nil {
		return m.needsRehash(hash)
	}
	return false
}

func TestUserBuilder_Build_Success(t *testing.T) {
	t.Parallel()

	user, err := NewUserBuilder().
		WithID(1).
		WithPrimaryEmail("user@example.com").
		WithPasswordHash("hashed_password").
		WithIsBlocked(false).
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, int64(1), user.ID().Value())
	assert.Equal(t, "user@example.com", user.PrimaryEmail().Value())
	assert.Equal(t, "hashed_password", user.PasswordHash().Value())
	assert.False(t, user.IsBlocked())
}

func TestUserBuilder_Build_InvalidEmail(t *testing.T) {
	t.Parallel()

	user, err := NewUserBuilder().
		WithID(1).
		WithPrimaryEmail("invalid-email").
		WithPasswordHash("hashed_password").
		Build()

	assert.Error(t, err)
	assert.Nil(t, user)
	fieldErr, ok := err.(*domain.FieldError)
	assert.True(t, ok)
	assert.Equal(t, "email", fieldErr.Field)
}

func TestUserBuilder_Build_EmptyEmail(t *testing.T) {
	t.Parallel()

	user, err := NewUserBuilder().
		WithID(1).
		WithPrimaryEmail("").
		WithPasswordHash("hashed_password").
		Build()

	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestUser_SetID_Success(t *testing.T) {
	t.Parallel()

	user, err := NewUserBuilder().
		WithPrimaryEmail("user@example.com").
		WithPasswordHash("hashed_password").
		Build()
	assert.NoError(t, err)

	err = user.SetID(1)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.ID().Value())
}

func TestUser_SetID_AlreadySet(t *testing.T) {
	t.Parallel()

	user, err := NewUserBuilder().
		WithID(1).
		WithPrimaryEmail("user@example.com").
		WithPasswordHash("hashed_password").
		Build()
	assert.NoError(t, err)

	err = user.SetID(2)

	assert.Error(t, err)
	fieldErr, ok := err.(*domain.FieldError)
	assert.True(t, ok)
	assert.Equal(t, "id", fieldErr.Field)
}

func TestUser_SetID_InvalidZero(t *testing.T) {
	t.Parallel()

	user, err := NewUserBuilder().
		WithPrimaryEmail("user@example.com").
		WithPasswordHash("hashed_password").
		Build()
	assert.NoError(t, err)

	err = user.SetID(0)

	assert.Error(t, err)
	fieldErr, ok := err.(*domain.FieldError)
	assert.True(t, ok)
	assert.Equal(t, "id", fieldErr.Field)
}

func TestUser_MatchesPassword_Success(t *testing.T) {
	t.Parallel()

	hasher := &MockPasswordHasher{}
	user, err := NewUserBuilder().
		WithID(1).
		WithPrimaryEmail("user@example.com").
		WithPasswordHash("hashed_Senha123").
		Build()
	assert.NoError(t, err)

	result := user.MatchesPassword("Senha123", hasher)

	assert.True(t, result)
}

func TestUser_MatchesPassword_Failure(t *testing.T) {
	t.Parallel()

	hasher := &MockPasswordHasher{}
	user, err := NewUserBuilder().
		WithID(1).
		WithPrimaryEmail("user@example.com").
		WithPasswordHash("hashed_Senha123").
		Build()
	assert.NoError(t, err)

	result := user.MatchesPassword("WrongPassword", hasher)

	assert.False(t, result)
}

func TestUser_ChangePassword_Success(t *testing.T) {
	t.Parallel()

	hasher := &MockPasswordHasher{}
	user, err := NewUserBuilder().
		WithID(1).
		WithPrimaryEmail("user@example.com").
		WithPasswordHash("hashed_OldPassword").
		WithIsBlocked(false).
		Build()
	assert.NoError(t, err)

	oldUpdatedAt := user.UpdatedAt()
	time.Sleep(10 * time.Millisecond) // Garante que o tempo passe

	err = user.ChangePassword("NewPassword123", hasher)

	assert.NoError(t, err)
	assert.Equal(t, "hashed_NewPassword123", user.PasswordHash().Value())
	assert.True(t, user.UpdatedAt().After(oldUpdatedAt))
}

func TestUser_ChangePassword_BlockedUser(t *testing.T) {
	t.Parallel()

	hasher := &MockPasswordHasher{}
	user, err := NewUserBuilder().
		WithID(1).
		WithPrimaryEmail("user@example.com").
		WithPasswordHash("hashed_OldPassword").
		WithIsBlocked(true).
		Build()
	assert.NoError(t, err)

	err = user.ChangePassword("NewPassword123", hasher)

	assert.Error(t, err)
	assert.Equal(t, ErrUserBlocked, err)
}

func TestUser_ChangePassword_EmptyPassword(t *testing.T) {
	t.Parallel()

	hasher := &MockPasswordHasher{}
	user, err := NewUserBuilder().
		WithID(1).
		WithPrimaryEmail("user@example.com").
		WithPasswordHash("hashed_OldPassword").
		Build()
	assert.NoError(t, err)

	err = user.ChangePassword("", hasher)

	assert.Error(t, err)
	fieldErr, ok := err.(*domain.FieldError)
	assert.True(t, ok)
	assert.Equal(t, "password", fieldErr.Field)
}

func TestUser_ChangePassword_HashError(t *testing.T) {
	t.Parallel()

	hasher := &MockPasswordHasher{
		hashFunc: func(string) (string, error) {
			return "", assert.AnError
		},
	}
	user, err := NewUserBuilder().
		WithID(1).
		WithPrimaryEmail("user@example.com").
		WithPasswordHash("hashed_OldPassword").
		Build()
	assert.NoError(t, err)

	err = user.ChangePassword("NewPassword123", hasher)

	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}

func TestUser_Block_Unblock(t *testing.T) {
	t.Parallel()

	user, err := NewUserBuilder().
		WithID(1).
		WithPrimaryEmail("user@example.com").
		WithPasswordHash("hashed_password").
		WithIsBlocked(false).
		Build()
	assert.NoError(t, err)

	assert.False(t, user.IsBlocked())

	user.Block()
	assert.True(t, user.IsBlocked())

	// Idempotência: bloquear novamente não causa erro
	user.Block()
	assert.True(t, user.IsBlocked())

	user.Unblock()
	assert.False(t, user.IsBlocked())

	// Idempotência: desbloquear novamente não causa erro
	user.Unblock()
	assert.False(t, user.IsBlocked())
}

func TestUser_Block_UpdatesTimestamp(t *testing.T) {
	t.Parallel()

	user, err := NewUserBuilder().
		WithID(1).
		WithPrimaryEmail("user@example.com").
		WithPasswordHash("hashed_password").
		Build()
	assert.NoError(t, err)

	oldUpdatedAt := user.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	user.Block()

	assert.True(t, user.UpdatedAt().After(oldUpdatedAt))
}

func TestUser_Unblock_UpdatesTimestamp(t *testing.T) {
	t.Parallel()

	user, err := NewUserBuilder().
		WithID(1).
		WithPrimaryEmail("user@example.com").
		WithPasswordHash("hashed_password").
		WithIsBlocked(true).
		Build()
	assert.NoError(t, err)

	oldUpdatedAt := user.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	user.Unblock()

	assert.True(t, user.UpdatedAt().After(oldUpdatedAt))
}

func TestUser_EmailEquals_True(t *testing.T) {
	t.Parallel()

	user, err := NewUserBuilder().
		WithID(1).
		WithPrimaryEmail("user@example.com").
		WithPasswordHash("hashed_password").
		Build()
	assert.NoError(t, err)

	result := user.EmailEquals("user@example.com")

	assert.True(t, result)
}

func TestUser_EmailEquals_False(t *testing.T) {
	t.Parallel()

	user, err := NewUserBuilder().
		WithID(1).
		WithPrimaryEmail("user@example.com").
		WithPasswordHash("hashed_password").
		Build()
	assert.NoError(t, err)

	result := user.EmailEquals("different@example.com")

	assert.False(t, result)
}

func TestUser_Timestamps_AreSet(t *testing.T) {
	t.Parallel()

	now := time.Now()
	user, err := NewUserBuilder().
		WithID(1).
		WithPrimaryEmail("user@example.com").
		WithPasswordHash("hashed_password").
		WithCreatedAt(now).
		WithUpdatedAt(now).
		Build()
	assert.NoError(t, err)

	assert.Equal(t, now, user.CreatedAt())
	assert.Equal(t, now, user.UpdatedAt())
}
