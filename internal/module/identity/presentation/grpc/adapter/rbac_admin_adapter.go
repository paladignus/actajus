// Package adapter
package adapter

import (
	"github.com/paladignus/actajus/internal/module/identity/application/dto"
	identityv1 "github.com/paladignus/actajus/proto/identity/v1"
)

func ProtoToAssignRoleCmd(in *identityv1.AssignRoleToUserRequest) dto.AssignRoleToUserCommand {
	return dto.AssignRoleToUserCommand{
		IDUser:     in.GetIdUser(),
		IDRole:     int16(in.GetIdRole()),
		AssignedBy: in.GetAssignedBy(),
	}
}

func ProtoToRemoveRoleCmd(in *identityv1.RemoveRoleFromUserRequest) dto.RemoveRoleFromUserCommand {
	return dto.RemoveRoleFromUserCommand{
		IDUser: in.GetIdUser(),
		IDRole: int16(in.GetIdRole()),
	}
}

func ProtoToGrantPermCmd(in *identityv1.GrantPermissionToRoleRequest) dto.GrantPermissionToRoleCommand {
	return dto.GrantPermissionToRoleCommand{
		IDRole:       int16(in.GetIdRole()),
		IDPermission: int16(in.GetIdPermission()),
	}
}

func ProtoToRevokePermCmd(in *identityv1.RevokePermissionFromRoleRequest) dto.RevokePermissionFromRoleCommand {
	return dto.RevokePermissionFromRoleCommand{
		IDRole:       int16(in.GetIdRole()),
		IDPermission: int16(in.GetIdPermission()),
	}
}
