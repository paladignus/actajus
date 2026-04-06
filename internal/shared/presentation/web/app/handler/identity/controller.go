package identityhandler

import (
	"net/http"
	"strings"
	"time"

	identitycmd "github.com/paladignus/actajus/internal/module/identity/application/command"
	identityread "github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	identityrepo "github.com/paladignus/actajus/internal/module/identity/application/repository"
	identityusecase "github.com/paladignus/actajus/internal/module/identity/application/usecase"
	webflash "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/flash"
	requestmeta "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/requestmeta"
	identitybinder "github.com/paladignus/actajus/internal/shared/presentation/web/app/binder/http/identity"
	querybinder "github.com/paladignus/actajus/internal/shared/presentation/web/app/binder/http/query"
	handlerctx "github.com/paladignus/actajus/internal/shared/presentation/web/app/handlerctx"
	identitymessage "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/message/identity"
	identitypage "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/page/identity"
	identityviewmodel "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/viewmodel/identity"
	httperror "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/http/error"
)

type Controller struct {
	Context                  handlerctx.Context
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

func (c Controller) SessionsPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := c.Context.RequireAuth(w, r)
		if err != nil {
			return
		}
		adminFilter := strings.TrimSpace(r.URL.Query().Get("email"))
		flash := c.Context.Navigator.ReadFlash(w, r, c.Context.FlashCookie)
		own, err := c.ViewSessions.Own(r.Context(), state.Claims.IDUser)
		if err != nil {
			c.Context.Logger.Error(r.Context(), "list own sessions", "error", err, "id_user", state.Claims.IDUser)
			c.renderSessions(w, r, http.StatusInternalServerError, state, adminFilter, flash, nil, nil, time.Now(), "Nao foi possivel carregar as suas sessoes.")
			return
		}
		var adminSessions []identityread.SessionReadModel
		if state.CanManageSessions {
			adminSessions, err = c.ViewSessions.Admin(r.Context(), identityrepo.SessionQueryFilter{UserEmail: adminFilter, Limit: 50})
			if err != nil {
				c.Context.Logger.Error(r.Context(), "list admin sessions", "error", err, "email", adminFilter)
				c.renderSessions(w, r, http.StatusInternalServerError, state, adminFilter, flash, own, nil, time.Now(), "Nao foi possivel carregar as sessoes administrativas.")
				return
			}
		}
		c.renderSessions(w, r, http.StatusOK, state, adminFilter, flash, own, adminSessions, time.Now(), "")
	})
}

func (c Controller) RevokeSessionAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := c.Context.RequireAuth(w, r)
		if err != nil {
			return
		}
		form, cmd, err := identitybinder.BindRevokeSessionForm(r)
		if err != nil {
			c.redirectSessions(w, r, identitymessage.InvalidForm())
			return
		}
		target, err := c.ViewSessions.ByID(r.Context(), form.IDSession)
		if err != nil {
			c.Context.Logger.Error(r.Context(), "find session", "error", err, "id_session", form.IDSession)
			c.redirectSessions(w, r, identitymessage.SessionLookupFailed())
			return
		}
		if target == nil {
			c.redirectSessions(w, r, identitymessage.SessionNotFound())
			return
		}
		if target.IDUser != state.Claims.IDUser && !state.CanManageSessions {
			if c.Context.AuditFailure != nil {
				c.Context.AuditFailure(r, "identity.session.revoke.forbidden", "id_user", state.Claims.IDUser, "id_session", form.IDSession)
			}
			c.Context.RespondForbidden(w, r)
			return
		}
		if err := c.RevokeSession.Execute(r.Context(), cmd); err != nil {
			c.Context.Logger.Error(r.Context(), "revoke session", "error", err, "id_session", form.IDSession)
			if c.Context.AuditFailure != nil {
				c.Context.AuditFailure(r, "identity.session.revoke.failed", "id_user", state.Claims.IDUser, "id_session", form.IDSession)
			}
			c.redirectSessions(w, r, identitymessage.SessionRevokeFailed())
			return
		}
		if c.Context.AuditSuccess != nil {
			c.Context.AuditSuccess(r, "identity.session.revoke.succeeded", "id_user", state.Claims.IDUser, "id_session", form.IDSession)
		}
		if form.IDSession == state.Claims.IDSession {
			c.Context.ClearAuthCookies(w)
			c.Context.Navigator.SeeOther(w, r, "/login")
			return
		}
		c.redirectSessions(w, r, identitymessage.SessionRevokeSucceeded())
	})
}

