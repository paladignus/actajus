// Package repository
package repository

import (
	"context"
)

// RoleUserQueryRepository é responsável por consultar usuários de uma role
// com estratégia de cache embutida (padrão Cache-Aside com fallback).
//
// Esta interface abstrai a complexidade de orquestração entre cache (index)
// e repositório (database), mantendo os use cases limpos e focados na regra
// de negócio.
type RoleUserQueryRepository interface {
	// ListUserIDsByRole retorna todos os usuários de uma role,
	// usando cache quando disponível e fazendo fallback para o banco.
	ListUserIDsByRole(ctx context.Context, rid int16) ([]int64, error)
}
