// Package command
package command

type CreateCompanyAddressCommand struct {
	ZIP          string
	Title        string
	Street       string
	Number       uint
	Complement   string
	Reference    string
	Neighborhood string
	City         string
	State        string
	Country      string
}

type UpdateCompanyAddressCommand struct {
	IDAddress    int64
	ZIP          string
	Title        string
	Street       string
	Number       uint
	Complement   string
	Reference    string
	Neighborhood string
	City         string
	State        string
	Country      string
}

type CreateCompanyPhoneCommand struct {
	Number     string
	Kind       string
	Department string
}

type UpdateCompanyPhoneCommand struct {
	IDPhone    int64
	Number     string
	Kind       string
	Department string
}

type CreateCompanyEmailCommand struct {
	Address string
}

type UpdateCompanyEmailCommand struct {
	IDEmail int64
	Address string
}

type CreateCompanySocialMediaCommand struct {
	Platform string
	URL      string
}

type UpdateCompanySocialMediaCommand struct {
	IDSocialMedia int64
	Platform      string
	URL           string
}