func (c Controller) RevokeAllAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := c.Context.RequireAuth(w, r)
		if err != nil {
			return
		}
		form, cmd, err := identitybinder.BindRevokeAllSessionsForm(r)
		if err != nil {
			c.redirectSessions(w, r, identitymessage.InvalidForm())
			return
		}
		targetUserID := state.Claims.IDUser
		if form.IDUser != nil {
			targetUserID = *form.IDUser
		}
		if targetUserID != state.Claims.IDUser && !state.CanManageSessions {
			if c.Context.AuditFailure != nil {
				c.Context.AuditFailure(r, "identity.session.revoke_all.forbidden", "id_user", state.Claims.IDUser, "target_id_user", targetUserID)
			}
			c.Context.RespondForbidden(w, r)
			return
		}
		if cmd == nil {
			cmd = &identitycmd.LogoutAllCommand{IDUser: targetUserID}
		}
		if err := c.LogoutAll.Execute(r.Context(), *cmd); err != nil {
			c.Context.Logger.Error(r.Context(), "revoke all sessions", "error", err, "id_user", targetUserID)
			if c.Context.AuditFailure != nil {
				c.Context.AuditFailure(r, "identity.session.revoke_all.failed", "id_user", state.Claims.IDUser, "target_id_user", targetUserID)
			}
			c.redirectSessions(w, r, identitymessage.SessionRevokeAllFailed())
			return
		}
		if c.Context.AuditSuccess != nil {
			c.Context.AuditSuccess(r, "identity.session.revoke_all.succeeded", "id_user", state.Claims.IDUser, "target_id_user", targetUserID)
		}
		if targetUserID == state.Claims.IDUser {
			c.Context.ClearAuthCookies(w)
			c.Context.Navigator.SeeOther(w, r, "/login")
			return
		}
		c.redirectSessions(w, r, identitymessage.SessionRevokeAllSucceeded())
	})
}

func (c Controller) UsersPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := c.Context.RequireAuth(w, r)
		if err != nil {
			return
		}
		filter := identitybinder.BindCatalogUserFilter(r, 50)
		items, err := c.ListUsers.Execute(r.Context(), filter)
		if err != nil {
			c.Context.Logger.Error(r.Context(), "list users", "error", err)
			httperror.InternalServerError(w)
			return
		}
		flash := c.Context.Navigator.ReadFlash(w, r, c.Context.FlashCookie)
		roles := []identityread.RoleListItemReadModel(nil)
		assignments := make(map[int64][]identityread.UserRoleAssignmentReadModel, len(items))
		if state.CanManageRBAC {
			roles, err = c.ListRoles.Execute(r.Context())
			if err == nil {
				for _, item := range items {
					assignments[item.ID], err = c.ViewUserRoles.Execute(r.Context(), item.ID)
					if err != nil {
						break
					}
				}
			}
		}
		if err != nil {
			c.Context.Logger.Error(r.Context(), "build users page data", "error", err)
			httperror.InternalServerError(w)
			return
		}
		c.renderUsers(w, http.StatusOK, items, filter.Query, state.CanManageRBAC, roles, assignments, flash)
	})
}

func (c Controller) PermissionsPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := c.Context.RequireAuth(w, r)
		if err != nil {
			return
		}
		filter := identitybinder.BindCatalogPermissionFilter(r, 100)
		items, err := c.ListPermissions.Execute(r.Context(), filter)
		if err != nil {
			c.Context.Logger.Error(r.Context(), "list permissions", "error", err)
			httperror.InternalServerError(w)
			return
		}
		flash := c.Context.Navigator.ReadFlash(w, r, c.Context.FlashCookie)
		roles := []identityread.RoleListItemReadModel(nil)
		assignments := map[int16][]identityread.RolePermissionAssignmentReadModel{}
		if state.CanManageRBAC {
			roles, err = c.ListRoles.Execute(r.Context())
			if err == nil {
				assignments = make(map[int16][]identityread.RolePermissionAssignmentReadModel, len(roles))
				for _, role := range roles {
					assignments[role.ID], err = c.ViewRolePermissions.Execute(r.Context(), role.ID)
					if err != nil {
						break
					}
				}
			}
		}
		if err != nil {
			c.Context.Logger.Error(r.Context(), "build permissions page data", "error", err)
			httperror.InternalServerError(w)
			return
		}
		c.renderPermissions(w, r, http.StatusOK, items, filter.Query, state.CanManageRBAC, roles, assignments, flash)
	})
}

