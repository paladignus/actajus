// Package web wires the WEB presentation module.
package web

import (
	"net/http"

	"github.com/paladignus/actajus/internal/module/company"
	"github.com/paladignus/actajus/internal/module/identity"
	"github.com/paladignus/actajus/internal/module/web/presentation/http/handler"
	sharedRepo "github.com/paladignus/actajus/internal/shared/application/repository"
)

type Dependencies struct {
	Logger   sharedRepo.Logger
	Company  company.Module
	Identity identity.Module
	Secure   bool
}

type Module struct {
	handler handler.AppHandler
}

func NewModule(dep Dependencies) (Module, error) {
	h, err := handler.NewAppHandler(
		dep.Logger,
		dep.Identity.ValidateAccess,
		dep.Identity.Login,
		dep.Identity.Refresh,
		dep.Identity.Logout,
		dep.Identity.LogoutAll,
		dep.Identity.RevokeSession,
		dep.Identity.ViewSessions,
		dep.Identity.ListUsers,
		dep.Identity.ListPermissions,
		dep.Identity.ListRoles,
		dep.Identity.ViewUserRoles,
		dep.Identity.ViewRolePermissions,
		dep.Identity.AssignRoleToUser,
		dep.Identity.RemoveRoleFromUser,
		dep.Identity.GrantPermissionToRole,
		dep.Identity.RevokePermissionFromRole,
		dep.Identity.CreateRole,
		dep.Identity.UpdateRole,
		dep.Identity.DeleteRole,
		dep.Identity.CreatePermission,
		dep.Identity.UpdatePermission,
		dep.Identity.DeletePermission,
		dep.Company.Create,
		dep.Company.List,
		dep.Company.Update,
		dep.Company.Delete,
		dep.Company.FindByID,
		dep.Identity.RBACChecker,
		dep.Secure,
	)
	if err != nil {
		return Module{}, err
	}
	return Module{handler: h}, nil
}

func (m Module) Mount(mux *http.ServeMux) {
	mux.Handle("GET /{$}", m.handler.Home())
	mux.Handle("GET /login", m.handler.LoginPage())
	mux.Handle("POST /login", m.handler.LoginAction())
	mux.Handle("POST /logout", m.handler.LogoutAction())
	mux.Handle("GET /sessions", m.handler.SessionsPage())
	mux.Handle("GET /users", m.handler.UsersPage())
	mux.Handle("GET /roles", m.handler.RolesPage())
	mux.Handle("GET /permissions", m.handler.PermissionsPage())
	mux.Handle("POST /roles", m.handler.CreateRoleAction())
	mux.Handle("POST /roles/{id}", m.handler.UpdateRoleAction())
	mux.Handle("POST /roles/{id}/delete", m.handler.DeleteRoleAction())
	mux.Handle("POST /permissions", m.handler.CreatePermissionAction())
	mux.Handle("POST /permissions/{id}", m.handler.UpdatePermissionAction())
	mux.Handle("POST /permissions/{id}/delete", m.handler.DeletePermissionAction())
	mux.Handle("POST /users/roles", m.handler.AssignRoleToUserAction())
	mux.Handle("POST /users/roles/delete", m.handler.RemoveRoleFromUserAction())
	mux.Handle("POST /roles/permissions", m.handler.GrantPermissionToRoleAction())
	mux.Handle("POST /roles/permissions/delete", m.handler.RevokePermissionFromRoleAction())
	mux.Handle("POST /sessions/revoke", m.handler.RevokeSessionAction())
	mux.Handle("POST /sessions/revoke-all", m.handler.RevokeAllAction())
	mux.Handle("GET /companies", m.handler.CompaniesPage())
	mux.Handle("GET /companies/new", m.handler.CompanyNewPage())
	mux.Handle("POST /companies", m.handler.CompanyCreateAction())
	mux.Handle("GET /companies/{id}", m.handler.CompanyDetailPage())
	mux.Handle("GET /companies/{id}/edit", m.handler.CompanyEditPage())
	mux.Handle("POST /companies/{id}", m.handler.CompanyUpdateAction())
	mux.Handle("POST /companies/{id}/delete", m.handler.CompanyDeleteAction())
	mux.Handle("GET /bootstrap", m.handler.Bootstrap())
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", m.handler.Assets()))
}

func (m Module) RenderNotFound(w http.ResponseWriter, r *http.Request) {
	m.handler.RenderNotFound(w, r)
}

func (m Module) RenderForbidden(w http.ResponseWriter, r *http.Request) {
	m.handler.RenderForbidden(w, r)
}
