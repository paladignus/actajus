// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/identity/application/dto"
	"github.com/paladignus/actajus/internal/module/identity/application/event"
	"github.com/paladignus/actajus/internal/module/identity/application/mapper"
	"github.com/paladignus/actajus/internal/module/identity/application/repository"
	"github.com/paladignus/actajus/internal/module/identity/application/service"
	"github.com/paladignus/actajus/internal/shared/application/messaging"
	sharedsvc "github.com/paladignus/actajus/internal/shared/application/service"
	"github.com/paladignus/actajus/internal/shared/application/uow"
	"github.com/paladignus/actajus/internal/shared/infrastructure/config"
)

type RequestPasswordReset struct {
	uow            uow.UnitOfWork
	repository     repository.Factory
	outbox         messaging.OutboxFactory
	refresh        service.RefreshTokenService
	clock          service.Clock
	cfg            config.PasswordResetConfig
	mapper         mapper.AuthMapper
	serializer     sharedsvc.MessageSerializer
	idGenerator    sharedsvc.IDGenerator
	revokePrevious bool
}

func NewRequestPasswordReset(
	uow uow.UnitOfWork,
	repository repository.Factory,
	outbox messaging.OutboxFactory,
	refresh service.RefreshTokenService,
	clock service.Clock,
	cfg config.PasswordResetConfig,
	mapper mapper.AuthMapper,
	serializer sharedsvc.MessageSerializer,
	idGenerator sharedsvc.IDGenerator,
	revokePrevious bool,
) RequestPasswordReset {
	return RequestPasswordReset{
		uow, repository, outbox,
		refresh, clock, cfg, mapper,
		serializer, idGenerator, revokePrevious,
	}
}

func (uc RequestPasswordReset) Execute(ctx context.Context, input dto.RequestPasswordResetCommand) (*dto.RequestPasswordResetReadModel, error) {
	norm, err := uc.mapper.RequestPasswordResetInputToNormalized(input)
	if err != nil {
		return nil, fmt.Errorf("invalid password reset request data: %w", err)
	}
	user, err := uc.repository.User().FindByEmail(ctx, norm.Email)
	if err != nil {
		return nil, err
	}
	resp := &dto.RequestPasswordResetReadModel{
		Message: "Se existir uma conta com esse e-mail, enviaremos instruções para redefinição de senha.",
	}
	if user == nil {
		return resp, nil
	}
	now := uc.clock.Now()
	expiresAt := now.Add(uc.cfg.ResetTTL)
	rawToken, hash, err := uc.refresh.Generate()
	if err != nil {
		return nil, err
	}
	err = uc.uow.Do(ctx, func(tx uow.Tx) error {
		r := uc.repository.WithTx(tx).PasswordReset()
		o := uc.outbox.WithTx(tx).Outbox()
		if uc.revokePrevious {
			if err := r.RevokeAllByUser(ctx, user.ID().Value(), now); err != nil {
				return err
			}
		}
		idReset, err := r.Create(ctx, &mapper.PasswordResetTokenCreate{
			IDUser:    user.ID().Value(),
			Hash:      hash,
			ExpiresAt: expiresAt,
			CreatedAt: now,
		})
		if err != nil {
			return err
		}
		idMessage := uc.idGenerator.NewString()
		payload, err := uc.serializer.Marshal(event.EmailSendRequested{
			IDMessage: idMessage,
			Template:  "password-reset",
			To:        user.PrimaryEmail().Value(),
			Data: map[string]any{
				"id_reset":   idReset,
				"token":      rawToken,
				"expires_at": expiresAt.Sub(now).Minutes(),
				"id_user":    user.ID().Value(),
			},
			Meta: map[string]string{
				"module":  "identity",
				"usecase": "request_password_reset",
			},
		})
		if err != nil {
			return err
		}
		return o.Add(ctx, messaging.OutboxMessage{
			ID:         idMessage,
			Subject:    event.SubjectEmailSendRequested,
			Payload:    payload,
			OccurredAt: now,
		})
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
