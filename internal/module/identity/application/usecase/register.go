// Package usecase
package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/paladignus/actajus/internal/module/identity/application/command"
	"github.com/paladignus/actajus/internal/module/identity/application/event"
	"github.com/paladignus/actajus/internal/module/identity/application/mapper"
	"github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	"github.com/paladignus/actajus/internal/module/identity/application/repository"
	"github.com/paladignus/actajus/internal/module/identity/application/service"
	"github.com/paladignus/actajus/internal/shared/application/messaging"
	sharedsvc "github.com/paladignus/actajus/internal/shared/application/service"
	"github.com/paladignus/actajus/internal/shared/application/uow"
)

type Register struct {
	uow         uow.UnitOfWork
	repository  repository.Factory
	outbox      messaging.OutboxFactory
	mapper      mapper.AuthMapper
	hasher      service.PasswordHasher
	refresh     service.RefreshTokenService
	clock       service.Clock
	serializer  sharedsvc.MessageSerializer
	idGenerator sharedsvc.IDGenerator
	ttl         time.Duration
}

const registerNeutralMessage = "Se o cadastro puder ser concluido, voce recebera um email com os proximos passos."

func NewRegister(
	uow uow.UnitOfWork,
	repository repository.Factory,
	outbox messaging.OutboxFactory,
	mapper mapper.AuthMapper,
	hasher service.PasswordHasher,
	refresh service.RefreshTokenService,
	clock service.Clock,
	serializer sharedsvc.MessageSerializer,
	idGenerator sharedsvc.IDGenerator,
	ttl time.Duration,
) Register {
	return Register{
		uow:         uow,
		repository:  repository,
		outbox:      outbox,
		mapper:      mapper,
		hasher:      hasher,
		refresh:     refresh,
		clock:       clock,
		serializer:  serializer,
		idGenerator: idGenerator,
		ttl:         ttl,
	}
}

func (uc Register) Execute(ctx context.Context, input command.RegisterCommand) (*readmodel.RegisterReadModel, error) {
	norm, err := uc.mapper.RegisterInputToNormalized(input)
	if err != nil {
		return nil, fmt.Errorf("invalid register data: %w", err)
	}
	user, err := uc.repository.User().FindByEmail(ctx, norm.Email)
	if err != nil {
		return nil, err
	}
	if user != nil {
		if user.IsPrimaryEmailVerified() {
			return &readmodel.RegisterReadModel{Message: registerNeutralMessage}, nil
		}
		if err := uc.resendVerification(ctx, user.ID().Value(), norm.Email); err != nil {
			return nil, err
		}
		return &readmodel.RegisterReadModel{Message: registerNeutralMessage}, nil
	}
	passwordHash, err := uc.hasher.Hash(norm.Password)
	if err != nil {
		return nil, err
	}
	now := uc.clock.Now()
	if err := uc.uow.Do(ctx, func(tx uow.Tx) error {
		r := uc.repository.WithTx(tx)
		created, err := r.User().Create(ctx, repository.CreateUserRegistration{
			FirstName:    norm.FirstName,
			LastName:     norm.LastName,
			Birthday:     norm.Birthday,
			GenderID:     norm.GenderID,
			Email:        norm.Email,
			PasswordHash: passwordHash,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
		if err != nil {
			return err
		}
		return uc.enqueueVerification(ctx, r, uc.outbox.WithTx(tx), created.IDUser, created.IDEmail, norm.Email, now, "register")
	}); err != nil {
		return nil, err
	}
	return &readmodel.RegisterReadModel{
		Message: registerNeutralMessage,
	}, nil
}

func (uc Register) resendVerification(ctx context.Context, idUser int64, email string) error {
	now := uc.clock.Now()
	return uc.uow.Do(ctx, func(tx uow.Tx) error {
		r := uc.repository.WithTx(tx)
		idEmail, err := r.User().FindPrimaryEmailIDByUser(ctx, idUser)
		if err != nil {
			return err
		}
		return uc.enqueueVerification(ctx, r, uc.outbox.WithTx(tx), idUser, idEmail, email, now, "register_resend")
	})
}

func (uc Register) enqueueVerification(
	ctx context.Context,
	repo repository.Factory,
	outbox messaging.OutboxFactory,
	idUser int64,
	idEmail int64,
	email string,
	now time.Time,
	usecaseName string,
) error {
	rawToken, hash, err := uc.refresh.Generate()
	if err != nil {
		return err
	}
	expiresAt := now.Add(uc.ttl)
	if err := repo.EmailVerification().RevokeAllByUser(ctx, idUser, now); err != nil {
		return err
	}
	idVerification, err := repo.EmailVerification().Create(ctx, &mapper.EmailVerificationTokenCreate{
		IDUser:    idUser,
		IDEmail:   idEmail,
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
		Template:  "email-verification",
		To:        email,
		Data: map[string]any{
			"id_verification": idVerification,
			"token":           rawToken,
			"expires_at":      int(expiresAt.Sub(now).Minutes()),
			"id_user":         idUser,
		},
		Meta: map[string]string{
			"module":  "identity",
			"usecase": usecaseName,
		},
	})
	if err != nil {
		return err
	}
	return outbox.Outbox().Add(ctx, messaging.OutboxMessage{
		ID:         idMessage,
		Subject:    event.SubjectEmailSendRequested,
		Payload:    payload,
		OccurredAt: now,
	})
}
