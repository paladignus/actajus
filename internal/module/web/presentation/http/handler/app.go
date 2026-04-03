// Package handler contains WEB HTTP handlers.
package handler

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	companycmd "github.com/paladignus/actajus/internal/module/company/application/command"
	companyread "github.com/paladignus/actajus/internal/module/company/application/readmodel"
	companyrepo "github.com/paladignus/actajus/internal/module/company/application/repository"
	companyusecase "github.com/paladignus/actajus/internal/module/company/application/usecase"
	identitycmd "github.com/paladignus/actajus/internal/module/identity/application/command"
	identityread "github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	identityrepo "github.com/paladignus/actajus/internal/module/identity/application/repository"
	identityusecase "github.com/paladignus/actajus/internal/module/identity/application/usecase"
	identitydomain "github.com/paladignus/actajus/internal/module/identity/domain"
	sharedRepo "github.com/paladignus/actajus/internal/shared/application/repository"
	sharedhttp "github.com/paladignus/actajus/internal/shared/infrastructure/http/handler"
	webtemplate "github.com/paladignus/actajus/internal/shared/presentation/web/template"
)

//go:embed content/templates/*.gohtml content/dist/.vite/manifest.json content/dist/assets/*
var contentFS embed.FS

const (
	cookieAccessToken  = "actajus_access_token"
	cookieRefreshToken = "actajus_refresh_token"
	cookieSessionID    = "actajus_session_id"
	cookieFlash        = "actajus_flash"
	permManageSessions = "session:manage"
	permManageRBAC     = "rbac:manage"
)

type AppHandler struct {
	logger                   sharedRepo.Logger
	renderer                 *webtemplate.Renderer
	assets                   http.Handler
	validateAccess           identityusecase.ValidateAccess
	login                    identityusecase.Login
	refresh                  identityusecase.Refresh
	logout                   identityusecase.Logout
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
	createCompany            companyusecase.CreateCompany
	listCompanies            companyusecase.ListCompanies
	updateCompany            companyusecase.UpdateCompany
	deleteCompany            companyusecase.DeleteCompany
	findCompanyByID          companyusecase.FindByID
	permissionChecker        sharedRepo.PermissionChecker
	secureCookies            bool
}

type authState struct {
	Claims            *identityread.AccessTokenClaims
	CanManageSessions bool
	CanManageRBAC     bool
}

type sessionItemView struct {
	IDSession int64
	IDUser    int64
	UserEmail string
	IP        string
	UserAgent string
	CreatedAt string
	ExpiresAt string
	UpdatedAt string
	Status    string
	IsCurrent bool
}

type sessionsPageData struct {
	ViewerID          int64
	CanManageSessions bool
	OwnSessions       []sessionItemView
	AdminSessions     []sessionItemView
	AdminFilter       string
	Error             string
	Flash             flashMessage
}

type companiesPageData struct {
	Page            int
	Limit           int
	RangeStart      int
	RangeEnd        int
	Items           []companyListItemView
	NextURL         string
	PreviousURL     string
	HasNextPage     bool
	HasPreviousPage bool
	FilterName      string
	FilterCNPJ      string
	Flash           flashMessage
}

type usersPageData struct {
	Filter        string
	CanManageRBAC bool
	Roles         []roleOptionView
	Flash         flashMessage
	Items         []userListItemView
}

type userListItemView struct {
	ID          int64
	FullName    string
	Email       string
	Status      string
	LastLoginAt string
	Roles       []userRoleView
}

type permissionsPageData struct {
	Filter        string
	CanManageRBAC bool
	Roles         []rolePermissionBlockView
	Permissions   []permissionListItemView
	Form          permissionFormView
	Flash         flashMessage
}

type permissionListItemView struct {
	ID          int16
	Resource    string
	Action      string
	Code        string
	Description string
	CreatedAt   string
}

type roleOptionView struct {
	ID   int16
	Name string
}

type userRoleView struct {
	ID         int16
	Name       string
	AssignedBy int64
}

type rolePermissionBlockView struct {
	ID          int16
	Name        string
	Description string
	Permissions []rolePermissionView
}

type rolePermissionView struct {
	ID          int16
	Code        string
	Description string
}

type rolesPageData struct {
	CanManageRBAC bool
	Roles         []roleListItemView
	Form          roleFormView
	Flash         flashMessage
}

type roleListItemView struct {
	ID          int16
	Name        string
	Description string
	CreatedAt   string
}

type roleFormView struct {
	ID          int16
	Name        string
	Description string
}

type permissionFormView struct {
	ID          int16
	Resource    string
	Action      string
	Description string
}

type flashMessage struct {
	Code    string
	Kind    string
	Message string
	Source  string
}

type companyListItemView struct {
	ID               int64
	Name             string
	TradeName        string
	CNPJ             string
	RegisteredByName string
	City             string
	State            string
}

type companyDetailPageData struct {
	Company *companyread.CompanyReadModel
}

type companyFormData struct {
	Mode            string
	Action          string
	SubmitLabel     string
	Error           string
	CompanyID       int64
	Name            string
	TradeName       string
	CNPJ            string
	AddressID       int64
	ZIP             string
	Title           string
	Street          string
	Number          string
	Complement      string
	Reference       string
	Neighborhood    string
	City            string
	State           string
	Country         string
	PhoneID         int64
	PhoneNumber     string
	PhoneKind       string
	PhoneDepartment string
	EmailID         int64
	EmailAddress    string
	SocialMedias    []companySocialMediaFormItem
}

type companySocialMediaFormItem struct {
	ID       int64
	Platform string
	URL      string
}

func NewAppHandler(
	logger sharedRepo.Logger,
	validateAccess identityusecase.ValidateAccess,
	login identityusecase.Login,
	refresh identityusecase.Refresh,
	logout identityusecase.Logout,
	logoutAll identityusecase.LogoutAll,
	revokeSession identityusecase.RevokeSession,
	viewSessions identityusecase.ViewSessions,
	listUsers identityusecase.ListUsers,
	listPermissions identityusecase.ListPermissions,
	listRoles identityusecase.ListRoles,
	viewUserRoles identityusecase.ViewUserRoles,
	viewRolePermissions identityusecase.ViewRolePermissions,
	assignRoleToUser identityusecase.AssignRoleToUser,
	removeRoleFromUser identityusecase.RemoveRoleFromUser,
	grantPermissionToRole identityusecase.GrantPermissionToRole,
	revokePermissionFromRole identityusecase.RevokePermissionFromRole,
	createRole identityusecase.CreateRole,
	updateRole identityusecase.UpdateRole,
	deleteRole identityusecase.DeleteRole,
	createPermission identityusecase.CreatePermission,
	updatePermission identityusecase.UpdatePermission,
	deletePermission identityusecase.DeletePermission,
	createCompany companyusecase.CreateCompany,
	listCompanies companyusecase.ListCompanies,
	updateCompany companyusecase.UpdateCompany,
	deleteCompany companyusecase.DeleteCompany,
	findCompanyByID companyusecase.FindByID,
	permissionChecker sharedRepo.PermissionChecker,
	secureCookies bool,
) (AppHandler, error) {
	renderer, err := webtemplate.NewRenderer(webtemplate.Source{
		TemplateFS:       contentFS,
		TemplatePatterns: []string{"content/templates/*.gohtml"},
		AssetFS:          contentFS,
		ManifestPath:     "content/dist/.vite/manifest.json",
		PublicPath:       "/assets",
	})
	if err != nil {
		return AppHandler{}, err
	}
	assets, err := renderer.AssetsHandler("content/dist")
	if err != nil {
		return AppHandler{}, err
	}
	return AppHandler{
		logger:                   logger,
		renderer:                 renderer,
		assets:                   assets,
		validateAccess:           validateAccess,
		login:                    login,
		refresh:                  refresh,
		logout:                   logout,
		logoutAll:                logoutAll,
		revokeSession:            revokeSession,
		viewSessions:             viewSessions,
		listUsers:                listUsers,
		listPermissions:          listPermissions,
		listRoles:                listRoles,
		viewUserRoles:            viewUserRoles,
		viewRolePermissions:      viewRolePermissions,
		assignRoleToUser:         assignRoleToUser,
		removeRoleFromUser:       removeRoleFromUser,
		grantPermissionToRole:    grantPermissionToRole,
		revokePermissionFromRole: revokePermissionFromRole,
		createRole:               createRole,
		updateRole:               updateRole,
		deleteRole:               deleteRole,
		createPermission:         createPermission,
		updatePermission:         updatePermission,
		deletePermission:         deletePermission,
		createCompany:            createCompany,
		listCompanies:            listCompanies,
		updateCompany:            updateCompany,
		deleteCompany:            deleteCompany,
		findCompanyByID:          findCompanyByID,
		permissionChecker:        permissionChecker,
		secureCookies:            secureCookies,
	}, nil
}

