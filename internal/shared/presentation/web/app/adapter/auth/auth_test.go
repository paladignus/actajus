package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	identityread "github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	identityusecase "github.com/paladignus/actajus/internal/module/identity/application/usecase"
	sharedrepo "github.com/paladignus/actajus/internal/shared/application/repository"
)

type permissionCheckerFunc func(context.Context, int64, string) (bool, error)

func (f permissionCheckerFunc) HasPermission(ctx context.Context, idUser int64, perm string) (bool, error) {
	return f(ctx, idUser, perm)
}

func TestBuildStateWithoutPermissionCheckerDefaultsToNoAdminPermissions(t *testing.T) {
	adapter := New(identityusecase.ValidateAccess{}, identityusecase.Refresh{}, nil, false, nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/sessions", nil)

	state, err := adapter.BuildState(req, &identityread.AccessTokenClaims{IDUser: 7, IDSession: 3}, "session:manage", "rbac:manage")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state == nil || state.Claims == nil {
		t.Fatalf("expected state")
	}
	if state.CanManageSessions || state.CanManageRBAC {
		t.Fatalf("expected no permissions when checker is nil")
	}
}

func TestBuildStateUsesPermissionChecker(t *testing.T) {
	adapter := New(
		identityusecase.ValidateAccess{},
		identityusecase.Refresh{},
		permissionCheckerFunc(func(_ context.Context, _ int64, perm string) (bool, error) {
			return perm == "session:manage", nil
		}),
		false,
		nil,
		nil,
	)
	req := httptest.NewRequest(http.MethodGet, "/sessions", nil)

	state, err := adapter.BuildState(req, &identityread.AccessTokenClaims{IDUser: 7, IDSession: 3}, "session:manage", "rbac:manage")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !state.CanManageSessions || state.CanManageRBAC {
		t.Fatalf("unexpected permissions %#v", state)
	}
}

func TestSetCookiesUsesStrictSameSiteWhenSecure(t *testing.T) {
	adapter := New(identityusecase.ValidateAccess{}, identityusecase.Refresh{}, nil, true, nil, nil)
	rec := httptest.NewRecorder()

	adapter.SetCookies(rec, &identityread.AuthTokensReadModel{
		IDSession:        9,
		IDUser:           7,
		AccessToken:      "access",
		RefreshToken:     "refresh",
		AccessExpiresAt:  time.Unix(1000, 0),
		RefreshExpiresAt: time.Unix(2000, 0),
	}, "access_token", "refresh_token", "session_id")

	cookies := rec.Result().Cookies()
	if len(cookies) != 3 {
		t.Fatalf("expected 3 cookies, got %d", len(cookies))
	}
	for _, cookie := range cookies {
		if !cookie.Secure {
			t.Fatalf("expected secure cookie %q", cookie.Name)
		}
		if cookie.SameSite != http.SameSiteStrictMode {
			t.Fatalf("expected strict same-site for %q, got %v", cookie.Name, cookie.SameSite)
		}
	}
}

func TestBootstrapIncludesAuthorizationContext(t *testing.T) {
	adapter := New(identityusecase.ValidateAccess{}, identityusecase.Refresh{}, sharedrepo.PermissionChecker(nil), false, func(*http.Request) string {
		return "127.0.0.1"
	}, func(*http.Request) string {
		return "http://localhost:8080"
	})
	req := httptest.NewRequest(http.MethodGet, "/bootstrap", nil)

	payload := adapter.Bootstrap(req, &State{
		Claims:            &identityread.AccessTokenClaims{IDUser: 11, IDSession: 4},
		CanManageSessions: true,
		CanManageRBAC:     true,
	})

	if payload["authenticated"] != true || payload["idUser"] != int64(11) || payload["idSession"] != int64(4) {
		t.Fatalf("unexpected payload %#v", payload)
	}
	if payload["canManageSessions"] != true || payload["canManageRBAC"] != true {
		t.Fatalf("expected authorization flags in payload %#v", payload)
	}
}
