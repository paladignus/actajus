// Pacakge rbac
package rbac

import (
	"context"

	"github.com/paladignus/actajus/internal/module/identity/domain"
	"github.com/paladignus/actajus/internal/module/identity/infrastructure/security"
)

type Checker struct {
	authz *security.AuthorizationService
}

func NewChecker(authz *security.AuthorizationService) *Checker {
	return &Checker{authz: authz}
}

func (c *Checker) HasPermission(ctx context.Context, idUser int64, permission string) (bool, error) {
	return c.authz.HasPermission(ctx, domain.IDUser(idUser), permission)
}
