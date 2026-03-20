// Package readmodel
package readmodel

import "github.com/paladignus/actajus/internal/shared/application/dto"

type CompanyListReadModel struct {
	Data     []CompanyReadModel
	PageInfo dto.PageInfo
}
