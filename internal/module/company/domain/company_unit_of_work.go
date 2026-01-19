// Package domain
package domain

import (
	uow "github.com/paladignus/actajus/internal/shared/domain/unit_of_work"
)

type CompanyUnitOfWork interface {
	uow.UnitOfWork
	Company() CompanyRepository
}
