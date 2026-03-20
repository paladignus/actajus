// Package command
package command

import (
	addrCommand "github.com/paladignus/actajus/internal/module/address/application/command"
	emailCommand "github.com/paladignus/actajus/internal/module/email/application/command"
	phoneCommand "github.com/paladignus/actajus/internal/module/phone/application/command"
	socialMediaCommand "github.com/paladignus/actajus/internal/module/social_media/application/command"
)

type UpdateCompanyCommand struct {
	IDCompany   int64
	Name        string
	TradeName   string
	CNPJ        string
	Address     addrCommand.UpdateAddressCommand
	Phone       phoneCommand.UpdatePhoneCommand
	Email       emailCommand.UpdateEmailCommand
	SocialMedia []socialMediaCommand.UpdateSocialMediaCommand
}