func (h AppHandler) Home() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := h.authFromRequest(w, r)
		if err == nil && state != nil && state.Claims != nil {
			http.Redirect(w, r, "/sessions", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
}

func (h AppHandler) LoginPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if state, err := h.authFromRequest(w, r); err == nil && state != nil && state.Claims != nil {
			http.Redirect(w, r, "/sessions", http.StatusSeeOther)
			return
		}
		_ = h.renderLogin(w, "", "")
	})
}

func (h AppHandler) LoginAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}
		email := strings.TrimSpace(r.FormValue("email"))
		input := identitycmd.LoginCommand{
			Email:     email,
			Password:  r.FormValue("password"),
			IP:        clientIP(r),
			UserAgent: r.UserAgent(),
		}
		tokens, err := h.login.Execute(r.Context(), input)
		if err != nil {
			status, message := mapLoginError(err)
			_ = h.renderLoginStatus(w, status, email, message)
			return
		}
		h.setAuthCookies(w, tokens)
		http.Redirect(w, r, "/sessions", http.StatusSeeOther)
	})
}

func (h AppHandler) LogoutAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := h.requireAuth(w, r)
		if err != nil {
			return
		}
		_ = h.logout.Execute(r.Context(), identitycmd.LogoutCommand{IDSession: state.Claims.IDSession})
		h.clearAuthCookies(w)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
}

func (h AppHandler) SessionsPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := h.requireAuth(w, r)
		if err != nil {
			return
		}
		data, status := h.buildSessionsPageData(w, r, state)
		if err := h.renderer.RenderHTTP(w, status, "pages/sessions", webtemplate.Page{
			Title:       "Sessões",
			Description: "Gerenciamento de sessões autenticadas.",
			NavKey:      "identity.sessions",
			Entry:       "src/main.ts",
			Bootstrap: map[string]any{
				"appName":           "Actajus",
				"apiBaseURL":        baseURL(r),
				"viewerID":          state.Claims.IDUser,
				"canManageSessions": state.CanManageSessions,
			},
			Data: data,
		}); err != nil {
			h.logger.Error(r.Context(), "render sessions page", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	})
}

func (h AppHandler) UsersPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := h.requireAuth(w, r)
		if err != nil {
			return
		}
		filter := identityrepo.CatalogUserFilter{
			Query: strings.TrimSpace(r.URL.Query().Get("q")),
			Limit: 50,
		}
		items, err := h.listUsers.Execute(r.Context(), filter)
		if err != nil {
			h.logger.Error(r.Context(), "list users", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		data, err := h.buildUsersPageData(w, r, items, filter.Query, state)
		if err != nil {
			h.logger.Error(r.Context(), "build users page data", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if err := h.renderer.RenderHTTP(w, http.StatusOK, "pages/users", webtemplate.Page{
			Title:       "Usuarios",
			Description: "Catalogo WEB de usuarios autenticaveis.",
			NavKey:      "identity.users",
			Entry:       "src/main.ts",
			Data:        data,
		}); err != nil {
			h.logger.Error(r.Context(), "render users page", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	})
}

func (h AppHandler) PermissionsPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := h.requireAuth(w, r)
		if err != nil {
			return
		}
		filter := identityrepo.CatalogPermissionFilter{
			Query: strings.TrimSpace(r.URL.Query().Get("q")),
			Limit: 100,
		}
		items, err := h.listPermissions.Execute(r.Context(), filter)
		if err != nil {
			h.logger.Error(r.Context(), "list permissions", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		data, err := h.buildPermissionsPageData(w, r, items, filter.Query, state)
		if err != nil {
			h.logger.Error(r.Context(), "build permissions page data", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if err := h.renderer.RenderHTTP(w, http.StatusOK, "pages/permissions", webtemplate.Page{
			Title:       "Permissoes",
			Description: "Catalogo WEB de permissoes RBAC.",
			NavKey:      "identity.permissions",
			Entry:       "src/main.ts",
			Data:        data,
		}); err != nil {
			h.logger.Error(r.Context(), "render permissions page", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	})
}

func (h AppHandler) RolesPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := h.requireAuth(w, r)
		if err != nil {
			return
		}
		roles, err := h.listRoles.Execute(r.Context())
		if err != nil {
			h.logger.Error(r.Context(), "list roles", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		data := h.buildRolesPageData(w, r, roles, state.CanManageRBAC)
		if err := h.renderer.RenderHTTP(w, http.StatusOK, "pages/roles", webtemplate.Page{
			Title:       "Roles",
			Description: "Catalogo WEB de roles RBAC.",
			NavKey:      "identity.roles",
			Entry:       "src/main.ts",
			Data:        data,
		}); err != nil {
			h.logger.Error(r.Context(), "render roles page", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	})
}

func (h AppHandler) CreateRoleAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.requireRBACAdmin(w, r); err != nil {
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}
		err := h.createRole.Execute(r.Context(), identitycmd.CreateRoleCommand{
			Name:        strings.TrimSpace(r.FormValue("name")),
			Description: strings.TrimSpace(r.FormValue("description")),
		})
		if err != nil {
			h.logger.Error(r.Context(), "create role", "error", err)
			h.redirectWithTypedFlash(w, r, "/roles", "roles.create", "roles_create_failed", "error", mapRBACCatalogError(err, "Nao foi possivel criar a role."))
			return
		}
		h.redirectWithTypedFlash(w, r, "/roles", "roles.create", "roles_create_succeeded", "success", "Role criada com sucesso.")
	})
}

func (h AppHandler) UpdateRoleAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.requireRBACAdmin(w, r); err != nil {
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}
		id, err := parseInt16Field(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid role id", http.StatusBadRequest)
			return
		}
		err = h.updateRole.Execute(r.Context(), identitycmd.UpdateRoleCommand{
			ID:          id,
			Name:        strings.TrimSpace(r.FormValue("name")),
			Description: strings.TrimSpace(r.FormValue("description")),
		})
		if err != nil {
			h.logger.Error(r.Context(), "update role", "error", err, "id_role", id)
			h.redirectWithTypedFlash(w, r, "/roles", "roles.update", "roles_update_failed", "error", mapRBACCatalogError(err, "Nao foi possivel atualizar a role."))
			return
		}
		h.redirectWithTypedFlash(w, r, "/roles", "roles.update", "roles_update_succeeded", "success", "Role atualizada com sucesso.")
	})
}

func (h AppHandler) DeleteRoleAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.requireRBACAdmin(w, r); err != nil {
			return
		}
		id, err := parseInt16Field(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid role id", http.StatusBadRequest)
			return
		}
		if err := h.deleteRole.Execute(r.Context(), identitycmd.DeleteRoleCommand{ID: id}); err != nil {
			h.logger.Error(r.Context(), "delete role", "error", err, "id_role", id)
			h.redirectWithTypedFlash(w, r, "/roles", "roles.delete", "roles_delete_failed", "error", mapRBACCatalogError(err, "Nao foi possivel excluir a role."))
			return
		}
		h.redirectWithTypedFlash(w, r, "/roles", "roles.delete", "roles_delete_succeeded", "success", "Role excluida com sucesso.")
	})
}

