package handler

import "net/http"

func (h AppHandler) authRoutes() authRoutes {
	return authRoutes{app: h}
}

func (h AppHandler) identityRoutes() identityRoutes {
	return identityRoutes{app: h}
}

func (h AppHandler) companyRoutes() companyRoutes {
	return companyRoutes{app: h}
}

func (h AppHandler) supportRoutes() supportRoutes {
	return supportRoutes{app: h}
}

func (h AppHandler) Home() http.Handler         { return h.authRoutes().Home() }
func (h AppHandler) LoginPage() http.Handler    { return h.authRoutes().LoginPage() }
func (h AppHandler) LoginAction() http.Handler  { return h.authRoutes().LoginAction() }
func (h AppHandler) LogoutAction() http.Handler { return h.authRoutes().LogoutAction() }
func (h AppHandler) Bootstrap() http.Handler    { return h.authRoutes().Bootstrap() }
func (h AppHandler) SessionsPage() http.Handler { return h.identityRoutes().SessionsPage() }
func (h AppHandler) RevokeSessionAction() http.Handler {
	return h.identityRoutes().RevokeSessionAction()
}
func (h AppHandler) RevokeAllAction() http.Handler  { return h.identityRoutes().RevokeAllAction() }
func (h AppHandler) UsersPage() http.Handler        { return h.identityRoutes().UsersPage() }
func (h AppHandler) PermissionsPage() http.Handler  { return h.identityRoutes().PermissionsPage() }
func (h AppHandler) RolesPage() http.Handler        { return h.identityRoutes().RolesPage() }
func (h AppHandler) CreateRoleAction() http.Handler { return h.identityRoutes().CreateRoleAction() }
func (h AppHandler) UpdateRoleAction() http.Handler { return h.identityRoutes().UpdateRoleAction() }
func (h AppHandler) DeleteRoleAction() http.Handler { return h.identityRoutes().DeleteRoleAction() }
func (h AppHandler) CreatePermissionAction() http.Handler {
	return h.identityRoutes().CreatePermissionAction()
}
func (h AppHandler) UpdatePermissionAction() http.Handler {
	return h.identityRoutes().UpdatePermissionAction()
}
func (h AppHandler) DeletePermissionAction() http.Handler {
	return h.identityRoutes().DeletePermissionAction()
}
func (h AppHandler) AssignRoleToUserAction() http.Handler {
	return h.identityRoutes().AssignRoleToUserAction()
}
func (h AppHandler) RemoveRoleFromUserAction() http.Handler {
	return h.identityRoutes().RemoveRoleFromUserAction()
}
func (h AppHandler) GrantPermissionToRoleAction() http.Handler {
	return h.identityRoutes().GrantPermissionToRoleAction()
}
func (h AppHandler) RevokePermissionFromRoleAction() http.Handler {
	return h.identityRoutes().RevokePermissionFromRoleAction()
}
func (h AppHandler) CompaniesPage() http.Handler  { return h.companyRoutes().CompaniesPage() }
func (h AppHandler) CompanyNewPage() http.Handler { return h.companyRoutes().CompanyNewPage() }
func (h AppHandler) CompanyCreateAction() http.Handler {
	return h.companyRoutes().CompanyCreateAction()
}
func (h AppHandler) CompanyDetailPage() http.Handler { return h.companyRoutes().CompanyDetailPage() }
func (h AppHandler) CompanyEditPage() http.Handler   { return h.companyRoutes().CompanyEditPage() }
func (h AppHandler) CompanyUpdateAction() http.Handler {
	return h.companyRoutes().CompanyUpdateAction()
}
func (h AppHandler) CompanyDeleteAction() http.Handler {
	return h.companyRoutes().CompanyDeleteAction()
}
func (h AppHandler) Assets() http.Handler                { return h.supportRoutes().Assets() }
func (h AppHandler) StaticFile(name string) http.Handler { return h.supportRoutes().StaticFile(name) }
func (h AppHandler) RenderNotFound(w http.ResponseWriter, r *http.Request) {
	h.supportRoutes().RenderNotFound(w, r)
}
func (h AppHandler) RenderForbidden(w http.ResponseWriter, r *http.Request) {
	h.supportRoutes().RenderForbidden(w, r)
}
func (h AppHandler) authFromRequest(w http.ResponseWriter, r *http.Request) (*authState, error) {
	return h.guard().authFromRequest(w, r)
}
func (h AppHandler) requireAuth(w http.ResponseWriter, r *http.Request) (*authState, error) {
	return h.guard().requireAuth(w, r)
}
func (h AppHandler) requireRBACAdmin(w http.ResponseWriter, r *http.Request) (*authState, error) {
	return h.guard().requireRBACAdmin(w, r)
}
func (h AppHandler) clearAuthCookies(w http.ResponseWriter) {
	h.guard().clearAuthCookies(w)
}
func (h AppHandler) respondForbidden(w http.ResponseWriter, r *http.Request) {
	h.supportRoutes().respondForbidden(w, r)
}
func (h AppHandler) respondNotFound(w http.ResponseWriter, r *http.Request) {
	h.supportRoutes().respondNotFound(w, r)
}
