package adapter_test

import (
	"testing"
	"time"

	"github.com/paladignus/actajus/internal/module/company/application/readmodel"
	"github.com/paladignus/actajus/internal/module/company/presentation/grpc/adapter"
	sharedpagination "github.com/paladignus/actajus/internal/shared/application/pagination"
	addressv1 "github.com/paladignus/actajus/proto/address/v1"
	companyv1 "github.com/paladignus/actajus/proto/company/v1"
	emailv1 "github.com/paladignus/actajus/proto/email/v1"
	phonev1 "github.com/paladignus/actajus/proto/phone/v1"
	socialmediav1 "github.com/paladignus/actajus/proto/social_media/v1"
)

func TestProtoToCompanyCreateCommand(t *testing.T) {
	t.Parallel()

	req := &companyv1.CreateCompanyRequest{
		Name:      "Acme SA",
		TradeName: "Acme",
		Cnpj:      "12345678000199",
		Address: &addressv1.CreateAddressRequest{
			Zip:          "79000000",
			Title:        "Matriz",
			Street:       "Rua A",
			Number:       42,
			Complement:   "Sala 1",
			Reference:    "Esquina",
			Neighborhood: "Centro",
			City:         "Campo Grande",
			State:        "MS",
			Country:      "BR",
		},
		Phone: &phonev1.CreatePhoneRequest{
			Number:     "6733334444",
			Kind:       "commercial",
			Department: "sales",
		},
		Email: &emailv1.CreateEmailRequest{
			Address: "contato@acme.com",
		},
		SocialMedia: []*socialmediav1.CreateSocialMediaRequest{
			{Platform: "linkedin", Url: "https://linkedin.com/company/acme"},
		},
	}

	got := adapter.ProtoToCompanyCreateCommand(req)

	if got.Name != req.Name || got.TradeName != req.TradeName || got.CNPJ != req.Cnpj {
		t.Fatalf("unexpected root mapping: %+v", got)
	}
	if got.Address.Number != uint(req.Address.Number) || got.Address.ZIP != req.Address.Zip {
		t.Fatalf("unexpected address mapping: %+v", got.Address)
	}
	if got.Phone.Number != req.Phone.Number || got.Phone.Kind != req.Phone.Kind {
		t.Fatalf("unexpected phone mapping: %+v", got.Phone)
	}
	if got.Email.Address != req.Email.Address {
		t.Fatalf("unexpected email mapping: %+v", got.Email)
	}
	if len(got.SocialMedia) != 1 || got.SocialMedia[0].Platform != "linkedin" || got.SocialMedia[0].URL != "https://linkedin.com/company/acme" {
		t.Fatalf("unexpected social media mapping: %+v", got.SocialMedia)
	}
}

func TestProtoToCompanyUpdateCommand(t *testing.T) {
	t.Parallel()

	req := &companyv1.UpdateCompanyRequest{
		Id:        7,
		Name:      "Acme SA",
		TradeName: "Acme",
		Cnpj:      "12345678000199",
		Address: &addressv1.UpdateAddressRequest{
			Id:           11,
			Zip:          "79000000",
			Title:        "Filial",
			Street:       "Rua B",
			Number:       77,
			Complement:   "Fundos",
			Reference:    "Praca",
			Neighborhood: "Centro",
			City:         "Campo Grande",
			State:        "MS",
			Country:      "BR",
		},
		Phone: &phonev1.UpdatePhoneRequest{
			Id:         22,
			Number:     "6799999999",
			Kind:       "support",
			Department: "ops",
		},
		Email: &emailv1.UpdateEmailRequest{
			Id:      33,
			Address: "suporte@acme.com",
		},
		SocialMedia: []*socialmediav1.UpdateSocialMediaRequest{
			{Id: 44, Platform: "instagram", Url: "https://instagram.com/acme"},
		},
	}

	got := adapter.ProtoToCompanyUpdateCommand(req)

	if got.IDCompany != req.Id || got.Address.IDAddress != req.Address.Id || got.Phone.IDPhone != req.Phone.Id || got.Email.IDEmail != req.Email.Id {
		t.Fatalf("unexpected id mapping: %+v", got)
	}
	if len(got.SocialMedia) != 1 || got.SocialMedia[0].IDSocialMedia != 44 || got.SocialMedia[0].URL != "https://instagram.com/acme" {
		t.Fatalf("unexpected social media mapping: %+v", got.SocialMedia)
	}
}

