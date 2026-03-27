// Package adapter
package adapter

import (
	"github.com/paladignus/actajus/internal/module/identity/application/command"
	identityv1 "github.com/paladignus/actajus/proto/identity/v1"
)

func ProtoToAssignRoleCmd(in *identityv1.AssignRoleToUserRequest) command.AssignRoleToUserCommand {
	return command.AssignRoleToUserCommand{
		IDUser:     in.GetIdUser(),
		IDRole:     int16(in.GetIdRole()),
		AssignedBy: in.GetAssignedBy(),
	}
}

func ProtoToRemoveRoleCmd(in *identityv1.RemoveRoleFromUserRequest) command.RemoveRoleFromUserCommand {
	return command.RemoveRoleFromUserCommand{
		IDUser: in.GetIdUser(),
		IDRole: int16(in.GetIdRole()),
	}
}

func ProtoToGrantPermCmd(in *identityv1.GrantPermissionToRoleRequest) command.GrantPermissionToRoleCommand {
	return command.GrantPermissionToRoleCommand{
		IDRole:       int16(in.GetIdRole()),
		IDPermission: int16(in.GetIdPermission()),
	}
}

func ProtoToRevokePermCmd(in *identityv1.RevokePermissionFromRoleRequest) command.RevokePermissionFromRoleCommand {
	return command.RevokePermissionFromRoleCommand{
		IDRole:       int16(in.GetIdRole()),
		IDPermission: int16(in.GetIdPermission()),
	}
}
