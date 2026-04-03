package csrf

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"strings"

	cookiepolicy "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/cookiepolicy"
)

const (
	CookieName = "actajus_csrf"
	FieldName  = "_csrf"
)

type Adapter struct {
	secureCookies bool
}

func New(secureCookies bool) *Adapter {
	return &Adapter{secureCookies: secureCookies}
}

func (a *Adapter) EnsureCookie(w http.ResponseWriter, r *http.Request) string {
	token := a.Token(r)
	if token != "" {
		return token
	}
	token = newToken()
	policy := cookiepolicy.New(a.secureCookies)
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Path:     "/",
		SameSite: policy.SameSite(),
		Secure:   policy.Secure(),
		MaxAge:   60 * 60 * 12,
	})
	return token
}

func (a *Adapter) Token(r *http.Request) string {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(cookie.Value)
}

func (a *Adapter) Validate(r *http.Request) bool {
	cookieToken := a.Token(r)
	if cookieToken == "" {
		return false
	}
	requestToken := strings.TrimSpace(r.Header.Get("X-CSRF-Token"))
	if requestToken == "" {
		requestToken = strings.TrimSpace(r.PostFormValue(FieldName))
	}
	if requestToken == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(cookieToken), []byte(requestToken)) == 1
}

func newToken() string {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(raw)
}
