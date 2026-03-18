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
	create     usecase.CreateCompany
	list       usecase.ListCompanies
	update     usecase.UpdateCompany
	delete     usecase.DeleteCompany
	findByID   usecase.FindByID
	findByCNPJ usecase.FindByCNPJ
}

func NewCompanyHandler(
	create usecase.CreateCompany,
	list usecase.ListCompanies,
	update usecase.UpdateCompany,
	delete usecase.DeleteCompany,
	findByID usecase.FindByID,
	findByCNPJ usecase.FindByCNPJ,
) CompanyHandler {
	return CompanyHandler{
		create,
		list,
		update,
		delete,
		findByID,
		findByCNPJ,
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
	baseURL := "http://localhost:50051"
	limit := 10
	if req.Msg.Limit != nil && *req.Msg.Limit > 0 && *req.Msg.Limit < 100 {
		limit = int(*req.Msg.Limit)
	}
	companies, err := h.list.Execute(ctx, req.Msg.After, req.Msg.Before, limit, baseURL)
	if err != nil {
		return nil, mapErr(err)
	}
	adapter.CompaniesReadModelToProto(companies)
	return connect.NewResponse(adapter.CompaniesReadModelToProto(companies)), nil
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

func (h CompanyHandler) FindCompanyByID(
	ctx context.Context,
	req *connect.Request[companyv1.FindCompanyByIDRequest],
) (*connect.Response[companyv1.FindCompanyByIDResponse], error) {
	company, err := h.findByID.Execute(ctx, req.Msg.Id)
	if err != nil {
		return nil, mapErr(err)
	}
	return connect.NewResponse(adapter.FindByIDCompanyReadModelToProto(company)), nil
}

func (h CompanyHandler) FindCompanyByCNPJ(
	ctx context.Context,
	req *connect.Request[companyv1.FindCompanyByCNPJRequest],
) (*connect.Response[companyv1.FindCompanyByCNPJResponse], error) {
	company, err := h.findByCNPJ.Execute(ctx, req.Msg.Cnpj)
	if err != nil {
		return nil, mapErr(err)
	}
	return connect.NewResponse(adapter.FindByCNPJCompanyReadModelToProto(company)), nil
}

func mapErr(err error) error {
	var ve sharedDomain.ValidationError
	if errors.As(err, &ve) {
		return sharedAdapter.ValidationErrorDetail(ve)
	}
	// Erros de identidade (regra/negócio/autenticação)
	return sharedAdapter.ToConnectIdentityError(err)
}
