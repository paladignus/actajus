// Package dto
package dto

type DeleteCompanyRequest struct {
	IDCompany int64 `json:"id_company" validate:"required|numeric"`
}
