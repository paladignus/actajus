// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/entity"
)

type Company struct {
	IDCompany   string
	CreateError error
}

func NewCompany() *Company {
	return &Company{}
}

func (e Company) Create(ctx context.Context, company entity.Company) (id string, err error) {
	return e.IDCompany, e.CreateError
}
