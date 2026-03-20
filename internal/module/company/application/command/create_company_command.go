// Package command
package command

import (
	addrCommand "github.com/paladignus/actajus/internal/module/address/application/command"
	emailCommand "github.com/paladignus/actajus/internal/module/email/application/command"
	phoneCommand "github.com/paladignus/actajus/internal/module/phone/application/command"
	socialMediaCommand "github.com/paladignus/actajus/internal/module/social_media/application/command"
)

type CreateCompanyCommand struct {
	Name         string
	TradeName    string
	CNPJ         string
	RegisteredBy int64
	Address      addrCommand.CreateAddressCommand
	Phone        phoneCommand.CreatePhoneCommand
	Email        emailCommand.CreateEmailCommand
	SocialMedia  []socialMediaCommand.CreateSocialMediaCommand
}
