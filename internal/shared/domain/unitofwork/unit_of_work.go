// Package unitofwork
package unitofwork

import "context"

// UnitOfWork define um contrato para gerenciamento de transações
// Esta é a interface base que deve ser estendida por módulos específicos
type UnitOfWork interface {
	// Begin inicia uma nova transação
	Begin(ctx context.Context) error
	// Commit confirma a transação
	Commit(ctx context.Context) error
	// Rollback desfaz a transação
	Rollback(ctx context.Context) error
}

// Transactional define um contrato para operações que podem ser executadas dentro de uma transação
// Esta interface é útil para repositórios que precisam operar dentro de um UnitOfWork
type Transactional interface {
	// GetPgxPool retorna o pool de conexões ou a transação atual
	// Se uma transação estiver ativa, retorna um adapter que usa a transação
	// Caso contrário, retorna o pool de conexões padrão
	GetPgxPool() any
}
