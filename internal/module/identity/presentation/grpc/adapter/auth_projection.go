// Package adapter
package adapter

import (
	identityDTO "github.com/paladignus/actajus/internal/module/identity/application/dto"
	identityv1 "github.com/paladignus/actajus/proto/identity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func LoginReadModelToProto(rm *identityDTO.AuthTokensReadModel) *identityv1.LoginResponse {
	return &identityv1.LoginResponse{
		AccessToken:      rm.AccessToken,
		AccessExpiresAt:  timestamppb.New(rm.AccessExpiresAt),
		RefreshToken:     rm.RefreshToken,
		RefreshExpiresAt: timestamppb.New(rm.RefreshExpiresAt),
		IdSession:        rm.IDSession,
	}
}

func RefreshReadModelToProto(rm *identityDTO.AuthTokensReadModel) *identityv1.RefreshResponse {
	return &identityv1.RefreshResponse{
		AccessToken:      rm.AccessToken,
		AccessExpiresAt:  timestamppb.New(rm.AccessExpiresAt),
		RefreshToken:     rm.RefreshToken,
		RefreshExpiresAt: timestamppb.New(rm.RefreshExpiresAt),
		IdSession:        rm.IDSession,
	}
}

func RequestPasswordResetReadModelToProto(rm *identityDTO.RequestPasswordResetReadModel) *identityv1.RequestPasswordResetResponse {
	resp := &identityv1.RequestPasswordResetResponse{
		Ok:      true,
		Message: rm.Message,
	}
	// if rm == nil {
	// 	return resp
	// }
	// resp.IdReset = rm.IDReset
	// resp.ResetToken = rm.ResetToken
	// if !rm.ExpiresAt.IsZero() {
	// 	resp.ExpiresAt = timestamppb.New(rm.ExpiresAt)
	// }
	return resp
}
