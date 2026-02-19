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
        c.idcompanies, c.name, c.trade_name, c.cnpj, c.created_at, c.updated_at,
        a.idaddresses, a.zip, a.title, a.street, a.complement, a.reference, a.number, a.neighborhood,
        a.city, a.state, a.country, a.created_at, a.updated_at,
        e.idemails, e.address, e.created_at, e.updated_at,
        p.idphones, p.number, p.kind, p.department, p.created_at, p.updated_at,
        sm.idsocial_media, sm.platform, sm.url, sm.created_at, sm.updated_at,
				CONCAT(pe.first_name, ' ', pe.last_name) AS name
    FROM (
        SELECT idcompanies, registered_by, name, trade_name, cnpj, created_at, updated_at
        FROM companies
        %s
        ORDER BY %s
        LIMIT $%d
    ) c
    LEFT JOIN company_address ca ON c.idcompanies = ca.id_companies AND ca.ended_at IS NULL
    LEFT JOIN addresses a ON ca.id_addresses = a.idaddresses
    LEFT JOIN company_email ce ON c.idcompanies = ce.id_companies
    LEFT JOIN emails e ON ce.id_emails = e.idemails AND e.deleted_at IS NULL
    LEFT JOIN company_phone cp ON c.idcompanies = cp.id_companies AND cp.ended_at IS NULL
    LEFT JOIN phones p ON cp.id_phones = p.idphones
    LEFT JOIN social_media sm ON sm.id_companies = c.idcompanies AND sm.deleted_at IS NULL
		LEFT JOIN people pe ON c.registered_by = pe.idpeople
    ORDER BY c.created_at ASC, c.idcompanies ASC;`
	var whereClause string
	var innerOrder string
	if after != nil {
		cursorData, err := postgres.DecodeCursor(*after)
		if err != nil {
			return nil, fmt.Errorf("invalid after cursor: %w", err)
		}
		whereClause = fmt.Sprintf(`WHERE (created_at, idcompanies) > ($%d, $%d)`, argPosition, argPosition+1)
		args = append(args, cursorData.Timestamp, cursorData.ID)
		argPosition += 2
		innerOrder = "created_at ASC, idcompanies ASC"
	} else if before != nil {
		cursorData, err := postgres.DecodeCursor(*before)
		if err != nil {
			return nil, fmt.Errorf("invalid before cursor: %w", err)
		}
		whereClause = fmt.Sprintf(`WHERE (created_at, idcompanies) < ($%d, $%d)`, argPosition, argPosition+1)
		args = append(args, cursorData.Timestamp, cursorData.ID)
		argPosition += 2
		innerOrder = "created_at DESC, idcompanies DESC" // pega os mais próximos do cursor
	} else {
		innerOrder = "created_at ASC, idcompanies ASC"
	}
	query = fmt.Sprintf(baseQuery, whereClause, innerOrder, argPosition)
	args = append(args, limit+1)
	fmt.Println(query, args)
	rows, err := c.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query companies: %w", err)
	}
	defer rows.Close()
	companiesMap := make(map[uint]*dto.CompanyReadModel)
	order := make([]uint, 0, limit)
	for rows.Next() {
		var (
			cp companyScan
			a  addressScan
			e  emailScan
			p  phoneScan
			sm socialMediaScan
		)
		if err := rows.Scan(
			&cp.id, &cp.name, &cp.tradeName, &cp.cnpj, &cp.createdAt, &cp.updatedAt,
			&a.id, &a.zip, &a.title, &a.street, &a.complement, &a.reference,
			&a.number, &a.neighborhood, &a.city, &a.state, &a.country, &a.createdAt, &a.updatedAt,
			&e.id, &e.address, &e.createdAt, &e.updatedAt,
			&p.id, &p.number, &p.kind, &p.department, &p.createdAt, &p.updatedAt,
			&sm.id, &sm.platform, &sm.url, &sm.createdAt, &sm.updatedAt, &cp.registeredByName,
		); err != nil {
			return nil, err
		}
		company := cp.companyToDTO()
		existing, seen := companiesMap[company.ID]
		if !seen {
			companiesMap[company.ID] = company
			order = append(order, company.ID)
			existing = company
		}
		if existing.Addresses == nil {
			if addr := a.addressToDTO(); addr != nil {
				existing.Addresses = addr
			}
		}
		if e.id != nil && !hasEmail(existing.Emails, *e.id) {
			existing.Emails = append(existing.Emails, &emailDTO.EmailReadModel{
				ID:        *e.id,
				Address:   *e.address,
				CreatedAt: *e.createdAt,
				UpdatedAt: *e.updatedAt,
			})
		}
		if p.id != nil && !hasPhone(existing.Phones, *p.id) {
			existing.Phones = append(existing.Phones, &phoneDTO.PhoneReadModel{
				ID:         *p.id,
				Number:     *p.number,
				Kind:       *p.kind,
				Department: *p.department,
				CreatedAt:  *p.createdAt,
				UpdatedAt:  *p.updatedAt,
			})
		}
		if sm.id != nil && !hasSocialMedia(existing.SocialMedia, *sm.id) {
			existing.SocialMedia = append(existing.SocialMedia, &socialMediaDTO.SocialMediaReadModel{
				ID:        *sm.id,
				IDCompany: existing.ID,
				Platform:  *sm.platform,
				URL:       *sm.url,
				CreatedAt: *sm.createdAt,
				UpdatedAt: *sm.updatedAt,
			})
		}
	}
	companies := make([]dto.CompanyReadModel, 0, len(order))
	for _, id := range order {
		companies = append(companies, *companiesMap[id])
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
