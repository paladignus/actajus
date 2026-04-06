// Package webapp wires the shared WEB presentation layer.
package webapp

import (
	"net/http"

	"github.com/paladignus/actajus/internal/module/company"
	"github.com/paladignus/actajus/internal/module/identity"
	sharedRepo "github.com/paladignus/actajus/internal/shared/application/repository"
	webcsrf "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/csrf"
	securityheaders "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/securityheaders"
	sharedhandler "github.com/paladignus/actajus/internal/shared/presentation/web/app/handler"
	webtemplate "github.com/paladignus/actajus/internal/shared/presentation/web/template"
)

type Dependencies struct {
	Logger         sharedRepo.Logger
	Company        company.Module
	Identity       identity.Module
	Secure         bool
	RendererSource webtemplate.Source
}

type Module struct {
	handler sharedhandler.AppHandler
	csrf    *webcsrf.Adapter
	headers *securityheaders.Adapter
}

func NewModule(dep Dependencies) (Module, error) {
	h, err := sharedhandler.NewAppHandler(sharedhandler.Dependencies{
		Logger:   dep.Logger,
		Renderer: dep.RendererSource,
		Auth: sharedhandler.AuthDependencies{
			ValidateAccess:           dep.Identity.ValidateAccess,
			Login:                    dep.Identity.Login,
			Register:                 dep.Identity.Register,
			RequestEmailVerification: dep.Identity.RequestEmailVerification,
			ConfirmEmailVerification: dep.Identity.ConfirmEmailVerification,
			Refresh:                  dep.Identity.Refresh,
			Logout:                   dep.Identity.Logout,
			RBACChecker:              dep.Identity.RBACChecker,
		},
		Identity: sharedhandler.IdentityDependencies{
			LogoutAll:                dep.Identity.LogoutAll,
			RevokeSession:            dep.Identity.RevokeSession,
			ViewSessions:             dep.Identity.ViewSessions,
			ListUsers:                dep.Identity.ListUsers,
			ListPermissions:          dep.Identity.ListPermissions,
			ListRoles:                dep.Identity.ListRoles,
			ViewUserRoles:            dep.Identity.ViewUserRoles,
			ViewRolePermissions:      dep.Identity.ViewRolePermissions,
			AssignRoleToUser:         dep.Identity.AssignRoleToUser,
			RemoveRoleFromUser:       dep.Identity.RemoveRoleFromUser,
			GrantPermissionToRole:    dep.Identity.GrantPermissionToRole,
			RevokePermissionFromRole: dep.Identity.RevokePermissionFromRole,
			CreateRole:               dep.Identity.CreateRole,
			UpdateRole:               dep.Identity.UpdateRole,
			DeleteRole:               dep.Identity.DeleteRole,
			CreatePermission:         dep.Identity.CreatePermission,
			UpdatePermission:         dep.Identity.UpdatePermission,
			DeletePermission:         dep.Identity.DeletePermission,
		},
		Company: sharedhandler.CompanyDependencies{
			Create:   dep.Company.Create,
			List:     dep.Company.List,
			Update:   dep.Company.Update,
			Delete:   dep.Company.Delete,
			FindByID: dep.Company.FindByID,
		},
		SecureCookies: dep.Secure,
	})
	if err != nil {
		return Module{}, err
	}
	return Module{
		handler: h,
		csrf:    webcsrf.New(dep.Secure),
		headers: securityheaders.New(dep.RendererSource.DevServerURL),
	}, nil
}

