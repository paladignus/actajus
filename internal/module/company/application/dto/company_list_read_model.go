// Package dto
package dto

type CompanyListReadModel struct {
	Companies []CompanyReadModel `json:"companies"`
	Page      uint               `json:"page"`
	PageSize  uint               `json:"page_size"`
	Total     uint64             `json:"total"`
}