func (h AppHandler) CreatePermissionAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.requireRBACAdmin(w, r); err != nil {
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}
		err := h.createPermission.Execute(r.Context(), identitycmd.CreatePermissionCommand{
			Resource:    strings.TrimSpace(r.FormValue("resource")),
			Action:      strings.TrimSpace(r.FormValue("action")),
			Description: strings.TrimSpace(r.FormValue("description")),
		})
		if err != nil {
			h.logger.Error(r.Context(), "create permission", "error", err)
			h.redirectWithTypedFlash(w, r, "/permissions", "permissions.create", "permissions_create_failed", "error", mapRBACCatalogError(err, "Nao foi possivel criar a permissao."))
			return
		}
		h.redirectWithTypedFlash(w, r, "/permissions", "permissions.create", "permissions_create_succeeded", "success", "Permissao criada com sucesso.")
	})
}

func (h AppHandler) UpdatePermissionAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.requireRBACAdmin(w, r); err != nil {
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}
		id, err := parseInt16Field(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid permission id", http.StatusBadRequest)
			return
		}
		err = h.updatePermission.Execute(r.Context(), identitycmd.UpdatePermissionCommand{
			ID:          id,
			Resource:    strings.TrimSpace(r.FormValue("resource")),
			Action:      strings.TrimSpace(r.FormValue("action")),
			Description: strings.TrimSpace(r.FormValue("description")),
		})
		if err != nil {
			h.logger.Error(r.Context(), "update permission", "error", err, "id_permission", id)
			h.redirectWithTypedFlash(w, r, "/permissions", "permissions.update", "permissions_update_failed", "error", mapRBACCatalogError(err, "Nao foi possivel atualizar a permissao."))
			return
		}
		h.redirectWithTypedFlash(w, r, "/permissions", "permissions.update", "permissions_update_succeeded", "success", "Permissao atualizada com sucesso.")
	})
}

func (h AppHandler) DeletePermissionAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.requireRBACAdmin(w, r); err != nil {
			return
		}
		id, err := parseInt16Field(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid permission id", http.StatusBadRequest)
			return
		}
		if err := h.deletePermission.Execute(r.Context(), identitycmd.DeletePermissionCommand{ID: id}); err != nil {
			h.logger.Error(r.Context(), "delete permission", "error", err, "id_permission", id)
			h.redirectWithTypedFlash(w, r, "/permissions", "permissions.delete", "permissions_delete_failed", "error", mapRBACCatalogError(err, "Nao foi possivel excluir a permissao."))
			return
		}
		h.redirectWithTypedFlash(w, r, "/permissions", "permissions.delete", "permissions_delete_succeeded", "success", "Permissao excluida com sucesso.")
	})
}

func (h AppHandler) AssignRoleToUserAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := h.requireRBACAdmin(w, r)
		if err != nil {
			return
		}
		if err := r.ParseForm(); err != nil {
			h.redirectWithFlash(w, r, "/users", "error", "Formulario invalido.")
			return
		}
		idUser, err := parseInt64Field(r.FormValue("id_user"))
		if err != nil {
			h.redirectWithFlash(w, r, "/users", "error", "Usuario invalido.")
			return
		}
		idRole, err := parseInt16Field(r.FormValue("id_role"))
		if err != nil {
			h.redirectWithFlash(w, r, "/users", "error", "Role invalida.")
			return
		}
		if err := h.assignRoleToUser.Execute(r.Context(), identitycmd.AssignRoleToUserCommand{
			IDUser:     idUser,
			IDRole:     idRole,
			AssignedBy: state.Claims.IDUser,
		}); err != nil {
			h.logger.Error(r.Context(), "assign role to user", "error", err, "id_user", idUser, "id_role", idRole)
			h.redirectWithTypedFlash(w, r, "/users", "users.roles.assign", "users_roles_assign_failed", "error", mapRBACCatalogError(err, "Nao foi possivel atribuir a role."))
			return
		}
		h.redirectWithTypedFlash(w, r, "/users", "users.roles.assign", "users_roles_assign_succeeded", "success", "Role atribuida com sucesso.")
	})
}

func (h AppHandler) RemoveRoleFromUserAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.requireRBACAdmin(w, r); err != nil {
			return
		}
		if err := r.ParseForm(); err != nil {
			h.redirectWithFlash(w, r, "/users", "error", "Formulario invalido.")
			return
		}
		idUser, err := parseInt64Field(r.FormValue("id_user"))
		if err != nil {
			h.redirectWithFlash(w, r, "/users", "error", "Usuario invalido.")
			return
		}
		idRole, err := parseInt16Field(r.FormValue("id_role"))
		if err != nil {
			h.redirectWithFlash(w, r, "/users", "error", "Role invalida.")
			return
		}
		if err := h.removeRoleFromUser.Execute(r.Context(), identitycmd.RemoveRoleFromUserCommand{
			IDUser: idUser,
			IDRole: idRole,
		}); err != nil {
			h.logger.Error(r.Context(), "remove role from user", "error", err, "id_user", idUser, "id_role", idRole)
			h.redirectWithTypedFlash(w, r, "/users", "users.roles.remove", "users_roles_remove_failed", "error", mapRBACCatalogError(err, "Nao foi possivel remover a role."))
			return
		}
		h.redirectWithTypedFlash(w, r, "/users", "users.roles.remove", "users_roles_remove_succeeded", "success", "Role removida com sucesso.")
	})
}

func (h AppHandler) GrantPermissionToRoleAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.requireRBACAdmin(w, r); err != nil {
			return
		}
		if err := r.ParseForm(); err != nil {
			h.redirectWithFlash(w, r, "/permissions", "error", "Formulario invalido.")
			return
		}
		idRole, err := parseInt16Field(r.FormValue("id_role"))
		if err != nil {
			h.redirectWithFlash(w, r, "/permissions", "error", "Role invalida.")
			return
		}
		idPermission, err := parseInt16Field(r.FormValue("id_permission"))
		if err != nil {
			h.redirectWithFlash(w, r, "/permissions", "error", "Permissao invalida.")
			return
		}
		if err := h.grantPermissionToRole.Execute(r.Context(), identitycmd.GrantPermissionToRoleCommand{
			IDRole:       idRole,
			IDPermission: idPermission,
		}); err != nil {
			h.logger.Error(r.Context(), "grant permission to role", "error", err, "id_role", idRole, "id_permission", idPermission)
			h.redirectWithTypedFlash(w, r, "/permissions", "roles.permissions.grant", "roles_permissions_grant_failed", "error", mapRBACCatalogError(err, "Nao foi possivel conceder a permissao."))
			return
		}
		h.redirectWithTypedFlash(w, r, "/permissions", "roles.permissions.grant", "roles_permissions_grant_succeeded", "success", "Permissao concedida com sucesso.")
	})
}

func (h AppHandler) RevokePermissionFromRoleAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.requireRBACAdmin(w, r); err != nil {
			return
		}
		if err := r.ParseForm(); err != nil {
			h.redirectWithFlash(w, r, "/permissions", "error", "Formulario invalido.")
			return
		}
		idRole, err := parseInt16Field(r.FormValue("id_role"))
		if err != nil {
			h.redirectWithFlash(w, r, "/permissions", "error", "Role invalida.")
			return
		}
		idPermission, err := parseInt16Field(r.FormValue("id_permission"))
		if err != nil {
			h.redirectWithFlash(w, r, "/permissions", "error", "Permissao invalida.")
			return
		}
		if err := h.revokePermissionFromRole.Execute(r.Context(), identitycmd.RevokePermissionFromRoleCommand{
			IDRole:       idRole,
			IDPermission: idPermission,
		}); err != nil {
			h.logger.Error(r.Context(), "revoke permission from role", "error", err, "id_role", idRole, "id_permission", idPermission)
			h.redirectWithTypedFlash(w, r, "/permissions", "roles.permissions.revoke", "roles_permissions_revoke_failed", "error", mapRBACCatalogError(err, "Nao foi possivel revogar a permissao."))
			return
		}
		h.redirectWithTypedFlash(w, r, "/permissions", "roles.permissions.revoke", "roles_permissions_revoke_succeeded", "success", "Permissao revogada com sucesso.")
	})
}

