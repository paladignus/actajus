// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/identity/application/command"
	"github.com/paladignus/actajus/internal/module/identity/application/mapper"
	"github.com/paladignus/actajus/internal/module/identity/application/readmodel"
)

type RequestEmailVerification struct {
	register Register
	mapper   mapper.AuthMapper
}

func NewRequestEmailVerification(register Register, mapper mapper.AuthMapper) RequestEmailVerification {
	return RequestEmailVerification{
		register: register,
		mapper:   mapper,
	}
}

func (uc RequestEmailVerification) Execute(ctx context.Context, input command.RequestEmailVerificationCommand) (*readmodel.RegisterReadModel, error) {
	norm, err := uc.mapper.RequestPasswordResetInputToNormalized(command.RequestPasswordResetCommand{
		Email: input.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("invalid request email verification data: %w", err)
	}
	user, err := uc.register.repository.User().FindByEmail(ctx, norm.Email)
	if err != nil {
		return nil, err
	}
	if user == nil || user.IsPrimaryEmailVerified() {
		return &readmodel.RegisterReadModel{Message: registerNeutralMessage}, nil
	}
	if err := uc.register.resendVerification(ctx, user.ID().Value(), norm.Email); err != nil {
		return nil, err
	}
	return &readmodel.RegisterReadModel{Message: registerNeutralMessage}, nil
}
