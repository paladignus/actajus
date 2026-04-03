package handler

import (
	"net/http"
	"strings"
	"time"

	identitycmd "github.com/paladignus/actajus/internal/module/identity/application/command"
	identityread "github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	identityrepo "github.com/paladignus/actajus/internal/module/identity/application/repository"
	requestmeta "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/requestmeta"
	identitybinder "github.com/paladignus/actajus/internal/shared/presentation/web/app/binder/http/identity"
	identitymessage "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/message/identity"
	identitypage "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/page/identity"
	identityviewmodel "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/viewmodel/identity"
	httperror "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/http/error"
)

type identityRoutes struct {
	app AppHandler
}

func (h identityRoutes) SessionsPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := h.app.requireAuth(w, r)
		if err != nil {
			return
		}
		adminFilter := strings.TrimSpace(r.URL.Query().Get("email"))
		flash := h.app.navigator.ReadFlash(w, r, cookieFlash)
		own, err := h.app.identity.viewSessions.Own(r.Context(), state.Claims.IDUser)
		if err != nil {
			h.app.logger.Error(r.Context(), "list own sessions", "error", err, "id_user", state.Claims.IDUser)
			pageData := identityviewmodel.BuildSessionsPage(state.Claims.IDUser, state.CanManageSessions, adminFilter, flash, nil, nil, state.Claims.IDSession, time.Now())
			pageData.Error = "Nao foi possivel carregar as suas sessoes."
			page := identitypage.SessionsPage(state.Claims.IDUser, state.CanManageSessions, adminFilter, flash, nil, nil, state.Claims.IDSession, time.Now())
			page.Bootstrap = map[string]any{"appName": "Actajus", "apiBaseURL": requestmeta.BaseURL(r), "viewerID": state.Claims.IDUser, "canManageSessions": state.CanManageSessions}
			page.Data = pageData
			if err := h.app.html.Render(w, http.StatusInternalServerError, "pages/sessions", page); err != nil {
				h.app.logger.Error(r.Context(), "render sessions page", "error", err)
				httperror.InternalServerError(w)
			}
			return
		}
		adminSessions := []identityread.SessionReadModel(nil)
		if state.CanManageSessions {
			adminSessions, err = h.app.identity.viewSessions.Admin(r.Context(), identityrepo.SessionQueryFilter{UserEmail: adminFilter, Limit: 50})
			if err != nil {
				h.app.logger.Error(r.Context(), "list admin sessions", "error", err, "email", adminFilter)
				pageData := identityviewmodel.BuildSessionsPage(state.Claims.IDUser, state.CanManageSessions, adminFilter, flash, own, nil, state.Claims.IDSession, time.Now())
				pageData.Error = "Nao foi possivel carregar as sessoes administrativas."
				page := identitypage.SessionsPage(state.Claims.IDUser, state.CanManageSessions, adminFilter, flash, own, nil, state.Claims.IDSession, time.Now())
				page.Bootstrap = map[string]any{"appName": "Actajus", "apiBaseURL": requestmeta.BaseURL(r), "viewerID": state.Claims.IDUser, "canManageSessions": state.CanManageSessions}
				page.Data = pageData
				if err := h.app.html.Render(w, http.StatusInternalServerError, "pages/sessions", page); err != nil {
					h.app.logger.Error(r.Context(), "render sessions page", "error", err)
					httperror.InternalServerError(w)
				}
				return
			}
		}
		page := identitypage.SessionsPage(state.Claims.IDUser, state.CanManageSessions, adminFilter, flash, own, adminSessions, state.Claims.IDSession, time.Now())
		page.Bootstrap = map[string]any{
			"appName":           "Actajus",
			"apiBaseURL":        requestmeta.BaseURL(r),
			"viewerID":          state.Claims.IDUser,
			"canManageSessions": state.CanManageSessions,
		}
		if err := h.app.html.Render(w, http.StatusOK, "pages/sessions", page); err != nil {
			h.app.logger.Error(r.Context(), "render sessions page", "error", err)
			httperror.InternalServerError(w)
		}
	})
}

func (h identityRoutes) RevokeSessionAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := h.app.requireAuth(w, r)
		if err != nil {
			return
		}
		form, cmd, err := identitybinder.BindRevokeSessionForm(r)
		if err != nil {
			h.app.navigator.SeeOtherWithFlash(w, r, "/sessions", cookieFlash, identitymessage.InvalidForm())
			return
		}
		target, err := h.app.identity.viewSessions.ByID(r.Context(), form.IDSession)
		if err != nil {
			h.app.logger.Error(r.Context(), "find session", "error", err, "id_session", form.IDSession)
			h.app.navigator.SeeOtherWithFlash(w, r, "/sessions", cookieFlash, identitymessage.SessionLookupFailed())
			return
		}
		if target == nil {
			h.app.navigator.SeeOtherWithFlash(w, r, "/sessions", cookieFlash, identitymessage.SessionNotFound())
			return
		}
		if target.IDUser != state.Claims.IDUser && !state.CanManageSessions {
			h.app.respondForbidden(w, r)
			return
		}
		if err := h.app.identity.revokeSession.Execute(r.Context(), cmd); err != nil {
			h.app.logger.Error(r.Context(), "revoke session", "error", err, "id_session", form.IDSession)
			h.app.navigator.SeeOtherWithFlash(w, r, "/sessions", cookieFlash, identitymessage.SessionRevokeFailed())
			return
		}
		if form.IDSession == state.Claims.IDSession {
			h.app.clearAuthCookies(w)
			h.app.navigator.SeeOther(w, r, "/login")
			return
		}
		h.app.navigator.SeeOtherWithFlash(w, r, "/sessions", cookieFlash, identitymessage.SessionRevokeSucceeded())
	})
}

func (h identityRoutes) RevokeAllAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := h.app.requireAuth(w, r)
		if err != nil {
			return
		}
		form, cmd, err := identitybinder.BindRevokeAllSessionsForm(r)
		if err != nil {
			h.app.navigator.SeeOtherWithFlash(w, r, "/sessions", cookieFlash, identitymessage.InvalidForm())
			return
		}
		targetUserID := state.Claims.IDUser
		if form.IDUser != nil {
			targetUserID = *form.IDUser
		}
		if targetUserID != state.Claims.IDUser && !state.CanManageSessions {
			h.app.respondForbidden(w, r)
			return
		}
		if cmd == nil {
			cmd = &identitycmd.LogoutAllCommand{IDUser: targetUserID}
		}
		if err := h.app.identity.logoutAll.Execute(r.Context(), *cmd); err != nil {
			h.app.logger.Error(r.Context(), "revoke all sessions", "error", err, "id_user", targetUserID)
			h.app.navigator.SeeOtherWithFlash(w, r, "/sessions", cookieFlash, identitymessage.SessionRevokeAllFailed())
			return
		}
		if targetUserID == state.Claims.IDUser {
			h.app.clearAuthCookies(w)
			h.app.navigator.SeeOther(w, r, "/login")
			return
		}
		h.app.navigator.SeeOtherWithFlash(w, r, "/sessions", cookieFlash, identitymessage.SessionRevokeAllSucceeded())
	})
}
