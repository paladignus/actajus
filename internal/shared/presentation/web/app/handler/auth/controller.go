package authhandler

import (
	"net/http"
	"strconv"

	identitycmd "github.com/paladignus/actajus/internal/module/identity/application/command"
	identityusecase "github.com/paladignus/actajus/internal/module/identity/application/usecase"
	sharedhttp "github.com/paladignus/actajus/internal/shared/infrastructure/http/handler"
	requestmeta "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/requestmeta"
	authbinder "github.com/paladignus/actajus/internal/shared/presentation/web/app/binder/http/auth"
	handlerctx "github.com/paladignus/actajus/internal/shared/presentation/web/app/handlerctx"
	authmessage "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/message/auth"
	authpage "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/page/auth"
	httperror "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/http/error"
)

type Controller struct {
	Context handlerctx.Context
	Login   identityusecase.Login
	Logout  identityusecase.Logout
}

func (c Controller) Home() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := c.Context.AuthFromRequest(w, r)
		if err == nil && state != nil && state.Claims != nil {
			c.Context.Navigator.SeeOther(w, r, "/sessions")
			return
		}
		c.Context.Navigator.SeeOther(w, r, "/login")
	})
}

func (c Controller) LoginPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if state, err := c.Context.AuthFromRequest(w, r); err == nil && state != nil && state.Claims != nil {
			c.Context.Navigator.SeeOther(w, r, "/sessions")
			return
		}
		_ = c.Context.HTML.Render(w, http.StatusOK, "pages/login", authpage.LoginPage("", ""))
	})
}

func (c Controller) LoginAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		form, err := authbinder.BindLoginForm(r)
		if err != nil {
			httperror.BadRequest(w, "invalid form")
			return
		}
		if c.Context.CheckLoginAttempt != nil {
			allowed, retryAfter := c.Context.CheckLoginAttempt(r, form.Email)
			if !allowed {
				if c.Context.AuditFailure != nil {
					c.Context.AuditFailure(r, "auth.login.rate_limited", "email", form.Email)
				}
				if retryAfter > 0 {
					w.Header().Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())+1))
				}
				_ = c.Context.HTML.Render(w, http.StatusTooManyRequests, "pages/login", authpage.LoginPage(form.Email, "Muitas tentativas de login. Tente novamente em instantes."))
				return
			}
		}
		input := identitycmd.LoginCommand{
			Email:     form.Email,
			Password:  form.Password,
			IP:        requestmeta.ClientIP(r),
			UserAgent: r.UserAgent(),
		}
		tokens, err := c.Login.Execute(r.Context(), input)
		if err != nil {
			if c.Context.RecordLoginFailure != nil {
				c.Context.RecordLoginFailure(r, form.Email)
			}
			if c.Context.AuditFailure != nil {
				c.Context.AuditFailure(r, "auth.login.failed", "email", form.Email)
			}
			status, message := authmessage.LoginError(err)
			_ = c.Context.HTML.Render(w, status, "pages/login", authpage.LoginPage(form.Email, message))
			return
		}
		if c.Context.RecordLoginSuccess != nil {
			c.Context.RecordLoginSuccess(r, form.Email)
		}
		if c.Context.AuditSuccess != nil {
			c.Context.AuditSuccess(r, "auth.login.succeeded", "id_user", tokens.IDUser, "id_session", tokens.IDSession)
		}
		c.Context.SetAuthCookies(w, tokens)
		c.Context.Navigator.SeeOther(w, r, "/sessions")
	})
}

func (c Controller) LogoutAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := c.Context.RequireAuth(w, r)
		if err != nil {
			return
		}
		_ = c.Logout.Execute(r.Context(), identitycmd.LogoutCommand{IDSession: state.Claims.IDSession})
		if c.Context.AuditSuccess != nil {
			c.Context.AuditSuccess(r, "auth.logout.succeeded", "id_user", state.Claims.IDUser, "id_session", state.Claims.IDSession)
		}
		c.Context.ClearAuthCookies(w)
		c.Context.Navigator.SeeOther(w, r, "/login")
	})
}

func (c Controller) Bootstrap() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, _ := c.Context.AuthFromRequest(w, r)
		payload := c.Context.AuthAdapter.Bootstrap(r, state)
		if err := sharedhttp.RespondJSON(w, http.StatusOK, payload); err != nil {
			c.Context.Logger.Error(r.Context(), "write web bootstrap", "error", err)
		}
	})
}
