// Package usecase
package usecase

import (
	"context"
	"fmt"
	"strconv"

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
	user         repository.IUser
	token        repository.IPasswordResetToken
	tokenService service.ITokenGenerator
	publisher    gateway.Publisher
}

func NewRequestPasswordReset(
	user repository.IUser,
	token repository.IPasswordResetToken,
	tokenService service.ITokenGenerator,
	publisher gateway.Publisher,
) RequestPasswordReset {
	return RequestPasswordReset{
		user,
		token,
		tokenService,
		publisher,
	}
}

func (r RequestPasswordReset) Execute(ctx context.Context, req dto.RequestPasswordResetInput) error {
	email := vo.Email(req.Email)
	if !email.IsValid() {
		return fmt.Errorf("invalid email format %s in request password reset use case: %w", req.Email, exception.ErrEmailNotFound)
	}
	idUser, err := r.user.FindIDUserByEmail(ctx, email.Value())
	if err != nil {
		return fmt.Errorf("request password reset use case failed to find user by email %s: %w", req.Email, err)
	}
	if err := r.token.InvalidateUserTokens(ctx, idUser); err != nil {
		return fmt.Errorf("request password reset use case failed to invalidate previous tokens for user ID %d: %w", idUser, err)
	}
	token, err := r.tokenService.Generate()
	if err != nil {
		return fmt.Errorf("request password reset use case failed to generate reset token for email %s: %w", req.Email, err)
	}
	resetToken := entity.NewPasswordResetToken(idUser, token)
	if err := r.token.Create(ctx, resetToken); err != nil {
		return fmt.Errorf("request password reset use case failed to save reset token for user ID %d: %w", idUser, err)
	}
	if err := r.publisher.Publish(ctx, event.NewPasswordResetRequestedEvent(
		strconv.Itoa(idUser),
		req.Email,
		"https://dynamicsolutions-hlg.com.br/alterar-senha/"+resetToken.Token,
	)); err != nil {
		return fmt.Errorf("request password reset use case failed to publish event for email %s: %w", req.Email, err)
	}
	return nil
}
