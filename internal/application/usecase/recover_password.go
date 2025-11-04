// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/internal/domain/gateway"
	"github.com/paladignus/actajus/internal/domain/repository"
	vo "github.com/paladignus/actajus/internal/domain/value_object"
)

type RecoverPassword struct {
	persistence repository.Authentication
	logger      repository.Logger
	token       gateway.Token
	smtp        gateway.SMTP
}

func NewRecoverPassword(
	persistence repository.Authentication,
	logger repository.Logger,
	token gateway.Token,
	smtp gateway.SMTP,
) RecoverPassword {
	return RecoverPassword{persistence, logger, token, smtp}
}

func (r RecoverPassword) Execute(ctx context.Context, req dto.RecoverPasswordInput) (dto.RecoverPasswordOutput, error) {
	r.logger.Info(ctx, "recover password", "email", req.Email)
	email := vo.Email(req.Email)
	if !email.IsValid() {
		r.logger.Warn(ctx, "invalid email format provided",
			"email", req.Email,
		)
		return dto.RecoverPasswordOutput{}, exception.ErrInvalidEmail
	}

	IDUser, err := r.persistence.AccountIsActive(ctx, req.Email)
	if err != nil {
		r.logger.Error(ctx, "account not found or is inactive", "error", err, "email", req.Email)
		return dto.RecoverPasswordOutput{}, err
	}

	r.logger.Info(ctx, "account is ative to email", "email", req.Email)
	token, err := r.token.GenerateResetToken(IDUser)
	if err != nil {
		r.logger.Error(ctx, "failed to generate reset token", "error", err, "email", req.Email)
		return dto.RecoverPasswordOutput{}, err
	}

	r.logger.Info(ctx, "generated reset token", "email", req.Email)
	if err = r.persistence.CreateRecoverPassword(ctx, IDUser, token.ResetToken); err != nil {
		r.logger.Error(ctx, "failed to create reset token record", "error", err, "email", req.Email)
		return dto.RecoverPasswordOutput{}, err
	}

	r.logger.Info(ctx, "created reset token record", "email", req.Email)
	if err = r.smtp.SendEmail(ctx, email.Value(), token.ResetToken); err != nil {
		// r.logger.Error(ctx, "failed to send email", "error", err, "email", req.Email)
		r.logger.Error(ctx, "failed to send email", "error", err, "email", "marcelo@marcelo.eti.br")
		return dto.RecoverPasswordOutput{}, err
	}

	r.logger.Info(ctx, "recover password successful", "email", req.Email)
	return dto.RecoverPasswordOutput{
		RecoverToken: "token",
	}, nil
}
