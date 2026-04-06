package companyviewmodel

import (
	companyread "github.com/paladignus/actajus/internal/module/company/application/readmodel"
	companyrepo "github.com/paladignus/actajus/internal/module/company/application/repository"
	webflash "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/flash"
)

type ListItem struct {
	ID               int64
	Name             string
	TradeName        string
	CNPJ             string
	RegisteredByName string
	City             string
	State            string
}

type CompaniesPage struct {
	Page            int
	Limit           int
	RangeStart      int
	RangeEnd        int
	Items           []ListItem
	NextURL         string
	PreviousURL     string
	HasNextPage     bool
	HasPreviousPage bool
	FilterName      string
	FilterCNPJ      string
	Flash           webflash.Message
}

type DetailPage struct {
	Company *companyread.CompanyReadModel
}

func BuildCompaniesPage(rm *companyread.CompanyListReadModel, filter companyrepo.CompanyListFilter, page, limit int, flash webflash.Message, nextURL, previousURL string) CompaniesPage {
	items := make([]ListItem, 0, len(rm.Data))
	for _, company := range rm.Data {
		city := "-"
		state := "-"
		if company.Address != nil {
			if company.Address.City != "" {
				city = company.Address.City
			}
			if company.Address.State != "" {
				state = company.Address.State
			}
		}
		items = append(items, ListItem{
			ID:               company.ID,
			Name:             company.Name,
			TradeName:        company.TradeName,
			CNPJ:             company.CNPJ,
			RegisteredByName: fallback(company.RegisteredByName, "-"),
			City:             city,
			State:            state,
		})
	}
	rangeStart := 0
	rangeEnd := 0
	if len(items) > 0 {
		rangeStart = ((page - 1) * limit) + 1
		rangeEnd = rangeStart + len(items) - 1
	}
	return CompaniesPage{
		Page:            page,
		Limit:           limit,
		RangeStart:      rangeStart,
		RangeEnd:        rangeEnd,
		Items:           items,
		HasNextPage:     rm.PageInfo.HasNextPage,
		HasPreviousPage: rm.PageInfo.HasPreviousPage,
		FilterName:      filter.Name,
		FilterCNPJ:      filter.CNPJ,
		Flash:           flash,
		NextURL:         nextURL,
		PreviousURL:     previousURL,
	}
}

func fallback(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
