package authhandler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	identitymapper "github.com/paladignus/actajus/internal/module/identity/application/mapper"
	identityread "github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	identityrepo "github.com/paladignus/actajus/internal/module/identity/application/repository"
	identityusecase "github.com/paladignus/actajus/internal/module/identity/application/usecase"
	identitydomain "github.com/paladignus/actajus/internal/module/identity/domain"
	"github.com/paladignus/actajus/internal/shared/presentation/validation"
	webauth "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/auth"
	handlerctx "github.com/paladignus/actajus/internal/shared/presentation/web/app/handlerctx"
	pagepresenter "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/page"
	webhtml "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/html"
	navigation "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/http/navigation"
	webtemplate "github.com/paladignus/actajus/internal/shared/presentation/web/template"
)

type stubSessionRepository struct {
	revokeCalled bool
	lastID       int64
}

func (s *stubSessionRepository) Create(context.Context, *identitydomain.Session) error { return nil }
func (s *stubSessionRepository) GetByID(context.Context, int64) (*identitydomain.Session, error) {
	return nil, nil
}
func (s *stubSessionRepository) RotateRefreshToken(context.Context, int64, [32]byte, time.Time) error {
	return nil
}
func (s *stubSessionRepository) Revoke(_ context.Context, id int64) error {
	s.revokeCalled = true
	s.lastID = id
	return nil
}
func (s *stubSessionRepository) RevokeAllByUser(context.Context, int64) error { return nil }
func (s *stubSessionRepository) CountActiveByUser(context.Context, int64) (int, error) {
	return 0, nil
}
func (s *stubSessionRepository) IsActive(context.Context, int64, int64, time.Time) (bool, error) {
	return true, nil
}
func (s *stubSessionRepository) RotateRefreshTokenAtomic(context.Context, int64, [32]byte, [32]byte, time.Time, time.Time) (bool, error) {
	return true, nil
}

func newTestHTMLResponder(t *testing.T) *webhtml.Responder {
	t.Helper()
	renderer, err := webtemplate.NewRenderer(webtemplate.Source{
		TemplateFS: fstest.MapFS{
			"templates/pages/login.gohtml": {
				Data: []byte(`{{define "pages/login"}}<html><body><h1>{{.Page.Title}}</h1><p>{{index .Data "Email"}}</p><p>{{index .Data "Error"}}</p></body></html>{{end}}`),
			},
		},
		TemplatePatterns: []string{"templates/pages/*.gohtml"},
		DevServerURL:     "http://localhost:5173",
	})
	if err != nil {
		t.Fatalf("new renderer: %v", err)
	}
	return webhtml.New(pagepresenter.New(renderer), nil)
}

func TestHomeRedirectsGuestToLogin(t *testing.T) {
	controller := Controller{
		Context: handlerctx.Context{
			Navigator: navigation.New(false),
			AuthFromRequest: func(http.ResponseWriter, *http.Request) (*handlerctx.State, error) {
				return nil, errors.New("missing auth")
			},
		},
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	controller.Home().ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, rec.Code)
	}
	if location := rec.Header().Get("Location"); location != "/login" {
		t.Fatalf("expected /login redirect, got %q", location)
	}
}

