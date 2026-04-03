package cookiepolicy

import (
	"net/http"
	"testing"
)

func TestSameSiteUsesLaxWithoutSecureCookies(t *testing.T) {
	policy := New(false)
	if policy.SameSite() != http.SameSiteLaxMode {
		t.Fatalf("expected lax same-site mode")
	}
}

func TestSameSiteUsesStrictWithSecureCookies(t *testing.T) {
	policy := New(true)
	if policy.SameSite() != http.SameSiteStrictMode {
		t.Fatalf("expected strict same-site mode")
	}
}
