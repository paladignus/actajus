// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/event"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/internal/domain/gateway"
	"github.com/paladignus/actajus/internal/domain/repository"
	vo "github.com/paladignus/actajus/internal/domain/value_object"
)

type RecoverPassword struct {
	persistence repository.Authentication
	logger      repository.Logger
	token       gateway.Token
	publisher   gateway.Publisher
}

func NewRecoverPassword(
	persistence repository.Authentication,
	logger repository.Logger,
	token gateway.Token,
	publisher gateway.Publisher,
) RecoverPassword {
	return RecoverPassword{persistence, logger, token, publisher}
}

func (r RecoverPassword) Execute(ctx context.Context, req dto.RecoverPasswordInput) error {
	email := vo.Email(req.Email)
	if !email.IsValid() {
		r.logger.Warn(ctx, "invalid email format provided",
			"email", req.Email,
		)
		return exception.ErrInvalidEmail
	}
	idUser, err := r.persistence.AccountIsActive(ctx, req.Email)
	if err != nil {
		r.logger.Error(ctx, "account not found or is inactive", "error", err, "email", req.Email)
		return err
	}
	token, err := r.token.GenerateResetToken(idUser)
	if err != nil {
		r.logger.Error(ctx, "failed to generate reset token", "error", err, "email", req.Email)
		return err
	}
	if err := r.persistence.InvalidAllTokensByIDUser(ctx, idUser); err != nil {
		r.logger.Error(ctx, "failed to invalidate previous token", "error", err, "id_user", idUser)
		return err
	}
	if err = r.persistence.CreateRecoverPassword(ctx, idUser, token.ResetToken); err != nil {
		r.logger.Error(ctx, "failed to create reset token record", "error", err, "email", req.Email)
		return err
	}
	if err := r.publisher.Publish(ctx, event.NewPasswordResetRequestedEvent(
		idUser,
		req.Email,
		"https://api.actajus.com.br/recover-password?token="+token.ResetToken,
	)); err != nil {
		r.logger.Error(ctx, "failed to publish recover password event", "error", err, "email", req.Email)
		return err
	}
	r.logger.Info(ctx, "recover password successfully", "email", req.Email)
	return nil
}
