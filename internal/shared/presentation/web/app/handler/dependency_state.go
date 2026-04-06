package handler

import (
	"net/http"
	"time"

	companyusecase "github.com/paladignus/actajus/internal/module/company/application/usecase"
	identityusecase "github.com/paladignus/actajus/internal/module/identity/application/usecase"
	auditadapter "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/audit"
	webauth "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/auth"
)

type loginLimiter interface {
	Check(*http.Request, string) (bool, time.Duration)
	RecordFailure(*http.Request, string)
	RecordSuccess(*http.Request, string)
}

type authDependencies struct {
	adapter      *webauth.Adapter
	audit        *auditadapter.Adapter
	login        identityusecase.Login
	registerUC   identityusecase.Register
	resendVerify identityusecase.RequestEmailVerification
	verifyEmail  identityusecase.ConfirmEmailVerification
	logout       identityusecase.Logout
	loginLimiter loginLimiter
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
