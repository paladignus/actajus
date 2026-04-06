// Package postgres
package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	socialMediaDomain "github.com/paladignus/actajus/internal/module/social_media/domain"
	sharedDomain "github.com/paladignus/actajus/internal/shared/domain"
	sharedpostgres "github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type SocialMediaRepository struct {
	exec sharedpostgres.Executor
}

func NewSocialMediaRepository(exec sharedpostgres.Executor) SocialMediaRepository {
	return SocialMediaRepository{exec: exec}
}

func (r SocialMediaRepository) Create(ctx context.Context, socialMedia *socialMediaDomain.SocialMedia) error {
	query := `INSERT INTO social_media (id_companies, platform, url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING idsocial_media`
	var id int64
	if err := r.exec.QueryRow(ctx, query,
		socialMedia.IDCompany(),
		socialMedia.Platform(),
		socialMedia.URL(),
		socialMedia.CreatedAt(),
		socialMedia.UpdatedAt(),
	).Scan(&id); err != nil {
		return err
	}
	return socialMedia.SetID(id)
}

func (r SocialMediaRepository) Update(ctx context.Context, socialMedia socialMediaDomain.SocialMedia) error {
	query := `UPDATE social_media SET platform = $1, url = $2, updated_at = $3 WHERE idsocial_media = $4 AND deleted_at IS NULL`
	_, err := r.exec.Exec(ctx, query,
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

func (r SocialMediaRepository) Delete(ctx context.Context, socialMedia socialMediaDomain.SocialMedia) error {
	query := `UPDATE social_media SET updated_at = $1, deleted_at = $2 WHERE idsocial_media = $3 AND deleted_at IS NULL`
	_, err := r.exec.Exec(ctx, query,
		socialMedia.UpdatedAt(),
		socialMedia.DeletedAt(),
		socialMedia.ID().Value(),
	)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return nil
	}
	return err
}

func (r SocialMediaRepository) DeleteByIDCompany(ctx context.Context, idCompany int64) error {
	now := time.Now()
	query := `UPDATE social_media SET updated_at = $1, deleted_at = $2 WHERE id_companies = $3 AND deleted_at IS NULL`
	_, err := r.exec.Exec(ctx, query, now, now, idCompany)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return nil
	}
	return err
}

func (r SocialMediaRepository) FindByIDCompany(ctx context.Context, idCompany int64) ([]*socialMediaDomain.SocialMedia, error) {
	query := `
		SELECT
			idsocial_media, platform, url, created_at, updated_at
		FROM social_media WHERE id_companies = $1 AND deleted_at IS NULL`
	rows, err := r.exec.Query(ctx, query, idCompany)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var socialMedia []*socialMediaDomain.SocialMedia
	for rows.Next() {
		var (
			id        int64
			platform  string
			url       string
			createdAt time.Time
			updatedAt time.Time
		)
		if err := rows.Scan(&id, &platform, &url, &createdAt, &updatedAt); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, nil
			}
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
				return nil, sharedDomain.NewFieldError("id_company", "idCompany")
			}
			return nil, err
		}
		item, err := socialMediaDomain.NewSocialMediaBuilder().
			WithID(id).
			WithIDCompany(idCompany).
			WithPlatform(platform).
			WithURL(url).
			WithCreatedAt(createdAt).
			WithUpdatedAt(updatedAt).
			Build()
		if err != nil {
			return nil, err
		}
		socialMedia = append(socialMedia, item)
	}
	return socialMedia, rows.Err()
}
