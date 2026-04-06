package authbinder

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBindLoginFormTrimsEmail(t *testing.T) {
	req := httptest.NewRequest("POST", "/login", strings.NewReader("email=%20user%40example.com%20&password=secret"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	form, err := BindLoginForm(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if form.Email != "user@example.com" {
		t.Fatalf("unexpected email %q", form.Email)
	}
	if form.Password != "secret" {
		t.Fatalf("unexpected password %q", form.Password)
	}
}
