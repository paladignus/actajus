package navigation

import (
	"net/http"
	"net/http/httptest"
	"testing"

	webflash "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/flash"
)

func TestNavigatorSeeOther(t *testing.T) {
	navigator := New(false)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/from", nil)

	navigator.SeeOther(rec, req, "/target")

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, rec.Code)
	}
	if rec.Header().Get("Location") != "/target" {
		t.Fatalf("unexpected location %q", rec.Header().Get("Location"))
	}
}

func TestNavigatorSeeOtherWithFlash(t *testing.T) {
	navigator := New(false)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/from", nil)

	navigator.SeeOtherWithFlash(rec, req, "/target", "actajus_flash", webflash.Message{
		Kind:    "success",
		Message: "Tudo certo.",
	})

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, rec.Code)
	}
	found := false
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == "actajus_flash" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected flash cookie")
	}
}
