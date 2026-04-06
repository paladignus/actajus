package auth

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	identitycmd "github.com/paladignus/actajus/internal/module/identity/application/command"
	identityread "github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	identityusecase "github.com/paladignus/actajus/internal/module/identity/application/usecase"
	identitydomain "github.com/paladignus/actajus/internal/module/identity/domain"
	sharedRepo "github.com/paladignus/actajus/internal/shared/application/repository"
	cookiepolicy "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/cookiepolicy"
)

type State struct {
	Claims            *identityread.AccessTokenClaims
	CanManageSessions bool
	CanManageRBAC     bool
}

type Adapter struct {
	validateAccess    identityusecase.ValidateAccess
	refresh           identityusecase.Refresh
	permissionChecker sharedRepo.PermissionChecker
	secureCookies     bool
	clientIP          func(*http.Request) string
	baseURL           func(*http.Request) string
}

func New(
	validateAccess identityusecase.ValidateAccess,
	refresh identityusecase.Refresh,
	permissionChecker sharedRepo.PermissionChecker,
	secureCookies bool,
	clientIP func(*http.Request) string,
	baseURL func(*http.Request) string,
) *Adapter {
	return &Adapter{
		validateAccess:    validateAccess,
		refresh:           refresh,
		permissionChecker: permissionChecker,
		secureCookies:     secureCookies,
		clientIP:          clientIP,
		baseURL:           baseURL,
	}
}

func (a *Adapter) FromRequest(w http.ResponseWriter, r *http.Request, accessCookieName, refreshCookieName, sessionCookieName, permManageSessions, permManageRBAC string) (*State, *identityread.AuthTokensReadModel, error) {
	accessCookie, err := r.Cookie(accessCookieName)
	if err == nil && strings.TrimSpace(accessCookie.Value) != "" {
		claims, verr := a.validateAccess.Execute(r.Context(), identitycmd.ValidateAccessCommand{
			AccessToken: accessCookie.Value,
		})
		if verr == nil {
			state, stateErr := a.BuildState(r, claims, permManageSessions, permManageRBAC)
			return state, nil, stateErr
		}
		if !ShouldAttemptRefresh(verr) {
			a.ClearCookies(w, accessCookieName, refreshCookieName, sessionCookieName)
			return nil, nil, verr
		}
	}

	refreshCookie, errRefresh := r.Cookie(refreshCookieName)
	sessionCookie, errSession := r.Cookie(sessionCookieName)
	if errRefresh != nil || errSession != nil {
		return nil, nil, errors.New("missing auth cookies")
	}
	sessionID, err := strconv.ParseInt(sessionCookie.Value, 10, 64)
	if err != nil || sessionID <= 0 {
		a.ClearCookies(w, accessCookieName, refreshCookieName, sessionCookieName)
		return nil, nil, errors.New("invalid session cookie")
	}
	tokens, err := a.refresh.Execute(r.Context(), identitycmd.RefreshCommand{
		IDSession:    sessionID,
		RefreshToken: refreshCookie.Value,
		IP:           a.clientIP(r),
		UserAgent:    r.UserAgent(),
	})
	if err != nil {
		a.ClearCookies(w, accessCookieName, refreshCookieName, sessionCookieName)
		return nil, nil, err
	}
	state, stateErr := a.BuildState(r, &identityread.AccessTokenClaims{
		IDUser:    tokens.IDUser,
		IDSession: tokens.IDSession,
	}, permManageSessions, permManageRBAC)
	if stateErr != nil {
		return nil, nil, stateErr
	}
	return state, tokens, nil
}

func (a *Adapter) BuildState(r *http.Request, claims *identityread.AccessTokenClaims, permManageSessions, permManageRBAC string) (*State, error) {
	if a.permissionChecker == nil {
		return &State{
			Claims:            claims,
			CanManageSessions: false,
			CanManageRBAC:     false,
		}, nil
	}
	canManage, err := a.permissionChecker.HasPermission(r.Context(), claims.IDUser, permManageSessions)
	if err != nil {
		return nil, err
	}
	canManageRBAC, err := a.permissionChecker.HasPermission(r.Context(), claims.IDUser, permManageRBAC)
	if err != nil {
		return nil, err
	}
	return &State{
		Claims:            claims,
		CanManageSessions: canManage,
		CanManageRBAC:     canManageRBAC,
	}, nil
}

func (a *Adapter) SetCookies(w http.ResponseWriter, tokens *identityread.AuthTokensReadModel, accessCookieName, refreshCookieName, sessionCookieName string) {
	policy := cookiepolicy.New(a.secureCookies)
	http.SetCookie(w, &http.Cookie{
		Name:     accessCookieName,
		Value:    tokens.AccessToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: policy.SameSite(),
		Secure:   policy.Secure(),
		Expires:  tokens.AccessExpiresAt,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    tokens.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: policy.SameSite(),
		Secure:   policy.Secure(),
		Expires:  tokens.RefreshExpiresAt,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    strconv.FormatInt(tokens.IDSession, 10),
		Path:     "/",
		HttpOnly: true,
		SameSite: policy.SameSite(),
		Secure:   policy.Secure(),
		Expires:  tokens.RefreshExpiresAt,
	})
}

func (a *Adapter) ClearCookies(w http.ResponseWriter, accessCookieName, refreshCookieName, sessionCookieName string) {
	policy := cookiepolicy.New(a.secureCookies)
	for _, name := range []string{accessCookieName, refreshCookieName, sessionCookieName} {
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			SameSite: policy.SameSite(),
			Secure:   policy.Secure(),
			MaxAge:   -1,
		})
	}
}

func (a *Adapter) Bootstrap(r *http.Request, state *State) map[string]any {
	payload := map[string]any{
		"appName":       "Actajus",
		"apiBaseURL":    a.baseURL(r),
		"authenticated": state != nil && state.Claims != nil,
	}
	if state != nil && state.Claims != nil {
		payload["idUser"] = state.Claims.IDUser
		payload["idSession"] = state.Claims.IDSession
		payload["canManageSessions"] = state.CanManageSessions
		payload["canManageRBAC"] = state.CanManageRBAC
	}
	return payload
}

func ShouldAttemptRefresh(err error) bool {
	return errors.Is(err, identitydomain.ErrMissingAccessToken) ||
		errors.Is(err, identitydomain.ErrInvalidToken) ||
		errors.Is(err, identitydomain.ErrSessionNotActive) ||
		errors.Is(err, identitydomain.ErrSessionExpired) ||
		errors.Is(err, identitydomain.ErrSessionRevoked)
}
