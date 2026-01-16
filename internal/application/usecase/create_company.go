// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/entity"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type Company struct {
	uow repository.UnitOfWorkCompany
}

func NewCompany(
	uow repository.UnitOfWorkCompany,
) Company {
	return Company{uow}
}

func (c Company) Execute(ctx context.Context, input dto.CompanyInputOutput) error {
	if err := c.uow.Begin(ctx); err != nil {
		return fmt.Errorf("database error for begin transaction: %w", err)
	}
	defer c.uow.Rollback(ctx)
	company := entity.NewCompany(input.Company)
	if err := company.Create(); err != nil {
		return fmt.Errorf("use case create company, invalid input: %w", err)
	}
	idCompany, err := c.uow.Company().Create(ctx, company)
	if err != nil {
		return fmt.Errorf("database error while saving company: %w", err)
	}
	address := entity.NewAddress(input.Address)
	if err = company.Create(); err != nil {
		return fmt.Errorf("use case create company, invalid input: %w", err)
	}
	idAddress, err := c.uow.Address().Create(ctx, address)
	if err != nil {
		return fmt.Errorf("database error while saving address company: %w", err)
	}
	if err = c.uow.CompanyAddress().Create(ctx, idCompany, idAddress); err != nil {
		return fmt.Errorf("database error while saving address company relation: %w", err)
	}
	phone := entity.NewPhone(input.Phone)
	idPhone, err := c.uow.Phone().Create(ctx, phone)
	if err != nil {
		return fmt.Errorf("database error while saving phone company: %w", err)
	}
	if err := c.uow.CompanyPhone().Create(ctx, idCompany, idPhone); err != nil {
		return fmt.Errorf("database error while saving company phone relation: %w", err)
	}
	email := entity.NewEmail(input.Email)
	if err := email.Create(); err != nil {
		return fmt.Errorf("use case create company, invalid input: %w", err)
	}
	idEmail, err := c.uow.Email().Create(ctx, email)
	if err != nil {
		return fmt.Errorf("database error while saving email company: %w", err)
	}
	if err := c.uow.CompanyEmail().Create(ctx, idCompany, idEmail); err != nil {
		return fmt.Errorf("database error while saving company email relation: %w", err)
	}
	for _, input := range input.SocialMedia {
		socialMedia := entity.NewSocialMedia(input)
		if err = socialMedia.Create(); err != nil {
			return fmt.Errorf("use case create company, invalid input: %w", err)
		}
		socialMedia.IDCompany = idCompany
		if err := c.uow.SocialMedia().Create(ctx, socialMedia); err != nil {
			return fmt.Errorf("database error while saving social media company: %w", err)
		}
	}
	return c.uow.Commit(ctx)
}