func (c Controller) RolesPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := c.Context.RequireAuth(w, r)
		if err != nil {
			return
		}
		roles, err := c.ListRoles.Execute(r.Context())
		if err != nil {
			c.Context.Logger.Error(r.Context(), "list roles", "error", err)
			httperror.InternalServerError(w)
			return
		}
		c.renderRoles(w, r, http.StatusOK, roles, state.CanManageRBAC, c.Context.Navigator.ReadFlash(w, r, c.Context.FlashCookie))
	})
}

func (c Controller) CreateRoleAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := c.Context.RequireRBACAdmin(w, r); err != nil {
			return
		}
		_, cmd, err := identitybinder.BindCreateRoleForm(r)
		if err != nil {
			httperror.BadRequest(w, "invalid form")
			return
		}
		err = c.CreateRole.Execute(r.Context(), cmd)
		if err != nil {
			c.Context.Logger.Error(r.Context(), "create role", "error", err)
			if c.Context.AuditFailure != nil {
				c.Context.AuditFailure(r, "identity.role.create.failed", "name", cmd.Name)
			}
			c.redirectRoles(w, r, identitymessage.RoleCreateFailed(err))
			return
		}
		if c.Context.AuditSuccess != nil {
			c.Context.AuditSuccess(r, "identity.role.create.succeeded", "name", cmd.Name)
		}
		c.redirectRoles(w, r, identitymessage.RoleCreateSucceeded())
	})
}

func (c Controller) UpdateRoleAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := c.Context.RequireRBACAdmin(w, r); err != nil {
			return
		}
		_, cmd, err := identitybinder.BindUpdateRoleForm(r)
		if err != nil {
			httperror.BadRequest(w, "invalid form")
			return
		}
		err = c.UpdateRole.Execute(r.Context(), cmd)
		if err != nil {
			c.Context.Logger.Error(r.Context(), "update role", "error", err, "id_role", cmd.ID)
			if c.Context.AuditFailure != nil {
				c.Context.AuditFailure(r, "identity.role.update.failed", "id_role", cmd.ID)
			}
			c.redirectRoles(w, r, identitymessage.RoleUpdateFailed(err))
			return
		}
		if c.Context.AuditSuccess != nil {
			c.Context.AuditSuccess(r, "identity.role.update.succeeded", "id_role", cmd.ID)
		}
		c.redirectRoles(w, r, identitymessage.RoleUpdateSucceeded())
	})
}

func (c Controller) DeleteRoleAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := c.Context.RequireRBACAdmin(w, r); err != nil {
			return
		}
		cmd, err := identitybinder.BindDeleteRolePath(r)
		if err != nil {
			httperror.BadRequest(w, "invalid role id")
			return
		}
		if err := c.DeleteRole.Execute(r.Context(), cmd); err != nil {
			c.Context.Logger.Error(r.Context(), "delete role", "error", err, "id_role", cmd.ID)
			if c.Context.AuditFailure != nil {
				c.Context.AuditFailure(r, "identity.role.delete.failed", "id_role", cmd.ID)
			}
			c.redirectRoles(w, r, identitymessage.RoleDeleteFailed(err))
			return
		}
		if c.Context.AuditSuccess != nil {
			c.Context.AuditSuccess(r, "identity.role.delete.succeeded", "id_role", cmd.ID)
		}
		c.redirectRoles(w, r, identitymessage.RoleDeleteSucceeded())
	})
}