func (h AppHandler) CompaniesPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.requireAuth(w, r); err != nil {
			return
		}
		page := parsePage(r.URL.Query().Get("page"))
		limit := parseLimit(r.URL.Query().Get("limit"))
		filter := companyrepo.CompanyListFilter{
			Name: strings.TrimSpace(r.URL.Query().Get("name")),
			CNPJ: strings.TrimSpace(r.URL.Query().Get("cnpj")),
		}
		after := optionalString(r.URL.Query().Get("after"))
		before := optionalString(r.URL.Query().Get("before"))
		rm, err := h.listCompanies.Execute(r.Context(), filter, after, before, limit, baseURL(r)+"/companies")
		if err != nil {
			h.logger.Error(r.Context(), "list companies", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		data := h.buildCompaniesPageData(w, r, rm, filter, page, limit)
		if err := h.renderer.RenderHTTP(w, http.StatusOK, "pages/companies", webtemplate.Page{
			Title:       "Companies",
			Description: "Listagem WEB de companies.",
			NavKey:      "companies.list",
			Entry:       "src/main.ts",
			Data:        data,
		}); err != nil {
			h.logger.Error(r.Context(), "render companies page", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	})
}

func (h AppHandler) CompanyNewPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.requireAuth(w, r); err != nil {
			return
		}
		_ = h.renderCompanyForm(w, http.StatusOK, newCompanyFormData(), "pages/company-form")
	})
}

func (h AppHandler) CompanyCreateAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := h.requireAuth(w, r)
		if err != nil {
			return
		}
		form, cmd, err := parseCreateCompanyForm(r, state.Claims.IDUser)
		if err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}
		company, execErr := h.createCompany.Execute(r.Context(), cmd)
		if execErr != nil {
			form.Error = execErr.Error()
			_ = h.renderCompanyForm(w, http.StatusBadRequest, form, "pages/company-form")
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/companies/%d", company.ID), http.StatusSeeOther)
	})
}

func (h AppHandler) CompanyDetailPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.requireAuth(w, r); err != nil {
			return
		}
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil || id <= 0 {
			h.respondNotFound(w, r)
			return
		}
		company, err := h.findCompanyByID.Execute(r.Context(), id)
		if err != nil {
			h.logger.Error(r.Context(), "find company by id", "error", err, "id_company", id)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if company == nil {
			h.respondNotFound(w, r)
			return
		}
		if err := h.renderer.RenderHTTP(w, http.StatusOK, "pages/company-detail", webtemplate.Page{
			Title:       "Company",
			Description: "Detalhe WEB de company.",
			NavKey:      "companies.list",
			Entry:       "src/main.ts",
			Data:        companyDetailPageData{Company: company},
		}); err != nil {
			h.logger.Error(r.Context(), "render company detail", "error", err, "id_company", id)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	})
}

func (h AppHandler) CompanyEditPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.requireAuth(w, r); err != nil {
			return
		}
		company, ok := h.loadCompanyForEdit(w, r)
		if !ok {
			return
		}
		_ = h.renderCompanyForm(w, http.StatusOK, companyToFormData(company), "pages/company-form")
	})
}

func (h AppHandler) CompanyUpdateAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.requireAuth(w, r); err != nil {
			return
		}
		company, ok := h.loadCompanyForEdit(w, r)
		if !ok {
			return
		}
		form, cmd, err := parseUpdateCompanyForm(r, company)
		if err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}
		updated, execErr := h.updateCompany.Execute(r.Context(), cmd)
		if execErr != nil {
			form.Error = execErr.Error()
			_ = h.renderCompanyForm(w, http.StatusBadRequest, form, "pages/company-form")
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/companies/%d", updated.ID), http.StatusSeeOther)
	})
}

func (h AppHandler) CompanyDeleteAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.requireAuth(w, r); err != nil {
			return
		}
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil || id <= 0 {
			h.respondNotFound(w, r)
			return
		}
		if err := h.deleteCompany.Execute(r.Context(), id); err != nil {
			h.logger.Error(r.Context(), "delete company", "error", err, "id_company", id)
			h.redirectWithTypedFlash(w, r, "/companies", "companies.delete", "companies_delete_failed", "error", "Nao foi possivel excluir a company.")
			return
		}
		h.redirectWithTypedFlash(w, r, "/companies", "companies.delete", "companies_delete_succeeded", "success", "Company excluida com sucesso.")
	})
}

func (h AppHandler) RevokeSessionAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := h.requireAuth(w, r)
		if err != nil {
			return
		}
		if err := r.ParseForm(); err != nil {
			h.redirectWithFlash(w, r, "/sessions", "error", "Formulario invalido.")
			return
		}
		idSession, err := strconv.ParseInt(r.FormValue("id_session"), 10, 64)
		if err != nil || idSession <= 0 {
			h.redirectWithFlash(w, r, "/sessions", "error", "Sessao invalida.")
			return
		}
		target, err := h.viewSessions.ByID(r.Context(), idSession)
		if err != nil {
			h.logger.Error(r.Context(), "find session", "error", err, "id_session", idSession)
			h.redirectWithFlash(w, r, "/sessions", "error", "Nao foi possivel localizar a sessao.")
			return
		}
		if target == nil {
			h.redirectWithFlash(w, r, "/sessions", "error", "Sessao nao encontrada.")
			return
		}
		if target.IDUser != state.Claims.IDUser && !state.CanManageSessions {
			h.respondForbidden(w, r)
			return
		}
		if err := h.revokeSession.Execute(r.Context(), identitycmd.RevokeSessionCommand{IDSession: idSession}); err != nil {
			h.logger.Error(r.Context(), "revoke session", "error", err, "id_session", idSession)
			h.redirectWithTypedFlash(w, r, "/sessions", "sessions.revoke", "sessions_revoke_failed", "error", "Nao foi possivel revogar a sessao.")
			return
		}
		if idSession == state.Claims.IDSession {
			h.clearAuthCookies(w)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		h.redirectWithTypedFlash(w, r, "/sessions", "sessions.revoke", "sessions_revoke_succeeded", "success", "Sessao revogada com sucesso.")
	})
}

func (h AppHandler) RevokeAllAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := h.requireAuth(w, r)
		if err != nil {
			return
		}
		if err := r.ParseForm(); err != nil {
			h.redirectWithFlash(w, r, "/sessions", "error", "Formulario invalido.")
			return
		}
		targetUserID := state.Claims.IDUser
		if raw := strings.TrimSpace(r.FormValue("id_user")); raw != "" {
			parsed, err := strconv.ParseInt(raw, 10, 64)
			if err != nil || parsed <= 0 {
				h.redirectWithFlash(w, r, "/sessions", "error", "Usuario invalido.")
				return
			}
			targetUserID = parsed
		}
		if targetUserID != state.Claims.IDUser && !state.CanManageSessions {
			h.respondForbidden(w, r)
			return
		}
		if err := h.logoutAll.Execute(r.Context(), identitycmd.LogoutAllCommand{IDUser: targetUserID}); err != nil {
			h.logger.Error(r.Context(), "revoke all sessions", "error", err, "id_user", targetUserID)
			h.redirectWithTypedFlash(w, r, "/sessions", "sessions.revoke_all", "sessions_revoke_all_failed", "error", "Nao foi possivel encerrar as sessoes.")
			return
		}
		if targetUserID == state.Claims.IDUser {
			h.clearAuthCookies(w)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		h.redirectWithTypedFlash(w, r, "/sessions", "sessions.revoke_all", "sessions_revoke_all_succeeded", "success", "Sessoes encerradas com sucesso.")
	})
}

