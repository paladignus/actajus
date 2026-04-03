package identitybinder

import (
	"net/http"
	"strings"

	identitycmd "github.com/paladignus/actajus/internal/module/identity/application/command"
	identityrepo "github.com/paladignus/actajus/internal/module/identity/application/repository"
	pathbinder "github.com/paladignus/actajus/internal/shared/presentation/web/app/binder/http/path"
)

type RoleCatalogForm struct {
	Name        string
	Description string
}

type PermissionCatalogForm struct {
	Resource    string
	Action      string
	Description string
}

type AssignUserRoleForm struct {
	IDUser int64
	IDRole int16
}

type RolePermissionForm struct {
	IDRole       int16
	IDPermission int16
}

type RevokeSessionForm struct {
	IDSession int64
}

type RevokeAllSessionsForm struct {
	IDUser *int64
}

func BindCatalogUserFilter(r *http.Request, limit int) identityrepo.CatalogUserFilter {
	return identityrepo.CatalogUserFilter{
		Query: strings.TrimSpace(r.URL.Query().Get("q")),
		Limit: limit,
	}
}

func BindCatalogPermissionFilter(r *http.Request, limit int) identityrepo.CatalogPermissionFilter {
	return identityrepo.CatalogPermissionFilter{
		Query: strings.TrimSpace(r.URL.Query().Get("q")),
		Limit: limit,
	}
}

func BindCreateRoleForm(r *http.Request) (RoleCatalogForm, identitycmd.CreateRoleCommand, error) {
	if err := r.ParseForm(); err != nil {
		return RoleCatalogForm{}, identitycmd.CreateRoleCommand{}, err
	}
	form := RoleCatalogForm{
		Name:        strings.TrimSpace(r.FormValue("name")),
		Description: strings.TrimSpace(r.FormValue("description")),
	}
	return form, identitycmd.CreateRoleCommand{
		Name:        form.Name,
		Description: form.Description,
	}, nil
}

func BindUpdateRoleForm(r *http.Request) (RoleCatalogForm, identitycmd.UpdateRoleCommand, error) {
	if err := r.ParseForm(); err != nil {
		return RoleCatalogForm{}, identitycmd.UpdateRoleCommand{}, err
	}
	id, err := pathbinder.BindInt16Path(r, "id")
	if err != nil {
		return RoleCatalogForm{}, identitycmd.UpdateRoleCommand{}, err
	}
	form := RoleCatalogForm{
		Name:        strings.TrimSpace(r.FormValue("name")),
		Description: strings.TrimSpace(r.FormValue("description")),
	}
	return form, identitycmd.UpdateRoleCommand{
		ID:          id,
		Name:        form.Name,
		Description: form.Description,
	}, nil
}

func BindDeleteRolePath(r *http.Request) (identitycmd.DeleteRoleCommand, error) {
	id, err := pathbinder.BindInt16Path(r, "id")
	if err != nil {
		return identitycmd.DeleteRoleCommand{}, err
	}
	return identitycmd.DeleteRoleCommand{ID: id}, nil
}

func BindCreatePermissionForm(r *http.Request) (PermissionCatalogForm, identitycmd.CreatePermissionCommand, error) {
	if err := r.ParseForm(); err != nil {
		return PermissionCatalogForm{}, identitycmd.CreatePermissionCommand{}, err
	}
	form := PermissionCatalogForm{
		Resource:    strings.TrimSpace(r.FormValue("resource")),
		Action:      strings.TrimSpace(r.FormValue("action")),
		Description: strings.TrimSpace(r.FormValue("description")),
	}
	return form, identitycmd.CreatePermissionCommand{
		Resource:    form.Resource,
		Action:      form.Action,
		Description: form.Description,
	}, nil
}

func BindUpdatePermissionForm(r *http.Request) (PermissionCatalogForm, identitycmd.UpdatePermissionCommand, error) {
	if err := r.ParseForm(); err != nil {
		return PermissionCatalogForm{}, identitycmd.UpdatePermissionCommand{}, err
	}
	id, err := pathbinder.BindInt16Path(r, "id")
	if err != nil {
		return PermissionCatalogForm{}, identitycmd.UpdatePermissionCommand{}, err
	}
	form := PermissionCatalogForm{
		Resource:    strings.TrimSpace(r.FormValue("resource")),
		Action:      strings.TrimSpace(r.FormValue("action")),
		Description: strings.TrimSpace(r.FormValue("description")),
	}
	return form, identitycmd.UpdatePermissionCommand{
		ID:          id,
		Resource:    form.Resource,
		Action:      form.Action,
		Description: form.Description,
	}, nil
}

