// Package postgres
package postgres

import (
	"context"
	"testing"

	"github.com/paladignus/actajus/internal/domain/entity"
	"github.com/paladignus/actajus/internal/infrastructure/database"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestPasswordResetToken_Create(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	db := &database.DB{Pool: mock}
	repo := NewPasswordResetToken(db)
	inputToken := entity.PasswordResetToken{
		IDUser: 1,
		Token:  "sample",
	}
	t.Run("successfully creates a password reset token", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO password_reset`).
			WithArgs(inputToken.IDUser, inputToken.Token, inputToken.ExpiresAt).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))
		err = repo.Create(ctx, inputToken)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
