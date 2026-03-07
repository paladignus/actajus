// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/identity/application/dto"
	"github.com/paladignus/actajus/internal/module/identity/application/event"
	"github.com/paladignus/actajus/internal/module/identity/application/mapper"
	"github.com/paladignus/actajus/internal/module/identity/application/model"
	"github.com/paladignus/actajus/internal/module/identity/application/repository"
	"github.com/paladignus/actajus/internal/module/identity/application/service"
	sharedsvc "github.com/paladignus/actajus/internal/shared/application/service"
	"github.com/paladignus/actajus/internal/shared/application/uow"
	"github.com/paladignus/actajus/internal/shared/infrastructure/config"
	"github.com/paladignus/actajus/shared/infrastructure/messaging"
)

// type RequestPasswordReset struct {
// 	user           repository.UserRepository
// 	reset          repository.PasswordResetRepository
// 	refresh        service.RefreshTokenService // reutiliza gerador/hash (mesmo do refresh!)
// 	clock          service.Clock
// 	cfg            config.PasswordResetConfig
// 	mapper         mapper.AuthMapper
// 	RevokePrevious bool
// }
//
// func NewRequestPasswordReset(
// 	user repository.UserRepository,
// 	reset repository.PasswordResetRepository,
// 	refresh service.RefreshTokenService,
// 	clock service.Clock,
// 	cfg config.PasswordResetConfig,
// 	mapper mapper.AuthMapper,
// 	revokePrevious bool,
// ) RequestPasswordReset {
// 	return RequestPasswordReset{
// 		user, reset, refresh, clock,
// 		cfg, mapper, revokePrevious,
// 	}
// }

type RequestPasswordReset struct {
	uow            uow.UnitOfWork
	user           repository.UserRepository
	reset          repository.PasswordResetRepository
	outbox         repository.OutboxRepository
	refresh        service.RefreshTokenService
	clock          service.Clock
	cfg            config.PasswordResetConfig
	mapper         mapper.AuthMapper
	serializer     sharedsvc.MessageSerializer
	idGenerator    sharedsvc.IDGenerator
	RevokePrevious bool
}

func NewRequestPasswordReset(
	uow uow.UnitOfWork,
	user repository.UserRepository,
	reset repository.PasswordResetRepository,
	outbox repository.OutboxRepository,
	refresh service.RefreshTokenService,
	clock service.Clock,
	cfg config.PasswordResetConfig,
	mapper mapper.AuthMapper,
	serializer sharedsvc.MessageSerializer,
	idGenerator sharedsvc.IDGenerator,
	RevokePrevious bool,
) RequestPasswordReset {
	return RequestPasswordReset{
		uow, user, reset, outbox,
		refresh, clock, cfg, mapper,
		serializer, idGenerator, RevokePrevious,
	}
}

func (uc RequestPasswordReset) Execute(ctx context.Context, input dto.RequestPasswordResetCommand) (*dto.RequestPasswordResetReadModel, error) {
	norm, err := uc.mapper.RequestPasswordResetInputToNormalized(input)
	if err != nil {
		return nil, fmt.Errorf("invalid password reset request data: %w", err)
	}
	user, err := uc.user.FindByEmail(ctx, norm.Email)
	if err != nil {
		return nil, err
	}
	resp := &dto.RequestPasswordResetReadModel{
		Message: "Se existir uma conta com esse e-mail, enviaremos instruções para redefinição de senha.",
	}
	if user == nil {
		return resp, nil
	}
	// if err != nil || user == nil {
	// 	return &dto.RequestPasswordResetReadModel{
	// 		IDReset:    0,
	// 		ResetToken: "",
	// 		ExpiresAt:  time.Time{},
	// 	}, nil
	// }
	now := uc.clock.Now()
	// if uc.RevokePrevious {
	// 	_ = uc.reset.RevokeAllByUser(ctx, user.ID().Value(), now)
	// }
	token, hash, err := uc.refresh.Generate()
	if err != nil {
		return nil, err
	}
	expiresAt := now.Add(uc.cfg.ResetTTL)
	err = uc.uow.Do(ctx, func(tx uow.Tx) error {
		if uc.RevokePrevious {
			if err := uc.reset.RevokeAllByUser(ctx, tx, user.ID(), now); err != nil {
				return err
			}
		}
		IDReset, err := uc.reset.Create(ctx, tx, &model.PasswordResetTokenCreate{
			IDUser:    user.ID().Value(),
			Hash:      hash,
			ExpiresAt: expiresAt,
			CreatedAt: now,
		})
		if err != nil {
			return err
		}
		IDMessage := uc.idGenerator.NewString()
		payload, err := uc.serializer.Marshal(event.EmailSendRequested{
			IDMessage: IDMessage,
			Template:  "password-reset",
			To:        user.PrimaryEmail().Value(),
			Data: map[string]any{
				"reset_id":   IDReset,
				"token":      token,
				"expires_at": expiresAt,
				"user_id":    user.ID(),
			},
			Meta: map[string]string{
				"module":  "identity",
				"usecase": "request_password_reset",
			},
		})
		if err != nil {
			return err
		}
		return uc.outbox.Add(ctx, tx, messaging.OutboxMessage{
			ID:         IDMessage,
			Subject:    event.SubjectEmailSendRequested,
			Payload:    payload,
			OccurredAt: now,
		})
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
	// in := model.PasswordResetTokenCreate{
	// 	IDUser:    user.ID().Value(),
	// 	Hash:      hash,
	// 	ExpiresAt: expiresAt,
	// 	CreatedAt: now,
	// }
	// rid, err := uc.reset.Create(ctx, &in)
	// if err != nil {
	// 	return nil, err
	// }
	// return &dto.RequestPasswordResetReadModel{
	// 	IDReset:    rid,
	// 	ResetToken: token,
	// 	ExpiresAt:  expiresAt,
	// }, nil
}
