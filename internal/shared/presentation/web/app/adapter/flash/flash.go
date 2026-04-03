package flash

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
)

type Message struct {
	Code    string
	Kind    string
	Message string
	Source  string
}

type Adapter struct {
	secureCookies bool
}

func New(secureCookies bool) *Adapter {
	return &Adapter{secureCookies: secureCookies}
}

func (a *Adapter) Read(w http.ResponseWriter, r *http.Request, cookieName string) Message {
	cookie, err := r.Cookie(cookieName)
	if err != nil || strings.TrimSpace(cookie.Value) == "" {
		return Message{}
	}
	a.Clear(w, cookieName)
	raw, err := base64.RawURLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return Message{}
	}
	var flash Message
	if err := json.Unmarshal(raw, &flash); err != nil {
		return Message{}
	}
	if strings.TrimSpace(flash.Message) == "" {
		return Message{}
	}
	if flash.Kind != "success" {
		flash.Kind = "error"
	}
	return flash
}

func (a *Adapter) Redirect(w http.ResponseWriter, r *http.Request, path, cookieName string, flash Message) {
	if strings.TrimSpace(flash.Message) == "" {
		http.Redirect(w, r, path, http.StatusSeeOther)
		return
	}
	if flash.Kind != "success" {
		flash.Kind = "error"
	}
	raw, err := json.Marshal(flash)
	if err != nil {
		http.Redirect(w, r, path, http.StatusSeeOther)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    base64.RawURLEncoding.EncodeToString(raw),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   a.secureCookies,
		MaxAge:   10,
	})
	http.Redirect(w, r, path, http.StatusSeeOther)
}

func (a *Adapter) Clear(w http.ResponseWriter, cookieName string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   a.secureCookies,
		MaxAge:   -1,
	})
}
