package company

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func RewritePageLink(r *http.Request, raw *string, page int) string {
	if raw == nil {
		return ""
	}
	return RewritePageLinkURL(*raw, page, r.URL.Query().Get("name"), r.URL.Query().Get("cnpj"))
}

func RewritePageLinkURL(raw string, page int, filters ...string) string {
	if strings.TrimSpace(raw) == "" {
		return ""
	}
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	q := u.Query()
	q.Set("page", strconv.Itoa(page))
	if len(filters) > 0 {
		q.Set("name", strings.TrimSpace(filters[0]))
	}
	if len(filters) > 1 {
		q.Set("cnpj", strings.TrimSpace(filters[1]))
	}
	u.RawQuery = q.Encode()
	return u.RequestURI()
}
