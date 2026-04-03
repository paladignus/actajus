package identitypage

import (
	"strings"
	"time"

	identityread "github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	webflash "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/flash"
	identityviewmodel "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/viewmodel/identity"
	webtemplate "github.com/paladignus/actajus/internal/shared/presentation/web/template"
)

func SessionsPage(viewerID int64, canManageSessions bool, adminFilter string, flash webflash.Message, own []identityread.SessionReadModel, admin []identityread.SessionReadModel, currentSessionID int64, now time.Time) webtemplate.Page {
	return webtemplate.Page{
		Title:       "Sessões",
		Description: "Gerenciamento de sessões autenticadas.",
		NavKey:      "identity.sessions",
		Entry:       "src/main.ts",
		Data:        identityviewmodel.BuildSessionsPage(viewerID, canManageSessions, adminFilter, flash, own, admin, currentSessionID, now),
	}
}

func UsersPage(items []identityread.UserListItemReadModel, filter string, canManageRBAC bool, roles []identityread.RoleListItemReadModel, assignments map[int64][]identityread.UserRoleAssignmentReadModel, flash webflash.Message) webtemplate.Page {
	return webtemplate.Page{
		Title:       "Usuarios",
		Description: "Catalogo WEB de usuarios autenticaveis.",
		NavKey:      "identity.users",
		Entry:       "src/main.ts",
		Data:        identityviewmodel.BuildUsersPage(items, filter, canManageRBAC, roles, assignments, flash),
	}
}

func PermissionsPage(items []identityread.PermissionListItemReadModel, filter string, canManageRBAC bool, roles []identityread.RoleListItemReadModel, assignments map[int16][]identityread.RolePermissionAssignmentReadModel, editID int16, resource, action, description string, flash webflash.Message) webtemplate.Page {
	form := identityviewmodel.PermissionForm{
		ID:          editID,
		Resource:    strings.TrimSpace(resource),
		Action:      strings.TrimSpace(action),
		Description: strings.TrimSpace(description),
	}
	return webtemplate.Page{
		Title:       "Permissoes",
		Description: "Catalogo WEB de permissoes RBAC.",
		NavKey:      "identity.permissions",
		Entry:       "src/main.ts",
		Data:        identityviewmodel.BuildPermissionsPage(items, filter, canManageRBAC, roles, assignments, form, flash),
	}
}

func RolesPage(items []identityread.RoleListItemReadModel, canManageRBAC bool, editID int16, name, description string, flash webflash.Message) webtemplate.Page {
	form := identityviewmodel.RoleForm{
		ID:          editID,
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
	}
	return webtemplate.Page{
		Title:       "Roles",
		Description: "Catalogo WEB de roles RBAC.",
		NavKey:      "identity.roles",
		Entry:       "src/main.ts",
		Data:        identityviewmodel.BuildRolesPage(items, canManageRBAC, form, flash),
	}
}
