// Package adapter
package adapter

import (
	addr "github.com/paladignus/actajus/internal/module/address/presentation/grpc/adapter"
	"github.com/paladignus/actajus/internal/module/company/application/dto"
	email "github.com/paladignus/actajus/internal/module/email/presentation/grpc/adapter"
	phone "github.com/paladignus/actajus/internal/module/phone/presentation/grpc/adapter"
	socialMedia "github.com/paladignus/actajus/internal/module/social_media/presentation/grpc/adapter"
	companyv1 "github.com/paladignus/actajus/proto/company/v1"
)

func ProtoToCompanyCreateCommand(in *companyv1.CreateCompanyRequest) dto.CreateCompanyRequest {
	return dto.CreateCompanyRequest{
		Name:        in.Name,
		TradeName:   in.TradeName,
		CNPJ:        in.Cnpj,
		Address:     addr.ProtoToAddressCreateCommand(in.Address),
		Phone:       phone.ProtoToPhoneCreateCommand(in.Phone),
		Email:       email.ProtoToEmailCreateCommand(in.Email),
		SocialMedia: socialMedia.ProtoToSocialMediaCreateCommands(in.SocialMedia),
	}
}

func ProtoToCompanyUpdateCommand(in *companyv1.UpdateCompanyRequest) dto.UpdateCompanyRequest {
	return dto.UpdateCompanyRequest{
		IDCompany:   in.Id,
		Name:        in.Name,
		TradeName:   in.TradeName,
		CNPJ:        in.Cnpj,
		Address:     addr.ProtoToAddressUpdateCommand(in.Address),
		Phone:       phone.ProtoToPhoneUpdateCommand(in.Phone),
		Email:       email.ProtoToEmailUpdateCommand(in.Email),
		SocialMedia: socialMedia.ProtoToSocialMediaUpdateCommands(in.SocialMedia),
	}
}

func CompaniesReadModelToProto(in *dto.CompanyListReadModel) *companyv1.ListCompaniesResponse {
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
			Address:     addr.AddressReadModelToProto(crm.Addresses),
			Phone:       phone.PhoneReadModelToProto(crm.Phones),
			Email:       email.EmailReadModelToProto(crm.Emails),
			SocialMedia: socialMedia.SocialMediaReadModelsToProto(crm.SocialMedia),
			CreatedAt:   crm.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   crm.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return companies
}

func FindByIDCompanyReadModelToProto(in *dto.CompanyReadModel) *companyv1.FindCompanyByIDResponse {
	return &companyv1.FindCompanyByIDResponse{
		Id:          int64(in.ID),
		Name:        in.Name,
		TradeName:   in.TradeName,
		Cnpj:        in.CNPJ,
		Address:     addr.AddressReadModelToProto(in.Addresses),
		Phone:       phone.PhoneReadModelToProto(in.Phones),
		Email:       email.EmailReadModelToProto(in.Emails),
		SocialMedia: socialMedia.SocialMediaReadModelsToProto(in.SocialMedia),
		CreatedAt:   in.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   in.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func FindByCNPJCompanyReadModelToProto(in *dto.CompanyReadModel) *companyv1.FindCompanyByCNPJResponse {
	return &companyv1.FindCompanyByCNPJResponse{
		Id:          int64(in.ID),
		Name:        in.Name,
		TradeName:   in.TradeName,
		Cnpj:        in.CNPJ,
		Address:     addr.AddressReadModelToProto(in.Addresses),
		Phone:       phone.PhoneReadModelToProto(in.Phones),
		Email:       email.EmailReadModelToProto(in.Emails),
		SocialMedia: socialMedia.SocialMediaReadModelsToProto(in.SocialMedia),
		CreatedAt:   in.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   in.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
