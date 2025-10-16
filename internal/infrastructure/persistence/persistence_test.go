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
