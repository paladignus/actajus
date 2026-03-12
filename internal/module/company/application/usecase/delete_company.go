// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/company/application/repository"
	"github.com/paladignus/actajus/internal/shared/application/uow"
	sharedDomain "github.com/paladignus/actajus/internal/shared/domain"
)

type DeleteCompany struct {
	uow        uow.UnitOfWork
	repository repository.Factory
}

func NewDeleteCompany(
	uow uow.UnitOfWork,
	repository repository.Factory,
) DeleteCompany {
	return DeleteCompany{uow, repository}
}

func (u DeleteCompany) Execute(ctx context.Context, id uint) error {
	company, err := u.repository.Company().FindByID(ctx, id)
	if err != nil {
		return err
	}
	if company == nil {
		return sharedDomain.NewFieldError("id", "company not found")
	}
	if err := company.Delete(); err != nil {
		return err
	}
	err = u.uow.Do(ctx, func(tx uow.Tx) error {
		r := u.repository.WithTx(tx)
		if err := r.CompanyAddress().DeleteByIDCompany(ctx, id); err != nil {
			return fmt.Errorf("failed to delete company address relationship: %w", err)
		}
		if err := r.CompanyPhone().DeleteByIDCompany(ctx, id); err != nil {
			return fmt.Errorf("failed to delete company phone relationship: %w", err)
		}
		if err := r.CompanyEmail().DeleteByIDCompany(ctx, id); err != nil {
			return fmt.Errorf("failed to delete company email relationship: %w", err)
		}
		if err := r.SocialMedia().DeleteByIDCompany(ctx, id); err != nil {
			return fmt.Errorf("failed to delete company social media relationship: %w", err)
		}
		if err := r.Company().Delete(ctx, company); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
