package handler

import (
	"net/http"
	"time"

	authhandler "github.com/paladignus/actajus/internal/shared/presentation/web/app/handler/auth"
	companyhandler "github.com/paladignus/actajus/internal/shared/presentation/web/app/handler/company"
	identityhandler "github.com/paladignus/actajus/internal/shared/presentation/web/app/handler/identity"
	handlerctx "github.com/paladignus/actajus/internal/shared/presentation/web/app/handlerctx"
)

func (h AppHandler) baseContext() handlerctx.Context {
	return handlerctx.Context{
		Logger:             h.logger,
		HTML:               h.html,
		Navigator:          h.navigator,
		AuthAdapter:        h.auth.adapter,
		FlashCookie:        cookieFlash,
		AuthFromRequest:    h.authFromRequest,
		RequireAuth:        h.requireAuth,
		RequireRBACAdmin:   h.requireRBACAdmin,
		SetAuthCookies:     h.guard.SetAuthCookies,
		ClearAuthCookies:   h.clearAuthCookies,
		RespondForbidden:   h.respondForbidden,
		RespondNotFound:    h.respondNotFound,
		CheckLoginAttempt:  h.checkLoginAttempt,
		RecordLoginFailure: h.recordLoginFailure,
		RecordLoginSuccess: h.recordLoginSuccess,
		AuditSuccess:       h.auditSuccess,
		AuditFailure:       h.auditFailure,
	}
}

func (h AppHandler) authController() authhandler.Controller {
	return authhandler.Controller{
		Context: h.baseContext(),
		Login:   h.auth.login,
		Logout:  h.auth.logout,
	}
}

func (h AppHandler) companyController() companyhandler.Controller {
	return companyhandler.Controller{
		Context:  h.baseContext(),
		Create:   h.company.create,
		List:     h.company.list,
		Update:   h.company.update,
		Delete:   h.company.delete,
		FindByID: h.company.findByID,
	}
}

func (h AppHandler) identityController() identityhandler.Controller {
	return identityhandler.Controller{
		Context:                  h.baseContext(),
		LogoutAll:                h.identity.logoutAll,
		RevokeSession:            h.identity.revokeSession,
		ViewSessions:             h.identity.viewSessions,
		ListUsers:                h.identity.listUsers,
		ListPermissions:          h.identity.listPermissions,
		ListRoles:                h.identity.listRoles,
		ViewUserRoles:            h.identity.viewUserRoles,
		ViewRolePermissions:      h.identity.viewRolePermissions,
		AssignRoleToUser:         h.identity.assignRoleToUser,
		RemoveRoleFromUser:       h.identity.removeRoleFromUser,
		GrantPermissionToRole:    h.identity.grantPermissionToRole,
		RevokePermissionFromRole: h.identity.revokePermissionFromRole,
		CreateRole:               h.identity.createRole,
		UpdateRole:               h.identity.updateRole,
		DeleteRole:               h.identity.deleteRole,
		CreatePermission:         h.identity.createPermission,
		UpdatePermission:         h.identity.updatePermission,
		DeletePermission:         h.identity.deletePermission,
	}
}

