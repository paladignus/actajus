// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/module/company/domain"
	sharedDomain "github.com/paladignus/actajus/internal/shared/domain"
)

type DeleteCompany struct {
	uow domain.CompanyUnitOfWork
}

func NewDeleteCompany(uow domain.CompanyUnitOfWork) DeleteCompany {
	return DeleteCompany{uow}
}

func (d DeleteCompany) Execute(ctx context.Context, id uint) error {
	if err := d.uow.Begin(ctx); err != nil {
		return err
	}
	defer d.uow.Rollback(ctx)
	company, err := d.uow.Company().FindByID(ctx, id)
	if err != nil {
		return err
	}
	if company == nil {
		return sharedDomain.NewFieldError("id", "company not found")
	}
	company.Delete()
	if err := d.uow.Company().Delete(ctx, company); err != nil {
		return err
	}
	if err := d.uow.Commit(ctx); err != nil {
		return err
	}
	return nil
}
