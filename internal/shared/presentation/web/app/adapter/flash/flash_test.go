package flash

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAdapterRedirectAndRead(t *testing.T) {
	adapter := New(false)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/from", nil)

	adapter.Redirect(rec, req, "/to", "actajus_flash", Message{
		Kind:    "success",
		Message: "Operacao concluida.",
		Source:  "test",
		Code:    "ok",
	})

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, rec.Code)
	}
	var flashCookie *http.Cookie
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == "actajus_flash" {
			flashCookie = cookie
			break
		}
	}
	if flashCookie == nil {
		t.Fatalf("expected flash cookie")
	}

	readReq := httptest.NewRequest(http.MethodGet, "/to", nil)
	readReq.AddCookie(flashCookie)
	readRec := httptest.NewRecorder()
	message := adapter.Read(readRec, readReq, "actajus_flash")

	if message.Message != "Operacao concluida." || message.Kind != "success" {
		t.Fatalf("unexpected flash message %#v", message)
	}
}

func TestAdapterReadNormalizesNonSuccessToError(t *testing.T) {
	adapter := New(false)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/from", nil)

	adapter.Redirect(rec, req, "/to", "actajus_flash", Message{
		Kind:    "info",
		Message: "Aviso",
	})

	var flashCookie *http.Cookie
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == "actajus_flash" {
			flashCookie = cookie
			break
		}
	}
	readReq := httptest.NewRequest(http.MethodGet, "/to", nil)
	readReq.AddCookie(flashCookie)
	readRec := httptest.NewRecorder()
	message := adapter.Read(readRec, readReq, "actajus_flash")

	if message.Kind != "error" {
		t.Fatalf("expected kind error, got %q", message.Kind)
	}
}