func (c Controller) CreatePermissionAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := c.Context.RequireRBACAdmin(w, r); err != nil {
			return
		}
		_, cmd, err := identitybinder.BindCreatePermissionForm(r)
		if err != nil {
			httperror.BadRequest(w, "invalid form")
			return
		}
		err = c.CreatePermission.Execute(r.Context(), cmd)
		if err != nil {
			c.Context.Logger.Error(r.Context(), "create permission", "error", err)
			if c.Context.AuditFailure != nil {
				c.Context.AuditFailure(r, "identity.permission.create.failed", "resource", cmd.Resource, "action_name", cmd.Action)
			}
			c.redirectPermissions(w, r, identitymessage.PermissionCreateFailed(err))
			return
		}
		if c.Context.AuditSuccess != nil {
			c.Context.AuditSuccess(r, "identity.permission.create.succeeded", "resource", cmd.Resource, "action_name", cmd.Action)
		}
		c.redirectPermissions(w, r, identitymessage.PermissionCreateSucceeded())
	})
}

func (c Controller) UpdatePermissionAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := c.Context.RequireRBACAdmin(w, r); err != nil {
			return
		}
		_, cmd, err := identitybinder.BindUpdatePermissionForm(r)
		if err != nil {
			httperror.BadRequest(w, "invalid form")
			return
		}
		err = c.UpdatePermission.Execute(r.Context(), cmd)
		if err != nil {
			c.Context.Logger.Error(r.Context(), "update permission", "error", err, "id_permission", cmd.ID)
			if c.Context.AuditFailure != nil {
				c.Context.AuditFailure(r, "identity.permission.update.failed", "id_permission", cmd.ID)
			}
			c.redirectPermissions(w, r, identitymessage.PermissionUpdateFailed(err))
			return
		}
		if c.Context.AuditSuccess != nil {
			c.Context.AuditSuccess(r, "identity.permission.update.succeeded", "id_permission", cmd.ID)
		}
		c.redirectPermissions(w, r, identitymessage.PermissionUpdateSucceeded())
	})
}

func (c Controller) DeletePermissionAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := c.Context.RequireRBACAdmin(w, r); err != nil {
			return
		}
		cmd, err := identitybinder.BindDeletePermissionPath(r)
		if err != nil {
			httperror.BadRequest(w, "invalid permission id")
			return
		}
		if err := c.DeletePermission.Execute(r.Context(), cmd); err != nil {
			c.Context.Logger.Error(r.Context(), "delete permission", "error", err, "id_permission", cmd.ID)
			if c.Context.AuditFailure != nil {
				c.Context.AuditFailure(r, "identity.permission.delete.failed", "id_permission", cmd.ID)
			}
			c.redirectPermissions(w, r, identitymessage.PermissionDeleteFailed(err))
			return
		}
		if c.Context.AuditSuccess != nil {
			c.Context.AuditSuccess(r, "identity.permission.delete.succeeded", "id_permission", cmd.ID)
		}
		c.redirectPermissions(w, r, identitymessage.PermissionDeleteSucceeded())
	})
}

func (c Controller) AssignRoleToUserAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := c.Context.RequireRBACAdmin(w, r)
		if err != nil {
			return
		}
		form, err := identitybinder.BindAssignUserRoleForm(r)
		if err != nil {
			c.redirectUsers(w, r, identitymessage.InvalidForm())
			return
		}
		if err := c.AssignRoleToUser.Execute(r.Context(), c.assignRoleCommand(form, state.Claims.IDUser)); err != nil {
			c.Context.Logger.Error(r.Context(), "assign role to user", "error", err, "id_user", form.IDUser, "id_role", form.IDRole)
			if c.Context.AuditFailure != nil {
				c.Context.AuditFailure(r, "identity.user_role.assign.failed", "id_user", form.IDUser, "id_role", form.IDRole)
			}
			c.redirectUsers(w, r, identitymessage.UserRoleAssignFailed(err))
			return
		}
		if c.Context.AuditSuccess != nil {
			c.Context.AuditSuccess(r, "identity.user_role.assign.succeeded", "id_user", form.IDUser, "id_role", form.IDRole)
		}
		c.redirectUsers(w, r, identitymessage.UserRoleAssignSucceeded())
	})
}

func (c Controller) RemoveRoleFromUserAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := c.Context.RequireRBACAdmin(w, r); err != nil {
			return
		}
		cmd, err := identitybinder.BindRemoveUserRoleForm(r)
		if err != nil {
			c.redirectUsers(w, r, identitymessage.InvalidForm())
			return
		}
		if err := c.RemoveRoleFromUser.Execute(r.Context(), cmd); err != nil {
			c.Context.Logger.Error(r.Context(), "remove role from user", "error", err, "id_user", cmd.IDUser, "id_role", cmd.IDRole)
			if c.Context.AuditFailure != nil {
				c.Context.AuditFailure(r, "identity.user_role.remove.failed", "id_user", cmd.IDUser, "id_role", cmd.IDRole)
			}
			c.redirectUsers(w, r, identitymessage.UserRoleRemoveFailed(err))
			return
		}
		if c.Context.AuditSuccess != nil {
			c.Context.AuditSuccess(r, "identity.user_role.remove.succeeded", "id_user", cmd.IDUser, "id_role", cmd.IDRole)
		}
		c.redirectUsers(w, r, identitymessage.UserRoleRemoveSucceeded())
	})
}

func (c Controller) GrantPermissionToRoleAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := c.Context.RequireRBACAdmin(w, r); err != nil {
			return
		}
		_, cmd, err := identitybinder.BindGrantRolePermissionForm(r)
		if err != nil {
			c.redirectPermissions(w, r, identitymessage.InvalidForm())
			return
		}
		if err := c.GrantPermissionToRole.Execute(r.Context(), cmd); err != nil {
			c.Context.Logger.Error(r.Context(), "grant permission to role", "error", err, "id_role", cmd.IDRole, "id_permission", cmd.IDPermission)
			if c.Context.AuditFailure != nil {
				c.Context.AuditFailure(r, "identity.role_permission.grant.failed", "id_role", cmd.IDRole, "id_permission", cmd.IDPermission)
			}
			c.redirectPermissions(w, r, identitymessage.RolePermissionGrantFailed(err))
			return
		}
		if c.Context.AuditSuccess != nil {
			c.Context.AuditSuccess(r, "identity.role_permission.grant.succeeded", "id_role", cmd.IDRole, "id_permission", cmd.IDPermission)
		}
		c.redirectPermissions(w, r, identitymessage.RolePermissionGrantSucceeded())
	})
}

func (c Controller) RevokePermissionFromRoleAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := c.Context.RequireRBACAdmin(w, r); err != nil {
			return
		}
		_, cmd, err := identitybinder.BindRevokeRolePermissionForm(r)
		if err != nil {
			c.redirectPermissions(w, r, identitymessage.InvalidForm())
			return
		}
		if err := c.RevokePermissionFromRole.Execute(r.Context(), cmd); err != nil {
			c.Context.Logger.Error(r.Context(), "revoke permission from role", "error", err, "id_role", cmd.IDRole, "id_permission", cmd.IDPermission)
			if c.Context.AuditFailure != nil {
				c.Context.AuditFailure(r, "identity.role_permission.revoke.failed", "id_role", cmd.IDRole, "id_permission", cmd.IDPermission)
			}
			c.redirectPermissions(w, r, identitymessage.RolePermissionRevokeFailed(err))
			return
		}
		if c.Context.AuditSuccess != nil {
			c.Context.AuditSuccess(r, "identity.role_permission.revoke.succeeded", "id_role", cmd.IDRole, "id_permission", cmd.IDPermission)
		}
		c.redirectPermissions(w, r, identitymessage.RolePermissionRevokeSucceeded())
	})
}

func (c Controller) renderSessions(w http.ResponseWriter, r *http.Request, status int, state *handlerctx.State, adminFilter string, flash webflash.Message, own []identityread.SessionReadModel, admin []identityread.SessionReadModel, now time.Time, pageError string) {
	page := identitypage.SessionsPage(state.Claims.IDUser, state.CanManageSessions, adminFilter, flash, own, admin, state.Claims.IDSession, now)
	if pageError != "" {
		page.Data = c.sessionsViewData(state, adminFilter, flash, own, admin, now, pageError)
	}
	page.Bootstrap = c.sessionsBootstrap(r, state)
	if err := c.Context.HTML.Render(w, status, "pages/sessions", page); err != nil {
		c.Context.Logger.Error(r.Context(), "render sessions page", "error", err)
		httperror.InternalServerError(w)
	}
}