func (h AppHandler) Bootstrap() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, _ := h.authFromRequest(w, r)
		payload := map[string]any{
			"appName":       "Actajus",
			"apiBaseURL":    baseURL(r),
			"authenticated": state != nil && state.Claims != nil,
		}
		if state != nil && state.Claims != nil {
			payload["idUser"] = state.Claims.IDUser
			payload["idSession"] = state.Claims.IDSession
			payload["canManageSessions"] = state.CanManageSessions
			payload["canManageRBAC"] = state.CanManageRBAC
		}
		if err := sharedhttp.RespondJSON(w, http.StatusOK, payload); err != nil {
			h.logger.Error(r.Context(), "write web bootstrap", "error", err)
		}
	})
}

func (h AppHandler) Assets() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.assets == nil {
			http.NotFound(w, r)
			return
		}
		h.assets.ServeHTTP(w, r)
	})
}

func (h AppHandler) RenderNotFound(w http.ResponseWriter, r *http.Request) {
	h.renderErrorPage(w, r, http.StatusNotFound, "404", "A pagina solicitada nao foi encontrada.")
}

func (h AppHandler) RenderForbidden(w http.ResponseWriter, r *http.Request) {
	h.renderErrorPage(w, r, http.StatusForbidden, "403", "Voce nao tem permissao para acessar esta pagina.")
}

func (h AppHandler) renderLogin(w http.ResponseWriter, email, message string) error {
	return h.renderLoginStatus(w, http.StatusOK, email, message)
}

func (h AppHandler) renderErrorPage(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	primaryHref := "/login"
	primaryLabel := "Ir para login"
	shortcuts := []map[string]string{
		{"Label": "Pagina inicial", "Href": "/"},
		{"Label": "Login", "Href": "/login"},
	}
	if state, err := h.authFromRequest(w, r); err == nil && state != nil && state.Claims != nil {
		primaryHref = "/sessions"
		primaryLabel = "Voltar ao painel"
		shortcuts = []map[string]string{
			{"Label": "Sessoes", "Href": "/sessions"},
			{"Label": "Companies", "Href": "/companies"},
			{"Label": "Usuarios", "Href": "/users"},
			{"Label": "Roles", "Href": "/roles"},
		}
	}
	if err := h.renderer.RenderHTTP(w, status, "pages/error", webtemplate.Page{
		Title:       code + " | Actajus",
		Description: "Erro de navegacao do painel WEB.",
		NavKey:      "login",
		Entry:       "src/main.ts",
		Data: map[string]any{
			"Code":         code,
			"Message":      message,
			"CurrentPath":  r.URL.Path,
			"Breadcrumbs":  buildErrorBreadcrumbs(r.URL.Path),
			"PrimaryHref":  primaryHref,
			"PrimaryLabel": primaryLabel,
			"Shortcuts":    shortcuts,
		},
	}); err != nil {
		http.Error(w, message, status)
	}
}

func buildErrorBreadcrumbs(path string) []map[string]string {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return []map[string]string{{"Label": "inicio", "Href": "/"}}
	}
	parts := strings.Split(trimmed, "/")
	items := make([]map[string]string, 0, len(parts)+1)
	items = append(items, map[string]string{"Label": "inicio", "Href": "/"})
	current := ""
	for _, part := range parts {
		current += "/" + part
		items = append(items, map[string]string{
			"Label": part,
			"Href":  current,
		})
	}
	return items
}

func (h AppHandler) renderLoginStatus(w http.ResponseWriter, status int, email, message string) error {
	return h.renderer.RenderHTTP(w, status, "pages/login", webtemplate.Page{
		Title:       "Entrar",
		Description: "Acesso WEB do Actajus.",
		NavKey:      "login",
		Entry:       "src/main.ts",
		Data: map[string]any{
			"Email": email,
			"Error": message,
		},
	})
}

func (h AppHandler) authFromRequest(w http.ResponseWriter, r *http.Request) (*authState, error) {
	accessCookie, err := r.Cookie(cookieAccessToken)
	if err == nil && strings.TrimSpace(accessCookie.Value) != "" {
		claims, verr := h.validateAccess.Execute(r.Context(), identitycmd.ValidateAccessCommand{
			AccessToken: accessCookie.Value,
		})
		if verr == nil {
			return h.buildAuthState(r, claims)
		}
		if !shouldAttemptRefresh(verr) {
			h.clearAuthCookies(w)
			return nil, verr
		}
	}

	refreshCookie, errRefresh := r.Cookie(cookieRefreshToken)
	sessionCookie, errSession := r.Cookie(cookieSessionID)
	if errRefresh != nil || errSession != nil {
		return nil, errors.New("missing auth cookies")
	}
	sessionID, err := strconv.ParseInt(sessionCookie.Value, 10, 64)
	if err != nil || sessionID <= 0 {
		h.clearAuthCookies(w)
		return nil, errors.New("invalid session cookie")
	}
	tokens, err := h.refresh.Execute(r.Context(), identitycmd.RefreshCommand{
		IDSession:    sessionID,
		RefreshToken: refreshCookie.Value,
		IP:           clientIP(r),
		UserAgent:    r.UserAgent(),
	})
	if err != nil {
		h.clearAuthCookies(w)
		return nil, err
	}
	h.setAuthCookies(w, tokens)
	return h.buildAuthState(r, &identityread.AccessTokenClaims{
		IDUser:    tokens.IDUser,
		IDSession: tokens.IDSession,
	})
}

func (h AppHandler) buildAuthState(r *http.Request, claims *identityread.AccessTokenClaims) (*authState, error) {
	canManage, err := h.permissionChecker.HasPermission(r.Context(), claims.IDUser, permManageSessions)
	if err != nil {
		return nil, err
	}
	canManageRBAC, err := h.permissionChecker.HasPermission(r.Context(), claims.IDUser, permManageRBAC)
	if err != nil {
		return nil, err
	}
	return &authState{
		Claims:            claims,
		CanManageSessions: canManage,
		CanManageRBAC:     canManageRBAC,
	}, nil
}

func (h AppHandler) requireAuth(w http.ResponseWriter, r *http.Request) (*authState, error) {
	state, err := h.authFromRequest(w, r)
	if err != nil || state == nil || state.Claims == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return nil, err
	}
	return state, nil
}

func (h AppHandler) requireRBACAdmin(w http.ResponseWriter, r *http.Request) (*authState, error) {
	state, err := h.requireAuth(w, r)
	if err != nil {
		return nil, err
	}
	if state == nil || !state.CanManageRBAC {
		h.respondForbidden(w, r)
		return nil, errors.New("missing rbac manage permission")
	}
	return state, nil
}

func (h AppHandler) respondForbidden(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.RenderForbidden(w, r)
		return
	}
	http.Error(w, "forbidden", http.StatusForbidden)
}

func (h AppHandler) respondNotFound(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.RenderNotFound(w, r)
		return
	}
	http.NotFound(w, r)
}

func (h AppHandler) buildSessionsPageData(w http.ResponseWriter, r *http.Request, state *authState) (sessionsPageData, int) {
	data := sessionsPageData{
		ViewerID:          state.Claims.IDUser,
		CanManageSessions: state.CanManageSessions,
		AdminFilter:       strings.TrimSpace(r.URL.Query().Get("email")),
		Flash:             h.readFlash(w, r),
	}
	own, err := h.viewSessions.Own(r.Context(), state.Claims.IDUser)
	if err != nil {
		h.logger.Error(r.Context(), "list own sessions", "error", err, "id_user", state.Claims.IDUser)
		data.Error = "Nao foi possivel carregar as suas sessoes."
		return data, http.StatusInternalServerError
	}
	data.OwnSessions = makeSessionViews(own, state.Claims.IDSession)
	if state.CanManageSessions {
		adminSessions, err := h.viewSessions.Admin(r.Context(), identityrepo.SessionQueryFilter{
			UserEmail: data.AdminFilter,
			Limit:     50,
		})
		if err != nil {
			h.logger.Error(r.Context(), "list admin sessions", "error", err, "email", data.AdminFilter)
			data.Error = "Nao foi possivel carregar as sessoes administrativas."
			return data, http.StatusInternalServerError
		}
		data.AdminSessions = makeSessionViews(adminSessions, state.Claims.IDSession)
	}
	return data, http.StatusOK
}

