// Package database
package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	addrDTO "github.com/paladignus/actajus/internal/module/address/application/dto"
	"github.com/paladignus/actajus/internal/module/company/application/dto"
	"github.com/paladignus/actajus/internal/module/company/domain"
	emailDTO "github.com/paladignus/actajus/internal/module/email/application/dto"
	phoneDTO "github.com/paladignus/actajus/internal/module/phone/application/dto"
	socialMediaDTO "github.com/paladignus/actajus/internal/module/social_media/application/dto"
	sharedDomain "github.com/paladignus/actajus/internal/shared/domain"
	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/postgres"
)

type Company struct {
	pool postgres.PgxPool
}

func NewCompany(pool postgres.PgxPool) *Company {
	return &Company{
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

func (c Company) List(ctx context.Context, page, pageSize uint) (*dto.CompanyListReadModel, error) {
	offset := (page - 1) * pageSize
	query := `
			SELECT
					c.idcompanies, c.registered_by, c.name, c.trade_name, c.cnpj, c.created_at, c.updated_at,
					a.idaddresses, a.zip, a.title, a.street, a.complement, a.reference, a.number, a.neighborhood,
					a.city, a.state, a.country, a.created_at, a.updated_at, e.idemails, e.address, e.created_at, e.updated_at,
					p.idphones, p.number, p.kind, p.department, p.created_at, p.updated_at
			FROM companies c 
			LEFT JOIN company_address ca ON c.idcompanies = ca.id_companies
			LEFT JOIN addresses a ON ca.id_addresses = a.idaddresses AND a.deleted_at IS NULL
			LEFT JOIN company_email ce ON c.idcompanies = ce.id_companies
			LEFT JOIN emails e ON ce.id_emails = e.idemails AND e.deleted_at IS NULL
			LEFT JOIN company_phone cp ON c.idcompanies = cp.id_companies
			LEFT JOIN phones p ON cp.id_phones = p.idphones AND p.deleted_at IS NULL
			WHERE c.deleted_at IS NULL
			ORDER BY c.idcompanies DESC
			LIMIT $1 OFFSET $2`
	rows, err := c.pool.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	companies := []dto.CompanyReadModel{}
	for rows.Next() {
		var (
			idCompany        uint
			registeredBy     uint
			name             string
			tradeName        string
			cnpj             string
			companyCreatedAt time.Time
			companyUpdatedAt time.Time

			idAddress       *uint
			zip             *string
			title           *string
			street          *string
			complement      *string
			reference       *string
			number          *uint
			neighborhood    *string
			city            *string
			state           *string
			country         *string
			addresCreatedAt *time.Time
			addresUpdatedAt *time.Time

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
			&idCompany, &registeredBy, &name, &tradeName, &cnpj, &companyCreatedAt, &companyUpdatedAt,
			&idAddress, &zip, &title, &street, &complement, &reference, &number, &neighborhood, &city,
			&state, &country, &addresCreatedAt, &addresUpdatedAt, &idEmail, &address, &emailCreatedAt, &emailUpdatedAt,
			&idPhone, &phoneNumber, &kind, &department, &phoneCreatedAt, &phoneUpdatedAt,
		); err != nil {
			return nil, err
		}
		company := dto.CompanyReadModel{
			ID:           idCompany,
			RegisteredBy: registeredBy,
			Name:         name,
			TradeName:    tradeName,
			CNPJ:         cnpj,
			CreatedAt:    companyCreatedAt.Format(time.RFC3339),
			UpdatedAt:    companyUpdatedAt.Format(time.RFC3339),
		}
		if idAddress != nil {
			company.Address = &addrDTO.AddressReadModel{
				ID:           *idAddress,
				ZIP:          *zip,
				Title:        *title,
				Street:       *street,
				Number:       *number,
				Complement:   *complement,
				Reference:    *reference,
				Neighborhood: *neighborhood,
				City:         *city,
				State:        *state,
				Country:      *country,
				CreatedAt:    addresCreatedAt.Format(time.RFC3339),
				UpdatedAt:    addresUpdatedAt.Format(time.RFC3339),
			}
		}
		if idEmail != nil {
			company.Email = &emailDTO.EmailReadModel{
				ID:        *idEmail,
				Address:   *address,
				CreatedAt: emailCreatedAt.Format(time.RFC3339),
				UpdatedAt: emailUpdatedAt.Format(time.RFC3339),
			}
		}
		if idPhone != nil {
			company.Phone = &phoneDTO.PhoneReadModel{
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
		rows, err := c.pool.Query(ctx, query, idCompany)
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
			company.SocialMedia = append(company.SocialMedia, &socialMediaDTO.SocialMediaReadModel{
				ID:        idSocialMedia,
				IDCompany: idCompany,
				Platform:  platform,
				URL:       url,
				CreatedAt: createdAt.Format(time.RFC3339),
				UpdatedAt: updatedAt.Format(time.RFC3339),
			})
		}
		companies = append(companies, company)
	}
	var total uint64
	if err := c.pool.QueryRow(ctx, "SELECT count(*) FROM companies WHERE deleted_at IS NULL").Scan(&total); err != nil {
		return nil, err
	}

	return &dto.CompanyListReadModel{
		Companies: companies,
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
	}, nil
}

func (c *Company) FindByCNPJ(ctx context.Context, cnpj string) (*domain.Company, error) {
	return c.find(ctx, "cnpj", cnpj)
}

func (c *Company) FindByID(ctx context.Context, id uint) (*domain.Company, error) {
	sid := fmt.Sprintf("%d", id)
	return c.find(ctx, "id", sid)
}

func (c *Company) find(ctx context.Context, by, search string) (*domain.Company, error) {
	var where string
	switch by {
	case "id":
		where = `idcompanies = $1 AND deleted_at IS NULL`
	default:
		where = `cnpj = $1 AND deleted_at IS NULL`
	}
	query := `
			SELECT
					idcompanies, registered_by, name, trade_name, cnpj, created_at, updated_at
			FROM companies
			WHERE `
	var (
		id           uint
		registeredBy uint
		name         string
		tradeName    string
		cnpj         string
		createdAt    time.Time
		updatedAt    time.Time
	)
	err := c.pool.QueryRow(ctx, query+where, search).
		Scan(
			&id, &registeredBy, &name, &tradeName, &cnpj, &createdAt, &updatedAt,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return nil, sharedDomain.NewFieldError(by, search)
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
