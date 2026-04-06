// Package database
package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/paladignus/actajus/internal/module/social_media/domain"
	sharedDomain "github.com/paladignus/actajus/internal/shared/domain"
	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type SocialMedia struct {
	db postgres.Executor
}

func NewSocialMedia(db postgres.Executor) SocialMedia {
	return SocialMedia{
		db,
	}
}

func (s SocialMedia) Create(ctx context.Context, socialMedia *domain.SocialMedia) error {
	query := `INSERT INTO social_media (id_companies, platform, url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING idsocial_media`
	var id int64
	if err := s.db.QueryRow(ctx, query,
		socialMedia.IDCompany(),
		socialMedia.Platform(),
		socialMedia.URL(),
		socialMedia.CreatedAt(),
		socialMedia.UpdatedAt()).Scan(&id); err != nil {
		return err
	}
	return socialMedia.SetID(id)
}

func (s SocialMedia) Update(ctx context.Context, socialMedia domain.SocialMedia) error {
	query := `UPDATE social_media SET platform = $1, url = $2, updated_at = $3 WHERE idsocial_media = $4 AND deleted_at IS NULL`
	_, err := s.db.Exec(ctx, query,
		socialMedia.Platform(),
		socialMedia.URL(),
		socialMedia.UpdatedAt(),
		socialMedia.ID().Value(),
	)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return sharedDomain.NewFieldError("social_media", "social media already exists")
	}
	return err
}

func (s SocialMedia) Delete(ctx context.Context, socialMedia domain.SocialMedia) error {
	query := `UPDATE social_media SET updated_at = $1, deleted_at = $2 WHERE idsocial_media = $3 AND deleted_at IS NULL`
	_, err := s.db.Exec(ctx, query,
		socialMedia.UpdatedAt(),
		socialMedia.DeletedAt(),
		socialMedia.ID().Value())
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return nil
	}
	return err
}

func (s SocialMedia) DeleteByIDCompany(ctx context.Context, idCompany int64) error {
	now := time.Now()
	query := `UPDATE social_media SET updated_at = $1, deleted_at = $2 WHERE id_companies = $3 AND deleted_at IS NULL`
	_, err := s.db.Exec(ctx, query, now, now, idCompany)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return nil
	}
	return err
}

func (s SocialMedia) FindByIDCompany(ctx context.Context, idCompany int64) ([]*domain.SocialMedia, error) {
	query := `
			SELECT
				idsocial_media, platform, url, created_at, updated_at
			FROM social_media WHERE id_companies = $1 AND deleted_at IS NULL`
	var (
		idsocialMedia int64
		platform      string
		url           string
		createdAt     time.Time
		updatedAt     time.Time
	)
	rows, err := s.db.Query(ctx, query, idCompany)
	if err != nil {
		return nil, err
	}
	var socialMedia []*domain.SocialMedia
	for rows.Next() {
		if err := rows.Scan(&idsocialMedia, &platform, &url, &createdAt, &updatedAt); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, nil
			}
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
				return nil, sharedDomain.NewFieldError("id_company", "idCompany")
			}
			return nil, err
		}
		sm, err := domain.NewSocialMediaBuilder().
			WithID(idsocialMedia).
			WithIDCompany(idCompany).
			WithPlatform(platform).
			WithURL(url).
			WithCreatedAt(createdAt).
			WithUpdatedAt(updatedAt).
			Build()
		if err != nil {
			return nil, err
		}
		socialMedia = append(socialMedia, sm)
	}
	return socialMedia, nil
}