func (h AppHandler) buildUsersPageData(w http.ResponseWriter, r *http.Request, items []identityread.UserListItemReadModel, filter string, state *authState) (usersPageData, error) {
	data := usersPageData{
		Filter:        filter,
		CanManageRBAC: state.CanManageRBAC,
		Flash:         h.readFlash(w, r),
		Items:         make([]userListItemView, 0, len(items)),
	}
	if state.CanManageRBAC {
		roles, err := h.listRoles.Execute(r.Context())
		if err != nil {
			return data, err
		}
		data.Roles = makeRoleOptions(roles)
	}
	for _, item := range items {
		status := "Ativo"
		if item.IsBlocked {
			status = "Bloqueado"
		}
		lastLogin := "Nunca"
		if item.LastLoginAt != nil {
			lastLogin = item.LastLoginAt.Format(time.DateTime)
		}
		view := userListItemView{
			ID:          item.ID,
			FullName:    fallback(item.FullName, "-"),
			Email:       fallback(item.Email, "-"),
			Status:      status,
			LastLoginAt: lastLogin,
		}
		if state.CanManageRBAC {
			assignments, err := h.viewUserRoles.Execute(r.Context(), item.ID)
			if err != nil {
				return data, err
			}
			view.Roles = makeUserRoleViews(assignments)
		}
		data.Items = append(data.Items, view)
	}
	return data, nil
}

func (h AppHandler) buildPermissionsPageData(w http.ResponseWriter, r *http.Request, items []identityread.PermissionListItemReadModel, filter string, state *authState) (permissionsPageData, error) {
	data := permissionsPageData{
		Filter:        filter,
		CanManageRBAC: state.CanManageRBAC,
		Permissions:   make([]permissionListItemView, 0, len(items)),
		Flash:         h.readFlash(w, r),
		Form: permissionFormView{
			ID:          parseOptionalInt16(r.URL.Query().Get("edit")),
			Resource:    strings.TrimSpace(r.URL.Query().Get("resource")),
			Action:      strings.TrimSpace(r.URL.Query().Get("action")),
			Description: strings.TrimSpace(r.URL.Query().Get("description")),
		},
	}
	for _, item := range items {
		data.Permissions = append(data.Permissions, permissionListItemView{
			ID:          item.ID,
			Resource:    item.Resource,
			Action:      item.Action,
			Code:        item.Resource + ":" + item.Action,
			Description: item.Description,
			CreatedAt:   item.CreatedAt.Format(time.DateTime),
		})
	}
	if !state.CanManageRBAC {
		return data, nil
	}
	roles, err := h.listRoles.Execute(r.Context())
	if err != nil {
		return data, err
	}
	data.Roles = make([]rolePermissionBlockView, 0, len(roles))
	for _, role := range roles {
		assignments, err := h.viewRolePermissions.Execute(r.Context(), role.ID)
		if err != nil {
			return data, err
		}
		data.Roles = append(data.Roles, rolePermissionBlockView{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Permissions: makeRolePermissionViews(assignments),
		})
	}
	return data, nil
}

func (h AppHandler) buildRolesPageData(w http.ResponseWriter, r *http.Request, items []identityread.RoleListItemReadModel, canManageRBAC bool) rolesPageData {
	data := rolesPageData{
		CanManageRBAC: canManageRBAC,
		Roles:         make([]roleListItemView, 0, len(items)),
		Flash:         h.readFlash(w, r),
	}
	for _, item := range items {
		data.Roles = append(data.Roles, roleListItemView{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			CreatedAt:   item.CreatedAt.Format(time.DateTime),
		})
	}
	return data
}

func makeSessionViews(items []identityread.SessionReadModel, currentSessionID int64) []sessionItemView {
	out := make([]sessionItemView, 0, len(items))
	now := time.Now()
	for _, item := range items {
		out = append(out, sessionItemView{
			IDSession: item.IDSession,
			IDUser:    item.IDUser,
			UserEmail: item.UserEmail,
			IP:        fallback(item.IP, "-"),
			UserAgent: fallback(item.UserAgent, "-"),
			CreatedAt: item.CreatedAt.Format("2006-01-02 15:04"),
			ExpiresAt: item.ExpiresAt.Format("2006-01-02 15:04"),
			UpdatedAt: item.UpdatedAt.Format("2006-01-02 15:04"),
			Status:    sessionStatus(item, now),
			IsCurrent: item.IDSession == currentSessionID,
		})
	}
	return out
}

func sessionStatus(item identityread.SessionReadModel, now time.Time) string {
	if item.RevokedAt != nil {
		return "revogada"
	}
	if !now.Before(item.ExpiresAt) {
		return "expirada"
	}
	return "ativa"
}

func shouldAttemptRefresh(err error) bool {
	return errors.Is(err, identitydomain.ErrMissingAccessToken) ||
		errors.Is(err, identitydomain.ErrInvalidToken) ||
		errors.Is(err, identitydomain.ErrSessionNotActive) ||
		errors.Is(err, identitydomain.ErrSessionExpired) ||
		errors.Is(err, identitydomain.ErrSessionRevoked)
}

func mapLoginError(err error) (int, string) {
	switch {
	case errors.Is(err, identitydomain.ErrInvalidCredentials):
		return http.StatusUnauthorized, "Email ou senha invalidos."
	case errors.Is(err, identitydomain.ErrUserBlocked):
		return http.StatusForbidden, "Usuario bloqueado."
	case errors.Is(err, identitydomain.ErrSessionLimit):
		return http.StatusTooManyRequests, "Limite de sessoes ativas atingido."
	default:
		return http.StatusBadRequest, "Nao foi possivel autenticar."
	}
}

func (h AppHandler) setAuthCookies(w http.ResponseWriter, tokens *identityread.AuthTokensReadModel) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieAccessToken,
		Value:    tokens.AccessToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.secureCookies,
		Expires:  tokens.AccessExpiresAt,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     cookieRefreshToken,
		Value:    tokens.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.secureCookies,
		Expires:  tokens.RefreshExpiresAt,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     cookieSessionID,
		Value:    strconv.FormatInt(tokens.IDSession, 10),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.secureCookies,
		Expires:  tokens.RefreshExpiresAt,
	})
}

func (h AppHandler) clearAuthCookies(w http.ResponseWriter) {
	for _, name := range []string{cookieAccessToken, cookieRefreshToken, cookieSessionID} {
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   h.secureCookies,
			MaxAge:   -1,
		})
	}
}

func clientIP(r *http.Request) string {
	if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil {
		return host
	}
	return r.RemoteAddr
}

func baseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}

func fallback(value, defaultValue string) string {
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}
	return value
}

func optionalString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func parsePage(raw string) int {
	page, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || page < 1 {
		return 1
	}
	return page
}

func parseLimit(raw string) int {
	limit, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || limit <= 0 || limit > 100 {
		return 10
	}
	return limit
}

func (h AppHandler) buildCompaniesPageData(w http.ResponseWriter, r *http.Request, rm *companyread.CompanyListReadModel, filter companyrepo.CompanyListFilter, page, limit int) companiesPageData {
	items := make([]companyListItemView, 0, len(rm.Data))
	for _, company := range rm.Data {
		city := "-"
		state := "-"
		if company.Address != nil {
			if company.Address.City != "" {
				city = company.Address.City
			}
			if company.Address.State != "" {
				state = company.Address.State
			}
		}
		items = append(items, companyListItemView{
			ID:               company.ID,
			Name:             company.Name,
			TradeName:        company.TradeName,
			CNPJ:             company.CNPJ,
			RegisteredByName: fallback(company.RegisteredByName, "-"),
			City:             city,
			State:            state,
		})
	}
	rangeStart := 0
	rangeEnd := 0
	if len(items) > 0 {
		rangeStart = ((page - 1) * limit) + 1
		rangeEnd = rangeStart + len(items) - 1
	}
	return companiesPageData{
		Page:            page,
		Limit:           limit,
		RangeStart:      rangeStart,
		RangeEnd:        rangeEnd,
		Items:           items,
		HasNextPage:     rm.PageInfo.HasNextPage,
		HasPreviousPage: rm.PageInfo.HasPreviousPage,
		FilterName:      filter.Name,
		FilterCNPJ:      filter.CNPJ,
		Flash:           h.readFlash(w, r),
		NextURL:         rewritePageLink(r, rm.PageInfo.NextURL, page+1),
		PreviousURL:     rewritePageLink(r, rm.PageInfo.PreviousURL, max(1, page-1)),
	}
}

