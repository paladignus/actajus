package companybinder

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	companycmd "github.com/paladignus/actajus/internal/module/company/application/command"
	companyread "github.com/paladignus/actajus/internal/module/company/application/readmodel"
	companyrepo "github.com/paladignus/actajus/internal/module/company/application/repository"
)

type ListQuery struct {
	Page   int
	Limit  int
	Filter companyrepo.CompanyListFilter
	After  *string
	Before *string
}

type SocialMediaFormItem struct {
	ID       int64
	Platform string
	URL      string
}

type FormData struct {
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
	SocialMedias    []SocialMediaFormItem
}

func BindListQuery(r *http.Request) ListQuery {
	return ListQuery{
		Page:  parsePage(r.URL.Query().Get("page")),
		Limit: parseLimit(r.URL.Query().Get("limit")),
		Filter: companyrepo.CompanyListFilter{
			Name: strings.TrimSpace(r.URL.Query().Get("name")),
			CNPJ: strings.TrimSpace(r.URL.Query().Get("cnpj")),
		},
		After:  optionalString(r.URL.Query().Get("after")),
		Before: optionalString(r.URL.Query().Get("before")),
	}
}

func BindCreateForm(r *http.Request, idUser int64) (FormData, companycmd.CreateCompanyCommand, error) {
	if err := r.ParseForm(); err != nil {
		return FormData{}, companycmd.CreateCompanyCommand{}, err
	}
	number, err := parseUintForm(r.FormValue("number"))
	if err != nil {
		return FormData{}, companycmd.CreateCompanyCommand{}, err
	}
	form := NewFormData()
	form.Name = strings.TrimSpace(r.FormValue("name"))
	form.TradeName = strings.TrimSpace(r.FormValue("trade_name"))
	form.CNPJ = strings.TrimSpace(r.FormValue("cnpj"))
	form.ZIP = strings.TrimSpace(r.FormValue("zip"))
	form.Title = strings.TrimSpace(r.FormValue("title"))
	form.Street = strings.TrimSpace(r.FormValue("street"))
	form.Number = strings.TrimSpace(r.FormValue("number"))
	form.Complement = strings.TrimSpace(r.FormValue("complement"))
	form.Reference = strings.TrimSpace(r.FormValue("reference"))
	form.Neighborhood = strings.TrimSpace(r.FormValue("neighborhood"))
	form.City = strings.TrimSpace(r.FormValue("city"))
	form.State = strings.TrimSpace(r.FormValue("state"))
	form.Country = strings.TrimSpace(r.FormValue("country"))
	form.PhoneNumber = strings.TrimSpace(r.FormValue("phone_number"))
	form.PhoneKind = strings.TrimSpace(r.FormValue("phone_kind"))
	form.PhoneDepartment = strings.TrimSpace(r.FormValue("phone_department"))
	form.EmailAddress = strings.TrimSpace(r.FormValue("email_address"))
	form.SocialMedias = bindSocialMediaForm(r)

	cmd := companycmd.CreateCompanyCommand{
		Name:         form.Name,
		TradeName:    form.TradeName,
		CNPJ:         form.CNPJ,
		RegisteredBy: idUser,
		Address: companycmd.CreateCompanyAddressCommand{
			ZIP:          form.ZIP,
			Title:        form.Title,
			Street:       form.Street,
			Number:       number,
			Complement:   form.Complement,
			Reference:    form.Reference,
			Neighborhood: form.Neighborhood,
			City:         form.City,
			State:        form.State,
			Country:      form.Country,
		},
		Phone: companycmd.CreateCompanyPhoneCommand{
			Number:     form.PhoneNumber,
			Kind:       form.PhoneKind,
			Department: form.PhoneDepartment,
		},
		Email: companycmd.CreateCompanyEmailCommand{
			Address: form.EmailAddress,
		},
	}
	for _, item := range form.SocialMedias {
		if item.Platform == "" && item.URL == "" {
			continue
		}
		cmd.SocialMedia = append(cmd.SocialMedia, companycmd.CreateCompanySocialMediaCommand{
			Platform: item.Platform,
			URL:      item.URL,
		})
	}
	return form, cmd, nil
}

