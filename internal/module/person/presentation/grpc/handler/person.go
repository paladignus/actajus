// Package handler
package handler

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/paladignus/actajus/internal/module/person/application/usecase"
	personAdapter "github.com/paladignus/actajus/internal/module/person/presentation/grpc/adapter"
	"github.com/paladignus/actajus/internal/shared/application/repository"
	"github.com/paladignus/actajus/internal/shared/domain"
	sharedAdapter "github.com/paladignus/actajus/internal/shared/presentation/adapter"
	personv1 "github.com/paladignus/actajus/proto/person/v1"
)

type PersonHandler struct {
	person usecase.CreatePerson
	logger repository.Logger
}

func NewPersonHandler(
	person usecase.CreatePerson,
	logger repository.Logger,
) PersonHandler {
	return PersonHandler{
		person,
		logger,
	}
}

func (p PersonHandler) CreatePerson(
	ctx context.Context,
	req *connect.Request[personv1.CreatePersonRequest],
) (*connect.Response[personv1.CreatePersonResponse], error) {
	input := personAdapter.ProtoToCreatePersonDTO(req.Msg)

	person, err := p.person.Execute(ctx, input)
	if err != nil {
		var ve domain.ValidationError
		if errors.As(err, &ve) {
			return nil, sharedAdapter.ValidationErrorDetail(ve)
		}
		return nil, connect.NewError(connect.CodeFailedPrecondition, err)
	}
	return connect.NewResponse(&personv1.CreatePersonResponse{Id: uint32(person.ID)}), nil
}
