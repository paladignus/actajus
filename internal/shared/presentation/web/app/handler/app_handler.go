// Package handler contains shared WEB HTTP handlers.
package handler

import (
	"net/http"

	companyusecase "github.com/paladignus/actajus/internal/module/company/application/usecase"
	identityusecase "github.com/paladignus/actajus/internal/module/identity/application/usecase"
	sharedRepo "github.com/paladignus/actajus/internal/shared/application/repository"
	webauth "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/auth"
	webhtml "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/html"
	navigation "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/http/navigation"
	webtemplate "github.com/paladignus/actajus/internal/shared/presentation/web/template"
)

const (
	cookieAccessToken  = "actajus_access_token"
	cookieRefreshToken = "actajus_refresh_token"
	cookieSessionID    = "actajus_session_id"
	cookieFlash        = "actajus_flash"
	permManageSessions = "session:manage"
	permManageRBAC     = "rbac:manage"
)

type AppHandler struct {
	logger    sharedRepo.Logger
	renderer  *webtemplate.Renderer
	assets    http.Handler
	html      *webhtml.Responder
	navigator *navigation.Navigator
	auth      authDependencies
	identity  identityDependencies
	company   companyDependencies
}

type authState = webauth.State

type AuthDependencies struct {
	ValidateAccess identityusecase.ValidateAccess
	Login          identityusecase.Login
	Refresh        identityusecase.Refresh
	Logout         identityusecase.Logout
	RBACChecker    sharedRepo.PermissionChecker
}

type IdentityDependencies struct {
	LogoutAll                identityusecase.LogoutAll
	RevokeSession            identityusecase.RevokeSession
	ViewSessions             identityusecase.ViewSessions
	ListUsers                identityusecase.ListUsers
	ListPermissions          identityusecase.ListPermissions
	ListRoles                identityusecase.ListRoles
	ViewUserRoles            identityusecase.ViewUserRoles
	ViewRolePermissions      identityusecase.ViewRolePermissions
	AssignRoleToUser         identityusecase.AssignRoleToUser
	RemoveRoleFromUser       identityusecase.RemoveRoleFromUser
	GrantPermissionToRole    identityusecase.GrantPermissionToRole
	RevokePermissionFromRole identityusecase.RevokePermissionFromRole
	CreateRole               identityusecase.CreateRole
	UpdateRole               identityusecase.UpdateRole
	DeleteRole               identityusecase.DeleteRole
	CreatePermission         identityusecase.CreatePermission
	UpdatePermission         identityusecase.UpdatePermission
	DeletePermission         identityusecase.DeletePermission
}

type CompanyDependencies struct {
	Create   companyusecase.CreateCompany
	List     companyusecase.ListCompanies
	Update   companyusecase.UpdateCompany
	Delete   companyusecase.DeleteCompany
	FindByID companyusecase.FindByID
}

type Dependencies struct {
	Logger        sharedRepo.Logger
	Renderer      webtemplate.Source
	Auth          AuthDependencies
	Identity      IdentityDependencies
	Company       CompanyDependencies
	SecureCookies bool
}

type authDependencies struct {
	adapter *webauth.Adapter
	login   identityusecase.Login
	logout  identityusecase.Logout
}

type identityDependencies struct {
	logoutAll                identityusecase.LogoutAll
	revokeSession            identityusecase.RevokeSession
	viewSessions             identityusecase.ViewSessions
	listUsers                identityusecase.ListUsers
	listPermissions          identityusecase.ListPermissions
	listRoles                identityusecase.ListRoles
	viewUserRoles            identityusecase.ViewUserRoles
	viewRolePermissions      identityusecase.ViewRolePermissions
	assignRoleToUser         identityusecase.AssignRoleToUser
	removeRoleFromUser       identityusecase.RemoveRoleFromUser
	grantPermissionToRole    identityusecase.GrantPermissionToRole
	revokePermissionFromRole identityusecase.RevokePermissionFromRole
	createRole               identityusecase.CreateRole
	updateRole               identityusecase.UpdateRole
	deleteRole               identityusecase.DeleteRole
	createPermission         identityusecase.CreatePermission
	updatePermission         identityusecase.UpdatePermission
	deletePermission         identityusecase.DeletePermission
}

type companyDependencies struct {
	create   companyusecase.CreateCompany
	list     companyusecase.ListCompanies
	update   companyusecase.UpdateCompany
	delete   companyusecase.DeleteCompany
	findByID companyusecase.FindByID
}
