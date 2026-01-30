// Package database
package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/paladignus/actajus/internal/infrastructure/persistence/postgres"
	"github.com/paladignus/actajus/internal/module/social_media/domain"
	sharedDomain "github.com/paladignus/actajus/internal/shared/domain"
)

type SocialMedia struct {
	pool postgres.PgxPool
}

func NewSocialMedia(pool postgres.PgxPool) *SocialMedia {
	return &SocialMedia{
		pool,
	}
}

func (e *SocialMedia) Create(ctx context.Context, socialMedia *domain.SocialMedia) error {
	query := `INSERT INTO social_media (id_companies, platform, url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING idsocial_media`
	var id uint
	if err := e.pool.QueryRow(ctx, query,
		socialMedia.IDCompany(),
		socialMedia.Platform(),
		socialMedia.URL(),
		socialMedia.CreatedAt(),
		socialMedia.UpdatedAt()).Scan(&id); err != nil {
		return err
	}
	return socialMedia.SetID(id)
}

func (e *SocialMedia) Update(ctx context.Context, socialMedia *domain.SocialMedia) error {
	query := `UPDATE social_media SET platform = $1, url = $2, updated_at = $3 WHERE idsocial_media = $4 AND deleted_at IS NULL`
	_, err := e.pool.Exec(ctx, query,
		socialMedia.Platform(),
		socialMedia.URL(),
		socialMedia.UpdatedAt(),
		socialMedia.ID(),
	)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return sharedDomain.NewFieldError("social_media", "social media already exists")
	}
	return err
}
