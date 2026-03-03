// Package usecase
package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/paladignus/actajus/internal/module/identity/application/dto"
	"github.com/paladignus/actajus/internal/module/identity/application/mapper"
	identityrepo "github.com/paladignus/actajus/internal/module/identity/application/repository"
	"github.com/paladignus/actajus/internal/module/identity/application/service"
	"github.com/paladignus/actajus/internal/shared/domain/dispatcher"
	"github.com/paladignus/actajus/internal/shared/domain/unitofwork"
)

type AssignRoleToUser struct {
	repo        identityrepo.RoleUserAdminRepository
	cache       service.RBACCacheInvalidator
	index       service.RBACRoleUsersIndex
	mapper      *mapper.RBACAdminMapper
	uow         unitofwork.UnitOfWork
	dispatcher  *dispatcher.SimpleEventDispatcher
}

func NewAssignRoleToUser(
	repo identityrepo.RoleUserAdminRepository,
	cache service.RBACCacheInvalidator,
	index service.RBACRoleUsersIndex,
	mapper *mapper.RBACAdminMapper,
	uow unitofwork.UnitOfWork,
	dispatcher *dispatcher.SimpleEventDispatcher,
) AssignRoleToUser {
	return AssignRoleToUser{
		repo:        repo,
		cache:       cache,
		index:       index,
		mapper:      mapper,
		uow:         uow,
		dispatcher:  dispatcher,
	}
}

func (uc AssignRoleToUser) Execute(ctx context.Context, input dto.AssignRoleToUserCommand) error {
	norm, err := uc.mapper.AssignRoleInputToNormalized(input)
	if err != nil {
		return fmt.Errorf("invalid assign role data: %w", err)
	}

	// Se UnitOfWork estiver disponível, usa transação
	if uc.uow != nil {
		return uc.executeWithTransaction(ctx, norm)
	}

	// Fallback sem transação
	if err := uc.repo.AssignRole(ctx, norm.IDUser, norm.IDRole, norm.AssignedBy); err != nil {
		return err
	}
	uc.index.AddUserToRole(ctx, norm.IDRole, norm.IDUser.Value())
	uc.cache.InvalidateUser(ctx, norm.IDUser)
	
	// Publicar evento de domínio
	uc.publishRoleAssignedEvent(norm)
	
	return nil
}

func (uc AssignRoleToUser) executeWithTransaction(ctx context.Context, norm mapper.AssignRoleNormalized) error {
	// Nota: Esta implementação requer que o UnitOfWork seja do tipo concreto
	// Para uma implementação mais limpa, o repositório deveria aceitar um contexto com transação
	// Por enquanto, mantemos a estrutura atual como placeholder para futura refatoração
	
	// A implementação real dependeria de uma interface mais rica que expusesse o pgx.Tx
	// ou de repositórios que aceitem transação como parâmetro
	if err := uc.repo.AssignRole(ctx, norm.IDUser, norm.IDRole, norm.AssignedBy); err != nil {
		return err
	}
	
	uc.index.AddUserToRole(ctx, norm.IDRole, norm.IDUser.Value())
	uc.cache.InvalidateUser(ctx, norm.IDUser)

	// Publicar evento de domínio após operação bem-sucedida
	uc.publishRoleAssignedEvent(norm)

	return nil
}

func (uc AssignRoleToUser) publishRoleAssignedEvent(norm mapper.AssignRoleNormalized) {
	if uc.dispatcher == nil {
		return
	}
	
	uc.dispatcher.Publish(dispatcher.RoleAssignedEvent{
		ID:         norm.IDUser.Value(),
		IDRole:     norm.IDRole,
		AssignedBy: norm.AssignedBy,
		AssignedAt: time.Now(),
	})
}
