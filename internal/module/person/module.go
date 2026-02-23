// Package person
package person

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paladignus/actajus/internal/module/person/application/mapper"
	"github.com/paladignus/actajus/internal/module/person/application/usecase"
	"github.com/paladignus/actajus/internal/module/person/infrastructure/persistence/database/postgres"
	"github.com/paladignus/actajus/internal/module/person/presentation/grpc/handler"
	"github.com/paladignus/actajus/internal/shared/domain/repository"
	"github.com/paladignus/actajus/proto/person/v1/personv1connect"
)

type Module struct {
	Handler handler.PersonHandler
}

func (m Module) Route(opts ...connect.HandlerOption) (string, http.Handler) {
	return personv1connect.NewPersonServiceHandler(m.Handler, opts...)
}

func NewModule(
	pool *pgxpool.Pool,
	logger repository.Logger,
) Module {
	repository := postgres.NewPerson(pool)
	projection := mapper.NewPersonProjection()
	mapper := mapper.NewPersonMapper()
	usecase := usecase.NewCreatePerson(
		repository,
		mapper,
		projection,
	)
	h := handler.NewPersonHandler(usecase, logger)
	return Module{h}
}
