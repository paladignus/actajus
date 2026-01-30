// Package dto
package dto

type DeleteCompanyRequest struct {
	IDCompany uint `json:"id_company" validate:"required|numeric"`
}