func BindDeletePermissionPath(r *http.Request) (identitycmd.DeletePermissionCommand, error) {
	id, err := pathbinder.BindInt16Path(r, "id")
	if err != nil {
		return identitycmd.DeletePermissionCommand{}, err
	}
	return identitycmd.DeletePermissionCommand{ID: id}, nil
}

func BindAssignUserRoleForm(r *http.Request) (AssignUserRoleForm, error) {
	if err := r.ParseForm(); err != nil {
		return AssignUserRoleForm{}, err
	}
	idUser, err := pathbinder.ParseInt64Value(r.FormValue("id_user"), "invalid user id")
	if err != nil {
		return AssignUserRoleForm{}, err
	}
	idRole, err := pathbinder.ParseInt16Value(r.FormValue("id_role"), "invalid role id")
	if err != nil {
		return AssignUserRoleForm{}, err
	}
	return AssignUserRoleForm{IDUser: idUser, IDRole: idRole}, nil
}

func BindRemoveUserRoleForm(r *http.Request) (identitycmd.RemoveRoleFromUserCommand, error) {
	form, err := BindAssignUserRoleForm(r)
	if err != nil {
		return identitycmd.RemoveRoleFromUserCommand{}, err
	}
	return identitycmd.RemoveRoleFromUserCommand{IDUser: form.IDUser, IDRole: form.IDRole}, nil
}

func BindGrantRolePermissionForm(r *http.Request) (RolePermissionForm, identitycmd.GrantPermissionToRoleCommand, error) {
	if err := r.ParseForm(); err != nil {
		return RolePermissionForm{}, identitycmd.GrantPermissionToRoleCommand{}, err
	}
	idRole, err := pathbinder.ParseInt16Value(r.FormValue("id_role"), "invalid role id")
	if err != nil {
		return RolePermissionForm{}, identitycmd.GrantPermissionToRoleCommand{}, err
	}
	idPermission, err := pathbinder.ParseInt16Value(r.FormValue("id_permission"), "invalid permission id")
	if err != nil {
		return RolePermissionForm{}, identitycmd.GrantPermissionToRoleCommand{}, err
	}
	form := RolePermissionForm{IDRole: idRole, IDPermission: idPermission}
	return form, identitycmd.GrantPermissionToRoleCommand{
		IDRole:       form.IDRole,
		IDPermission: form.IDPermission,
	}, nil
}

func BindRevokeRolePermissionForm(r *http.Request) (RolePermissionForm, identitycmd.RevokePermissionFromRoleCommand, error) {
	if err := r.ParseForm(); err != nil {
		return RolePermissionForm{}, identitycmd.RevokePermissionFromRoleCommand{}, err
	}
	idRole, err := pathbinder.ParseInt16Value(r.FormValue("id_role"), "invalid role id")
	if err != nil {
		return RolePermissionForm{}, identitycmd.RevokePermissionFromRoleCommand{}, err
	}
	idPermission, err := pathbinder.ParseInt16Value(r.FormValue("id_permission"), "invalid permission id")
	if err != nil {
		return RolePermissionForm{}, identitycmd.RevokePermissionFromRoleCommand{}, err
	}
	form := RolePermissionForm{IDRole: idRole, IDPermission: idPermission}
	return form, identitycmd.RevokePermissionFromRoleCommand{
		IDRole:       form.IDRole,
		IDPermission: form.IDPermission,
	}, nil
}

func BindRevokeSessionForm(r *http.Request) (RevokeSessionForm, identitycmd.RevokeSessionCommand, error) {
	if err := r.ParseForm(); err != nil {
		return RevokeSessionForm{}, identitycmd.RevokeSessionCommand{}, err
	}
	idSession, err := pathbinder.ParseInt64Value(r.FormValue("id_session"), "invalid session id")
	if err != nil {
		return RevokeSessionForm{}, identitycmd.RevokeSessionCommand{}, err
	}
	form := RevokeSessionForm{IDSession: idSession}
	return form, identitycmd.RevokeSessionCommand{IDSession: idSession}, nil
}

func BindRevokeAllSessionsForm(r *http.Request) (RevokeAllSessionsForm, *identitycmd.LogoutAllCommand, error) {
	if err := r.ParseForm(); err != nil {
		return RevokeAllSessionsForm{}, nil, err
	}
	raw := strings.TrimSpace(r.FormValue("id_user"))
	if raw == "" {
		return RevokeAllSessionsForm{}, nil, nil
	}
	idUser, err := pathbinder.ParseInt64Value(raw, "invalid user id")
	if err != nil {
		return RevokeAllSessionsForm{}, nil, err
	}
	form := RevokeAllSessionsForm{IDUser: &idUser}
	cmd := identitycmd.LogoutAllCommand{IDUser: idUser}
	return form, &cmd, nil
}
