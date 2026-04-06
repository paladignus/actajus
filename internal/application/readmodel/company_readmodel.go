// Package readmodel
package readmodel

type CompanyReadModel struct {
	IDCompany    int
	RegisteredBy int
	Name         string
	TradeName    string
	CNPJ         string
	IDAddress    uint
	Zip          string
	Title        string
	Street       string
	Complement   string
	Reference    string
	Neighborhood string
	City         string
	State        string
	Country      string
}
