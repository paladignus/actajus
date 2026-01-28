// Package database
package database

import (
	"context"

	"github.com/paladignus/actajus/internal/infrastructure/persistence/postgres"
	"github.com/paladignus/actajus/internal/module/social_media/domain"
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
	query := `INSERT INTO social_media (id_companies, name, url, created_at, updated_at)
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
