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
	list   usecase.ListCompanies
	update usecase.UpdateCompany
	delete usecase.DeleteCompany
}

func NewCompanyHandler(
	create usecase.CreateCompany,
	list usecase.ListCompanies,
	update usecase.UpdateCompany,
	delete usecase.DeleteCompany,
) CompanyHandler {
	return CompanyHandler{
		create,
		list,
		update,
		delete,
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

func (h CompanyHandler) ListCompanies(
	ctx context.Context,
	req *connect.Request[companyv1.ListCompaniesRequest],
) (*connect.Response[companyv1.ListCompaniesResponse], error) {
	return *connect.NewResponse(&companyv1.ListCompaniesResponse{}), nil
}

func (h CompanyHandler) UpdateCompany(
	ctx context.Context,
	req *connect.Request[companyv1.UpdateCompanyRequest],
) (*connect.Response[companyv1.UpdateCompanyResponse], error) {
	cmd := adapter.ProtoToCompanyUpdateCommand(req.Msg)
	rm, err := h.update.Execute(ctx, cmd)
	if err != nil {
		return nil, mapErr(err)
	}
	return connect.NewResponse(&companyv1.UpdateCompanyResponse{
		Id: int64(rm.ID),
	}), nil
}

func (h CompanyHandler) DeleteCompany(
	ctx context.Context,
	req *connect.Request[companyv1.DeleteCompanyRequest],
) (*connect.Response[companyv1.DeleteCompanyResponse], error) {
	err := h.delete.Execute(ctx, req.Msg.Id)
	if err != nil {
		return nil, mapErr(err)
	}
	return connect.NewResponse(&companyv1.DeleteCompanyResponse{}), nil
}

func mapErr(err error) error {
	var ve sharedDomain.ValidationError
	if errors.As(err, &ve) {
		return sharedAdapter.ValidationErrorDetail(ve)
	}
	// Erros de identidade (regra/negócio/autenticação)
	return sharedAdapter.ToConnectIdentityError(err)
}
