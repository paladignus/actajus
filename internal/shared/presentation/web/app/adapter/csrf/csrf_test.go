package csrf

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestEnsureCookieCreatesToken(t *testing.T) {
	adapter := New(false)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/login", nil)

	token := adapter.EnsureCookie(rec, req)
	if token == "" {
		t.Fatalf("expected token")
	}

	res := rec.Result()
	cookies := res.Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}
	if cookies[0].Name != CookieName {
		t.Fatalf("expected cookie %q, got %q", CookieName, cookies[0].Name)
	}
	if cookies[0].Value != token {
		t.Fatalf("expected cookie token %q, got %q", token, cookies[0].Value)
	}
}

func TestEnsureCookieReusesExistingToken(t *testing.T) {
	adapter := New(false)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	req.AddCookie(&http.Cookie{Name: CookieName, Value: "existing-token"})

	token := adapter.EnsureCookie(rec, req)
	if token != "existing-token" {
		t.Fatalf("expected existing token, got %q", token)
	}
	if len(rec.Result().Cookies()) != 0 {
		t.Fatalf("expected no cookie rewrite")
	}
}

func TestValidateAcceptsFormToken(t *testing.T) {
	adapter := New(false)
	form := url.Values{}
	form.Set(FieldName, "valid-token")

	req := httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: CookieName, Value: "valid-token"})

	if !adapter.Validate(req) {
		t.Fatalf("expected token validation to succeed")
	}
}

func TestValidateAcceptsHeaderToken(t *testing.T) {
	adapter := New(false)
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req.Header.Set("X-CSRF-Token", "valid-token")
	req.AddCookie(&http.Cookie{Name: CookieName, Value: "valid-token"})

	if !adapter.Validate(req) {
		t.Fatalf("expected header token validation to succeed")
	}
}

func TestValidateRejectsMissingOrDifferentToken(t *testing.T) {
	adapter := New(false)

	missing := httptest.NewRequest(http.MethodPost, "/logout", nil)
	if adapter.Validate(missing) {
		t.Fatalf("expected missing token to fail")
	}

	form := url.Values{}
	form.Set(FieldName, "wrong-token")
	mismatch := httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader(form.Encode()))
	mismatch.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	mismatch.AddCookie(&http.Cookie{Name: CookieName, Value: "valid-token"})

	if adapter.Validate(mismatch) {
		t.Fatalf("expected mismatched token to fail")
	}
}
