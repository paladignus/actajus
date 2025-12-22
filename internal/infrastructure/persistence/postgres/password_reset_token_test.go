// Package postgres
package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/paladignus/actajus/internal/domain/entity"
	"github.com/paladignus/actajus/internal/infrastructure/database"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
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
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error when creating password reset token", func(t *testing.T) {
		expectdErr := errors.New("database error")
		mock.ExpectExec(`INSERT INTO password_reset`).
			WithArgs(inputToken.IDUser, inputToken.Token, inputToken.ExpiresAt).
			WillReturnError(expectdErr)
		err = repo.Create(ctx, inputToken)
		assert.Error(t, expectdErr)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPasswordResetToken_FindByToken(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	db := &database.DB{Pool: mock}
	repo := NewPasswordResetToken(db)
	token := "okokopkpodjklsjiodjsoi"
	now := time.Now()
	inputToken := entity.PasswordResetToken{
		IDPasswordReset: 1,
		IDUser:          1,
		Token:           token,
		ExpiresAt:       now.Add(1 * time.Hour),
		UsedAt:          nil,
		CreatedAt:       now,
	}

	t.Run("successfully finds a password reset token", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"idpassword_reset", "id_users", "token", "expires_at", "used_at"}).
			AddRow(inputToken.IDPasswordReset, inputToken.IDUser, inputToken.Token, inputToken.ExpiresAt, inputToken.UsedAt)
		mock.ExpectQuery(`SELECT idpassword_reset, id_users, token, expires_at, used_at FROM password_reset`).
			WithArgs(token).
			WillReturnRows(rows)
		result, err := repo.FindByToken(ctx, token)
		assert.NoError(t, err)
		assert.Equal(t, inputToken.IDPasswordReset, result.IDPasswordReset)
		assert.Equal(t, inputToken.IDUser, result.IDUser)
		assert.Equal(t, inputToken.Token, result.Token)
		assert.Equal(t, inputToken.ExpiresAt, result.ExpiresAt)
		assert.Equal(t, inputToken.UsedAt, result.UsedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when token is not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT idpassword_reset, id_users, token, expires_at, used_at FROM password_reset`).
			WithArgs(token).
			WillReturnError(pgx.ErrNoRows)
		result, err := repo.FindByToken(ctx, token)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "invalid token okokopkpodjklsjiodjsoi")
		assert.Equal(t, entity.PasswordResetToken{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should retturn error from database", func(t *testing.T) {
		expectedErr := errors.New("database error")
		mock.ExpectQuery(`SELECT idpassword_reset, id_users, token, expires_at, used_at FROM password_reset`).
			WithArgs(token).
			WillReturnError(expectedErr)
		result, err := repo.FindByToken(ctx, token)
		assert.Error(t, expectedErr)
		assert.ErrorContains(t, err, "database error while finding token okokopkpodjklsjiodjsoi")
		assert.Equal(t, entity.PasswordResetToken{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPasswordResetToken_MarkAsUsed(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	db := &database.DB{Pool: mock}
	repo := NewPasswordResetToken(db)
	token := "okokopkpodjklsjiodjsoi"

	t.Run("should to mark the token as used", func(t *testing.T) {
		mock.ExpectExec(`UPDATE password_reset`).
			WithArgs(token).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		err := repo.MarkAsUsed(ctx, token)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("sdhould return error when token is not found", func(t *testing.T) {
		mock.ExpectExec(`UPDATE password_reset`).
			WithArgs(token).
			WillReturnError(pgx.ErrNoRows)
		err := repo.MarkAsUsed(ctx, token)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "database error while marking token okokopkpodjklsjiodjsoi as used")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPasswordResetToken_InvalidateUserTokens(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	db := &database.DB{Pool: mock}
	repo := NewPasswordResetToken(db)

	t.Run("should invalidate the token", func(t *testing.T) {
		mock.ExpectExec(`UPDATE password_reset`).
			WithArgs(1).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		err := repo.InvalidateUserTokens(ctx, 1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error", func(t *testing.T) {
		mockErr := errors.New("database error")
		mock.ExpectExec(`UPDATE password_reset`).
			WithArgs(1).
			WillReturnError(mockErr)
		err := repo.InvalidateUserTokens(ctx, 1)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "database error while invalidating tokens for user ID 1")
		assert.ErrorContains(t, err, mockErr.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
