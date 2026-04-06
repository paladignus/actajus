package identitybinder

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBindRevokeAllSessionsFormWithEmptyUser(t *testing.T) {
	req := httptest.NewRequest("POST", "/sessions/revoke-all", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	form, cmd, err := BindRevokeAllSessionsForm(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if form.IDUser != nil {
		t.Fatalf("expected nil user id")
	}
	if cmd != nil {
		t.Fatalf("expected nil command")
	}
}

func TestBindAssignUserRoleFormParsesIDs(t *testing.T) {
	req := httptest.NewRequest("POST", "/users/roles", strings.NewReader("id_user=7&id_role=3"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	form, err := BindAssignUserRoleForm(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if form.IDUser != 7 || form.IDRole != 3 {
		t.Fatalf("unexpected form %#v", form)
	}
}
