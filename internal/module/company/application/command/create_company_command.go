// Package command
package command

type CreateCompanyCommand struct {
	Name         string
	TradeName    string
	CNPJ         string
	RegisteredBy int64
	Address      CreateCompanyAddressCommand
	Phone        CreateCompanyPhoneCommand
	Email        CreateCompanyEmailCommand
	SocialMedia  []CreateCompanySocialMediaCommand
}
