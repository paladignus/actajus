// Package handler
package handler

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/paladignus/actajus/internal/module/company/application/usecase"
	"github.com/paladignus/actajus/internal/module/company/presentation/grpc/adapter"
	sharedDomain "github.com/paladignus/actajus/internal/shared/domain"
	sharedAdapter "github.com/paladignus/actajus/internal/shared/presentation/adapter"
	"github.com/paladignus/actajus/internal/shared/presentation/authctx"
	companyv1 "github.com/paladignus/actajus/proto/company/v1"
)

type CompanyHandler struct {
	create usecase.CreateCompany
}

func NewCompanyHandler(
	create usecase.CreateCompany,
) CompanyHandler {
	return CompanyHandler{
		create,
	}
}

func (h CompanyHandler) CreateCompany(
	ctx context.Context,
	req *connect.Request[companyv1.CreateCompanyRequest],
) (*connect.Response[companyv1.CreateCompanyResponse], error) {
	c := authctx.MustGetClaims(ctx)
	cmd := adapter.ProtoToCompanyCreateCommand(req.Msg)
	cmd.RegisteredBy = c.IDUser
	rm, err := h.create.Execute(ctx, cmd)
	if err != nil {
		return nil, mapErr(err)
	}
	return connect.NewResponse(&companyv1.CreateCompanyResponse{
		Id: int64(rm.ID),
	}), nil
}

func mapErr(err error) error {
	var ve sharedDomain.ValidationError
	if errors.As(err, &ve) {
		return sharedAdapter.ValidationErrorDetail(ve)
	}
	// Erros de identidade (regra/negócio/autenticação)
	return sharedAdapter.ToConnectIdentityError(err)
}
