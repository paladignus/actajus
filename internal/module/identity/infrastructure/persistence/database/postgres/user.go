// Package postgres
package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	identityrepo "github.com/paladignus/actajus/internal/module/identity/application/repository"
	"github.com/paladignus/actajus/internal/module/identity/domain"
	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type User struct {
	db postgres.Executor
}

func NewUser(db postgres.Executor) *User {
	return &User{db}
}

func (r User) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	const query = `
		SELECT
			u.idusers, u.password_hash, u.is_blocked, e.verified_at
		FROM users u
		LEFT JOIN email_person ep ON u.idusers = ep.id_people
		JOIN emails e ON e.idemails = ep.id_emails AND e.is_primary = TRUE AND e.deleted_at IS NULL
		WHERE e.address = $1 AND u.deleted_at IS NULL LIMIT 1;`
	var (
		id           int64
		passwordHash string
		isBlocked    bool
		verifiedAt   *time.Time
	)
	err := r.db.QueryRow(ctx, query, email).Scan(
		&id, &passwordHash, &isBlocked, &verifiedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return domain.NewUserBuilder().
		WithID(id).
		WithPrimaryEmail(email).
		WithPasswordHash(passwordHash).
		WithIsBlocked(isBlocked).
		WithPrimaryEmailVerifiedAt(verifiedAt).
		Build()
}

func (r User) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	const query = `
		SELECT
			u.password_hash, u.is_blocked, e.address, e.verified_at
		FROM users u
		LEFT JOIN email_person ep ON u.idusers = ep.id_people
		JOIN emails e ON e.idemails = ep.id_emails AND e.is_primary = TRUE AND e.deleted_at IS NULL
		WHERE u.idusers = $1 AND u.deleted_at IS NULL LIMIT 1;`
	var (
		passwordHash string
		isBlocked    bool
		email        string
		verifiedAt   *time.Time
	)
	err := r.db.QueryRow(ctx, query, id).Scan(
		&passwordHash, &isBlocked, &email, &verifiedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return domain.NewUserBuilder().
		WithID(id).
		WithPrimaryEmail(email).
		WithPasswordHash(passwordHash).
		WithIsBlocked(isBlocked).
		WithPrimaryEmailVerifiedAt(verifiedAt).
		Build()
}

func (r User) Create(ctx context.Context, input identityrepo.CreateUserRegistration) (identityrepo.CreatedUserRegistration, error) {
	const personQuery = `
		INSERT INTO people
			(first_name, last_name, birthday, id_gender, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING idpeople;`
	const userQuery = `
		INSERT INTO users
			(idusers, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4);`
	const emailQuery = `
		INSERT INTO emails
			(address, is_primary, created_at, updated_at)
		VALUES ($1, TRUE, $2, $3)
		RETURNING idemails;`
	const linkQuery = `
		INSERT INTO email_person (id_emails, id_people)
		VALUES ($1, $2);`

	var idUser int64
	if err := r.db.QueryRow(ctx, personQuery,
		input.FirstName,
		input.LastName,
		input.Birthday,
		input.GenderID,
		input.CreatedAt,
		input.UpdatedAt,
	).Scan(&idUser); err != nil {
		return identityrepo.CreatedUserRegistration{}, err
	}
	if _, err := r.db.Exec(ctx, userQuery, idUser, input.PasswordHash, input.CreatedAt, input.UpdatedAt); err != nil {
		return identityrepo.CreatedUserRegistration{}, err
	}
	var idEmail int64
	if err := r.db.QueryRow(ctx, emailQuery, input.Email, input.CreatedAt, input.UpdatedAt).Scan(&idEmail); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return identityrepo.CreatedUserRegistration{}, domain.ErrEmailAlreadyExists
		}
		return identityrepo.CreatedUserRegistration{}, err
	}
	if _, err := r.db.Exec(ctx, linkQuery, idEmail, idUser); err != nil {
		return identityrepo.CreatedUserRegistration{}, err
	}
	return identityrepo.CreatedUserRegistration{IDUser: idUser, IDEmail: idEmail}, nil
}

func (r User) FindPrimaryEmailIDByUser(ctx context.Context, id int64) (int64, error) {
	const query = `
		SELECT e.idemails
		FROM users u
		JOIN email_person ep ON ep.id_people = u.idusers
		JOIN emails e ON e.idemails = ep.id_emails AND e.is_primary = TRUE AND e.deleted_at IS NULL
		WHERE u.idusers = $1 AND u.deleted_at IS NULL
		LIMIT 1;`
	var idEmail int64
	if err := r.db.QueryRow(ctx, query, id).Scan(&idEmail); err != nil {
		return 0, err
	}
	return idEmail, nil
}

func (r User) UpdatePasswordHash(ctx context.Context, id int64, passwordHash string) error {
	const query = `UPDATE users SET password_hash = $1, updated_at = NOW() WHERE idusers = $2 AND deleted_at IS NULL`
	ct, err := r.db.Exec(ctx, query, passwordHash, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r User) UpdateLastLoginAt(ctx context.Context, id int64, lastLoginAt time.Time) error {
	const query = `UPDATE users SET last_login_at = $1, updated_at = NOW() WHERE idusers = $2 AND deleted_at IS NULL`
	ct, err := r.db.Exec(ctx, query, lastLoginAt, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r User) MarkPrimaryEmailVerified(ctx context.Context, idEmail int64, verifiedAt time.Time) error {
	const query = `UPDATE emails SET verified_at = $1, updated_at = NOW() WHERE idemails = $2 AND deleted_at IS NULL AND verified_at IS NULL`
	_, err := r.db.Exec(ctx, query, verifiedAt, idEmail)
	return err
}
