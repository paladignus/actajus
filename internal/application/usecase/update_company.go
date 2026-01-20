// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/entity"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type UpdateCompany struct {
	uow repository.UnitOfWorkCompany
}

func NewUpdateCompany(uow repository.UnitOfWorkCompany) UpdateCompany {
	return UpdateCompany{uow}
}

func (u UpdateCompany) Execute(ctx context.Context, input dto.CompanyInputOutput) error {
	if err := u.uow.Begin(ctx); err != nil {
		return fmt.Errorf("database error for begin transaction: %w", err)
	}
	defer u.uow.Rollback(ctx)
	company := entity.NewCompany(input.Company)
	if err := company.Update(); err != nil {
		return fmt.Errorf("use case update company, invalid input: %w", err)
	}
	err := u.uow.Company().Update(ctx, company)
	if err != nil {
		return fmt.Errorf("database error while update company: %w", err)
	}
	address := entity.NewAddress(input.Address)
	if err := address.Update(); err != nil {
		return fmt.Errorf("use case update company, invalid input: %w", err)
	}
	err = u.uow.Address().Update(ctx, address)
	if err != nil {
		return fmt.Errorf("database error while updating address company: %w", err)
	}
	phone := entity.NewPhone(input.Phone)
	if err := phone.Update(); err != nil {
		return fmt.Errorf("use case update company, invalid input: %w", err)
	}
	if err = u.uow.Phone().Update(ctx, phone); err != nil {
		return fmt.Errorf("database error while updating phone company: %w", err)
	}
	email := entity.NewEmail(input.Email)
	if err = email.Update(); err != nil {
		return fmt.Errorf("use case update company, invalid input: %w", err)
	}
	if err = u.uow.Email().Update(ctx, email); err != nil {
		return fmt.Errorf("database error while saving email company: %w", err)
	}
	for _, input := range input.SocialMedia {
		socialMedia := entity.NewSocialMedia(input)
		if err = socialMedia.Update(); err != nil {
			return fmt.Errorf("use case update company, invalid input: %w", err)
		}
		if err := u.uow.SocialMedia().Update(ctx, socialMedia); err != nil {
			return fmt.Errorf("database error while updating social media company: %w", err)
		}
	}
	return u.uow.Commit(ctx)
}
