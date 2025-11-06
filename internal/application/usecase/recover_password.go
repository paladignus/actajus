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
	// smtp        gateway.SMTP
	publisher gateway.Publisher
}

func NewRecoverPassword(
	persistence repository.Authentication,
	logger repository.Logger,
	token gateway.Token,
	// smtp gateway.SMTP,
	publisher gateway.Publisher,
) RecoverPassword {
	return RecoverPassword{persistence, logger, token, publisher}
}

func (r RecoverPassword) Execute(ctx context.Context, req dto.RecoverPasswordInput) error {
	r.logger.Info(ctx, "recover password", "email", req.Email)
	email := vo.Email(req.Email)
	if !email.IsValid() {
		r.logger.Warn(ctx, "invalid email format provided",
			"email", req.Email,
		)
		return exception.ErrInvalidEmail
	}
	IDUser, err := r.persistence.AccountIsActive(ctx, req.Email)
	if err != nil {
		r.logger.Error(ctx, "account not found or is inactive", "error", err, "email", req.Email)
		return err
	}
	r.logger.Info(ctx, "account is ative to email", "email", req.Email)
	token, err := r.token.GenerateResetToken(IDUser)
	if err != nil {
		r.logger.Error(ctx, "failed to generate reset token", "error", err, "email", req.Email)
		return err
	}
	r.logger.Info(ctx, "invalidating previous token", "id_user", IDUser)
	if err := r.persistence.InvalidAllTokensByIDUser(ctx, IDUser); err != nil {
		r.logger.Error(ctx, "failed to invalidate previous token", "error", err, "id_user", IDUser)
		return err
	}
	r.logger.Info(ctx, "generated reset token", "email", req.Email)
	if err = r.persistence.CreateRecoverPassword(ctx, IDUser, token.ResetToken); err != nil {
		r.logger.Error(ctx, "failed to create reset token record", "error", err, "email", req.Email)
		return err
	}
	// URL := "https://api.actajus.com.br/recover-password?token=" + token.ResetToken
	// event := map[string]interface{}{
	// 	"id_user": IDUser,
	// 	"email":   email.Value(),
	// 	"name":    "marcelo",
	// 	"url":     "https://api.actajus.com.br/recover-password?token=" + token.ResetToken,
	// 	"event":   "user.created",
	// }
	// payload, _ := json.Marshal(event)
	event := event.RecoveredPassword{
		Email: email.Value(),
		URL:   "https://api.actajus.com.br/recover-password?token=" + token.ResetToken,
	}
	if err = r.publisher.Publish(ctx, event); err != nil {
		r.logger.Error(ctx, "failed to publish recover password event", "error", err, "email", req.Email)
		return err
	}
	// URL := "https://api.actajus.com.br/recover-password?token=" + token.ResetToken
	// r.logger.Info(ctx, "created reset token record", "email", req.Email)
	// if err = r.smtp.SendEmail(ctx, email.Value(), URL); err != nil {
	// 	r.logger.Error(ctx, "failed to send email", "error", err, "email", req.Email)
	// 	return err
	// }
	r.logger.Info(ctx, "recover password successful", "email", req.Email)
	return nil
}
