package handler

import (
	"net/http"

	identitycmd "github.com/paladignus/actajus/internal/module/identity/application/command"
	identityread "github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	identitybinder "github.com/paladignus/actajus/internal/shared/presentation/web/app/binder/http/identity"
	querybinder "github.com/paladignus/actajus/internal/shared/presentation/web/app/binder/http/query"
	identitymessage "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/message/identity"
	identitypage "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/page/identity"
	httperror "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/http/error"
)

func (h identityRoutes) UsersPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := h.app.requireAuth(w, r)
		if err != nil {
			return
		}
		filter := identitybinder.BindCatalogUserFilter(r, 50)
		items, err := h.app.identity.listUsers.Execute(r.Context(), filter)
		if err != nil {
			h.app.logger.Error(r.Context(), "list users", "error", err)
			httperror.InternalServerError(w)
			return
		}
		flash := h.app.navigator.ReadFlash(w, r, cookieFlash)
		roles := []identityread.RoleListItemReadModel(nil)
		assignments := make(map[int64][]identityread.UserRoleAssignmentReadModel, len(items))
		if state.CanManageRBAC {
			roles, err = h.app.identity.listRoles.Execute(r.Context())
			if err == nil {
				for _, item := range items {
					assignments[item.ID], err = h.app.identity.viewUserRoles.Execute(r.Context(), item.ID)
					if err != nil {
						break
					}
				}
			}
		}
		if err != nil {
			h.app.logger.Error(r.Context(), "build users page data", "error", err)
			httperror.InternalServerError(w)
			return
		}
		if err := h.app.html.Render(w, http.StatusOK, "pages/users", identitypage.UsersPage(items, filter.Query, state.CanManageRBAC, roles, assignments, flash)); err != nil {
			h.app.logger.Error(r.Context(), "render users page", "error", err)
			httperror.InternalServerError(w)
		}
	})
}

func (h identityRoutes) PermissionsPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := h.app.requireAuth(w, r)
		if err != nil {
			return
		}
		filter := identitybinder.BindCatalogPermissionFilter(r, 100)
		items, err := h.app.identity.listPermissions.Execute(r.Context(), filter)
		if err != nil {
			h.app.logger.Error(r.Context(), "list permissions", "error", err)
			httperror.InternalServerError(w)
			return
		}
		flash := h.app.navigator.ReadFlash(w, r, cookieFlash)
		roles := []identityread.RoleListItemReadModel(nil)
		assignments := map[int16][]identityread.RolePermissionAssignmentReadModel{}
		if state.CanManageRBAC {
			roles, err = h.app.identity.listRoles.Execute(r.Context())
			if err == nil {
				assignments = make(map[int16][]identityread.RolePermissionAssignmentReadModel, len(roles))
				for _, role := range roles {
					assignments[role.ID], err = h.app.identity.viewRolePermissions.Execute(r.Context(), role.ID)
					if err != nil {
						break
					}
				}
			}
		}
		if err != nil {
			h.app.logger.Error(r.Context(), "build permissions page data", "error", err)
			httperror.InternalServerError(w)
			return
		}
		if err := h.app.html.Render(w, http.StatusOK, "pages/permissions", identitypage.PermissionsPage(
			items,
			filter.Query,
			state.CanManageRBAC,
			roles,
			assignments,
			querybinder.OptionalInt16(r.URL.Query().Get("edit")),
			r.URL.Query().Get("resource"),
			r.URL.Query().Get("action"),
			r.URL.Query().Get("description"),
			flash,
		)); err != nil {
			h.app.logger.Error(r.Context(), "render permissions page", "error", err)
			httperror.InternalServerError(w)
		}
	})
}

func (h identityRoutes) RolesPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := h.app.requireAuth(w, r)
		if err != nil {
			return
		}
		roles, err := h.app.identity.listRoles.Execute(r.Context())
		if err != nil {
			h.app.logger.Error(r.Context(), "list roles", "error", err)
			httperror.InternalServerError(w)
			return
		}
		if err := h.app.html.Render(w, http.StatusOK, "pages/roles", identitypage.RolesPage(
			roles,
			state.CanManageRBAC,
			querybinder.OptionalInt16(r.URL.Query().Get("edit")),
			r.URL.Query().Get("name"),
			r.URL.Query().Get("description"),
			h.app.navigator.ReadFlash(w, r, cookieFlash),
		)); err != nil {
			h.app.logger.Error(r.Context(), "render roles page", "error", err)
			httperror.InternalServerError(w)
		}
	})
}

