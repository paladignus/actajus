// Package adapter
package adapter

import (
	"github.com/paladignus/actajus/internal/module/company/application/command"
	"github.com/paladignus/actajus/internal/module/company/application/readmodel"
	addressv1 "github.com/paladignus/actajus/proto/address/v1"
	companyv1 "github.com/paladignus/actajus/proto/company/v1"
	emailv1 "github.com/paladignus/actajus/proto/email/v1"
	phonev1 "github.com/paladignus/actajus/proto/phone/v1"
	socialmediav1 "github.com/paladignus/actajus/proto/social_media/v1"
)

func ProtoToCompanyCreateCommand(in *companyv1.CreateCompanyRequest) command.CreateCompanyCommand {
	return command.CreateCompanyCommand{
		Name:        in.Name,
		TradeName:   in.TradeName,
		CNPJ:        in.Cnpj,
		Address:     protoToCompanyAddressCreateCommand(in.Address),
		Phone:       protoToCompanyPhoneCreateCommand(in.Phone),
		Email:       protoToCompanyEmailCreateCommand(in.Email),
		SocialMedia: protoToCompanySocialMediaCreateCommands(in.SocialMedia),
	}
}

func ProtoToCompanyUpdateCommand(in *companyv1.UpdateCompanyRequest) command.UpdateCompanyCommand {
	return command.UpdateCompanyCommand{
		IDCompany:   in.Id,
		Name:        in.Name,
		TradeName:   in.TradeName,
		CNPJ:        in.Cnpj,
		Address:     protoToCompanyAddressUpdateCommand(in.Address),
		Phone:       protoToCompanyPhoneUpdateCommand(in.Phone),
		Email:       protoToCompanyEmailUpdateCommand(in.Email),
		SocialMedia: protoToCompanySocialMediaUpdateCommands(in.SocialMedia),
	}
}