func (c Controller) renderUsers(w http.ResponseWriter, status int, items []identityread.UserListItemReadModel, filter string, canManageRBAC bool, roles []identityread.RoleListItemReadModel, assignments map[int64][]identityread.UserRoleAssignmentReadModel, flash webflash.Message) {
	if err := c.Context.HTML.Render(w, status, "pages/users", identitypage.UsersPage(items, filter, canManageRBAC, roles, assignments, flash)); err != nil {
		c.Context.Logger.Error(nil, "render users page", "error", err)
		httperror.InternalServerError(w)
	}
}

func (c Controller) renderPermissions(w http.ResponseWriter, r *http.Request, status int, items []identityread.PermissionListItemReadModel, filter string, canManageRBAC bool, roles []identityread.RoleListItemReadModel, assignments map[int16][]identityread.RolePermissionAssignmentReadModel, flash webflash.Message) {
	if err := c.Context.HTML.Render(w, status, "pages/permissions", identitypage.PermissionsPage(
		items,
		filter,
		canManageRBAC,
		roles,
		assignments,
		querybinder.OptionalInt16(r.URL.Query().Get("edit")),
		r.URL.Query().Get("resource"),
		r.URL.Query().Get("action"),
		r.URL.Query().Get("description"),
		flash,
	)); err != nil {
		c.Context.Logger.Error(r.Context(), "render permissions page", "error", err)
		httperror.InternalServerError(w)
	}
}

func (c Controller) renderRoles(w http.ResponseWriter, r *http.Request, status int, items []identityread.RoleListItemReadModel, canManageRBAC bool, flash webflash.Message) {
	if err := c.Context.HTML.Render(w, status, "pages/roles", identitypage.RolesPage(
		items,
		canManageRBAC,
		querybinder.OptionalInt16(r.URL.Query().Get("edit")),
		r.URL.Query().Get("name"),
		r.URL.Query().Get("description"),
		flash,
	)); err != nil {
		c.Context.Logger.Error(r.Context(), "render roles page", "error", err)
		httperror.InternalServerError(w)
	}
}

func (c Controller) sessionsBootstrap(r *http.Request, state *handlerctx.State) map[string]any {
	return map[string]any{
		"appName":           "Actajus",
		"apiBaseURL":        requestmeta.BaseURL(r),
		"viewerID":          state.Claims.IDUser,
		"canManageSessions": state.CanManageSessions,
	}
}

func (c Controller) sessionsViewData(state *handlerctx.State, adminFilter string, flash webflash.Message, own []identityread.SessionReadModel, admin []identityread.SessionReadModel, now time.Time, pageError string) identityviewmodel.SessionsPage {
	pageData := identityviewmodel.BuildSessionsPage(state.Claims.IDUser, state.CanManageSessions, adminFilter, flash, own, admin, state.Claims.IDSession, now)
	pageData.Error = pageError
	return pageData
}

func (c Controller) redirectUsers(w http.ResponseWriter, r *http.Request, message webflash.Message) {
	c.Context.Navigator.SeeOtherWithFlash(w, r, "/users", c.Context.FlashCookie, message)
}

func (c Controller) redirectPermissions(w http.ResponseWriter, r *http.Request, message webflash.Message) {
	c.Context.Navigator.SeeOtherWithFlash(w, r, "/permissions", c.Context.FlashCookie, message)
}

func (c Controller) redirectRoles(w http.ResponseWriter, r *http.Request, message webflash.Message) {
	c.Context.Navigator.SeeOtherWithFlash(w, r, "/roles", c.Context.FlashCookie, message)
}

func (c Controller) redirectSessions(w http.ResponseWriter, r *http.Request, message webflash.Message) {
	c.Context.Navigator.SeeOtherWithFlash(w, r, "/sessions", c.Context.FlashCookie, message)
}

func (c Controller) assignRoleCommand(form identitybinder.AssignUserRoleForm, assignedBy int64) identitycmd.AssignRoleToUserCommand {
	return identitycmd.AssignRoleToUserCommand{
		IDUser:     form.IDUser,
		IDRole:     form.IDRole,
		AssignedBy: assignedBy,
	}
}
