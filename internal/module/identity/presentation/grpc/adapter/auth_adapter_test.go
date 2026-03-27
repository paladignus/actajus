package adapter_test

import (
	"testing"
	"time"

	"github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	"github.com/paladignus/actajus/internal/module/identity/presentation/grpc/adapter"
	identityv1 "github.com/paladignus/actajus/proto/identity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestProtoToLoginCommand(t *testing.T) {
	t.Parallel()

	req := &identityv1.LoginRequest{
		Email:     "user@example.com",
		Password:  "secret",
		Ip:        "127.0.0.1",
		UserAgent: "test-agent",
	}

	got := adapter.ProtoToLoginCommand(req)

	if got.Email != req.Email || got.Password != req.Password || got.IP != req.Ip || got.UserAgent != req.UserAgent {
		t.Fatalf("unexpected login mapping: %+v", got)
	}
}

func TestProtoToRefreshCommand(t *testing.T) {
	t.Parallel()

	req := &identityv1.RefreshRequest{
		IdSession:    99,
		RefreshToken: "refresh-token",
		Ip:           "10.0.0.1",
		UserAgent:    "mobile",
	}

	got := adapter.ProtoToRefreshCommand(req)

	if got.IDSession != 99 || got.RefreshToken != req.RefreshToken || got.IP != req.Ip || got.UserAgent != req.UserAgent {
		t.Fatalf("unexpected refresh mapping: %+v", got)
	}
}

func TestProtoToPasswordCommands(t *testing.T) {
	t.Parallel()

	changeReq := &identityv1.ChangePasswordRequest{
		IdUser:          7,
		CurrentPassword: "old",
		NewPassword:     "new",
	}
	changeGot := adapter.ProtoToChangePasswordCommand(changeReq)
	if changeGot.IDUser != 7 || changeGot.CurrentPassword != "old" || changeGot.NewPassword != "new" {
		t.Fatalf("unexpected change password mapping: %+v", changeGot)
	}

	resetReq := &identityv1.RequestPasswordResetRequest{Email: "user@example.com"}
	resetGot := adapter.ProtoToRequestPasswordResetCommand(resetReq)
	if resetGot.Email != "user@example.com" {
		t.Fatalf("unexpected request reset mapping: %+v", resetGot)
	}

	confirmReq := &identityv1.ConfirmPasswordResetRequest{
		IdReset:     11,
		ResetToken:  "token",
		NewPassword: "new-secret",
	}
	confirmGot := adapter.ProtoToConfirmPasswordResetCommand(confirmReq)
	if confirmGot.IDReset != 11 || confirmGot.ResetToken != "token" || confirmGot.NewPassword != "new-secret" {
		t.Fatalf("unexpected confirm reset mapping: %+v", confirmGot)
	}
}

func TestAuthReadModelsToProto(t *testing.T) {
	t.Parallel()

	accessExp := time.Date(2026, 3, 27, 12, 0, 0, 0, time.UTC)
	refreshExp := time.Date(2026, 3, 28, 12, 0, 0, 0, time.UTC)
	rm := &readmodel.AuthTokensReadModel{
		IDSession:        123,
		AccessToken:      "access",
		RefreshToken:     "refresh",
		AccessExpiresAt:  accessExp,
		RefreshExpiresAt: refreshExp,
	}

	loginResp := adapter.LoginReadModelToProto(rm)
	assertAuthTokensResponse(t, loginResp.AccessToken, loginResp.RefreshToken, loginResp.IdSession, loginResp.AccessExpiresAt, loginResp.RefreshExpiresAt, accessExp, refreshExp)

	refreshResp := adapter.RefreshReadModelToProto(rm)
	assertAuthTokensResponse(t, refreshResp.AccessToken, refreshResp.RefreshToken, refreshResp.IdSession, refreshResp.AccessExpiresAt, refreshResp.RefreshExpiresAt, accessExp, refreshExp)
}

func TestRequestPasswordResetReadModelToProto(t *testing.T) {
	t.Parallel()

	rm := &readmodel.RequestPasswordResetReadModel{Message: "If the account exists, an email will be sent."}
	got := adapter.RequestPasswordResetReadModelToProto(rm)

	if !got.Ok || got.Message != rm.Message {
		t.Fatalf("unexpected request reset response mapping: %+v", got)
	}
}

func assertAuthTokensResponse(
	t *testing.T,
	accessToken string,
	refreshToken string,
	idSession int64,
	accessExpiresAt *timestamppb.Timestamp,
	refreshExpiresAt *timestamppb.Timestamp,
	wantAccess time.Time,
	wantRefresh time.Time,
) {
	t.Helper()

	if accessToken != "access" || refreshToken != "refresh" || idSession != 123 {
		t.Fatalf("unexpected token payload: access=%s refresh=%s idSession=%d", accessToken, refreshToken, idSession)
	}
	if accessExpiresAt == nil || refreshExpiresAt == nil {
		t.Fatal("expected protobuf timestamps to be set")
	}
	if !accessExpiresAt.AsTime().Equal(wantAccess) || !refreshExpiresAt.AsTime().Equal(wantRefresh) {
		t.Fatalf("unexpected timestamp mapping: access=%v refresh=%v", accessExpiresAt.AsTime(), refreshExpiresAt.AsTime())
	}
}
