package companybinder

import (
	"net/http/httptest"
	"testing"
)

func TestBindListQuery(t *testing.T) {
	req := httptest.NewRequest("GET", "/companies?page=3&limit=20&name=Acme&cnpj=123&after=abc", nil)

	query := BindListQuery(req)
	if query.Page != 3 || query.Limit != 20 {
		t.Fatalf("unexpected paging %#v", query)
	}
	if query.Filter.Name != "Acme" || query.Filter.CNPJ != "123" {
		t.Fatalf("unexpected filter %#v", query.Filter)
	}
	if query.After == nil || *query.After != "abc" {
		t.Fatalf("unexpected after %#v", query.After)
	}
}
