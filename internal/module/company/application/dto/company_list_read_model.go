// Package dto
package dto

import "github.com/paladignus/actajus/internal/shared/application/dto"

type CompanyListReadModel struct {
	Data     []CompanyReadModel `json:"data"`
	PageInfo dto.PageInfo       `json:"page_info"`
}
