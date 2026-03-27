// Package readmodel
package readmodel

import "github.com/paladignus/actajus/internal/shared/application/pagination"

type CompanyListReadModel struct {
	Data     []CompanyReadModel
	PageInfo pagination.PageInfo
}