func (h AppHandler) Home() http.Handler         { return h.authController().Home() }
func (h AppHandler) LoginPage() http.Handler    { return h.authController().LoginPage() }
func (h AppHandler) LoginAction() http.Handler  { return h.authController().LoginAction() }
func (h AppHandler) LogoutAction() http.Handler { return h.authController().LogoutAction() }
func (h AppHandler) Bootstrap() http.Handler    { return h.authController().Bootstrap() }
func (h AppHandler) SessionsPage() http.Handler { return h.identityController().SessionsPage() }
func (h AppHandler) RevokeSessionAction() http.Handler {
	return h.identityController().RevokeSessionAction()
}
func (h AppHandler) RevokeAllAction() http.Handler  { return h.identityController().RevokeAllAction() }
func (h AppHandler) UsersPage() http.Handler        { return h.identityController().UsersPage() }
func (h AppHandler) PermissionsPage() http.Handler  { return h.identityController().PermissionsPage() }
func (h AppHandler) RolesPage() http.Handler        { return h.identityController().RolesPage() }
func (h AppHandler) CreateRoleAction() http.Handler { return h.identityController().CreateRoleAction() }
func (h AppHandler) UpdateRoleAction() http.Handler { return h.identityController().UpdateRoleAction() }
func (h AppHandler) DeleteRoleAction() http.Handler { return h.identityController().DeleteRoleAction() }
func (h AppHandler) CreatePermissionAction() http.Handler {
	return h.identityController().CreatePermissionAction()
}
func (h AppHandler) UpdatePermissionAction() http.Handler {
	return h.identityController().UpdatePermissionAction()
}
func (h AppHandler) DeletePermissionAction() http.Handler {
	return h.identityController().DeletePermissionAction()
}
func (h AppHandler) AssignRoleToUserAction() http.Handler {
	return h.identityController().AssignRoleToUserAction()
}
func (h AppHandler) RemoveRoleFromUserAction() http.Handler {
	return h.identityController().RemoveRoleFromUserAction()
}
func (h AppHandler) GrantPermissionToRoleAction() http.Handler {
	return h.identityController().GrantPermissionToRoleAction()
}
func (h AppHandler) RevokePermissionFromRoleAction() http.Handler {
	return h.identityController().RevokePermissionFromRoleAction()
}
func (h AppHandler) CompaniesPage() http.Handler  { return h.companyController().CompaniesPage() }
func (h AppHandler) CompanyNewPage() http.Handler { return h.companyController().CompanyNewPage() }
func (h AppHandler) CompanyCreateAction() http.Handler {
	return h.companyController().CompanyCreateAction()
}
func (h AppHandler) CompanyDetailPage() http.Handler {
	return h.companyController().CompanyDetailPage()
}
func (h AppHandler) CompanyEditPage() http.Handler { return h.companyController().CompanyEditPage() }
func (h AppHandler) CompanyUpdateAction() http.Handler {
	return h.companyController().CompanyUpdateAction()
}
func (h AppHandler) CompanyDeleteAction() http.Handler {
	return h.companyController().CompanyDeleteAction()
}
func (h AppHandler) Assets() http.Handler                { return h.support.Assets() }
func (h AppHandler) StaticFile(name string) http.Handler { return h.support.StaticFile(name) }
func (h AppHandler) RenderNotFound(w http.ResponseWriter, r *http.Request) {
	h.support.RenderNotFound(w, r)
}
func (h AppHandler) RenderForbidden(w http.ResponseWriter, r *http.Request) {
	h.support.RenderForbidden(w, r)
}
func (h AppHandler) authFromRequest(w http.ResponseWriter, r *http.Request) (*handlerctx.State, error) {
	return h.guard.AuthFromRequest(w, r)
}
func (h AppHandler) requireAuth(w http.ResponseWriter, r *http.Request) (*handlerctx.State, error) {
	return h.guard.RequireAuth(w, r)
}
func (h AppHandler) requireRBACAdmin(w http.ResponseWriter, r *http.Request) (*handlerctx.State, error) {
	return h.guard.RequireRBACAdmin(w, r)
}
func (h AppHandler) clearAuthCookies(w http.ResponseWriter) {
	h.guard.ClearAuthCookies(w)
}
func (h AppHandler) checkLoginAttempt(r *http.Request, email string) (bool, time.Duration) {
	if h.auth.loginLimiter == nil {
		return true, 0
	}
	return h.auth.loginLimiter.Check(r, email)
}
func (h AppHandler) recordLoginFailure(r *http.Request, email string) {
	if h.auth.loginLimiter == nil {
		return
	}
	h.auth.loginLimiter.RecordFailure(r, email)
}
func (h AppHandler) recordLoginSuccess(r *http.Request, email string) {
	if h.auth.loginLimiter == nil {
		return
	}
	h.auth.loginLimiter.RecordSuccess(r, email)
}
func (h AppHandler) auditSuccess(r *http.Request, action string, attrs ...any) {
	if h.auth.audit == nil {
		return
	}
	h.auth.audit.Success(r.Context(), r, action, attrs...)
}
func (h AppHandler) auditFailure(r *http.Request, action string, attrs ...any) {
	if h.auth.audit == nil {
		return
	}
	h.auth.audit.Failure(r.Context(), r, action, attrs...)
}
func (h AppHandler) respondForbidden(w http.ResponseWriter, r *http.Request) {
	h.support.RespondForbidden(w, r)
}
func (h AppHandler) respondNotFound(w http.ResponseWriter, r *http.Request) {
	h.support.RespondNotFound(w, r)
}
