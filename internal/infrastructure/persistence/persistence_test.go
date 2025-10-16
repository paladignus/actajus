package persistence

import (
	"testing"

	"github.com/paladignus/actajus/internal/infrastructure/database"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPersistence(t *testing.T) {
	t.Run("ensure creates persistence instance successfully", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()
		db := &database.DB{Pool: mock}
		persistence := NewPersistence(db)
		assert.NotNil(t, persistence)
		assert.NotNil(t, persistence.db)
		assert.Equal(t, db, persistence.db)
	})

	t.Run("ensure creates persistence with nil db", func(t *testing.T) {
		persistence := NewPersistence(nil)
		assert.NotNil(t, persistence)
		assert.Nil(t, persistence.db)
	})
}

func TestPersistence_Account(t *testing.T) {
	t.Run("should returns account repository instance", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()
		db := &database.DB{Pool: mock}
		persistence := NewPersistence(db)
		accountRepo := persistence.Account()
		assert.NotNil(t, accountRepo)
		assert.Equal(t, db, accountRepo.db)
	})

	t.Run("should returns new account instance on each call", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()
		db := &database.DB{Pool: mock}
		persistence := NewPersistence(db)
		account1 := persistence.Account()
		account2 := persistence.Account()
		// Verifica que ambas instâncias usam o mesmo db
		assert.Equal(t, account1.db, account2.db)
		// Note: Como Account é um struct (não ponteiro), cada chamada
		// retorna uma nova cópia, mas isso não é um problema pois
		// compartilham a mesma referência de db
		assert.NotNil(t, account1)
		assert.NotNil(t, account2)
	})

	t.Run("should account repository shares same db connection", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()
		db := &database.DB{Pool: mock}
		persistence := NewPersistence(db)
		account := persistence.Account()
		// Verifica que o repositório de account tem acesso ao mesmo pool de conexões
		assert.Same(t, persistence.db, account.db)
		assert.Same(t, db.Pool, account.db.Pool)
	})
}
