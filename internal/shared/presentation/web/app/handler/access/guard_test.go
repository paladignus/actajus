package access

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	identityread "github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	identityusecase "github.com/paladignus/actajus/internal/module/identity/application/usecase"
	identitydomain "github.com/paladignus/actajus/internal/module/identity/domain"
	webauth "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/auth"
	navigation "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/http/navigation"
)

type stubAccessTokenService struct {
	claims identityread.AccessTokenClaims
	err    error
}

func (s stubAccessTokenService) Sign(claims identityread.AccessTokenClaims) (string, error) {
	return "", nil
}

func (s stubAccessTokenService) Verify(token string) (identityread.AccessTokenClaims, error) {
	if s.err != nil {
		return identityread.AccessTokenClaims{}, s.err
	}
	return s.claims, nil
}

type stubSessionRepository struct{}

func (stubSessionRepository) Create(context.Context, *identitydomain.Session) error { return nil }
func (stubSessionRepository) GetByID(context.Context, int64) (*identitydomain.Session, error) {
	return nil, nil
}
func (stubSessionRepository) RotateRefreshToken(context.Context, int64, [32]byte, time.Time) error {
	return nil
}
func (stubSessionRepository) Revoke(context.Context, int64) error          { return nil }
func (stubSessionRepository) RevokeAllByUser(context.Context, int64) error { return nil }
func (stubSessionRepository) CountActiveByUser(context.Context, int64) (int, error) {
	return 0, nil
}
func (stubSessionRepository) IsActive(context.Context, int64, int64, time.Time) (bool, error) {
	return true, nil
}
func (stubSessionRepository) RotateRefreshTokenAtomic(context.Context, int64, [32]byte, [32]byte, time.Time, time.Time) (bool, error) {
	return true, nil
}

type stubClock struct{}

func (stubClock) Now() time.Time { return time.Unix(0, 0) }

type permissionCheckerFunc func(context.Context, int64, string) (bool, error)

func (f permissionCheckerFunc) HasPermission(ctx context.Context, idUser int64, perm string) (bool, error) {
	return f(ctx, idUser, perm)
}

func newTestGuard(t *testing.T, checker permissionCheckerFunc) Guard {
	t.Helper()
	validateAccess := identityusecase.NewValidateAccess(
		stubAccessTokenService{claims: identityread.AccessTokenClaims{IDUser: 42, IDSession: 9}},
		stubSessionRepository{},
		stubClock{},
		false,
	)
	adapter := webauth.New(
		validateAccess,
		identityusecase.Refresh{},
		checker,
		false,
		func(*http.Request) string { return "127.0.0.1" },
		func(*http.Request) string { return "http://localhost:8080" },
	)
	return Guard{
		Adapter:            adapter,
		Navigator:          navigation.New(false),
		CookieAccessToken:  "access",
		CookieRefreshToken: "refresh",
		CookieSessionID:    "session",
		ManageSessionsPerm: "session:manage",
		ManageRBACPerm:     "rbac:manage",
	}
}

func TestGuardRequireAuthRedirectsToLoginWhenMissingCookies(t *testing.T) {
	guard := newTestGuard(t, permissionCheckerFunc(func(context.Context, int64, string) (bool, error) {
		return false, nil
	}))

	req := httptest.NewRequest(http.MethodGet, "/sessions", nil)
	rec := httptest.NewRecorder()

	state, err := guard.RequireAuth(rec, req)
	if state != nil {
		t.Fatalf("expected nil state")
	}
	if err == nil {
		t.Fatalf("expected auth error")
	}
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect status, got %d", rec.Code)
	}
	if location := rec.Header().Get("Location"); location != "/login" {
		t.Fatalf("expected redirect to /login, got %q", location)
	}
}

func TestGuardAuthFromRequestReturnsStateFromAccessCookie(t *testing.T) {
	guard := newTestGuard(t, permissionCheckerFunc(func(_ context.Context, _ int64, perm string) (bool, error) {
		return perm == "session:manage", nil
	}))

	req := httptest.NewRequest(http.MethodGet, "/sessions", nil)
	req.AddCookie(&http.Cookie{Name: "access", Value: "valid-token"})
	rec := httptest.NewRecorder()

	state, err := guard.AuthFromRequest(rec, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state == nil || state.Claims == nil {
		t.Fatalf("expected authenticated state")
	}
	if state.Claims.IDUser != 42 || state.Claims.IDSession != 9 {
		t.Fatalf("unexpected claims %#v", state.Claims)
	}
	if !state.CanManageSessions {
		t.Fatalf("expected session permission")
	}
	if state.CanManageRBAC {
		t.Fatalf("did not expect rbac permission")
	}
}

func TestGuardRequireRBACAdminCallsForbiddenResponder(t *testing.T) {
	called := false
	guard := newTestGuard(t, permissionCheckerFunc(func(_ context.Context, _ int64, perm string) (bool, error) {
		return perm == "session:manage", nil
	}))
	guard.RespondForbidden = func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusForbidden)
	}

	req := httptest.NewRequest(http.MethodGet, "/roles", nil)
	req.AddCookie(&http.Cookie{Name: "access", Value: "valid-token"})
	rec := httptest.NewRecorder()

	state, err := guard.RequireRBACAdmin(rec, req)
	if state != nil {
		t.Fatalf("expected nil state")
	}
	if err == nil || err.Error() != "missing rbac manage permission" {
		t.Fatalf("expected missing permission error, got %v", err)
	}
	if !called {
		t.Fatalf("expected forbidden responder to be called")
	}
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}
}

func TestGuardSetAndClearAuthCookies(t *testing.T) {
	guard := newTestGuard(t, permissionCheckerFunc(func(context.Context, int64, string) (bool, error) {
		return false, nil
	}))
	rec := httptest.NewRecorder()

	tokens := &identityread.AuthTokensReadModel{
		IDSession:        33,
		IDUser:           7,
		AccessToken:      "a",
		RefreshToken:     "r",
		AccessExpiresAt:  time.Unix(1000, 0),
		RefreshExpiresAt: time.Unix(2000, 0),
	}
	guard.SetAuthCookies(rec, tokens)
	guard.ClearAuthCookies(rec)

	cookies := rec.Result().Cookies()
	if len(cookies) != 6 {
		t.Fatalf("expected 6 cookies written, got %d", len(cookies))
	}
}
