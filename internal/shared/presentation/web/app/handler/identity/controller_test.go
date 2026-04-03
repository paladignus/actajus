package identityhandler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	identityread "github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	handlerctx "github.com/paladignus/actajus/internal/shared/presentation/web/app/handlerctx"
	navigation "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/http/navigation"
)

func TestSessionsPageRedirectsWhenUnauthenticated(t *testing.T) {
	controller := Controller{
		Context: handlerctx.Context{
			Navigator: navigation.New(false),
			RequireAuth: func(http.ResponseWriter, *http.Request) (*handlerctx.State, error) {
				return nil, errors.New("missing auth")
			},
		},
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/sessions", nil)
	controller.Context.RequireAuth = func(w http.ResponseWriter, r *http.Request) (*handlerctx.State, error) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return nil, errors.New("missing auth")
	}
	controller.SessionsPage().ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, rec.Code)
	}
	if location := rec.Header().Get("Location"); location != "/login" {
		t.Fatalf("expected redirect to /login, got %q", location)
	}
}

func TestRevokeAllActionRedirectsWithFlashOnInvalidUserID(t *testing.T) {
	controller := Controller{
		Context: handlerctx.Context{
			Navigator:   navigation.New(false),
			FlashCookie: "actajus_flash",
			RequireAuth: func(http.ResponseWriter, *http.Request) (*handlerctx.State, error) {
				return &handlerctx.State{
					Claims: &identityread.AccessTokenClaims{IDUser: 7, IDSession: 9},
				}, nil
			},
		},
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/sessions/revoke-all", strings.NewReader("id_user=abc"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	controller.RevokeAllAction().ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, rec.Code)
	}
	if location := rec.Header().Get("Location"); location != "/sessions" {
		t.Fatalf("expected redirect to /sessions, got %q", location)
	}
	foundFlash := false
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == "actajus_flash" {
			foundFlash = true
			break
		}
	}
	if !foundFlash {
		t.Fatalf("expected flash cookie to be written")
	}
}

func TestRevokeSessionActionRedirectsWithFlashOnInvalidSessionID(t *testing.T) {
	controller := Controller{
		Context: handlerctx.Context{
			Navigator:   navigation.New(false),
			FlashCookie: "actajus_flash",
			RequireAuth: func(http.ResponseWriter, *http.Request) (*handlerctx.State, error) {
				return &handlerctx.State{
					Claims: &identityread.AccessTokenClaims{IDUser: 7, IDSession: 9},
				}, nil
			},
		},
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/sessions/revoke", strings.NewReader("id_session=abc"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	controller.RevokeSessionAction().ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, rec.Code)
	}
	if location := rec.Header().Get("Location"); location != "/sessions" {
		t.Fatalf("expected redirect to /sessions, got %q", location)
	}
}
