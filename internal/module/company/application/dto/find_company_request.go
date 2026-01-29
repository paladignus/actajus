// Package dto
package dto

type FindCompanyByCNPJRequest struct {
	CNPJ string `json:"cnpj"`
}

type FindCompanyByIDRequest struct {
	ID uint `json:"id"`
}
