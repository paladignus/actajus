// Package database
package database

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	addrDTO "github.com/paladignus/actajus/internal/module/address/application/dto"
	"github.com/paladignus/actajus/internal/module/company/application/dto"
	"github.com/paladignus/actajus/internal/module/company/domain"
	emailDTO "github.com/paladignus/actajus/internal/module/email/application/dto"
	phoneDTO "github.com/paladignus/actajus/internal/module/phone/application/dto"
	socialMediaDTO "github.com/paladignus/actajus/internal/module/social_media/application/dto"
	sharedDto "github.com/paladignus/actajus/internal/shared/application/dto"
	sharedDomain "github.com/paladignus/actajus/internal/shared/domain"
	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/postgres"
)

type Company struct {
	pool postgres.PgxPool
}

func NewCompany(pool postgres.PgxPool) Company {
	return Company{
		pool,
	}
}

func (c Company) Create(ctx context.Context, company *domain.Company) error {
	query := `
		INSERT INTO companies (registered_by, name, trade_name, cnpj, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING idcompanies`
	var id uint
	if err := c.pool.QueryRow(ctx, query,
		company.RegisteredBy(),
		company.Name().Value(),
		company.TradeName().Value(),
		company.CNPJ().Value(),
		company.CreatedAt(),
		company.UpdatedAt(),
	).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return sharedDomain.NewFieldError("cnpj", "cnpj already exists")
		}
		return err
	}
	return company.SetID(id)
}

func (c Company) Update(ctx context.Context, company *domain.Company) error {
	query := `
		UPDATE companies
		SET name = $1, trade_name = $2, cnpj = $3, updated_at = $4 
		WHERE idcompanies = $5 AND deleted_at IS NULL`
	_, err := c.pool.Exec(ctx, query,
		company.Name().Value(),
		company.TradeName().Value(),
		company.CNPJ().Value(),
		company.UpdatedAt(),
		company.ID(),
	)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return sharedDomain.NewFieldError("cnpj", "cnpj already exists")
	}
	return err
}

func (c Company) Delete(ctx context.Context, company *domain.Company) error {
	query := `UPDATE companies SET updated_at = $1, deleted_at = $2 WHERE idcompanies = $3 AND deleted_at IS NULL`
	_, err := c.pool.Exec(ctx, query,
		company.UpdatedAt(),
		company.DeletedAt(),
		company.ID(),
	)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return sharedDomain.NewFieldError("deleted_at", "company already deleted")
	}
	return err
}

