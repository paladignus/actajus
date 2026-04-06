// Package handler
package handler

import (
	"context"

	"connectrpc.com/connect"

	"github.com/paladignus/actajus/internal/module/identity/application/usecase"
	"github.com/paladignus/actajus/internal/module/identity/presentation/grpc/adapter"
	identityv1 "github.com/paladignus/actajus/proto/identity/v1"
)

type RbacAdminHandler struct {
	assign usecase.AssignRoleToUser
	remove usecase.RemoveRoleFromUser
	grant  usecase.GrantPermissionToRole
	revoke usecase.RevokePermissionFromRole
}

func NewRbacAdminHandler(
	assign usecase.AssignRoleToUser,
	remove usecase.RemoveRoleFromUser,
	grant usecase.GrantPermissionToRole,
	revoke usecase.RevokePermissionFromRole,
) *RbacAdminHandler {
	return &RbacAdminHandler{assign, remove, grant, revoke}
}

func (h *RbacAdminHandler) AssignRoleToUser(
	ctx context.Context,
	req *connect.Request[identityv1.AssignRoleToUserRequest],
) (*connect.Response[identityv1.AssignRoleToUserResponse], error) {
	cmd := adapter.ProtoToAssignRoleCmd(req.Msg)
	if err := h.assign.Execute(ctx, cmd); err != nil {
		return nil, mapErr(err)
	}
	return connect.NewResponse(&identityv1.AssignRoleToUserResponse{Ok: true}), nil
}

func (h *RbacAdminHandler) RemoveRoleFromUser(
	ctx context.Context,
	req *connect.Request[identityv1.RemoveRoleFromUserRequest],
) (*connect.Response[identityv1.RemoveRoleFromUserResponse], error) {
	cmd := adapter.ProtoToRemoveRoleCmd(req.Msg)
	if err := h.remove.Execute(ctx, cmd); err != nil {
		return nil, mapErr(err)
	}
	return connect.NewResponse(&identityv1.RemoveRoleFromUserResponse{Ok: true}), nil
}

func (h *RbacAdminHandler) GrantPermissionToRole(
	ctx context.Context,
	req *connect.Request[identityv1.GrantPermissionToRoleRequest],
) (*connect.Response[identityv1.GrantPermissionToRoleResponse], error) {
	cmd := adapter.ProtoToGrantPermCmd(req.Msg)
	if err := h.grant.Execute(ctx, cmd); err != nil {
		return nil, mapErr(err)
	}
	return connect.NewResponse(&identityv1.GrantPermissionToRoleResponse{Ok: true}), nil
}

func (h *RbacAdminHandler) RevokePermissionFromRole(
	ctx context.Context,
	req *connect.Request[identityv1.RevokePermissionFromRoleRequest],
) (*connect.Response[identityv1.RevokePermissionFromRoleResponse], error) {
	cmd := adapter.ProtoToRevokePermCmd(req.Msg)
	if err := h.revoke.Execute(ctx, cmd); err != nil {
		return nil, mapErr(err)
	}
	return connect.NewResponse(&identityv1.RevokePermissionFromRoleResponse{Ok: true}), nil
}

// func mapErr(err error) error {
// 	var ve sharedDomain.ValidationError
// 	if errors.As(err, &ve) {
// 		return sharedAdapter.ValidationErrorDetail(ve)
// 	}
// 	return sharedAdapter.ToConnectIdentityError(err)
// }