func BindUpdateForm(r *http.Request, current *companyread.CompanyReadModel) (FormData, companycmd.UpdateCompanyCommand, error) {
	if err := r.ParseForm(); err != nil {
		return FormData{}, companycmd.UpdateCompanyCommand{}, err
	}
	number, err := parseUintForm(r.FormValue("number"))
	if err != nil {
		return FormData{}, companycmd.UpdateCompanyCommand{}, err
	}
	form := FormDataFromReadModel(current)
	form.Name = strings.TrimSpace(r.FormValue("name"))
	form.TradeName = strings.TrimSpace(r.FormValue("trade_name"))
	form.CNPJ = strings.TrimSpace(r.FormValue("cnpj"))
	form.ZIP = strings.TrimSpace(r.FormValue("zip"))
	form.Title = strings.TrimSpace(r.FormValue("title"))
	form.Street = strings.TrimSpace(r.FormValue("street"))
	form.Number = strings.TrimSpace(r.FormValue("number"))
	form.Complement = strings.TrimSpace(r.FormValue("complement"))
	form.Reference = strings.TrimSpace(r.FormValue("reference"))
	form.Neighborhood = strings.TrimSpace(r.FormValue("neighborhood"))
	form.City = strings.TrimSpace(r.FormValue("city"))
	form.State = strings.TrimSpace(r.FormValue("state"))
	form.Country = strings.TrimSpace(r.FormValue("country"))
	form.PhoneNumber = strings.TrimSpace(r.FormValue("phone_number"))
	form.PhoneKind = strings.TrimSpace(r.FormValue("phone_kind"))
	form.PhoneDepartment = strings.TrimSpace(r.FormValue("phone_department"))
	form.EmailAddress = strings.TrimSpace(r.FormValue("email_address"))
	form.SocialMedias = bindSocialMediaForm(r)

	cmd := companycmd.UpdateCompanyCommand{
		IDCompany: current.ID,
		Name:      form.Name,
		TradeName: form.TradeName,
		CNPJ:      form.CNPJ,
		Address: companycmd.UpdateCompanyAddressCommand{
			IDAddress:    form.AddressID,
			ZIP:          form.ZIP,
			Title:        form.Title,
			Street:       form.Street,
			Number:       number,
			Complement:   form.Complement,
			Reference:    form.Reference,
			Neighborhood: form.Neighborhood,
			City:         form.City,
			State:        form.State,
			Country:      form.Country,
		},
		Phone: companycmd.UpdateCompanyPhoneCommand{
			IDPhone:    form.PhoneID,
			Number:     form.PhoneNumber,
			Kind:       form.PhoneKind,
			Department: form.PhoneDepartment,
		},
		Email: companycmd.UpdateCompanyEmailCommand{
			IDEmail: form.EmailID,
			Address: form.EmailAddress,
		},
	}
	for _, item := range form.SocialMedias {
		if item.Platform == "" && item.URL == "" {
			continue
		}
		cmd.SocialMedia = append(cmd.SocialMedia, companycmd.UpdateCompanySocialMediaCommand{
			IDSocialMedia: item.ID,
			Platform:      item.Platform,
			URL:           item.URL,
		})
	}
	return form, cmd, nil
}

func NewFormData() FormData {
	return FormData{
		Mode:        "create",
		Action:      "/companies",
		SubmitLabel: "Criar company",
		PhoneKind:   "commercial",
		Country:     "BR",
		SocialMedias: []SocialMediaFormItem{{
			Platform: "linkedin",
		}},
	}
}

func FormDataFromReadModel(company *companyread.CompanyReadModel) FormData {
	data := FormData{
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
		data.SocialMedias = make([]SocialMediaFormItem, 0, len(company.SocialMedia))
		for _, item := range company.SocialMedia {
			data.SocialMedias = append(data.SocialMedias, SocialMediaFormItem{
				ID:       item.ID,
				Platform: item.Platform,
				URL:      item.URL,
			})
		}
	}
	if len(data.SocialMedias) == 0 {
		data.SocialMedias = []SocialMediaFormItem{{}}
	}
	return data
}

func parsePage(raw string) int {
	page, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || page < 1 {
		return 1
	}
	return page
}

func parseLimit(raw string) int {
	limit, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || limit <= 0 || limit > 100 {
		return 10
	}
	return limit
}

func optionalString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func parseUintForm(raw string) (uint, error) {
	value, err := strconv.ParseUint(strings.TrimSpace(raw), 10, 64)
	if err != nil || value == 0 {
		return 0, errors.New("invalid number")
	}
	return uint(value), nil
}

func bindSocialMediaForm(r *http.Request) []SocialMediaFormItem {
	ids := r.Form["social_media_id"]
	platforms := r.Form["social_platform"]
	urls := r.Form["social_url"]
	size := len(ids)
	if len(platforms) > size {
		size = len(platforms)
	}
	if len(urls) > size {
		size = len(urls)
	}
	if size == 0 {
		return []SocialMediaFormItem{{}}
	}
	items := make([]SocialMediaFormItem, 0, size)
	for i := 0; i < size; i++ {
		var id int64
		if i < len(ids) {
			id, _ = strconv.ParseInt(strings.TrimSpace(ids[i]), 10, 64)
		}
		platform := ""
		if i < len(platforms) {
			platform = strings.TrimSpace(platforms[i])
		}
		link := ""
		if i < len(urls) {
			link = strings.TrimSpace(urls[i])
		}
		items = append(items, SocialMediaFormItem{ID: id, Platform: platform, URL: link})
	}
	return items
}