func makeRoleOptions(items []identityread.RoleListItemReadModel) []roleOptionView {
	out := make([]roleOptionView, 0, len(items))
	for _, item := range items {
		out = append(out, roleOptionView{ID: item.ID, Name: item.Name})
	}
	return out
}

func makeUserRoleViews(items []identityread.UserRoleAssignmentReadModel) []userRoleView {
	out := make([]userRoleView, 0, len(items))
	for _, item := range items {
		out = append(out, userRoleView{
			ID:         item.IDRole,
			Name:       item.RoleName,
			AssignedBy: item.AssignedBy,
		})
	}
	return out
}

func makeRolePermissionViews(items []identityread.RolePermissionAssignmentReadModel) []rolePermissionView {
	out := make([]rolePermissionView, 0, len(items))
	for _, item := range items {
		out = append(out, rolePermissionView{
			ID:          item.IDPermission,
			Code:        item.Resource + ":" + item.Action,
			Description: item.Description,
		})
	}
	return out
}

func parseInt64Field(value string) (int64, error) {
	parsed, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil || parsed <= 0 {
		return 0, errors.New("invalid int64 field")
	}
	return parsed, nil
}

func parseInt16Field(value string) (int16, error) {
	parsed, err := strconv.ParseInt(strings.TrimSpace(value), 10, 16)
	if err != nil || parsed <= 0 {
		return 0, errors.New("invalid int16 field")
	}
	return int16(parsed), nil
}

func parseOptionalInt16(value string) int16 {
	parsed, err := strconv.ParseInt(strings.TrimSpace(value), 10, 16)
	if err != nil || parsed <= 0 {
		return 0
	}
	return int16(parsed)
}

func (h AppHandler) readFlash(w http.ResponseWriter, r *http.Request) flashMessage {
	cookie, err := r.Cookie(cookieFlash)
	if err != nil || strings.TrimSpace(cookie.Value) == "" {
		return flashMessage{}
	}
	h.clearFlashCookie(w)
	raw, err := base64.RawURLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return flashMessage{}
	}
	var flash flashMessage
	if err := json.Unmarshal(raw, &flash); err != nil {
		return flashMessage{}
	}
	if strings.TrimSpace(flash.Message) == "" {
		return flashMessage{}
	}
	if flash.Kind != "success" {
		flash.Kind = "error"
	}
	return flash
}

func (h AppHandler) redirectWithFlash(w http.ResponseWriter, r *http.Request, path, kind, message string) {
	h.redirectWithTypedFlash(w, r, path, "", "", kind, message)
}

func (h AppHandler) redirectWithTypedFlash(w http.ResponseWriter, r *http.Request, path, source, code, kind, message string) {
	h.redirectWithFlashPayload(w, r, path, flashMessage{
		Source:  source,
		Code:    code,
		Kind:    kind,
		Message: message,
	})
}

func (h AppHandler) redirectWithFlashPayload(w http.ResponseWriter, r *http.Request, path string, flash flashMessage) {
	if strings.TrimSpace(flash.Message) == "" {
		http.Redirect(w, r, path, http.StatusSeeOther)
		return
	}
	if flash.Kind != "success" {
		flash.Kind = "error"
	}
	raw, err := json.Marshal(flash)
	if err != nil {
		http.Redirect(w, r, path, http.StatusSeeOther)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     cookieFlash,
		Value:    base64.RawURLEncoding.EncodeToString(raw),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.secureCookies,
		MaxAge:   10,
	})
	http.Redirect(w, r, path, http.StatusSeeOther)
}

func successFlash(source, code, message string) flashMessage {
	return flashMessage{
		Kind:    "success",
		Source:  source,
		Code:    code,
		Message: message,
	}
}

func errorFlash(source, code, message string) flashMessage {
	return flashMessage{
		Kind:    "error",
		Source:  source,
		Code:    code,
		Message: message,
	}
}

func (h AppHandler) clearFlashCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieFlash,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.secureCookies,
		MaxAge:   -1,
	})
}

func mapRBACCatalogError(err error, fallback string) string {
	if err == nil {
		return fallback
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23503":
			return "Nao foi possivel excluir porque o registro ainda esta em uso por outros relacionamentos."
		case "23505":
			return "Nao foi possivel salvar porque ja existe um registro com esses dados."
		}
	}
	return fallback
}

func rewritePageLink(r *http.Request, raw *string, page int) string {
	if raw == nil || strings.TrimSpace(*raw) == "" {
		return ""
	}
	u, err := url.Parse(*raw)
	if err != nil {
		return ""
	}
	q := u.Query()
	q.Set("page", strconv.Itoa(page))
	q.Set("name", strings.TrimSpace(r.URL.Query().Get("name")))
	q.Set("cnpj", strings.TrimSpace(r.URL.Query().Get("cnpj")))
	u.RawQuery = q.Encode()
	return u.RequestURI()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (h AppHandler) renderCompanyForm(w http.ResponseWriter, status int, data companyFormData, templateName string) error {
	return h.renderer.RenderHTTP(w, status, templateName, webtemplate.Page{
		Title:       "Company Form",
		Description: "Formulario WEB de company.",
		NavKey:      companyFormNavKey(data.Mode),
		Entry:       "src/main.ts",
		Data:        data,
	})
}

func companyFormNavKey(mode string) string {
	if mode == "create" {
		return "companies.new"
	}
	return "companies.list"
}

func (h AppHandler) loadCompanyForEdit(w http.ResponseWriter, r *http.Request) (*companyread.CompanyReadModel, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		h.respondNotFound(w, r)
		return nil, false
	}
	company, err := h.findCompanyByID.Execute(r.Context(), id)
	if err != nil {
		h.logger.Error(r.Context(), "find company for edit", "error", err, "id_company", id)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return nil, false
	}
	if company == nil {
		h.respondNotFound(w, r)
		return nil, false
	}
	return company, true
}

func newCompanyFormData() companyFormData {
	return companyFormData{
		Mode:        "create",
		Action:      "/companies",
		SubmitLabel: "Criar company",
		PhoneKind:   "commercial",
		Country:     "BR",
		SocialMedias: []companySocialMediaFormItem{{
			Platform: "linkedin",
		}},
	}
}

func companyToFormData(company *companyread.CompanyReadModel) companyFormData {
	data := companyFormData{
		Mode:        "edit",
		Action:      fmt.Sprintf("/companies/%d", company.ID),
		SubmitLabel: "Salvar alteracoes",
		CompanyID:   company.ID,
		Name:        company.Name,
		TradeName:   company.TradeName,
		CNPJ:        company.CNPJ,
	}
	if company.Address != nil {
		data.AddressID = company.Address.ID
		data.ZIP = company.Address.ZIP
		data.Title = company.Address.Title
		data.Street = company.Address.Street
		data.Number = strconv.FormatUint(uint64(company.Address.Number), 10)
		if company.Address.Complement != nil {
			data.Complement = *company.Address.Complement
		}
		if company.Address.Reference != nil {
			data.Reference = *company.Address.Reference
		}
		data.Neighborhood = company.Address.Neighborhood
		data.City = company.Address.City
		data.State = company.Address.State
		data.Country = company.Address.Country
	}
	if len(company.Phones) > 0 {
		data.PhoneID = company.Phones[0].ID
		data.PhoneNumber = company.Phones[0].Number
		data.PhoneKind = company.Phones[0].Kind
		data.PhoneDepartment = company.Phones[0].Department
	}
	if len(company.Emails) > 0 {
		data.EmailID = company.Emails[0].ID
		data.EmailAddress = company.Emails[0].Address
	}
	if len(company.SocialMedia) > 0 {
		data.SocialMedias = make([]companySocialMediaFormItem, 0, len(company.SocialMedia))
		for _, item := range company.SocialMedia {
			data.SocialMedias = append(data.SocialMedias, companySocialMediaFormItem{
				ID:       item.ID,
				Platform: item.Platform,
				URL:      item.URL,
			})
		}
	}
	if len(data.SocialMedias) == 0 {
		data.SocialMedias = []companySocialMediaFormItem{{}}
	}
	return data
}

