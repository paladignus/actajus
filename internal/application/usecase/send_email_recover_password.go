// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/gateway"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type SendEmailRecoverPassword struct {
	logger     repository.Logger
	smtp       gateway.SMTP
	subscriber gateway.Subscriber
}

func NewSendEmailRecoverPassword(
	repository repository.Logger,
	gateway gateway.SMTP,
	subscriber gateway.Subscriber,
) SendEmailRecoverPassword {
	return SendEmailRecoverPassword{
		repository,
		gateway,
		subscriber,
	}
}

func (s SendEmailRecoverPassword) Execute(ctx context.Context, req dto.SendEmailRecoverPasswordInput) {
	s.logger.Info(ctx, "Send email recover password")
}
