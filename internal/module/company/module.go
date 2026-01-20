// Package company
package company

import (
	"github.com/jackc/pgx/v5/pgxpool"
	address "github.com/paladignus/actajus/internal/module/address/application/mapper"
	"github.com/paladignus/actajus/internal/module/company/application/mapper"
	"github.com/paladignus/actajus/internal/module/company/application/usecase"
	"github.com/paladignus/actajus/internal/module/company/infrastructure/persistence/database"
	"github.com/paladignus/actajus/internal/module/company/presentation/handler"
)

type Module struct {
	Handler handler.CompanyHandler
}

func NewModule(pool *pgxpool.Pool) Module {
	uow := database.NewCompanyUnitOfWork(pool)
	address := address.NewAddressMapper()
	mapper := mapper.NewCompanyMapper(address)
	createUC := usecase.NewCreateCompany(&uow, *mapper)
	handler := handler.NewCompanyHandler(createUC)
	return Module{handler}
}
