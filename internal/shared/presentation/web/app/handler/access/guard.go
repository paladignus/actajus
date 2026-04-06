package access

import (
	"errors"
	"net/http"

	identityread "github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	webauth "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/auth"
	handlerctx "github.com/paladignus/actajus/internal/shared/presentation/web/app/handlerctx"
	navigation "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/http/navigation"
)

type Guard struct {
	Adapter            *webauth.Adapter
	Navigator          *navigation.Navigator
	CookieAccessToken  string
	CookieRefreshToken string
	CookieSessionID    string
	ManageSessionsPerm string
	ManageRBACPerm     string
	RespondForbidden   func(http.ResponseWriter, *http.Request)
}

func (g Guard) AuthFromRequest(w http.ResponseWriter, r *http.Request) (*handlerctx.State, error) {
	state, refreshed, err := g.Adapter.FromRequest(w, r, g.CookieAccessToken, g.CookieRefreshToken, g.CookieSessionID, g.ManageSessionsPerm, g.ManageRBACPerm)
	if refreshed != nil {
		g.Adapter.SetCookies(w, refreshed, g.CookieAccessToken, g.CookieRefreshToken, g.CookieSessionID)
	}
	return state, err
}

func (g Guard) RequireAuth(w http.ResponseWriter, r *http.Request) (*handlerctx.State, error) {
	state, err := g.AuthFromRequest(w, r)
	if err != nil || state == nil || state.Claims == nil {
		g.Navigator.SeeOther(w, r, "/login")
		return nil, err
	}
	return state, nil
}

func (g Guard) RequireRBACAdmin(w http.ResponseWriter, r *http.Request) (*handlerctx.State, error) {
	state, err := g.RequireAuth(w, r)
	if err != nil {
		return nil, err
	}
	if state == nil || !state.CanManageRBAC {
		g.RespondForbidden(w, r)
		return nil, errors.New("missing rbac manage permission")
	}
	return state, nil
}

func (g Guard) SetAuthCookies(w http.ResponseWriter, tokens *identityread.AuthTokensReadModel) {
	g.Adapter.SetCookies(w, tokens, g.CookieAccessToken, g.CookieRefreshToken, g.CookieSessionID)
}

func (g Guard) ClearAuthCookies(w http.ResponseWriter) {
	g.Adapter.ClearCookies(w, g.CookieAccessToken, g.CookieRefreshToken, g.CookieSessionID)
}
