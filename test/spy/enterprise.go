// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/entity"
)

type Enterprise struct {
	CreateError error
}

func NewEnterprise() *Enterprise {
	return &Enterprise{}
}

func (e Enterprise) Create(ctx context.Context, enterprise entity.Enterprise) error {
	return e.CreateError
}
