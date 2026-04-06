package adapter_test

import (
	"testing"

	"github.com/paladignus/actajus/internal/module/identity/presentation/grpc/adapter"
	identityv1 "github.com/paladignus/actajus/proto/identity/v1"
)

func TestProtoToAssignRoleCmd(t *testing.T) {
	t.Parallel()

	req := &identityv1.AssignRoleToUserRequest{
		IdUser:     10,
		IdRole:     3,
		AssignedBy: 99,
	}

	got := adapter.ProtoToAssignRoleCmd(req)

	if got.IDUser != 10 || got.IDRole != 3 || got.AssignedBy != 99 {
		t.Fatalf("unexpected assign role mapping: %+v", got)
	}
}

func TestProtoToRemoveRoleCmd(t *testing.T) {
	t.Parallel()

	req := &identityv1.RemoveRoleFromUserRequest{
		IdUser: 11,
		IdRole: 4,
	}

	got := adapter.ProtoToRemoveRoleCmd(req)

	if got.IDUser != 11 || got.IDRole != 4 {
		t.Fatalf("unexpected remove role mapping: %+v", got)
	}
}

func TestProtoToPermissionCommands(t *testing.T) {
	t.Parallel()

	grantReq := &identityv1.GrantPermissionToRoleRequest{
		IdRole:       5,
		IdPermission: 7,
	}
	grantGot := adapter.ProtoToGrantPermCmd(grantReq)
	if grantGot.IDRole != 5 || grantGot.IDPermission != 7 {
		t.Fatalf("unexpected grant permission mapping: %+v", grantGot)
	}

	revokeReq := &identityv1.RevokePermissionFromRoleRequest{
		IdRole:       8,
		IdPermission: 9,
	}
	revokeGot := adapter.ProtoToRevokePermCmd(revokeReq)
	if revokeGot.IDRole != 8 || revokeGot.IDPermission != 9 {
		t.Fatalf("unexpected revoke permission mapping: %+v", revokeGot)
	}
}
