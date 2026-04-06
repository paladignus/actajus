package pathbinder

import (
	"net/http/httptest"
	"testing"
)

func TestBindInt64Path(t *testing.T) {
	req := httptest.NewRequest("GET", "/companies/42", nil)
	req.SetPathValue("id", "42")

	got, err := BindInt64Path(req, "id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 42 {
		t.Fatalf("expected 42, got %d", got)
	}
}

func TestParseInt16ValueRejectsZero(t *testing.T) {
	_, err := ParseInt16Value("0", "invalid")
	if err == nil {
		t.Fatalf("expected error")
	}
}
