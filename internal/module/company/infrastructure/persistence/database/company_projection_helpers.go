// Package database
package database

import (
	companyDTO "github.com/paladignus/actajus/internal/module/company/application/dto"
	emailDTO "github.com/paladignus/actajus/internal/module/email/application/dto"
	phoneDTO "github.com/paladignus/actajus/internal/module/phone/application/dto"
	socialMediaDTO "github.com/paladignus/actajus/internal/module/social_media/application/dto"
)

func mergeCompanyScanWithRelations(
	cp *companyScan,
	a *addressScan,
	e *emailScan,
	p *phoneScan,
	sm *socialMediaScan,
	companiesMap map[uint]*companyDTO.CompanyReadModel,
	order *[]uint,
) *companyDTO.CompanyReadModel {
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
	if e.id != nil && !hasEmail(existing.Emails, *e.id) {
		existing.Emails = append(existing.Emails, e.emailToDTO())
	}
	if p.id != nil && !hasPhone(existing.Phones, *p.id) {
		existing.Phones = append(existing.Phones, p.phoneToDTO())
	}
	if sm.id != nil && !hasSocialMedia(existing.SocialMedia, *sm.id) {
		existing.SocialMedia = append(existing.SocialMedia, sm.socialMediaToDTO())
	}
	return existing
}

func hasEmail(emails []*emailDTO.EmailReadModel, id uint) bool {
	for _, e := range emails {
		if e.ID == id {
			return true
		}
	}
	return false
}

func hasPhone(phones []*phoneDTO.PhoneReadModel, id uint) bool {
	for _, p := range phones {
		if p.ID == id {
			return true
		}
	}
	return false
}

func hasSocialMedia(sms []*socialMediaDTO.SocialMediaReadModel, id uint) bool {
	for _, s := range sms {
		if s.ID == id {
			return true
		}
	}
	return false
}

func reverseSlice(items []companyDTO.CompanyReadModel) {
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}
}
