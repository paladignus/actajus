// Package handler
package handler

import (
	"context"
	"errors"
	"log"

	"connectrpc.com/connect"
	"github.com/paladignus/actajus/internal/module/identity/application/usecase"
	"github.com/paladignus/actajus/internal/module/identity/presentation/grpc/adapter"
	sharedDomain "github.com/paladignus/actajus/internal/shared/domain"
	sharedAdapter "github.com/paladignus/actajus/internal/shared/presentation/adapter"
	identityv1 "github.com/paladignus/actajus/proto/identity/v1"
)

type AuthHandler struct {
	login          usecase.Login
	refresh        usecase.Refresh
	logout         usecase.Logout
	logoutAll      usecase.LogoutAll
	changePassword usecase.ChangePassword
	requestReset   usecase.RequestPasswordReset
	confirmReset   usecase.ConfirmPasswordReset
}

func NewAuthHandler(
	login usecase.Login,
	refresh usecase.Refresh,
	logout usecase.Logout,
	logoutAll usecase.LogoutAll,
	changePassword usecase.ChangePassword,
	requestReset usecase.RequestPasswordReset,
	confirmReset usecase.ConfirmPasswordReset,
) AuthHandler {
	return AuthHandler{
		login, refresh,
		logout, logoutAll,
		changePassword,
		requestReset, confirmReset,
	}
}

func (h AuthHandler) Login(
	ctx context.Context,
	req *connect.Request[identityv1.LoginRequest],
) (*connect.Response[identityv1.LoginResponse], error) {
	cmd := adapter.ProtoToLoginCommand(req.Msg)
	rm, err := h.login.Execute(ctx, cmd)
	if err != nil {
		return nil, mapErr(err)
	}
	return connect.NewResponse(adapter.LoginReadModelToProto(rm)), nil
}

func (h AuthHandler) Refresh(
	ctx context.Context,
	req *connect.Request[identityv1.RefreshRequest],
) (*connect.Response[identityv1.RefreshResponse], error) {
	log.Printf("refresh sid=%d token_prefix=%s", req.Msg.IdSession, req.Msg.RefreshToken[:8])
	cmd := adapter.ProtoToRefreshCommand(req.Msg)
	log.Println(req.Msg)
	rm, err := h.refresh.Execute(ctx, cmd)
	log.Println(rm)
	if err != nil {
		log.Println(err)
		return nil, mapErr(err)
	}
	return connect.NewResponse(adapter.RefreshReadModelToProto(rm)), nil
}

func (h AuthHandler) Logout(
	ctx context.Context,
	req *connect.Request[identityv1.LogoutRequest],
) (*connect.Response[identityv1.LogoutResponse], error) {
	cmd := adapter.ProtoToLogoutCommand(req.Msg)
	if err := h.logout.Execute(ctx, cmd); err != nil {
		return nil, mapErr(err)
	}
	return connect.NewResponse(&identityv1.LogoutResponse{Ok: true}), nil
}

func (h AuthHandler) LogoutAll(
	ctx context.Context,
	req *connect.Request[identityv1.LogoutAllRequest],
) (*connect.Response[identityv1.LogoutAllResponse], error) {
	cmd := adapter.ProtoToLogoutAllCommand(req.Msg)
	if err := h.logoutAll.Execute(ctx, cmd); err != nil {
		return nil, mapErr(err)
	}
	return connect.NewResponse(&identityv1.LogoutAllResponse{Ok: true}), nil
}

func (h AuthHandler) ChangePassword(
	ctx context.Context,
	req *connect.Request[identityv1.ChangePasswordRequest],
) (*connect.Response[identityv1.ChangePasswordResponse], error) {
	cmd := adapter.ProtoToChangePasswordCommand(req.Msg)
	if err := h.changePassword.Execute(ctx, cmd); err != nil {
		return nil, mapErr(err)
	}
	return connect.NewResponse(&identityv1.ChangePasswordResponse{Ok: true}), nil
}

func (h AuthHandler) RequestPasswordReset(
	ctx context.Context,
	req *connect.Request[identityv1.RequestPasswordResetRequest],
) (*connect.Response[identityv1.RequestPasswordResetResponse], error) {
	cmd := adapter.ProtoToRequestPasswordResetCommand(req.Msg)
	rm, err := h.requestReset.Execute(ctx, cmd)
	if err != nil {
		return nil, mapErr(err)
	}
	return connect.NewResponse(adapter.RequestPasswordResetReadModelToProto(rm)), nil
}

func (h AuthHandler) ConfirmPasswordReset(
	ctx context.Context,
	req *connect.Request[identityv1.ConfirmPasswordResetRequest],
) (*connect.Response[identityv1.ConfirmPasswordResetResponse], error) {
	cmd := adapter.ProtoToConfirmPasswordResetCommand(req.Msg)
	if err := h.confirmReset.Execute(ctx, cmd); err != nil {
		return nil, mapErr(err)
	}
	return connect.NewResponse(&identityv1.ConfirmPasswordResetResponse{Ok: true}), nil
}

func mapErr(err error) error {
	var ve sharedDomain.ValidationError
	if errors.As(err, &ve) {
		return sharedAdapter.ValidationErrorDetail(ve)
	}
	// Erros de identidade (regra/negócio/autenticação)
	return sharedAdapter.ToConnectIdentityError(err)
}

// Só pra evitar import não usado, se Login/Refresh use dto
// var _ = dto.LoginCommand{}
