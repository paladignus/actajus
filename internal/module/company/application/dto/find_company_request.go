// Package dto
package dto

type FindCompanyByCNPJRequest struct {
	CNPJ string `json:"cnpj"`
}

type FindCompanyByIDRequest struct {
	ID string `json:"id"`
}