func TestHomeRedirectsAuthenticatedUserToSessions(t *testing.T) {
	controller := Controller{
		Context: handlerctx.Context{
			Navigator: navigation.New(false),
			AuthFromRequest: func(http.ResponseWriter, *http.Request) (*handlerctx.State, error) {
				return &handlerctx.State{Claims: &identityread.AccessTokenClaims{IDUser: 1, IDSession: 2}}, nil
			},
		},
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	controller.Home().ServeHTTP(rec, req)

	if location := rec.Header().Get("Location"); location != "/sessions" {
		t.Fatalf("expected /sessions redirect, got %q", location)
	}
}

func TestLoginPageRendersForGuest(t *testing.T) {
	controller := Controller{
		Context: handlerctx.Context{
			HTML:      newTestHTMLResponder(t),
			Navigator: navigation.New(false),
			AuthFromRequest: func(http.ResponseWriter, *http.Request) (*handlerctx.State, error) {
				return nil, errors.New("missing auth")
			},
		},
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	controller.LoginPage().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Entrar") {
		t.Fatalf("expected rendered login page, got %q", body)
	}
}

func TestLoginPageRedirectsAuthenticatedUser(t *testing.T) {
	controller := Controller{
		Context: handlerctx.Context{
			HTML:      newTestHTMLResponder(t),
			Navigator: navigation.New(false),
			AuthFromRequest: func(http.ResponseWriter, *http.Request) (*handlerctx.State, error) {
				return &handlerctx.State{Claims: &identityread.AccessTokenClaims{IDUser: 1, IDSession: 2}}, nil
			},
		},
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	controller.LoginPage().ServeHTTP(rec, req)

	if rec.Header().Get("Location") != "/sessions" {
		t.Fatalf("expected /sessions redirect, got %q", rec.Header().Get("Location"))
	}
}

func TestBootstrapReturnsJSON(t *testing.T) {
	adapter := webauth.New(
		identityusecase.ValidateAccess{},
		identityusecase.Refresh{},
		nil,
		false,
		func(*http.Request) string { return "127.0.0.1" },
		func(*http.Request) string { return "http://localhost:8080" },
	)
	controller := Controller{
		Context: handlerctx.Context{
			Navigator:   navigation.New(false),
			AuthAdapter: adapter,
			AuthFromRequest: func(http.ResponseWriter, *http.Request) (*handlerctx.State, error) {
				return &handlerctx.State{Claims: &identityread.AccessTokenClaims{IDUser: 9, IDSession: 3}}, nil
			},
		},
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/bootstrap", nil)
	controller.Bootstrap().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"authenticated":true`) {
		t.Fatalf("expected bootstrap payload, got %q", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"idUser":9`) {
		t.Fatalf("expected user payload, got %q", rec.Body.String())
	}
}

func TestLoginActionReturnsBadRequestOnMalformedForm(t *testing.T) {
	controller := Controller{}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("%zz"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	controller.LoginAction().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestLoginActionReturnsTooManyRequestsWhenRateLimited(t *testing.T) {
	auditCalled := false
	controller := Controller{
		Context: handlerctx.Context{
			HTML: newTestHTMLResponder(t),
			CheckLoginAttempt: func(*http.Request, string) (bool, time.Duration) {
				return false, 30 * time.Second
			},
			AuditFailure: func(_ *http.Request, action string, _ ...any) {
				auditCalled = action == "auth.login.rate_limited"
			},
		},
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("email=test%40mail.com&password=123"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	controller.LoginAction().ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected status %d, got %d", http.StatusTooManyRequests, rec.Code)
	}
	if retryAfter := rec.Header().Get("Retry-After"); retryAfter == "" {
		t.Fatalf("expected retry-after header")
	}
	if !strings.Contains(rec.Body.String(), "Muitas tentativas") {
		t.Fatalf("expected rate limit message, got %q", rec.Body.String())
	}
	if !auditCalled {
		t.Fatalf("expected audit callback for rate limit")
	}
}

func TestLogoutActionClearsCookiesRedirectsAndAudits(t *testing.T) {
	repo := &stubSessionRepository{}
	var _ identityrepo.SessionRepository = repo

	clearCalled := false
	auditCalled := false
	controller := Controller{
		Context: handlerctx.Context{
			Navigator: navigation.New(false),
			RequireAuth: func(http.ResponseWriter, *http.Request) (*handlerctx.State, error) {
				return &handlerctx.State{Claims: &identityread.AccessTokenClaims{IDUser: 7, IDSession: 11}}, nil
			},
			ClearAuthCookies: func(http.ResponseWriter) {
				clearCalled = true
			},
			AuditSuccess: func(_ *http.Request, action string, attrs ...any) {
				if action == "auth.logout.succeeded" {
					auditCalled = true
				}
			},
		},
		Logout: identityusecase.NewLogout(repo, *identitymapper.NewAuthMapper(validation.New())),
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	controller.LogoutAction().ServeHTTP(rec, req)

	if !repo.revokeCalled || repo.lastID != 11 {
		t.Fatalf("expected session revoke for id 11, got called=%v id=%d", repo.revokeCalled, repo.lastID)
	}
	if !clearCalled {
		t.Fatalf("expected auth cookies to be cleared")
	}
	if !auditCalled {
		t.Fatalf("expected logout audit callback")
	}
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, rec.Code)
	}
	if location := rec.Header().Get("Location"); location != "/login" {
		t.Fatalf("expected redirect to /login, got %q", location)
	}
}
