package handler

import (
	"net/http"

	identitycmd "github.com/paladignus/actajus/internal/module/identity/application/command"
	sharedhttp "github.com/paladignus/actajus/internal/shared/infrastructure/http/handler"
	requestmeta "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/requestmeta"
	authbinder "github.com/paladignus/actajus/internal/shared/presentation/web/app/binder/http/auth"
	authmessage "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/message/auth"
	authpage "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/page/auth"
	httperror "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/http/error"
)

type authRoutes struct {
	app AppHandler
}

func (h authRoutes) Home() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := h.app.guard().authFromRequest(w, r)
		if err == nil && state != nil && state.Claims != nil {
			h.app.navigator.SeeOther(w, r, "/sessions")
			return
		}
		h.app.navigator.SeeOther(w, r, "/login")
	})
}

func (h authRoutes) LoginPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if state, err := h.app.guard().authFromRequest(w, r); err == nil && state != nil && state.Claims != nil {
			h.app.navigator.SeeOther(w, r, "/sessions")
			return
		}
		_ = h.app.html.Render(w, http.StatusOK, "pages/login", authpage.LoginPage("", ""))
	})
}

func (h authRoutes) LoginAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		form, err := authbinder.BindLoginForm(r)
		if err != nil {
			httperror.BadRequest(w, "invalid form")
			return
		}
		input := identitycmd.LoginCommand{
			Email:     form.Email,
			Password:  form.Password,
			IP:        requestmeta.ClientIP(r),
			UserAgent: r.UserAgent(),
		}
		tokens, err := h.app.auth.login.Execute(r.Context(), input)
		if err != nil {
			status, message := authmessage.LoginError(err)
			_ = h.app.html.Render(w, status, "pages/login", authpage.LoginPage(form.Email, message))
			return
		}
		h.app.guard().setAuthCookies(w, tokens)
		h.app.navigator.SeeOther(w, r, "/sessions")
	})
}

func (h authRoutes) LogoutAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := h.app.guard().requireAuth(w, r)
		if err != nil {
			return
		}
		_ = h.app.auth.logout.Execute(r.Context(), h.app.guard().logoutCommand(state))
		h.app.guard().clearAuthCookies(w)
		h.app.navigator.SeeOther(w, r, "/login")
	})
}

func (h authRoutes) Bootstrap() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, _ := h.app.guard().authFromRequest(w, r)
		payload := h.app.auth.adapter.Bootstrap(r, state)
		if err := sharedhttp.RespondJSON(w, http.StatusOK, payload); err != nil {
			h.app.logger.Error(r.Context(), "write web bootstrap", "error", err)
		}
	})
}
