// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/entity"
	"github.com/paladignus/actajus/internal/domain/event"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/internal/domain/gateway"
	"github.com/paladignus/actajus/internal/domain/repository"
	"github.com/paladignus/actajus/internal/domain/service"
	vo "github.com/paladignus/actajus/internal/domain/value_object"
)

type RequestPasswordReset struct {
	auth         repository.Authentication
	token        repository.IPasswordResetToken
	tokenService service.IToken
	logger       repository.Logger
	// token     gateway.Token
	publisher gateway.Publisher
}

func NewRequestPasswordReset(
	auth repository.Authentication,
	token repository.IPasswordResetToken,
	tokenService service.IToken,
	logger repository.Logger,
	// token gateway.Token,
	publisher gateway.Publisher,
) RequestPasswordReset {
	return RequestPasswordReset{
		auth,
		token,
		tokenService,
		logger,
		// token,
		publisher,
	}
}

func (r RequestPasswordReset) Execute(ctx context.Context, req dto.RecoverPasswordInput) error {
	email := vo.Email(req.Email)
	if !email.IsValid() {
		r.logger.Warn(ctx, "invalid email format provided",
			"email", email.Value(),
		)
		return exception.ErrInvalidEmail
	}
	authenticate, err := r.auth.FindByEmail(ctx, email.Value())
	if err != nil {
		r.logger.Warn(ctx, "person does not have an account", "error", err, "email", email.Value())
		return err
	}
	if authenticate.IsDeleted() {
		r.logger.Warn(ctx, "person with inactive authentication", "email", email.Value())
	}
	if err := r.token.InvalidateUserTokens(ctx, authenticate.IDPerson); err != nil {
		r.logger.Error(ctx, "failed to invalidate previous token", "error", err, "id_user", authenticate.IDPerson)
		return err
	}
	token, err := r.tokenService.GenerateToken()
	if err != nil {
		r.logger.Error(ctx, "failed to generate reset token", "error", err, "email", email.Value())
		return err
	}
	resetToken := entity.NewPasswordResetToken(authenticate.IDPerson, token)
	if err := r.token.Create(ctx, resetToken); err != nil {
		r.logger.Error(ctx, "failed to create reset token record", "error", err, "email", email.Value())
		return err
	}
	if err := r.publisher.Publish(ctx, event.NewPasswordResetRequestedEvent(
		string(authenticate.IDPerson),
		req.Email,
		"https://api.actajus.com.br/auth/recover/"+resetToken.Token,
	)); err != nil {
		r.logger.Error(ctx, "failed to publish recover password event", "error", err, "email", req.Email)
		return err
	}
	r.logger.Info(ctx, "recover password successfully", "email", req.Email)
	return nil
	// idUser, err := r.auth.AccountIsActive(ctx, req.Email)
	// if err != nil {
	// 	r.logger.Error(ctx, "account not found or is inactive", "error", err, "email", req.Email)
	// 	return err
	// }
	// token, err := r.token.GenerateResetToken(idUser)
	// if err != nil {
	// 	r.logger.Error(ctx, "failed to generate reset token", "error", err, "email", req.Email)
	// 	return err
	// }
	// if err := r.auth.InvalidAllTokensByIDUser(ctx, idUser); err != nil {
	// 	r.logger.Error(ctx, "failed to invalidate previous token", "error", err, "id_user", idUser)
	// 	return err
	// }
	// if err = r.auth.CreateRecoverPassword(ctx, idUser, token.ResetToken); err != nil {
	// 	r.logger.Error(ctx, "failed to create reset token record", "error", err, "email", req.Email)
	// 	return err
	// }
	// if err := r.publisher.Publish(ctx, event.NewPasswordResetRequestedEvent(
	// 	idUser,
	// 	req.Email,
	// 	"https://api.actajus.com.br/auth/recover/"+token.ResetToken,
	// )); err != nil {
	// 	r.logger.Error(ctx, "failed to publish recover password event", "error", err, "email", req.Email)
	// 	return err
	// }
	// r.logger.Info(ctx, "recover password successfully", "email", req.Email)
	return nil
}
