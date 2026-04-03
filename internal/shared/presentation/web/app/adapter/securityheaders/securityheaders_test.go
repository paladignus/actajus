package securityheaders

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestApplyDocumentSetsHeaders(t *testing.T) {
	adapter := New("")
	rec := httptest.NewRecorder()

	adapter.ApplyDocument(rec)

	if rec.Header().Get("X-Frame-Options") != "DENY" {
		t.Fatalf("expected frame options header")
	}
	csp := rec.Header().Get("Content-Security-Policy")
	if !strings.Contains(csp, "default-src 'self'") {
		t.Fatalf("expected base csp, got %q", csp)
	}
}

func TestApplyDocumentIncludesDevServerWhenConfigured(t *testing.T) {
	adapter := New("http://localhost:5173")
	rec := httptest.NewRecorder()

	adapter.ApplyDocument(rec)

	csp := rec.Header().Get("Content-Security-Policy")
	if !strings.Contains(csp, "http://localhost:5173") {
		t.Fatalf("expected dev server in csp, got %q", csp)
	}
}

func TestApplyAssetSetsAssetHeaders(t *testing.T) {
	adapter := New("")
	rec := httptest.NewRecorder()

	adapter.ApplyAsset(rec)

	if rec.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("expected nosniff header")
	}
	if rec.Header().Get("Cross-Origin-Resource-Policy") != "same-origin" {
		t.Fatalf("expected cross-origin-resource-policy header")
	}
}
