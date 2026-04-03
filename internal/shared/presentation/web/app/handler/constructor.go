package handler

import (
	webauth "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/auth"
	requestmeta "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/requestmeta"
	pagepresenter "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/page"
	webhtml "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/html"
	navigation "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/http/navigation"
	webtemplate "github.com/paladignus/actajus/internal/shared/presentation/web/template"
)

func NewAppHandler(dep Dependencies) (AppHandler, error) {
	renderer, err := webtemplate.NewRenderer(dep.Renderer)
	if err != nil {
		return AppHandler{}, err
	}
	assets, err := renderer.AssetsHandler(dep.Renderer.AssetRoot)
	if err != nil {
		return AppHandler{}, err
	}
	return AppHandler{
		logger:    dep.Logger,
		renderer:  renderer,
		assets:    assets,
		html:      webhtml.New(pagepresenter.New(renderer)),
		navigator: navigation.New(dep.SecureCookies),
		auth: authDependencies{
			adapter: webauth.New(dep.Auth.ValidateAccess, dep.Auth.Refresh, dep.Auth.RBACChecker, dep.SecureCookies, requestmeta.ClientIP, requestmeta.BaseURL),
			login:   dep.Auth.Login,
			logout:  dep.Auth.Logout,
		},
		identity: identityDependencies{
			logoutAll:                dep.Identity.LogoutAll,
			revokeSession:            dep.Identity.RevokeSession,
			viewSessions:             dep.Identity.ViewSessions,
			listUsers:                dep.Identity.ListUsers,
			listPermissions:          dep.Identity.ListPermissions,
			listRoles:                dep.Identity.ListRoles,
			viewUserRoles:            dep.Identity.ViewUserRoles,
			viewRolePermissions:      dep.Identity.ViewRolePermissions,
			assignRoleToUser:         dep.Identity.AssignRoleToUser,
			removeRoleFromUser:       dep.Identity.RemoveRoleFromUser,
			grantPermissionToRole:    dep.Identity.GrantPermissionToRole,
			revokePermissionFromRole: dep.Identity.RevokePermissionFromRole,
			createRole:               dep.Identity.CreateRole,
			updateRole:               dep.Identity.UpdateRole,
			deleteRole:               dep.Identity.DeleteRole,
			createPermission:         dep.Identity.CreatePermission,
			updatePermission:         dep.Identity.UpdatePermission,
			deletePermission:         dep.Identity.DeletePermission,
		},
		company: companyDependencies{
			create:   dep.Company.Create,
			list:     dep.Company.List,
			update:   dep.Company.Update,
			delete:   dep.Company.Delete,
			findByID: dep.Company.FindByID,
		},
	}, nil
}
