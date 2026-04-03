package securityheaders

import (
	"net/http"
	"strings"
)

type Adapter struct {
	devServerURL string
}

func New(devServerURL string) *Adapter {
	return &Adapter{devServerURL: strings.TrimSpace(devServerURL)}
}

func (a *Adapter) ApplyDocument(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
	w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
	w.Header().Set("Content-Security-Policy", a.documentCSP())
}

func (a *Adapter) ApplyAsset(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
}

func (a *Adapter) documentCSP() string {
	connectSrc := "'self'"
	scriptSrc := "'self' 'unsafe-inline'"
	styleSrc := "'self' 'unsafe-inline'"

	if a.devServerURL != "" {
		connectSrc += " " + a.devServerURL + " ws://localhost:5173"
		scriptSrc += " " + a.devServerURL
		styleSrc += " " + a.devServerURL
	}

	return strings.Join([]string{
		"default-src 'self'",
		"base-uri 'self'",
		"frame-ancestors 'none'",
		"form-action 'self'",
		"img-src 'self' data:",
		"font-src 'self' data:",
		"script-src " + scriptSrc,
		"style-src " + styleSrc,
		"connect-src " + connectSrc,
		"object-src 'none'",
	}, "; ")
}