func (c Company) List(ctx context.Context, after, before *string, limit int, baseURL string) (*dto.CompanyListReadModel, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	var query string
	var args []any
	argPosition := 1
	baseQuery := `
		SELECT
			c.idcompanies, c.registered_by, c.name, c.trade_name, c.cnpj, c.created_at, c.updated_at,
			a.idaddresses, a.zip, a.title, a.street, a.complement, a.reference, a.number, a.neighborhood,
			a.city, a.state, a.country, a.created_at, a.updated_at, e.idemails, e.address, e.created_at, e.updated_at,
			p.idphones, p.number, p.kind, p.department, p.created_at, p.updated_at
		FROM companies c
		LEFT JOIN company_address ca ON c.idcompanies = ca.id_companies AND ca.ended_at IS NULL
		LEFT JOIN addresses a ON ca.id_addresses = a.idaddresses
		LEFT JOIN company_email ce ON c.idcompanies = ce.id_companies
		LEFT JOIN emails e ON ce.id_emails = e.idemails AND e.deleted_at IS NULL
		LEFT JOIN company_phone cp ON c.idcompanies = cp.id_companies
		LEFT JOIN phones p ON cp.id_phones = p.idphones
	`
	var whereClause string
	var orderClause string
	if after != nil {
		cursorData, err := postgres.DecodeCursor(*after)
		if err != nil {
			return nil, fmt.Errorf("invalid after cursor: %w", err)
		}
		whereClause = fmt.Sprintf(`
            WHERE (c.created_at > $%d AND c.idcompanies > $%d)
        `, argPosition, argPosition+1)
		args = append(args, cursorData.Timestamp, cursorData.ID)
		argPosition += 2
		orderClause = "ORDER BY c.created_at ASC, c.idcompanies ASC"
	} else if before != nil {
		cursorData, err := postgres.DecodeCursor(*before)
		if err != nil {
			return nil, fmt.Errorf("invalid before cursor: %w", err)
		}
		whereClause = fmt.Sprintf(`
            WHERE (c.created_at < $%d AND c.idcompanies < $%d)
        `, argPosition, argPosition+1)
		args = append(args, cursorData.Timestamp, cursorData.ID)
		argPosition += 2
		orderClause = "ORDER BY c.created_at DESC, c.idcompanies DESC"
	} else {
		whereClause = ""
		orderClause = "ORDER BY c.created_at ASC, c.idcompanies ASC"
	}
	query = fmt.Sprintf("%s %s %s LIMIT $%d",
		baseQuery, whereClause, orderClause, argPosition)
	args = append(args, limit+1)
	rows, err := c.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query companies: %w", err)
	}
	defer rows.Close()
	companies := make([]dto.CompanyReadModel, 0, limit)
	for rows.Next() {
		var (
			cp     dto.CompanyReadModel
			a      addrDTO.AddressReadModel
			addrID *uint
			// idAddress       *uint
			// zip             *string
			// title           *string
			// street          *string
			// complement      *string
			// reference       *string
			// number          *uint
			// neighborhood    *string
			// city            *string
			// state           *string
			// country         *string
			// addresCreatedAt *time.Time
			// addresUpdatedAt *time.Time

			idEmail        *uint
			address        *string
			emailCreatedAt *time.Time
			emailUpdatedAt *time.Time

			idPhone        *uint
			phoneNumber    *string
			kind           *string
			department     *string
			phoneCreatedAt *time.Time
			phoneUpdatedAt *time.Time
		)
		if err := rows.Scan(
			&cp.ID, &cp.RegisteredBy, &cp.Name, &cp.TradeName, &cp.CNPJ, &cp.CreatedAt, &cp.UpdatedAt,
			&addrID, &a.ZIP, &a.Title, &a.Street, &a.Complement, &a.Reference, &a.Number, &a.Neighborhood, &a.City,
			&a.State, &a.Country, &a.CreatedAt, &a.UpdatedAt, &idEmail, &address, &emailCreatedAt, &emailUpdatedAt,
			&idPhone, &phoneNumber, &kind, &department, &phoneCreatedAt, &phoneUpdatedAt,
		); err != nil {
			return nil, err
		}
		if addrID != nil {
			a.ID = *addrID
			cp.Address = &a
		}
		// if idAddress != nil {
		// 	cp.Address = &addrDTO.AddressReadModel{
		// 		ID:           *idAddress,
		// 		ZIP:          *zip,
		// 		Title:        *title,
		// 		Street:       *street,
		// 		Number:       *number,
		// 		Complement:   *complement,
		// 		Reference:    *reference,
		// 		Neighborhood: *neighborhood,
		// 		City:         *city,
		// 		State:        *state,
		// 		Country:      *country,
		// 		CreatedAt:    addresCreatedAt.Format(time.RFC3339),
		// 		UpdatedAt:    addresUpdatedAt.Format(time.RFC3339),
		// 	}
		// }
		if idEmail != nil {
			cp.Email = &emailDTO.EmailReadModel{
				ID:        *idEmail,
				Address:   *address,
				CreatedAt: emailCreatedAt.Format(time.RFC3339),
				UpdatedAt: emailUpdatedAt.Format(time.RFC3339),
			}
		}
		if idPhone != nil {
			cp.Phone = &phoneDTO.PhoneReadModel{
				ID:         *idPhone,
				Number:     *phoneNumber,
				Kind:       *kind,
				Department: *department,
				CreatedAt:  phoneCreatedAt.Format(time.RFC3339),
				UpdatedAt:  phoneUpdatedAt.Format(time.RFC3339),
			}
		}
		query := `
				SELECT idsocial_media, platform, url, created_at, updated_at
				FROM social_media WHERE id_companies = $1 AND deleted_at IS NULL`
		rows, err := c.pool.Query(ctx, query, cp.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var (
				idSocialMedia uint
				platform      string
				url           string
				createdAt     time.Time
				updatedAt     time.Time
			)
			if err := rows.Scan(&idSocialMedia, &platform, &url, &createdAt, &updatedAt); err != nil {
				return nil, err
			}
			cp.SocialMedia = append(cp.SocialMedia, &socialMediaDTO.SocialMediaReadModel{
				ID:        idSocialMedia,
				IDCompany: cp.ID,
				Platform:  platform,
				URL:       url,
				CreatedAt: createdAt.Format(time.RFC3339),
				UpdatedAt: updatedAt.Format(time.RFC3339),
			})
		}
		companies = append(companies, cp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	hasNextPage := len(companies) > limit
	hasPreviousPage := after != nil || before != nil
	if hasNextPage {
		companies = companies[:limit]
	}
	if before != nil {
		reversePersons(companies)
	}
	pageInfo := sharedDto.PageInfo{
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
	}
	if len(companies) > 0 {
		id := strconv.FormatUint(uint64(companies[0].ID), 10)
		idLen := strconv.FormatUint(uint64(companies[len(companies)-1].ID), 10)
		if hasNextPage {
			endCursor := postgres.EncodeCursor(
				companies[len(companies)-1].CreatedAt,
				idLen,
			)
			nextURL := postgres.BuildPaginationURL(baseURL, endCursor, "next", limit)
			pageInfo.NextURL = &nextURL
		}
		if hasPreviousPage {
			startCursor := postgres.EncodeCursor(
				companies[0].CreatedAt,
				id,
			)
			prevURL := postgres.BuildPaginationURL(baseURL, startCursor, "prev", limit)
			pageInfo.PreviousURL = &prevURL
		}
	}
	return &dto.CompanyListReadModel{
		Data:     companies,
		PageInfo: pageInfo,
	}, nil
}

func reversePersons(persons []dto.CompanyReadModel) {
	for i, j := 0, len(persons)-1; i < j; i, j = i+1, j-1 {
		fmt.Println(persons[i], persons[j])
		persons[i], persons[j] = persons[j], persons[i]
		fmt.Println(persons[i], persons[j])
	}
}

func (c Company) FindByCNPJ(ctx context.Context, cnpj string) (*domain.Company, error) {
	query := `
			SELECT
					idcompanies, registered_by, name, trade_name, created_at, updated_at
			FROM companies
			WHERE cnpj = $1 AND deleted_at IS NULL`
	var (
		id           uint
		registeredBy uint
		name         string
		tradeName    string
		createdAt    time.Time
		updatedAt    time.Time
	)
	err := c.pool.QueryRow(ctx, query, cnpj).
		Scan(
			&id, &registeredBy, &name, &tradeName, &createdAt, &updatedAt,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sharedDomain.NewFieldError("cnpj", "cnpj not found")
		}
		return nil, err
	}
	return domain.NewCompanyBuilder().
		WithID(id).
		WithRegisteredBy(registeredBy).
		WithName(name).
		WithTradeName(tradeName).
		WithCNPJ(cnpj).
		WithCreatedAt(createdAt).
		WithUpdatedAt(updatedAt).
		Build()
}

func (c Company) FindByID(ctx context.Context, id uint) (*domain.Company, error) {
	query := `
			SELECT
					registered_by, name, trade_name, cnpj, created_at, updated_at
			FROM companies
			WHERE idcompanies = $1 AND deleted_at IS NULL`
	var (
		registeredBy uint
		name         string
		tradeName    string
		cnpj         string
		createdAt    time.Time
		updatedAt    time.Time
	)
	err := c.pool.QueryRow(ctx, query, id).
		Scan(
			&registeredBy, &name, &tradeName, &cnpj, &createdAt, &updatedAt,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sharedDomain.NewFieldError("cnpj", "cnpj not found")
		}
		return nil, err
	}
	return domain.NewCompanyBuilder().
		WithID(id).
		WithRegisteredBy(registeredBy).
		WithName(name).
		WithTradeName(tradeName).
		WithCNPJ(cnpj).
		WithCreatedAt(createdAt).
		WithUpdatedAt(updatedAt).
		Build()
}
