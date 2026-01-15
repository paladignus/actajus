// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/entity"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type EnterpriseUpdate struct {
	uow repository.UnitOfWorkEnterprise
}

func NewEnterpriseUpdate(uow repository.UnitOfWorkEnterprise) EnterpriseUpdate {
	return EnterpriseUpdate{uow}
}

func (e EnterpriseUpdate) Execute(ctx context.Context, input dto.EnterpriseInput) error {
	if err := e.uow.Begin(ctx); err != nil {
		return fmt.Errorf("database error for begin transaction: %w", err)
	}
	defer e.uow.Rollback(ctx)
	enterprise := entity.NewEnterprise(input.Enterprise)
	if err := enterprise.Update(); err != nil {
		return fmt.Errorf("use case update enterprise, invalid input: %w", err)
	}
	err := e.uow.Enterprise().Update(ctx, enterprise)
	if err != nil {
		return fmt.Errorf("database error while update enterprise: %w", err)
	}

	address := entity.NewAddress(input.Address)
	if err := address.Update(); err != nil {
		return fmt.Errorf("use case update enterprise, invalid input: %w", err)
	}
	err = e.uow.Address().Update(ctx, address)
	if err != nil {
		return fmt.Errorf("database error while updating address: %w", err)
	}

	phone := entity.NewPhone(input.Phone)
	if err := phone.Update(); err != nil {
		return fmt.Errorf("use case update enterprise, invalid input: %w", err)
	}
	if err = e.uow.Phone().Update(ctx, phone); err != nil {
		return fmt.Errorf("database error while updating phone: %w", err)
	}
	email := entity.NewEmail(input.Email)
	if err = email.Update(); err != nil {
		return fmt.Errorf("use case update enterprise, invalid input: %w", err)
	}
	if err = e.uow.Email().Update(ctx, email); err != nil {
		return fmt.Errorf("database error while saving email: %w", err)
	}

	for _, input := range input.SocialMedia {
		socialMedia := entity.NewSocialMedia(input)
		if err = socialMedia.Update(); err != nil {
			return fmt.Errorf("use case update enterprise, invalid input: %w", err)
		}
		if err := e.uow.SocialMedia().Update(ctx, socialMedia); err != nil {
			return fmt.Errorf("database error while updating social media: %w", err)
		}
	}
	return e.uow.Commit(ctx)
}
