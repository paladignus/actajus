// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/identity/application/dto"
	"github.com/paladignus/actajus/internal/module/identity/application/mapper"
	"github.com/paladignus/actajus/internal/module/identity/application/service"
)

type Login struct {
	mapper     mapper.AuthMapper
	projection mapper.AuthProjectionMapper
	authn      service.AuthnService
}

func NewLogin(
	mapper mapper.AuthMapper,
	projection mapper.AuthProjectionMapper,
	authn service.AuthnService,
) Login {
	return Login{
		mapper:     mapper,
		projection: projection,
		authn:      authn,
	}
}

func (uc Login) Execute(ctx context.Context, input dto.LoginCommand) (*dto.AuthTokensReadModel, error) {
	norm, err := uc.mapper.LoginInputToNormalized(input)
	if err != nil {
		return nil, fmt.Errorf("invalid login data: %w", err)
	}

	result, err := uc.authn.Authenticate(
		ctx,
		norm.Email,
		norm.Password,
		norm.IP,
		norm.UserAgent,
	)
	if err != nil {
		return nil, err
	}

	return uc.projection.ProjectTokens(
		result.Session.ID().Value(),
		result.User.ID().Value(),
		result.AccessToken,
		result.RefreshToken,
		result.AccessExp,
		result.RefreshExp,
	), nil
}
