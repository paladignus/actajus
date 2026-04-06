// Package adapter
package adapter

import (
	"github.com/paladignus/actajus/internal/module/identity/application/command"
	identityv1 "github.com/paladignus/actajus/proto/identity/v1"
)

func ProtoToLoginCommand(in *identityv1.LoginRequest) command.LoginCommand {
	return command.LoginCommand{
		Email:     in.GetEmail(),
		Password:  in.GetPassword(),
		IP:        in.GetIp(),
		UserAgent: in.GetUserAgent(),
	}
}

func ProtoToRefreshCommand(in *identityv1.RefreshRequest) command.RefreshCommand {
	return command.RefreshCommand{
		IDSession:    in.GetIdSession(),
		RefreshToken: in.GetRefreshToken(),
		IP:           in.GetIp(),
		UserAgent:    in.GetUserAgent(),
	}
}

func ProtoToLogoutCommand(in *identityv1.LogoutRequest) command.LogoutCommand {
	return command.LogoutCommand{
		IDSession: in.GetIdSession(),
	}
}

func ProtoToLogoutAllCommand(in *identityv1.LogoutAllRequest) command.LogoutAllCommand {
	return command.LogoutAllCommand{
		IDUser: in.GetIdUser(),
	}
}

func ProtoToChangePasswordCommand(in *identityv1.ChangePasswordRequest) command.ChangePasswordCommand {
	return command.ChangePasswordCommand{
		IDUser:          in.GetIdUser(),
		CurrentPassword: in.GetCurrentPassword(),
		NewPassword:     in.GetNewPassword(),
	}
}

func ProtoToRequestPasswordResetCommand(in *identityv1.RequestPasswordResetRequest) command.RequestPasswordResetCommand {
	return command.RequestPasswordResetCommand{
		Email: in.GetEmail(),
	}
}

func ProtoToConfirmPasswordResetCommand(in *identityv1.ConfirmPasswordResetRequest) command.ConfirmPasswordResetCommand {
	return command.ConfirmPasswordResetCommand{
		IDReset:     in.GetIdReset(),
		ResetToken:  in.GetResetToken(),
		NewPassword: in.GetNewPassword(),
	}
}
