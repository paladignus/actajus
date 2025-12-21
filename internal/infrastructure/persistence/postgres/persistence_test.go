// Package postgres
package postgres

import (
	"testing"

	"github.com/paladignus/actajus/internal/infrastructure/database"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPersistence(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	db := &database.DB{Pool: mock}
	t.Run("ensure creates persistence instance successfully", func(t *testing.T) {
		persistence := NewPersistence(db)
		assert.NotNil(t, persistence)
		assert.NotNil(t, persistence.db)
		assert.Equal(t, db, persistence.db)
	})

	t.Run("ensure creates persistence with nil db", func(t *testing.T) {
		persistence := NewPersistence(nil)
		assert.NotNil(t, persistence)
		assert.Nil(t, persistence.db)
	})
}

func TestPersistence_User(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	db := &database.DB{Pool: mock}
	persistence := NewPersistence(db)
	t.Run("should returns user repository instance", func(t *testing.T) {
		userRepo := persistence.User()
		assert.NotNil(t, userRepo)
		assert.Equal(t, db, userRepo.db)
	})

	t.Run("should returns new account instance on each call", func(t *testing.T) {
		user1 := persistence.User()
		user2 := persistence.User()
		assert.Equal(t, user1.db, user2.db)
		assert.NotNil(t, user1)
		assert.NotNil(t, user2)
	})

	t.Run("should account repository shares same db connection", func(t *testing.T) {
		user := persistence.User()
		assert.Same(t, persistence.db, user.db)
		assert.Same(t, db.Pool, user.db.Pool)
	})
}

func TestPersistence_PasswordResetToken(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	db := &database.DB{Pool: mock}
	persistence := NewPersistence(db)
	t.Run("should returns password reset token repository instance", func(t *testing.T) {
		passwordResetTokenRepo := persistence.PasswordResetToken()
		assert.NotNil(t, passwordResetTokenRepo)
		assert.Equal(t, db, passwordResetTokenRepo.db)
	})

	t.Run("should returns new password reset token instance on each call", func(t *testing.T) {
		passwordResetToken1 := persistence.PasswordResetToken()
		passwordResetToken2 := persistence.PasswordResetToken()
		assert.Equal(t, passwordResetToken1.db, passwordResetToken2.db)
		assert.NotNil(t, passwordResetToken1)
		assert.NotNil(t, passwordResetToken2)
	})

	t.Run("should password reset token repository shares same db connection", func(t *testing.T) {
		passwordResetToken := persistence.PasswordResetToken()
		assert.Same(t, persistence.db, passwordResetToken.db)
		assert.Same(t, db.Pool, passwordResetToken.db.Pool)
	})
}