func (m Module) Mount(mux *http.ServeMux) {
	mux.Handle("GET /{$}", m.withPageSecurity(m.withCSRFCookie(m.handler.Home())))
	mux.Handle("GET /login", m.withPageSecurity(m.withCSRFCookie(m.handler.LoginPage())))
	mux.Handle("GET /register", m.withPageSecurity(m.withCSRFCookie(m.handler.RegisterPage())))
	mux.Handle("GET /verification-pending", m.withPageSecurity(m.withCSRFCookie(m.handler.VerificationPendingPage())))
	mux.Handle("POST /login", m.withCSRFProtection(m.handler.LoginAction()))
	mux.Handle("POST /register", m.withCSRFProtection(m.handler.RegisterAction()))
	mux.Handle("POST /request-email-verification", m.withCSRFProtection(m.handler.RequestEmailVerificationAction()))
	mux.Handle("GET /verify-email", m.withPageSecurity(m.handler.VerifyEmailAction()))
	mux.Handle("POST /logout", m.withCSRFProtection(m.handler.LogoutAction()))
	mux.Handle("GET /sessions", m.withPageSecurity(m.withCSRFCookie(m.handler.SessionsPage())))
	mux.Handle("GET /users", m.withPageSecurity(m.withCSRFCookie(m.handler.UsersPage())))
	mux.Handle("GET /roles", m.withPageSecurity(m.withCSRFCookie(m.handler.RolesPage())))
	mux.Handle("GET /permissions", m.withPageSecurity(m.withCSRFCookie(m.handler.PermissionsPage())))
	mux.Handle("POST /roles", m.withCSRFProtection(m.handler.CreateRoleAction()))
	mux.Handle("POST /roles/{id}", m.withCSRFProtection(m.handler.UpdateRoleAction()))
	mux.Handle("POST /roles/{id}/delete", m.withCSRFProtection(m.handler.DeleteRoleAction()))
	mux.Handle("POST /permissions", m.withCSRFProtection(m.handler.CreatePermissionAction()))
	mux.Handle("POST /permissions/{id}", m.withCSRFProtection(m.handler.UpdatePermissionAction()))
	mux.Handle("POST /permissions/{id}/delete", m.withCSRFProtection(m.handler.DeletePermissionAction()))
	mux.Handle("POST /users/roles", m.withCSRFProtection(m.handler.AssignRoleToUserAction()))
	mux.Handle("POST /users/roles/delete", m.withCSRFProtection(m.handler.RemoveRoleFromUserAction()))
	mux.Handle("POST /roles/permissions", m.withCSRFProtection(m.handler.GrantPermissionToRoleAction()))
	mux.Handle("POST /roles/permissions/delete", m.withCSRFProtection(m.handler.RevokePermissionFromRoleAction()))
	mux.Handle("POST /sessions/revoke", m.withCSRFProtection(m.handler.RevokeSessionAction()))
	mux.Handle("POST /sessions/revoke-all", m.withCSRFProtection(m.handler.RevokeAllAction()))
	mux.Handle("GET /companies", m.withPageSecurity(m.withCSRFCookie(m.handler.CompaniesPage())))
	mux.Handle("GET /companies/new", m.withPageSecurity(m.withCSRFCookie(m.handler.CompanyNewPage())))
	mux.Handle("POST /companies", m.withCSRFProtection(m.handler.CompanyCreateAction()))
	mux.Handle("GET /companies/{id}", m.withPageSecurity(m.withCSRFCookie(m.handler.CompanyDetailPage())))
	mux.Handle("GET /companies/{id}/edit", m.withPageSecurity(m.withCSRFCookie(m.handler.CompanyEditPage())))
	mux.Handle("POST /companies/{id}", m.withCSRFProtection(m.handler.CompanyUpdateAction()))
	mux.Handle("POST /companies/{id}/delete", m.withCSRFProtection(m.handler.CompanyDeleteAction()))
	mux.Handle("GET /bootstrap", m.handler.Bootstrap())
	mux.Handle("GET /favicon.svg", m.withAssetSecurity(m.handler.StaticFile("favicon.svg")))
	mux.Handle("GET /favicon.ico", http.RedirectHandler("/favicon.svg", http.StatusPermanentRedirect))
	mux.Handle("GET /icons.svg", m.withAssetSecurity(m.handler.StaticFile("icons.svg")))
	mux.Handle("GET /assets/", m.withAssetSecurity(http.StripPrefix("/assets/", m.handler.Assets())))
}

func (m Module) withCSRFCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.csrf.EnsureCookie(w, r)
		next.ServeHTTP(w, r)
	})
}

func (m Module) withCSRFProtection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.csrf.Validate(r) {
			m.handler.RenderForbidden(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m Module) withPageSecurity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.headers.ApplyDocument(w)
		next.ServeHTTP(w, r)
	})
}

func (m Module) withAssetSecurity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.headers.ApplyAsset(w)
		next.ServeHTTP(w, r)
	})
}

func (m Module) RenderNotFound(w http.ResponseWriter, r *http.Request) {
	m.handler.RenderNotFound(w, r)
}

func (m Module) RenderForbidden(w http.ResponseWriter, r *http.Request) {
	m.handler.RenderForbidden(w, r)
}
