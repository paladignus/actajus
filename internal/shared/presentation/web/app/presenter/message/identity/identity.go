package identitymessage

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	webflash "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/flash"
)

func InvalidForm() webflash.Message {
	return webflash.Message{Kind: "error", Message: "Formulario invalido."}
}

func SessionLookupFailed() webflash.Message {
	return webflash.Message{Kind: "error", Message: "Nao foi possivel localizar a sessao."}
}

func SessionNotFound() webflash.Message {
	return webflash.Message{Kind: "error", Message: "Sessao nao encontrada."}
}

func SessionRevokeFailed() webflash.Message {
	return webflash.Message{Source: "sessions.revoke", Code: "sessions_revoke_failed", Kind: "error", Message: "Nao foi possivel revogar a sessao."}
}

func SessionRevokeSucceeded() webflash.Message {
	return webflash.Message{Source: "sessions.revoke", Code: "sessions_revoke_succeeded", Kind: "success", Message: "Sessao revogada com sucesso."}
}

func SessionRevokeAllFailed() webflash.Message {
	return webflash.Message{Source: "sessions.revoke_all", Code: "sessions_revoke_all_failed", Kind: "error", Message: "Nao foi possivel encerrar as sessoes."}
}

func SessionRevokeAllSucceeded() webflash.Message {
	return webflash.Message{Source: "sessions.revoke_all", Code: "sessions_revoke_all_succeeded", Kind: "success", Message: "Sessoes encerradas com sucesso."}
}

func RoleCreateFailed(err error) webflash.Message {
	return catalogFailure("roles.create", "roles_create_failed", err, "Nao foi possivel criar a role.")
}

func RoleCreateSucceeded() webflash.Message {
	return webflash.Message{Source: "roles.create", Code: "roles_create_succeeded", Kind: "success", Message: "Role criada com sucesso."}
}

func RoleUpdateFailed(err error) webflash.Message {
	return catalogFailure("roles.update", "roles_update_failed", err, "Nao foi possivel atualizar a role.")
}

func RoleUpdateSucceeded() webflash.Message {
	return webflash.Message{Source: "roles.update", Code: "roles_update_succeeded", Kind: "success", Message: "Role atualizada com sucesso."}
}

func RoleDeleteFailed(err error) webflash.Message {
	return catalogFailure("roles.delete", "roles_delete_failed", err, "Nao foi possivel excluir a role.")
}

func RoleDeleteSucceeded() webflash.Message {
	return webflash.Message{Source: "roles.delete", Code: "roles_delete_succeeded", Kind: "success", Message: "Role excluida com sucesso."}
}

func PermissionCreateFailed(err error) webflash.Message {
	return catalogFailure("permissions.create", "permissions_create_failed", err, "Nao foi possivel criar a permissao.")
}

func PermissionCreateSucceeded() webflash.Message {
	return webflash.Message{Source: "permissions.create", Code: "permissions_create_succeeded", Kind: "success", Message: "Permissao criada com sucesso."}
}

func PermissionUpdateFailed(err error) webflash.Message {
	return catalogFailure("permissions.update", "permissions_update_failed", err, "Nao foi possivel atualizar a permissao.")
}

func PermissionUpdateSucceeded() webflash.Message {
	return webflash.Message{Source: "permissions.update", Code: "permissions_update_succeeded", Kind: "success", Message: "Permissao atualizada com sucesso."}
}

func PermissionDeleteFailed(err error) webflash.Message {
	return catalogFailure("permissions.delete", "permissions_delete_failed", err, "Nao foi possivel excluir a permissao.")
}

func PermissionDeleteSucceeded() webflash.Message {
	return webflash.Message{Source: "permissions.delete", Code: "permissions_delete_succeeded", Kind: "success", Message: "Permissao excluida com sucesso."}
}

func UserRoleAssignFailed(err error) webflash.Message {
	return catalogFailure("users.roles.assign", "users_roles_assign_failed", err, "Nao foi possivel atribuir a role.")
}

func UserRoleAssignSucceeded() webflash.Message {
	return webflash.Message{Source: "users.roles.assign", Code: "users_roles_assign_succeeded", Kind: "success", Message: "Role atribuida com sucesso."}
}

func UserRoleRemoveFailed(err error) webflash.Message {
	return catalogFailure("users.roles.remove", "users_roles_remove_failed", err, "Nao foi possivel remover a role.")
}

func UserRoleRemoveSucceeded() webflash.Message {
	return webflash.Message{Source: "users.roles.remove", Code: "users_roles_remove_succeeded", Kind: "success", Message: "Role removida com sucesso."}
}

func RolePermissionGrantFailed(err error) webflash.Message {
	return catalogFailure("roles.permissions.grant", "roles_permissions_grant_failed", err, "Nao foi possivel conceder a permissao.")
}

func RolePermissionGrantSucceeded() webflash.Message {
	return webflash.Message{Source: "roles.permissions.grant", Code: "roles_permissions_grant_succeeded", Kind: "success", Message: "Permissao concedida com sucesso."}
}

func RolePermissionRevokeFailed(err error) webflash.Message {
	return catalogFailure("roles.permissions.revoke", "roles_permissions_revoke_failed", err, "Nao foi possivel revogar a permissao.")
}

func RolePermissionRevokeSucceeded() webflash.Message {
	return webflash.Message{Source: "roles.permissions.revoke", Code: "roles_permissions_revoke_succeeded", Kind: "success", Message: "Permissao revogada com sucesso."}
}

func catalogFailure(source, code string, err error, fallback string) webflash.Message {
	return webflash.Message{
		Source:  source,
		Code:    code,
		Kind:    "error",
		Message: mapCatalogError(err, fallback),
	}
}

func mapCatalogError(err error, fallback string) string {
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
