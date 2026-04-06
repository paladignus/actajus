// Package postgres
package postgres

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	"github.com/paladignus/actajus/internal/module/identity/application/repository"
	sharedpostgres "github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type SessionQuery struct {
	db sharedpostgres.Executor
}

func NewSessionQuery(db sharedpostgres.Executor) SessionQuery {
	return SessionQuery{db: db}
}

func (r SessionQuery) GetByID(ctx context.Context, id int64) (*readmodel.SessionReadModel, error) {
	const query = `
		SELECT
			s.idsessions,
			s.id_users,
			COALESCE(e.address, ''),
			COALESCE(s.ip, ''),
			COALESCE(s.user_agent, ''),
			s.expires_at,
			s.revoked_at,
			s.rotated_at,
			s.created_at,
			s.updated_at
		FROM sessions s
		LEFT JOIN users u ON u.idusers = s.id_users
		LEFT JOIN email_person ep ON ep.id_people = u.idusers
		LEFT JOIN emails e ON e.idemails = ep.id_emails AND e.is_primary = TRUE AND e.deleted_at IS NULL
		WHERE s.idsessions = $1
		LIMIT 1;`

	rm, err := scanSessionReadModel(r.db.QueryRow(ctx, query, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return rm, nil
}

func (r SessionQuery) ListByUser(ctx context.Context, uid int64) ([]readmodel.SessionReadModel, error) {
	const query = `
		SELECT
			s.idsessions,
			s.id_users,
			COALESCE(e.address, ''),
			COALESCE(s.ip, ''),
			COALESCE(s.user_agent, ''),
			s.expires_at,
			s.revoked_at,
			s.rotated_at,
			s.created_at,
			s.updated_at
		FROM sessions s
		LEFT JOIN users u ON u.idusers = s.id_users
		LEFT JOIN email_person ep ON ep.id_people = u.idusers
		LEFT JOIN emails e ON e.idemails = ep.id_emails AND e.is_primary = TRUE AND e.deleted_at IS NULL
		WHERE s.id_users = $1
		ORDER BY s.created_at DESC, s.idsessions DESC;`
	rows, err := r.db.Query(ctx, query, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return collectSessionReadModels(rows)
}

func (r SessionQuery) ListAll(ctx context.Context, filter repository.SessionQueryFilter) ([]readmodel.SessionReadModel, error) {
	limit := filter.Limit
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	emailQuery := "%" + strings.TrimSpace(strings.ToLower(filter.UserEmail)) + "%"
	const query = `
		SELECT
			s.idsessions,
			s.id_users,
			COALESCE(e.address, ''),
			COALESCE(s.ip, ''),
			COALESCE(s.user_agent, ''),
			s.expires_at,
			s.revoked_at,
			s.rotated_at,
			s.created_at,
			s.updated_at
		FROM sessions s
		LEFT JOIN users u ON u.idusers = s.id_users
		LEFT JOIN email_person ep ON ep.id_people = u.idusers
		LEFT JOIN emails e ON e.idemails = ep.id_emails AND e.is_primary = TRUE AND e.deleted_at IS NULL
		WHERE ($1 = '%%' OR LOWER(COALESCE(e.address, '')) LIKE $1)
		ORDER BY s.created_at DESC, s.idsessions DESC
		LIMIT $2;`
	rows, err := r.db.Query(ctx, query, emailQuery, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return collectSessionReadModels(rows)
}

type sessionScanner interface {
	Scan(dest ...any) error
}

func scanSessionReadModel(scanner sessionScanner) (*readmodel.SessionReadModel, error) {
	var (
		rm        readmodel.SessionReadModel
		revokedAt *time.Time
		rotatedAt *time.Time
	)
	if err := scanner.Scan(
		&rm.IDSession,
		&rm.IDUser,
		&rm.UserEmail,
		&rm.IP,
		&rm.UserAgent,
		&rm.ExpiresAt,
		&revokedAt,
		&rotatedAt,
		&rm.CreatedAt,
		&rm.UpdatedAt,
	); err != nil {
		return nil, err
	}
	rm.RevokedAt = revokedAt
	rm.RotatedAt = rotatedAt
	return &rm, nil
}

func collectSessionReadModels(rows pgx.Rows) ([]readmodel.SessionReadModel, error) {
	items := make([]readmodel.SessionReadModel, 0, 16)
	for rows.Next() {
		rm, err := scanSessionReadModel(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *rm)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
