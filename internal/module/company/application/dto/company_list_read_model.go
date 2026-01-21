// Package dto
package dto

type CompanyListReadModel struct {
	Companies []CompanyReadModel `json:"companies"`
	Page      int                `json:"page"`
	PageSize  int                `json:"page_size"`
	Total     int64              `json:"total"`
}
