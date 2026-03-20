// Package postgres
package postgres

import (
	companyReadModel "github.com/paladignus/actajus/internal/module/company/application/readmodel"
	emailReadModel "github.com/paladignus/actajus/internal/module/email/application/readmodel"
	phoneReadModel "github.com/paladignus/actajus/internal/module/phone/application/readmodel"
	socialMediaReadModel "github.com/paladignus/actajus/internal/module/social_media/application/readmodel"
)

func mergeCompanyScanWithRelations(
	cp *companyScan,
	a *addressScan,
	e *emailScan,
	p *phoneScan,
	sm *socialMediaScan,
	companiesMap map[int64]*companyReadModel.CompanyReadModel,
	order *[]int64,
) *companyReadModel.CompanyReadModel {
	company := cp.companyToDTO()
	if company == nil {
		return nil
	}
	existing, seen := companiesMap[company.ID]
	if !seen {
		companiesMap[company.ID] = company
		*order = append(*order, company.ID)
		existing = company
	}
	if existing.Addresses == nil {
		if addr := a.addressToDTO(); addr != nil {
			existing.Addresses = addr
		}
	}
	if e.id != nil {
		if email := e.emailToDTO(); email != nil && !hasEmail(existing.Emails, *e.id) {
			existing.Emails = append(existing.Emails, email)
		}
	}
	if p.id != nil {
		if phone := p.phoneToDTO(); phone != nil && !hasPhone(existing.Phones, *p.id) {
			existing.Phones = append(existing.Phones, phone)
		}
	}
	if sm.id != nil {
		if socialMedia := sm.socialMediaToDTO(); socialMedia != nil && !hasSocialMedia(existing.SocialMedia, *sm.id) {
			existing.SocialMedia = append(existing.SocialMedia, socialMedia)
		}
	}
	return existing
}

func hasEmail(emails []*emailReadModel.EmailReadModel, id int64) bool {
	for _, e := range emails {
		if e.ID == id {
			return true
		}
	}
	return false
}

func hasPhone(phones []*phoneReadModel.PhoneReadModel, id int64) bool {
	for _, p := range phones {
		if p.ID == id {
			return true
		}
	}
	return false
}

func hasSocialMedia(sms []*socialMediaReadModel.SocialMediaReadModel, id int64) bool {
	for _, s := range sms {
		if s.ID == id {
			return true
		}
	}
	return false
}

func reverseSlice(items []companyReadModel.CompanyReadModel) {
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}
}
