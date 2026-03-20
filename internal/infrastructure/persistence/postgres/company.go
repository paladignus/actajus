// Package postgres
package postgres

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/application/readmodel"
	"github.com/paladignus/actajus/internal/domain/entity"
)

type Company struct {
	pool PgxPool
}

func NewCompany(pool PgxPool) Company {
	return Company{pool}
}

func (c Company) Create(ctx context.Context, company entity.Company) (id uint, err error) {
	sql := `INSERT INTO companies (registered_by, name, trade_name, cnpj) VALUES ($1, $2, $3, $4) returning idcompanies;`
	if err = c.pool.QueryRow(ctx, sql, company.RegisteredBy, company.Name, company.TradeName, company.CNPJ).Scan(&id); err != nil {
		return 0, fmt.Errorf("database error while saving company: %w", err)
	}
	return id, nil
}

func (c Company) Update(ctx context.Context, company entity.Company) error {
	sql := `UPDATE companies SET name = $1, trade_name = $2, cnpj = $3, updated_at = now() WHERE idcompanies = $4;`
	if _, err := c.pool.Exec(ctx, sql, company.Name, company.TradeName, company.CNPJ, company.IDCompany); err != nil {
		return fmt.Errorf("database error while update company: %w", err)
	}
	return nil
}

func (c Company) GetAll(ctx context.Context) (companies []readmodel.CompanyReadModel, err error) {
	sql := `SELECT c.idcompanies, c.registered_by, c.name, c.trade_name, c.cnpj,
		a.idaddresses, a.zip, a.title, a.street,  a.complement, a.reference, a.neighborhood, a.city, a.state, a.country
		FROM companies c
		LEFT JOIN company_address ca ON c.idcompanies = ca.id_companies
		LEFT JOIN addresses a ON a.idaddresses = ca.id_addresses AND a.deleted_at IS NULL
		WHERE c.deleted_at IS NULL;`
	rows, err := c.pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("database error while getting all companies: %w", err)
	}
	for rows.Next() {
		var company readmodel.CompanyReadModel
		if err = rows.Scan(&company.IDCompany, &company.RegisteredBy, &company.Name, &company.TradeName, &company.CNPJ,
			&company.IDAddress, &company.Zip, &company.Title, &company.Street, &company.Complement, &company.Reference, &company.Neighborhood, &company.City, &company.State, &company.Country,
		); err != nil {
			return nil, fmt.Errorf("database error while getting all companies: %w", err)
		}
		companies = append(companies, company)
	}
	return companies, nil
}

func (c Company) Delete(ctx context.Context, idCompany uint) error {
	sql := `UPDATE companies SET deleted_at = now() WHERE idcompanies = $1;`
	if _, err := c.pool.Exec(ctx, sql, idCompany); err != nil {
		return fmt.Errorf("database error while delete company: %w", err)
	}
	return nil
}
