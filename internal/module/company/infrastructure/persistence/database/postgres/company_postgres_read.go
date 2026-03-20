// Package postgres
package postgres

import (
	"context"
	"fmt"
	"strconv"

	"github.com/paladignus/actajus/internal/module/company/application/readmodel"
	sharedDto "github.com/paladignus/actajus/internal/shared/application/dto"
	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type CompanyReadRepository struct {
	db postgres.Executor
}

func NewCompanyReadRepository(db postgres.Executor) CompanyReadRepository {
	return CompanyReadRepository{
		db: db,
	}
}

func (r CompanyReadRepository) List(ctx context.Context, after, before *string, limit int, baseURL string) (*readmodel.CompanyListReadModel, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	baseQuery := `
	SELECT
        c.idcompanies, c.name, c.trade_name, c.cnpj, c.created_at, c.updated_at,
        a.idaddresses, a.zip, a.title, a.street, a.complement, a.reference, a.number, a.neighborhood,
        a.city, a.state, a.country, a.created_at, a.updated_at,
        e.idemails, e.address, e.created_at, e.updated_at,
        p.idphones, p.number, p.kind, p.department, p.created_at, p.updated_at,
        sm.idsocial_media, sm.platform, sm.url, sm.created_at, sm.updated_at,
				CONCAT(pe.first_name, ' ', pe.last_name) name
    FROM (
        SELECT idcompanies, registered_by, name, trade_name, cnpj, created_at, updated_at
        FROM companies
        WHERE deleted_at IS NULL %s
        ORDER BY %s
        LIMIT $%d
    ) c
    LEFT JOIN company_address ca ON c.idcompanies = ca.id_companies AND ca.ended_at IS NULL
    LEFT JOIN addresses a ON ca.id_addresses = a.idaddresses
    LEFT JOIN company_email ce ON c.idcompanies = ce.id_companies AND ce.ended_at IS NULL
    LEFT JOIN emails e ON ce.id_emails = e.idemails
    LEFT JOIN company_phone cp ON c.idcompanies = cp.id_companies AND cp.ended_at IS NULL
    LEFT JOIN phones p ON cp.id_phones = p.idphones
    LEFT JOIN social_media sm ON sm.id_companies = c.idcompanies AND sm.deleted_at IS NULL
		LEFT JOIN people pe ON c.registered_by = pe.idpeople;`
	var whereClause string
	var innerOrder string
	var args []any
	argPosition := 1
	if after != nil {
		cursorData, err := postgres.DecodeCursor(*after)
		if err != nil {
			return nil, fmt.Errorf("invalid after cursor: %w", err)
		}
		whereClause = fmt.Sprintf(`AND (created_at, idcompanies) > ($%d, $%d)`, argPosition, argPosition+1)
		args = append(args, cursorData.Timestamp, cursorData.ID)
		argPosition += 2
		innerOrder = "created_at ASC, idcompanies ASC"
	} else if before != nil {
		cursorData, err := postgres.DecodeCursor(*before)
		if err != nil {
			return nil, fmt.Errorf("invalid before cursor: %w", err)
		}
		whereClause = fmt.Sprintf(`AND (created_at, idcompanies) < ($%d, $%d)`, argPosition, argPosition+1)
		args = append(args, cursorData.Timestamp, cursorData.ID)
		argPosition += 2
		innerOrder = "created_at DESC, idcompanies DESC"
	} else {
		innerOrder = "created_at ASC, idcompanies ASC"
	}
	query := fmt.Sprintf(baseQuery, whereClause, innerOrder, argPosition)
	args = append(args, limit+1)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query companies: %w", err)
	}
	defer rows.Close()
	companiesMap := make(map[int64]*readmodel.CompanyReadModel)
	order := make([]int64, 0, limit)
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
		mergeCompanyScanWithRelations(&cp, &a, &e, &p, &sm, companiesMap, &order)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	companies := make([]readmodel.CompanyReadModel, 0, len(order))
	for _, id := range order {
		companies = append(companies, *companiesMap[id])
	}
	hasNextPage := len(companies) > limit
	hasPreviousPage := after != nil || before != nil
	if hasNextPage {
		companies = companies[:limit]
	}
	if before != nil {
		reverseSlice(companies)
	}
	pageInfo := sharedDto.PageInfo{
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
	}
	if len(companies) > 0 {
		id := strconv.FormatInt(int64(companies[0].ID), 10)
		idLen := strconv.FormatInt(int64(companies[len(companies)-1].ID), 10)
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
	return &readmodel.CompanyListReadModel{
		Data:     companies,
		PageInfo: pageInfo,
	}, nil
}
