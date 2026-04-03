package handler

import (
	"errors"
	"net/http"

	identitycmd "github.com/paladignus/actajus/internal/module/identity/application/command"
	identityread "github.com/paladignus/actajus/internal/module/identity/application/readmodel"
)

type accessGuard struct {
	app AppHandler
}

func (h AppHandler) guard() accessGuard {
	return accessGuard{app: h}
}

func (g accessGuard) authFromRequest(w http.ResponseWriter, r *http.Request) (*authState, error) {
	state, refreshed, err := g.app.auth.adapter.FromRequest(w, r, cookieAccessToken, cookieRefreshToken, cookieSessionID, permManageSessions, permManageRBAC)
	if refreshed != nil {
		g.app.auth.adapter.SetCookies(w, refreshed, cookieAccessToken, cookieRefreshToken, cookieSessionID)
	}
	return state, err
}

func (g accessGuard) buildAuthState(r *http.Request, claims *identityread.AccessTokenClaims) (*authState, error) {
	return g.app.auth.adapter.BuildState(r, claims, permManageSessions, permManageRBAC)
}

func (g accessGuard) requireAuth(w http.ResponseWriter, r *http.Request) (*authState, error) {
	state, err := g.authFromRequest(w, r)
	if err != nil || state == nil || state.Claims == nil {
		g.app.navigator.SeeOther(w, r, "/login")
		return nil, err
	}
	return state, nil
}

func (g accessGuard) requireRBACAdmin(w http.ResponseWriter, r *http.Request) (*authState, error) {
	state, err := g.requireAuth(w, r)
	if err != nil {
		return nil, err
	}
	if state == nil || !state.CanManageRBAC {
		g.app.support().respondForbidden(w, r)
		return nil, errors.New("missing rbac manage permission")
	}
	return state, nil
}

func (g accessGuard) setAuthCookies(w http.ResponseWriter, tokens *identityread.AuthTokensReadModel) {
	g.app.auth.adapter.SetCookies(w, tokens, cookieAccessToken, cookieRefreshToken, cookieSessionID)
}

func (g accessGuard) clearAuthCookies(w http.ResponseWriter) {
	g.app.auth.adapter.ClearCookies(w, cookieAccessToken, cookieRefreshToken, cookieSessionID)
}

func (g accessGuard) logoutCommand(state *authState) identitycmd.LogoutCommand {
	return identitycmd.LogoutCommand{IDSession: state.Claims.IDSession}
}
