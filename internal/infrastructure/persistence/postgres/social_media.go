// Package postgres
package postgres

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/domain/entity"
)

type SocialMedia struct {
	db PgxPool
}

func NewSocialMedia(db PgxPool) SocialMedia {
	return SocialMedia{db: db}
}

func (s SocialMedia) Create(ctx context.Context, input entity.SocialMedia) error {
	sql := `INSERT INTO social_media (id_companies, name, url) VALUES ($1, $2, $3)`
	if _, err := s.db.Exec(ctx, sql, input.IDEnterprise, input.Name, input.URL); err != nil {
		return fmt.Errorf("failed to create social media: %w", err)
	}
	return nil
}

func (s SocialMedia) Update(ctx context.Context, input entity.SocialMedia) error {
	sql := `UPDATE social_media SET name = $1, url = $2, updated_at = now() WHERE idsocial_media = $3`
	if _, err := s.db.Exec(ctx, sql, input.Name, input.URL, input.IDSocialMedia); err != nil {
		return fmt.Errorf("failed to update social media: %w", err)
	}
	return nil
}
