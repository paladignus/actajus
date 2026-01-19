// Package company
package company

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paladignus/actajus/internal/module/company/application/mapper"
	"github.com/paladignus/actajus/internal/module/company/application/usecase"
	"github.com/paladignus/actajus/internal/module/company/infrastructure/persistence"
	"github.com/paladignus/actajus/internal/module/company/presentation/handler"
)

type Module struct {
	Handler handler.CompanyHandler
}

func NewModule(pool *pgxpool.Pool) Module {
	persistence := persistence.NewCompanyRepositoryPostgres(pool)
	mapper := mapper.NewCompanyMapper()
	createUC := usecase.NewCreateCompany(persistence, *mapper)
	handler := handler.NewCompanyHandler(createUC)
	return Module{handler}
}
