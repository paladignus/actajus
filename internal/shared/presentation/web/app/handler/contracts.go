package handler

import (
	companyusecase "github.com/paladignus/actajus/internal/module/company/application/usecase"
	identityusecase "github.com/paladignus/actajus/internal/module/identity/application/usecase"
	sharedRepo "github.com/paladignus/actajus/internal/shared/application/repository"
	webtemplate "github.com/paladignus/actajus/internal/shared/presentation/web/template"
)

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