func parseCreateCompanyForm(r *http.Request, idUser int64) (companyFormData, companycmd.CreateCompanyCommand, error) {
	if err := r.ParseForm(); err != nil {
		return companyFormData{}, companycmd.CreateCompanyCommand{}, err
	}
	number, err := parseUintForm(r.FormValue("number"))
	if err != nil {
		return companyFormData{}, companycmd.CreateCompanyCommand{}, err
	}
	form := newCompanyFormData()
	form.Name = strings.TrimSpace(r.FormValue("name"))
	form.TradeName = strings.TrimSpace(r.FormValue("trade_name"))
	form.CNPJ = strings.TrimSpace(r.FormValue("cnpj"))
	form.ZIP = strings.TrimSpace(r.FormValue("zip"))
	form.Title = strings.TrimSpace(r.FormValue("title"))
	form.Street = strings.TrimSpace(r.FormValue("street"))
	form.Number = strings.TrimSpace(r.FormValue("number"))
	form.Complement = strings.TrimSpace(r.FormValue("complement"))
	form.Reference = strings.TrimSpace(r.FormValue("reference"))
	form.Neighborhood = strings.TrimSpace(r.FormValue("neighborhood"))
	form.City = strings.TrimSpace(r.FormValue("city"))
	form.State = strings.TrimSpace(r.FormValue("state"))
	form.Country = strings.TrimSpace(r.FormValue("country"))
	form.PhoneNumber = strings.TrimSpace(r.FormValue("phone_number"))
	form.PhoneKind = strings.TrimSpace(r.FormValue("phone_kind"))
	form.PhoneDepartment = strings.TrimSpace(r.FormValue("phone_department"))
	form.EmailAddress = strings.TrimSpace(r.FormValue("email_address"))
	form.SocialMedias = parseSocialMediaForm(r)

	cmd := companycmd.CreateCompanyCommand{
		Name:         form.Name,
		TradeName:    form.TradeName,
		CNPJ:         form.CNPJ,
		RegisteredBy: idUser,
		Address: companycmd.CreateCompanyAddressCommand{
			ZIP:          form.ZIP,
			Title:        form.Title,
			Street:       form.Street,
			Number:       number,
			Complement:   form.Complement,
			Reference:    form.Reference,
			Neighborhood: form.Neighborhood,
			City:         form.City,
			State:        form.State,
			Country:      form.Country,
		},
		Phone: companycmd.CreateCompanyPhoneCommand{
			Number:     form.PhoneNumber,
			Kind:       form.PhoneKind,
			Department: form.PhoneDepartment,
		},
		Email: companycmd.CreateCompanyEmailCommand{
			Address: form.EmailAddress,
		},
	}
	for _, item := range form.SocialMedias {
		if item.Platform == "" && item.URL == "" {
			continue
		}
		cmd.SocialMedia = append(cmd.SocialMedia, companycmd.CreateCompanySocialMediaCommand{
			Platform: item.Platform,
			URL:      item.URL,
		})
	}
	return form, cmd, nil
}

func parseUpdateCompanyForm(r *http.Request, current *companyread.CompanyReadModel) (companyFormData, companycmd.UpdateCompanyCommand, error) {
	if err := r.ParseForm(); err != nil {
		return companyFormData{}, companycmd.UpdateCompanyCommand{}, err
	}
	number, err := parseUintForm(r.FormValue("number"))
	if err != nil {
		return companyFormData{}, companycmd.UpdateCompanyCommand{}, err
	}
	form := companyToFormData(current)
	form.Name = strings.TrimSpace(r.FormValue("name"))
	form.TradeName = strings.TrimSpace(r.FormValue("trade_name"))
	form.CNPJ = strings.TrimSpace(r.FormValue("cnpj"))
	form.ZIP = strings.TrimSpace(r.FormValue("zip"))
	form.Title = strings.TrimSpace(r.FormValue("title"))
	form.Street = strings.TrimSpace(r.FormValue("street"))
	form.Number = strings.TrimSpace(r.FormValue("number"))
	form.Complement = strings.TrimSpace(r.FormValue("complement"))
	form.Reference = strings.TrimSpace(r.FormValue("reference"))
	form.Neighborhood = strings.TrimSpace(r.FormValue("neighborhood"))
	form.City = strings.TrimSpace(r.FormValue("city"))
	form.State = strings.TrimSpace(r.FormValue("state"))
	form.Country = strings.TrimSpace(r.FormValue("country"))
	form.PhoneNumber = strings.TrimSpace(r.FormValue("phone_number"))
	form.PhoneKind = strings.TrimSpace(r.FormValue("phone_kind"))
	form.PhoneDepartment = strings.TrimSpace(r.FormValue("phone_department"))
	form.EmailAddress = strings.TrimSpace(r.FormValue("email_address"))
	form.SocialMedias = parseSocialMediaForm(r)

	cmd := companycmd.UpdateCompanyCommand{
		IDCompany: current.ID,
		Name:      form.Name,
		TradeName: form.TradeName,
		CNPJ:      form.CNPJ,
		Address: companycmd.UpdateCompanyAddressCommand{
			IDAddress:    form.AddressID,
			ZIP:          form.ZIP,
			Title:        form.Title,
			Street:       form.Street,
			Number:       number,
			Complement:   form.Complement,
			Reference:    form.Reference,
			Neighborhood: form.Neighborhood,
			City:         form.City,
			State:        form.State,
			Country:      form.Country,
		},
		Phone: companycmd.UpdateCompanyPhoneCommand{
			IDPhone:    form.PhoneID,
			Number:     form.PhoneNumber,
			Kind:       form.PhoneKind,
			Department: form.PhoneDepartment,
		},
		Email: companycmd.UpdateCompanyEmailCommand{
			IDEmail: form.EmailID,
			Address: form.EmailAddress,
		},
	}
	for _, item := range form.SocialMedias {
		if item.Platform == "" && item.URL == "" {
			continue
		}
		cmd.SocialMedia = append(cmd.SocialMedia, companycmd.UpdateCompanySocialMediaCommand{
			IDSocialMedia: item.ID,
			Platform:      item.Platform,
			URL:           item.URL,
		})
	}
	return form, cmd, nil
}

func parseUintForm(raw string) (uint, error) {
	value, err := strconv.ParseUint(strings.TrimSpace(raw), 10, 64)
	if err != nil || value == 0 {
		return 0, errors.New("invalid number")
	}
	return uint(value), nil
}

func parseSocialMediaForm(r *http.Request) []companySocialMediaFormItem {
	ids := r.Form["social_media_id"]
	platforms := r.Form["social_platform"]
	urls := r.Form["social_url"]
	size := max(len(ids), len(platforms))
	size = max(size, len(urls))
	if size == 0 {
		return []companySocialMediaFormItem{{}}
	}
	items := make([]companySocialMediaFormItem, 0, size)
	for i := 0; i < size; i++ {
		var id int64
		if i < len(ids) {
			id, _ = strconv.ParseInt(strings.TrimSpace(ids[i]), 10, 64)
		}
		platform := ""
		if i < len(platforms) {
			platform = strings.TrimSpace(platforms[i])
		}
		link := ""
		if i < len(urls) {
			link = strings.TrimSpace(urls[i])
		}
		items = append(items, companySocialMediaFormItem{
			ID:       id,
			Platform: platform,
			URL:      link,
		})
	}
	return items
}