func TestFindByIDCompanyReadModelToProto(t *testing.T) {
	t.Parallel()

	complement := "Sala 3"
	reference := "Ao lado do forum"
	now := time.Date(2026, 3, 27, 14, 30, 0, 0, time.UTC)
	in := &readmodel.CompanyReadModel{
		ID:        9,
		Name:      "Acme SA",
		TradeName: "Acme",
		CNPJ:      "12345678000199",
		Address: &readmodel.CompanyAddressReadModel{
			ID:           1,
			ZIP:          "79000000",
			Title:        "Matriz",
			Street:       "Rua A",
			Number:       55,
			Complement:   &complement,
			Reference:    &reference,
			Neighborhood: "Centro",
			City:         "Campo Grande",
			State:        "MS",
			Country:      "BR",
		},
		Phones: []*readmodel.CompanyPhoneReadModel{
			{ID: 2, Number: "6733334444", Kind: "commercial", Department: "sales"},
		},
		Emails: []*readmodel.CompanyEmailReadModel{
			{ID: 3, Address: "contato@acme.com"},
		},
		SocialMedia: []*readmodel.CompanySocialMediaReadModel{
			{ID: 4, Platform: "linkedin", URL: "https://linkedin.com/company/acme"},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	got := adapter.FindByIDCompanyReadModelToProto(in)

	if got.Id != 9 || got.Cnpj != in.CNPJ || got.Address == nil {
		t.Fatalf("unexpected root mapping: %+v", got)
	}
	if got.Address.Number != 55 || got.Address.Complement != complement || got.Address.Reference != reference {
		t.Fatalf("unexpected address mapping: %+v", got.Address)
	}
	if len(got.Phone) != 1 || got.Phone[0].Id != 2 {
		t.Fatalf("unexpected phone mapping: %+v", got.Phone)
	}
	if len(got.Email) != 1 || got.Email[0].Address != "contato@acme.com" {
		t.Fatalf("unexpected email mapping: %+v", got.Email)
	}
	if len(got.SocialMedia) != 1 || got.SocialMedia[0].Url != "https://linkedin.com/company/acme" {
		t.Fatalf("unexpected social media mapping: %+v", got.SocialMedia)
	}
	if got.CreatedAt != "2026-03-27 14:30:00" || got.UpdatedAt != "2026-03-27 14:30:00" {
		t.Fatalf("unexpected timestamp mapping: %+v", got)
	}
}

func TestCompaniesReadModelToProto(t *testing.T) {
	t.Parallel()

	nextURL := "http://localhost:8080/companies?after=abc"
	prevURL := "http://localhost:8080/companies?before=def"
	now := time.Date(2026, 3, 27, 15, 0, 0, 0, time.UTC)
	in := &readmodel.CompanyListReadModel{
		Data: []readmodel.CompanyReadModel{
			{
				ID:        1,
				Name:      "Acme SA",
				TradeName: "Acme",
				CNPJ:      "12345678000199",
				Address: &readmodel.CompanyAddressReadModel{
					ID:           10,
					ZIP:          "79000000",
					Title:        "Matriz",
					Street:       "Rua A",
					Number:       10,
					Neighborhood: "Centro",
					City:         "Campo Grande",
					State:        "MS",
					Country:      "BR",
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		PageInfo: sharedpagination.PageInfo{
			HasNextPage:     true,
			HasPreviousPage: true,
			NextURL:         &nextURL,
			PreviousURL:     &prevURL,
		},
	}

	got := adapter.CompaniesReadModelToProto(in)

	if got.PageInfo == nil || got.PageInfo.NextUrl == nil || got.PageInfo.PreviousUrl == nil {
		t.Fatalf("unexpected page info mapping: %+v", got.PageInfo)
	}
	if *got.PageInfo.NextUrl != nextURL || *got.PageInfo.PreviousUrl != prevURL {
		t.Fatalf("unexpected page info mapping: %+v", got.PageInfo)
	}
	if len(got.Companies) != 1 || got.Companies[0].Address == nil {
		t.Fatalf("unexpected companies mapping: %+v", got.Companies)
	}
	if got.Companies[0].Address.Number != 10 {
		t.Fatalf("unexpected nested address mapping: %+v", got.Companies[0].Address)
	}
}
