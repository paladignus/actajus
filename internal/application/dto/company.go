// Package dto
package dto

type Company struct {
	IDCompany    int    `json:"id_company"`
	RegisteredBy int    `json:"registered_by"`
	Name         string `json:"name"`
	TradeName    string `json:"trade_name"`
	CNPJ         string `json:"cnpj"`
}

type CompanyInputOutput struct {
	Company     `json:"company"`
	Address     `json:"address,omitzero"`
	Phone       `json:"phone,omitzero"`
	Email       `json:"email,omitzero"`
	SocialMedia []SocialMedia `json:"social_media,omitzero"`
}
