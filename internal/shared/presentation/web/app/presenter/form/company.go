package form

import (
	"strconv"

	companyread "github.com/paladignus/actajus/internal/module/company/application/readmodel"
	companybinder "github.com/paladignus/actajus/internal/shared/presentation/web/app/binder/http/company"
)

type CompanySocialMediaItem struct {
	ID       int64
	Platform string
	URL      string
}

type CompanyData struct {
	Mode            string
	Action          string
	SubmitLabel     string
	Error           string
	CompanyID       int64
	Name            string
	TradeName       string
	CNPJ            string
	AddressID       int64
	ZIP             string
	Title           string
	Street          string
	Number          string
	Complement      string
	Reference       string
	Neighborhood    string
	City            string
	State           string
	Country         string
	PhoneID         int64
	PhoneNumber     string
	PhoneKind       string
	PhoneDepartment string
	EmailID         int64
	EmailAddress    string
	SocialMedias    []CompanySocialMediaItem
}

func NewCompanyData() CompanyData {
	return fromBinder(companybinder.NewFormData())
}

func CompanyDataFromBinder(input companybinder.FormData) CompanyData {
	return fromBinder(input)
}

func CompanyDataFromReadModel(company *companyread.CompanyReadModel) CompanyData {
	data := CompanyData{
		Mode:        "edit",
		Action:      "/companies/" + strconv.FormatInt(company.ID, 10),
		SubmitLabel: "Salvar alteracoes",
		CompanyID:   company.ID,
		Name:        company.Name,
		TradeName:   company.TradeName,
		CNPJ:        company.CNPJ,
	}
	if company.Address != nil {
		data.AddressID = company.Address.ID
		data.ZIP = company.Address.ZIP
		data.Title = company.Address.Title
		data.Street = company.Address.Street
		data.Number = strconv.FormatUint(uint64(company.Address.Number), 10)
		if company.Address.Complement != nil {
			data.Complement = *company.Address.Complement
		}
		if company.Address.Reference != nil {
			data.Reference = *company.Address.Reference
		}
		data.Neighborhood = company.Address.Neighborhood
		data.City = company.Address.City
		data.State = company.Address.State
		data.Country = company.Address.Country
	}
	if len(company.Phones) > 0 {
		data.PhoneID = company.Phones[0].ID
		data.PhoneNumber = company.Phones[0].Number
		data.PhoneKind = company.Phones[0].Kind
		data.PhoneDepartment = company.Phones[0].Department
	}
	if len(company.Emails) > 0 {
		data.EmailID = company.Emails[0].ID
		data.EmailAddress = company.Emails[0].Address
	}
	if len(company.SocialMedia) > 0 {
		data.SocialMedias = make([]CompanySocialMediaItem, 0, len(company.SocialMedia))
		for _, item := range company.SocialMedia {
			data.SocialMedias = append(data.SocialMedias, CompanySocialMediaItem{
				ID:       item.ID,
				Platform: item.Platform,
				URL:      item.URL,
			})
		}
	}
	if len(data.SocialMedias) == 0 {
		data.SocialMedias = []CompanySocialMediaItem{{}}
	}
	return data
}

func fromBinder(input companybinder.FormData) CompanyData {
	data := CompanyData{
		Mode:            input.Mode,
		Action:          input.Action,
		SubmitLabel:     input.SubmitLabel,
		Error:           input.Error,
		CompanyID:       input.CompanyID,
		Name:            input.Name,
		TradeName:       input.TradeName,
		CNPJ:            input.CNPJ,
		AddressID:       input.AddressID,
		ZIP:             input.ZIP,
		Title:           input.Title,
		Street:          input.Street,
		Number:          input.Number,
		Complement:      input.Complement,
		Reference:       input.Reference,
		Neighborhood:    input.Neighborhood,
		City:            input.City,
		State:           input.State,
		Country:         input.Country,
		PhoneID:         input.PhoneID,
		PhoneNumber:     input.PhoneNumber,
		PhoneKind:       input.PhoneKind,
		PhoneDepartment: input.PhoneDepartment,
		EmailID:         input.EmailID,
		EmailAddress:    input.EmailAddress,
		SocialMedias:    make([]CompanySocialMediaItem, 0, len(input.SocialMedias)),
	}
	for _, item := range input.SocialMedias {
		data.SocialMedias = append(data.SocialMedias, CompanySocialMediaItem{
			ID:       item.ID,
			Platform: item.Platform,
			URL:      item.URL,
		})
	}
	return data
}