func (h identityRoutes) CreateRoleAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.app.requireRBACAdmin(w, r); err != nil {
			return
		}
		_, cmd, err := identitybinder.BindCreateRoleForm(r)
		if err != nil {
			httperror.BadRequest(w, "invalid form")
			return
		}
		err = h.app.identity.createRole.Execute(r.Context(), cmd)
		if err != nil {
			h.app.logger.Error(r.Context(), "create role", "error", err)
			h.app.navigator.SeeOtherWithFlash(w, r, "/roles", cookieFlash, identitymessage.RoleCreateFailed(err))
			return
		}
		h.app.navigator.SeeOtherWithFlash(w, r, "/roles", cookieFlash, identitymessage.RoleCreateSucceeded())
	})
}

func (h identityRoutes) UpdateRoleAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.app.requireRBACAdmin(w, r); err != nil {
			return
		}
		_, cmd, err := identitybinder.BindUpdateRoleForm(r)
		if err != nil {
			httperror.BadRequest(w, "invalid form")
			return
		}
		err = h.app.identity.updateRole.Execute(r.Context(), cmd)
		if err != nil {
			h.app.logger.Error(r.Context(), "update role", "error", err, "id_role", cmd.ID)
			h.app.navigator.SeeOtherWithFlash(w, r, "/roles", cookieFlash, identitymessage.RoleUpdateFailed(err))
			return
		}
		h.app.navigator.SeeOtherWithFlash(w, r, "/roles", cookieFlash, identitymessage.RoleUpdateSucceeded())
	})
}

func (h identityRoutes) DeleteRoleAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.app.requireRBACAdmin(w, r); err != nil {
			return
		}
		cmd, err := identitybinder.BindDeleteRolePath(r)
		if err != nil {
			httperror.BadRequest(w, "invalid role id")
			return
		}
		if err := h.app.identity.deleteRole.Execute(r.Context(), cmd); err != nil {
			h.app.logger.Error(r.Context(), "delete role", "error", err, "id_role", cmd.ID)
			h.app.navigator.SeeOtherWithFlash(w, r, "/roles", cookieFlash, identitymessage.RoleDeleteFailed(err))
			return
		}
		h.app.navigator.SeeOtherWithFlash(w, r, "/roles", cookieFlash, identitymessage.RoleDeleteSucceeded())
	})
}

func (h identityRoutes) CreatePermissionAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.app.requireRBACAdmin(w, r); err != nil {
			return
		}
		_, cmd, err := identitybinder.BindCreatePermissionForm(r)
		if err != nil {
			httperror.BadRequest(w, "invalid form")
			return
		}
		err = h.app.identity.createPermission.Execute(r.Context(), cmd)
		if err != nil {
			h.app.logger.Error(r.Context(), "create permission", "error", err)
			h.app.navigator.SeeOtherWithFlash(w, r, "/permissions", cookieFlash, identitymessage.PermissionCreateFailed(err))
			return
		}
		h.app.navigator.SeeOtherWithFlash(w, r, "/permissions", cookieFlash, identitymessage.PermissionCreateSucceeded())
	})
}

func (h identityRoutes) UpdatePermissionAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.app.requireRBACAdmin(w, r); err != nil {
			return
		}
		_, cmd, err := identitybinder.BindUpdatePermissionForm(r)
		if err != nil {
			httperror.BadRequest(w, "invalid form")
			return
		}
		err = h.app.identity.updatePermission.Execute(r.Context(), cmd)
		if err != nil {
			h.app.logger.Error(r.Context(), "update permission", "error", err, "id_permission", cmd.ID)
			h.app.navigator.SeeOtherWithFlash(w, r, "/permissions", cookieFlash, identitymessage.PermissionUpdateFailed(err))
			return
		}
		h.app.navigator.SeeOtherWithFlash(w, r, "/permissions", cookieFlash, identitymessage.PermissionUpdateSucceeded())
	})
}

func (h identityRoutes) DeletePermissionAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.app.requireRBACAdmin(w, r); err != nil {
			return
		}
		cmd, err := identitybinder.BindDeletePermissionPath(r)
		if err != nil {
			httperror.BadRequest(w, "invalid permission id")
			return
		}
		if err := h.app.identity.deletePermission.Execute(r.Context(), cmd); err != nil {
			h.app.logger.Error(r.Context(), "delete permission", "error", err, "id_permission", cmd.ID)
			h.app.navigator.SeeOtherWithFlash(w, r, "/permissions", cookieFlash, identitymessage.PermissionDeleteFailed(err))
			return
		}
		h.app.navigator.SeeOtherWithFlash(w, r, "/permissions", cookieFlash, identitymessage.PermissionDeleteSucceeded())
	})
}

