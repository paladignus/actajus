// Package domain
package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPerson(t *testing.T) {
	now := time.Now()
	p, err := NewPersonBuilder().
		WithFirstName("John").
		WithLastName("Doe").
		WithGender(1).
		WithBirthday("01/01/1990").
		WithCreatedAt(now).
		WithUpdatedAt(now).
		Build()
	t.Run("should create a person", func(t *testing.T) {
		assert.NoError(t, err)
		assert.Equal(t, "John", p.FirstName().Value())
		assert.Equal(t, "Doe", p.LastName().Value())
		assert.Equal(t, uint(1), p.IDGender())
		assert.Equal(t, "01/01/1990", p.Birthday().Value())
		assert.Equal(t, now, p.CreatedAt())
		assert.Equal(t, now, p.UpdatedAt())
	})
	t.Run("should set id and delete a person", func(t *testing.T) {
		err := p.SetID(1)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), p.ID())
		err = p.Delete()
		assert.NoError(t, err)
		assert.True(t, p.IsDeleted())
	})
	t.Run("should error when set id", func(t *testing.T) {
		err := p.SetID(2)
		assert.Error(t, err)
		assert.Equal(t, "id is already set", err.Error())
		// Criar uma nova pessoa válida para testar o caso de ID zero
		newPerson, err := NewPersonBuilder().
			WithFirstName("Jane").
			WithLastName("Doe").
			WithGender(2).
			WithBirthday("15/05/1995").
			Build()
		assert.NoError(t, err)
		err = newPerson.SetID(0)
		assert.Error(t, err)
		assert.Equal(t, "id must be greater than 0", err.Error())
	})
	t.Run("should error when deleting a deleted person", func(t *testing.T) {
		// Criar uma nova pessoa com ID mas não deletada
		personWithID, err := NewPersonBuilder().
			WithID(1).
			WithFirstName("Jane").
			WithLastName("Smith").
			WithGender(2).
			WithBirthday("01/01/1990").
			Build()
		assert.NoError(t, err)

		// Agora deletar para simular o estado deletado
		err = personWithID.Delete()
		assert.NoError(t, err)
		assert.True(t, personWithID.IsDeleted())

		// Tentar deletar novamente deve retornar "person already deleted"
		err = personWithID.Delete()
		assert.Error(t, err)
		assert.Equal(t, "person already deleted", err.Error())

		// Criar uma nova pessoa completamente válida com DeletedAt já definido
		validPerson, err := NewPersonBuilder().
			WithID(1).
			WithFirstName("Jane").
			WithLastName("Doe").
			WithGender(2).
			WithBirthday("01/01/1990").
			WithDeletedAt(&now).
			Build()
		assert.NoError(t, err)
		assert.Equal(t, now, *validPerson.DeletedAt())
		assert.True(t, validPerson.IsDeleted())
		err = validPerson.Delete()
		assert.Error(t, err)
		assert.Equal(t, "person already deleted", err.Error())
	})
}
