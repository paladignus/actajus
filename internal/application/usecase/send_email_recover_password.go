// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/domain/gateway"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type SendEmailRecoverPassword struct {
	logger     repository.Logger
	handler    repository.IHandler
	subscriber gateway.Subscriber
}

func NewSendEmailRecoverPassword(
	repository repository.Logger,
	handler repository.IHandler,
	subscriber gateway.Subscriber,
) SendEmailRecoverPassword {
	return SendEmailRecoverPassword{
		repository,
		handler,
		subscriber,
	}
}

func (s SendEmailRecoverPassword) Execute(ctx context.Context) error {
	s.logger.Info(ctx, "send email recover password")
	eventName := "user.password_reset_requested"
	err := s.subscriber.Subscribe(ctx, eventName, s.handler)
	if err != nil {
		s.logger.Error(ctx, "failed to subscribe to event", "event", eventName, "error", err)
		return fmt.Errorf("failed to subscribe handler for %s: %w", eventName, err)
	}
	s.logger.Info(ctx, "handler registered for event", "event", eventName)
	return nil
}
