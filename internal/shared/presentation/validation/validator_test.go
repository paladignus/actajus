package validation

import "testing"

func TestValidatorUsesCleanJSONPath(t *testing.T) {
	type input struct {
		Email string `json:"email,omitempty" validate:"required|email"`
	}

	violations := New().ValidateStruct(input{})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Path != "email" {
		t.Fatalf("expected path email, got %q", violations[0].Path)
	}
}

func TestValidatorEmailRuleKeepsFieldPath(t *testing.T) {
	type input struct {
		Contact string `json:"contact_email" validate:"email"`
	}

	violations := New().ValidateStruct(input{Contact: "invalid"})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Path != "contact_email" {
		t.Fatalf("expected path contact_email, got %q", violations[0].Path)
	}
	if violations[0].Code != CodeEmail {
		t.Fatalf("expected code %q, got %q", CodeEmail, violations[0].Code)
	}
}

func TestValidatorIgnoresNonStructInput(t *testing.T) {
	violations := New().ValidateStruct("not-a-struct")
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}
