// Package person
package person

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/paladignus/actajus/internal/module/person/application/mapper"
	"github.com/paladignus/actajus/internal/module/person/application/usecase"
	"github.com/paladignus/actajus/internal/module/person/infrastructure/persistence/database/postgres"
	"github.com/paladignus/actajus/internal/module/person/presentation/grpc/handler"
	"github.com/paladignus/actajus/internal/shared/application/repository"
	postgresShared "github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
	"github.com/paladignus/actajus/proto/person/v1/personv1connect"
)

type Module struct {
	handler handler.PersonHandler
}

func NewModule(
	db postgresShared.PgxPool,
	logger repository.Logger,
) Module {
	personRepo := postgres.NewPerson(db)
	personProjection := mapper.NewPersonProjection()
	personMapper := mapper.NewPersonMapper()
	createPersonUC := usecase.NewCreatePerson(
		personRepo,
		personMapper,
		personProjection,
	)
	h := handler.NewPersonHandler(createPersonUC, logger)
	return Module{h}
}

func (m Module) Mount(mux *http.ServeMux, opts ...connect.HandlerOption) {
	path, h := personv1connect.NewPersonServiceHandler(m.handler, opts...)
	mux.Handle(path, h)
}

// func (m Module) Route(opts ...connect.HandlerOption) (string, http.Handler) {
// 	return personv1connect.NewPersonServiceHandler(m.handler, opts...)
// }
