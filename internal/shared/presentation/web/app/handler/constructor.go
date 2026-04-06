package handler

import (
	auditadapter "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/audit"
	webauth "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/auth"
	ratelimit "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/ratelimit"
	requestmeta "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/requestmeta"
	securityheaders "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/securityheaders"
	accessguard "github.com/paladignus/actajus/internal/shared/presentation/web/app/handler/access"
	supporthandler "github.com/paladignus/actajus/internal/shared/presentation/web/app/handler/support"
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
	navigator := navigation.New(dep.SecureCookies)
	authAdapter := webauth.New(dep.Auth.ValidateAccess, dep.Auth.Refresh, dep.Auth.RBACChecker, dep.SecureCookies, requestmeta.ClientIP, requestmeta.BaseURL)
	audit := auditadapter.New(dep.Logger, requestmeta.ClientIP)
	loginLimiter := ratelimit.NewLoginShield(requestmeta.ClientIP)
	headers := securityheaders.New(dep.Renderer.DevServerURL)
	html := webhtml.New(pagepresenter.New(renderer), headers)
	support := supporthandler.Service{
		AssetsHandler: assets,
		HTML:          html,
		Headers:       headers,
	}
	guard := accessguard.Guard{
		Adapter:            authAdapter,
		Navigator:          navigator,
		CookieAccessToken:  cookieAccessToken,
		CookieRefreshToken: cookieRefreshToken,
		CookieSessionID:    cookieSessionID,
		ManageSessionsPerm: permManageSessions,
		ManageRBACPerm:     permManageRBAC,
		RespondForbidden:   support.RespondForbidden,
	}
	support.AuthFromRequest = guard.AuthFromRequest
	return AppHandler{
		logger:    dep.Logger,
		assets:    assets,
		html:      html,
		navigator: navigator,
		guard:     guard,
		support:   support,
		auth: authDependencies{
			adapter:      authAdapter,
			audit:        audit,
			login:        dep.Auth.Login,
			registerUC:   dep.Auth.Register,
			resendVerify: dep.Auth.RequestEmailVerification,
			verifyEmail:  dep.Auth.ConfirmEmailVerification,
			logout:       dep.Auth.Logout,
			loginLimiter: loginLimiter,
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
