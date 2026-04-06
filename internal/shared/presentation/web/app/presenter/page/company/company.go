package companypage

import (
	"strconv"

	companyread "github.com/paladignus/actajus/internal/module/company/application/readmodel"
	companyrepo "github.com/paladignus/actajus/internal/module/company/application/repository"
	webflash "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/flash"
	companybinder "github.com/paladignus/actajus/internal/shared/presentation/web/app/binder/http/company"
	webform "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/form"
	companynav "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/navigation/company"
	companyviewmodel "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/viewmodel/company"
	webtemplate "github.com/paladignus/actajus/internal/shared/presentation/web/template"
)

func CompaniesPage(rm *companyread.CompanyListReadModel, filter companyrepo.CompanyListFilter, page, limit int, flash webflash.Message, nextURL, previousURL string) webtemplate.Page {
	return webtemplate.Page{
		Title:       "Companies",
		Description: "Listagem WEB de companies.",
		NavKey:      "companies.list",
		Entry:       "src/main.ts",
		Data: companyviewmodel.BuildCompaniesPage(
			rm,
			filter,
			page,
			limit,
			flash,
			companynav.RewritePageLinkURL(nextURL, page+1, filter.Name, filter.CNPJ),
			companynav.RewritePageLinkURL(previousURL, previousPage(page), filter.Name, filter.CNPJ),
		),
	}
}

func DetailPage(company *companyread.CompanyReadModel) webtemplate.Page {
	return webtemplate.Page{
		Title:       "Company",
		Description: "Detalhe WEB de company.",
		NavKey:      "companies.list",
		Entry:       "src/main.ts",
		Data:        companyviewmodel.DetailPage{Company: company},
	}
}

func FormPage(data webform.CompanyData, templateName string) (string, webtemplate.Page) {
	return templateName, webtemplate.Page{
		Title:       "Company Form",
		Description: "Formulario WEB de company.",
		NavKey:      formNavKey(data.Mode),
		Entry:       "src/main.ts",
		Data:        data,
	}
}

func NewFormData() webform.CompanyData {
	return webform.NewCompanyData()
}

func FormDataFromBinder(input companybinder.FormData) webform.CompanyData {
	return webform.CompanyDataFromBinder(input)
}

func FormDataFromReadModel(company *companyread.CompanyReadModel) webform.CompanyData {
	return webform.CompanyDataFromReadModel(company)
}

func ActionPath(id int64) string {
	return "/companies/" + strconv.FormatInt(id, 10)
}

func previousPage(page int) int {
	if page <= 1 {
		return 1
	}
	return page - 1
}

func formNavKey(mode string) string {
	if mode == "create" {
		return "companies.new"
	}
	return "companies.list"
}
