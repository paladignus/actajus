package cookiepolicy

import "net/http"

type Policy struct {
	secure bool
}

func New(secure bool) Policy {
	return Policy{secure: secure}
}

func (p Policy) SameSite() http.SameSite {
	if p.secure {
		return http.SameSiteStrictMode
	}
	return http.SameSiteLaxMode
}

func (p Policy) Secure() bool {
	return p.secure
}