func (h identityRoutes) AssignRoleToUserAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := h.app.requireRBACAdmin(w, r)
		if err != nil {
			return
		}
		form, err := identitybinder.BindAssignUserRoleForm(r)
		if err != nil {
			h.app.navigator.SeeOtherWithFlash(w, r, "/users", cookieFlash, identitymessage.InvalidForm())
			return
		}
		if err := h.app.identity.assignRoleToUser.Execute(r.Context(), identitybinderToAssignRoleCommand(form, state.Claims.IDUser)); err != nil {
			h.app.logger.Error(r.Context(), "assign role to user", "error", err, "id_user", form.IDUser, "id_role", form.IDRole)
			h.app.navigator.SeeOtherWithFlash(w, r, "/users", cookieFlash, identitymessage.UserRoleAssignFailed(err))
			return
		}
		h.app.navigator.SeeOtherWithFlash(w, r, "/users", cookieFlash, identitymessage.UserRoleAssignSucceeded())
	})
}

func (h identityRoutes) RemoveRoleFromUserAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.app.requireRBACAdmin(w, r); err != nil {
			return
		}
		cmd, err := identitybinder.BindRemoveUserRoleForm(r)
		if err != nil {
			h.app.navigator.SeeOtherWithFlash(w, r, "/users", cookieFlash, identitymessage.InvalidForm())
			return
		}
		if err := h.app.identity.removeRoleFromUser.Execute(r.Context(), cmd); err != nil {
			h.app.logger.Error(r.Context(), "remove role from user", "error", err, "id_user", cmd.IDUser, "id_role", cmd.IDRole)
			h.app.navigator.SeeOtherWithFlash(w, r, "/users", cookieFlash, identitymessage.UserRoleRemoveFailed(err))
			return
		}
		h.app.navigator.SeeOtherWithFlash(w, r, "/users", cookieFlash, identitymessage.UserRoleRemoveSucceeded())
	})
}

func (h identityRoutes) GrantPermissionToRoleAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.app.requireRBACAdmin(w, r); err != nil {
			return
		}
		_, cmd, err := identitybinder.BindGrantRolePermissionForm(r)
		if err != nil {
			h.app.navigator.SeeOtherWithFlash(w, r, "/permissions", cookieFlash, identitymessage.InvalidForm())
			return
		}
		if err := h.app.identity.grantPermissionToRole.Execute(r.Context(), cmd); err != nil {
			h.app.logger.Error(r.Context(), "grant permission to role", "error", err, "id_role", cmd.IDRole, "id_permission", cmd.IDPermission)
			h.app.navigator.SeeOtherWithFlash(w, r, "/permissions", cookieFlash, identitymessage.RolePermissionGrantFailed(err))
			return
		}
		h.app.navigator.SeeOtherWithFlash(w, r, "/permissions", cookieFlash, identitymessage.RolePermissionGrantSucceeded())
	})
}

func (h identityRoutes) RevokePermissionFromRoleAction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.app.requireRBACAdmin(w, r); err != nil {
			return
		}
		_, cmd, err := identitybinder.BindRevokeRolePermissionForm(r)
		if err != nil {
			h.app.navigator.SeeOtherWithFlash(w, r, "/permissions", cookieFlash, identitymessage.InvalidForm())
			return
		}
		if err := h.app.identity.revokePermissionFromRole.Execute(r.Context(), cmd); err != nil {
			h.app.logger.Error(r.Context(), "revoke permission from role", "error", err, "id_role", cmd.IDRole, "id_permission", cmd.IDPermission)
			h.app.navigator.SeeOtherWithFlash(w, r, "/permissions", cookieFlash, identitymessage.RolePermissionRevokeFailed(err))
			return
		}
		h.app.navigator.SeeOtherWithFlash(w, r, "/permissions", cookieFlash, identitymessage.RolePermissionRevokeSucceeded())
	})
}

func identitybinderToAssignRoleCommand(form identitybinder.AssignUserRoleForm, assignedBy int64) identitycmd.AssignRoleToUserCommand {
	return identitycmd.AssignRoleToUserCommand{
		IDUser:     form.IDUser,
		IDRole:     form.IDRole,
		AssignedBy: assignedBy,
	}
}
