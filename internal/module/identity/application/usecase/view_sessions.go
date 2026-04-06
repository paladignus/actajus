// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	"github.com/paladignus/actajus/internal/module/identity/application/repository"
)

type ViewSessions struct {
	query repository.SessionQueryRepository
}

func NewViewSessions(query repository.SessionQueryRepository) ViewSessions {
	return ViewSessions{query: query}
}

func (uc ViewSessions) ByID(ctx context.Context, id int64) (*readmodel.SessionReadModel, error) {
	return uc.query.GetByID(ctx, id)
}

func (uc ViewSessions) Own(ctx context.Context, idUser int64) ([]readmodel.SessionReadModel, error) {
	return uc.query.ListByUser(ctx, idUser)
}

func (uc ViewSessions) Admin(ctx context.Context, filter repository.SessionQueryFilter) ([]readmodel.SessionReadModel, error) {
	return uc.query.ListAll(ctx, filter)
}
