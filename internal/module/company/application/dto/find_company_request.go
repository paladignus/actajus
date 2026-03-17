// Package dto
package dto

type FindCompanyByCNPJRequest struct {
	CNPJ string `json:"cnpj"`
}

type FindCompanyByIDRequest struct {
	ID int64 `json:"id"`
}
