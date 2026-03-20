// Package adapter
package adapter

import (
	"github.com/paladignus/actajus/internal/module/identity/application/command"
	identityv1 "github.com/paladignus/actajus/proto/identity/v1"
)

func ProtoToLoginCommand(in *identityv1.LoginRequest) dto.LoginCommand {
	return dto.LoginCommand{
		Email:     in.GetEmail(),
		Password:  in.GetPassword(),
		IP:        in.GetIp(),
		UserAgent: in.GetUserAgent(),
	}
}

func ProtoToRefreshCommand(in *identityv1.RefreshRequest) dto.RefreshCommand {
	return dto.RefreshCommand{
		IDSession:    in.GetIdSession(),
		RefreshToken: in.GetRefreshToken(),
		IP:           in.GetIp(),
		UserAgent:    in.GetUserAgent(),
	}
}

func ProtoToLogoutCommand(in *identityv1.LogoutRequest) dto.LogoutCommand {
	return dto.LogoutCommand{
		IDSession: in.GetIdSession(),
	}
}

func ProtoToLogoutAllCommand(in *identityv1.LogoutAllRequest) dto.LogoutAllCommand {
	return dto.LogoutAllCommand{
		IDUser: in.GetIdUser(),
	}
}

func ProtoToChangePasswordCommand(in *identityv1.ChangePasswordRequest) dto.ChangePasswordCommand {
	return dto.ChangePasswordCommand{
		IDUser:          in.GetIdUser(),
		CurrentPassword: in.GetCurrentPassword(),
		NewPassword:     in.GetNewPassword(),
	}
}

func ProtoToRequestPasswordResetCommand(in *identityv1.RequestPasswordResetRequest) dto.RequestPasswordResetCommand {
	return dto.RequestPasswordResetCommand{
		Email: in.GetEmail(),
	}
}

func ProtoToConfirmPasswordResetCommand(in *identityv1.ConfirmPasswordResetRequest) dto.ConfirmPasswordResetCommand {
	return dto.ConfirmPasswordResetCommand{
		IDReset:     in.GetIdReset(),
		ResetToken:  in.GetResetToken(),
		NewPassword: in.GetNewPassword(),
	}
}
