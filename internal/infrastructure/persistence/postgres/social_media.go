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