func CompaniesReadModelToProto(in *readmodel.CompanyListReadModel) *companyv1.ListCompaniesResponse {
	companies := &companyv1.ListCompaniesResponse{
		PageInfo: &companyv1.PageInfo{
			HasNextPage:     in.PageInfo.HasNextPage,
			HasPreviousPage: in.PageInfo.HasPreviousPage,
			PreviousUrl:     in.PageInfo.PreviousURL,
			NextUrl:         in.PageInfo.NextURL,
		},
	}
	for _, crm := range in.Data {
		companies.Companies = append(companies.Companies, &companyv1.Company{
			Id:          int64(crm.ID),
			Name:        crm.Name,
			TradeName:   crm.TradeName,
			Cnpj:        crm.CNPJ,
			Address:     companyAddressReadModelToProto(crm.Address),
			Phone:       companyPhoneReadModelToProto(crm.Phones),
			Email:       companyEmailReadModelToProto(crm.Emails),
			SocialMedia: companySocialMediaReadModelsToProto(crm.SocialMedia),
			CreatedAt:   crm.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   crm.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return companies
}

func FindByIDCompanyReadModelToProto(in *readmodel.CompanyReadModel) *companyv1.FindCompanyByIDResponse {
	return &companyv1.FindCompanyByIDResponse{
		Id:          int64(in.ID),
		Name:        in.Name,
		TradeName:   in.TradeName,
		Cnpj:        in.CNPJ,
		Address:     companyAddressReadModelToProto(in.Address),
		Phone:       companyPhoneReadModelToProto(in.Phones),
		Email:       companyEmailReadModelToProto(in.Emails),
		SocialMedia: companySocialMediaReadModelsToProto(in.SocialMedia),
		CreatedAt:   in.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   in.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func FindByCNPJCompanyReadModelToProto(in *readmodel.CompanyReadModel) *companyv1.FindCompanyByCNPJResponse {
	return &companyv1.FindCompanyByCNPJResponse{
		Id:          int64(in.ID),
		Name:        in.Name,
		TradeName:   in.TradeName,
		Cnpj:        in.CNPJ,
		Address:     companyAddressReadModelToProto(in.Address),
		Phone:       companyPhoneReadModelToProto(in.Phones),
		Email:       companyEmailReadModelToProto(in.Emails),
		SocialMedia: companySocialMediaReadModelsToProto(in.SocialMedia),
		CreatedAt:   in.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   in.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func protoToCompanyAddressCreateCommand(in *addressv1.CreateAddressRequest) command.CreateCompanyAddressCommand {
	if in == nil {
		return command.CreateCompanyAddressCommand{}
	}
	return command.CreateCompanyAddressCommand{
		ZIP:          in.Zip,
		Title:        in.Title,
		Street:       in.Street,
		Number:       uint(in.Number),
		Complement:   in.Complement,
		Reference:    in.Reference,
		Neighborhood: in.Neighborhood,
		City:         in.City,
		State:        in.State,
		Country:      in.Country,
	}
}

func protoToCompanyAddressUpdateCommand(in *addressv1.UpdateAddressRequest) command.UpdateCompanyAddressCommand {
	if in == nil {
		return command.UpdateCompanyAddressCommand{}
	}
	return command.UpdateCompanyAddressCommand{
		IDAddress:    in.Id,
		ZIP:          in.Zip,
		Title:        in.Title,
		Street:       in.Street,
		Number:       uint(in.Number),
		Complement:   in.Complement,
		Reference:    in.Reference,
		Neighborhood: in.Neighborhood,
		City:         in.City,
		State:        in.State,
		Country:      in.Country,
	}
}

func protoToCompanyPhoneCreateCommand(in *phonev1.CreatePhoneRequest) command.CreateCompanyPhoneCommand {
	if in == nil {
		return command.CreateCompanyPhoneCommand{}
	}
	return command.CreateCompanyPhoneCommand{
		Number:     in.Number,
		Kind:       in.Kind,
		Department: in.Department,
	}
}

func protoToCompanyPhoneUpdateCommand(in *phonev1.UpdatePhoneRequest) command.UpdateCompanyPhoneCommand {
	if in == nil {
		return command.UpdateCompanyPhoneCommand{}
	}
	return command.UpdateCompanyPhoneCommand{
		IDPhone:    in.Id,
		Number:     in.Number,
		Kind:       in.Kind,
		Department: in.Department,
	}
}

func protoToCompanyEmailCreateCommand(in *emailv1.CreateEmailRequest) command.CreateCompanyEmailCommand {
	if in == nil {
		return command.CreateCompanyEmailCommand{}
	}
	return command.CreateCompanyEmailCommand{Address: in.Address}
}

func protoToCompanyEmailUpdateCommand(in *emailv1.UpdateEmailRequest) command.UpdateCompanyEmailCommand {
	if in == nil {
		return command.UpdateCompanyEmailCommand{}
	}
	return command.UpdateCompanyEmailCommand{
		IDEmail: in.Id,
		Address: in.Address,
	}
}

func protoToCompanySocialMediaCreateCommands(in []*socialmediav1.CreateSocialMediaRequest) []command.CreateCompanySocialMediaCommand {
	if len(in) == 0 {
		return nil
	}
	items := make([]command.CreateCompanySocialMediaCommand, len(in))
	for i, item := range in {
		items[i] = command.CreateCompanySocialMediaCommand{
			Platform: item.Platform,
			URL:      item.Url,
		}
	}
	return items
}

func protoToCompanySocialMediaUpdateCommands(in []*socialmediav1.UpdateSocialMediaRequest) []command.UpdateCompanySocialMediaCommand {
	if len(in) == 0 {
		return nil
	}
	items := make([]command.UpdateCompanySocialMediaCommand, len(in))
	for i, item := range in {
		items[i] = command.UpdateCompanySocialMediaCommand{
			IDSocialMedia: item.Id,
			Platform:      item.Platform,
			URL:           item.Url,
		}
	}
	return items
}

func companyAddressReadModelToProto(in *readmodel.CompanyAddressReadModel) *addressv1.AddressResponse {
	if in == nil {
		return nil
	}
	out := &addressv1.AddressResponse{
		Id:           in.ID,
		Zip:          in.ZIP,
		Title:        in.Title,
		Street:       in.Street,
		Number:       uint32(in.Number),
		Neighborhood: in.Neighborhood,
		City:         in.City,
		State:        in.State,
		Country:      in.Country,
	}
	if in.Complement != nil {
		out.Complement = *in.Complement
	}
	if in.Reference != nil {
		out.Reference = *in.Reference
	}
	return out
}

func companyPhoneReadModelToProto(in []*readmodel.CompanyPhoneReadModel) []*phonev1.PhoneResponse {
	if len(in) == 0 {
		return nil
	}
	items := make([]*phonev1.PhoneResponse, 0, len(in))
	for _, item := range in {
		items = append(items, &phonev1.PhoneResponse{
			Id:         item.ID,
			Number:     item.Number,
			Kind:       item.Kind,
			Department: item.Department,
		})
	}
	return items
}

func companyEmailReadModelToProto(in []*readmodel.CompanyEmailReadModel) []*emailv1.EmailResponse {
	if len(in) == 0 {
		return nil
	}
	items := make([]*emailv1.EmailResponse, 0, len(in))
	for _, item := range in {
		items = append(items, &emailv1.EmailResponse{
			Id:      item.ID,
			Address: item.Address,
		})
	}
	return items
}

func companySocialMediaReadModelsToProto(in []*readmodel.CompanySocialMediaReadModel) []*socialmediav1.SocialMediaResponse {
	if len(in) == 0 {
		return nil
	}
	items := make([]*socialmediav1.SocialMediaResponse, 0, len(in))
	for _, item := range in {
		items = append(items, &socialmediav1.SocialMediaResponse{
			Id:       item.ID,
			Platform: item.Platform,
			Url:      item.URL,
		})
	}
	return items
}
