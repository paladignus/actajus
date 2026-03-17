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
