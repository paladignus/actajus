// Package command
package command

type UpdateCompanyCommand struct {
	IDCompany   int64
	Name        string
	TradeName   string
	CNPJ        string
	Address     UpdateCompanyAddressCommand
	Phone       UpdateCompanyPhoneCommand
	Email       UpdateCompanyEmailCommand
	SocialMedia []UpdateCompanySocialMediaCommand
}
