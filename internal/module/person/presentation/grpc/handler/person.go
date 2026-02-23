// Package handler
package handler

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/paladignus/actajus/internal/module/person/application/usecase"
	"github.com/paladignus/actajus/internal/module/person/presentation/grpc/adapter"
	sharedAdapter "github.com/paladignus/actajus/internal/shared/application/adapter"
	"github.com/paladignus/actajus/internal/shared/domain"
	"github.com/paladignus/actajus/internal/shared/domain/repository"
	"github.com/paladignus/actajus/internal/shared/presentation/validation"
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
	// input := mapper.ProtoToCreatePersonDTO(req.Msg)
	input := adapter.ProtoToCreatePersonDTO(req.Msg)
	vs := validation.New().ValidateStruct(input)
	// if err := mapper.ValidateCreatePersonDTO(input); err != nil {
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		var ve domain.ValidationError
		if errors.As(err, &ve) {
			return nil, err
			// return nil, sharedAdapter.ViolationsToDomainError(ve)
		}
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	// person, err := mapper.PersonInputToDomain(dtoIn)
	// if err != nil {
	// 	var ve domain.ValidationError
	// 	if errors.As(err, &ve) {
	// 		return nil, connecterr.FromDomainValidation(ve)
	// 	}
	// 	return nil, connect.NewError(connect.CodeFailedPrecondition, err)
	// }
	person, err := p.person.Execute(ctx, input)
	if err != nil {
		var ve domain.ValidationError
		if errors.As(err, &ve) {
			return nil, err
			// return nil, sharedAdapter.ViolationsToDomainError(ve)
		}
		return nil, connect.NewError(connect.CodeFailedPrecondition, err)
	}
	// if err != nil {
	// 	return nil, connect.NewError(connect.CodeUnknown, err)
	// }
	return connect.NewResponse(&personv1.CreatePersonResponse{Id: uint32(person.ID)}), nil
}
