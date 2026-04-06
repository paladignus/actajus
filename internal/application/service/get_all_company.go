// Package service
package service

import (
	"context"

	"github.com/paladignus/actajus/internal/application/readmodel"
)

type IGetAllCompany interface {
	Execute(ctx context.Context) (output []readmodel.CompanyReadModel, err error)
}
