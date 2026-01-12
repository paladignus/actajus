// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/entity"
)

type Enterprise struct {
	IDEnterprise string
	CreateError  error
}

func NewEnterprise() *Enterprise {
	return &Enterprise{}
}

func (e Enterprise) Create(ctx context.Context, enterprise entity.Enterprise) (id string, err error) {
	return e.IDEnterprise, e.CreateError
}
