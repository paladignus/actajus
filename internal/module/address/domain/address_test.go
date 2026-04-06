// Package domain
package domain

import (
	"testing"
	"time"

	"github.com/paladignus/actajus/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

func TestAddressBuilder_Build_Success(t *testing.T) {
	t.Parallel()

	now := time.Now()
	address, err := NewAddressBuilder().
		WithID(1).
		WithZIP("01000-000").
		WithTitle("Casa").
		WithStreet("Rua das Flores").
		WithNumber(123).
		WithComplement("Apto 101").
		WithReference("Próximo à praça").
		WithNeighborhood("Centro").
		WithCity("São Paulo").
		WithState("SP").
		WithCountry("Brasil").
		WithCreatedAt(now).
		WithUpdatedAt(now).
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, address)
	assert.Equal(t, int64(1), address.ID().Value())
	assert.Equal(t, "01000-000", address.ZIP().Value())
	assert.Equal(t, "Casa", address.Title().Value())
	assert.Equal(t, "Rua das Flores", address.Street().Value())
	assert.Equal(t, uint(123), address.Number())
	assert.False(t, address.IsDeleted())
}

func TestAddressBuilder_Build_InvalidZIP(t *testing.T) {
	t.Parallel()

	address, err := NewAddressBuilder().
		WithID(1).
		WithZIP("invalid").
		WithTitle("Casa").
		WithStreet("Rua das Flores").
		WithNumber(123).
		WithNeighborhood("Centro").
		WithCity("São Paulo").
		WithState("SP").
		WithCountry("Brasil").
		Build()

	assert.Error(t, err)
	assert.Nil(t, address)
	fieldErr, ok := err.(*domain.FieldError)
	assert.True(t, ok)
	assert.Equal(t, "zip", fieldErr.Field)
}

func TestAddressBuilder_Build_MissingNumber(t *testing.T) {
	t.Parallel()

	address, err := NewAddressBuilder().
		WithID(1).
		WithZIP("01000-000").
		WithTitle("Casa").
		WithStreet("Rua das Flores").
		WithNumber(0).
		WithNeighborhood("Centro").
		WithCity("São Paulo").
		WithState("SP").
		WithCountry("Brasil").
		Build()

	assert.Error(t, err)
	assert.Nil(t, address)
	fieldErr, ok := err.(*domain.FieldError)
	assert.True(t, ok)
	assert.Equal(t, "number", fieldErr.Field)
}

func TestAddress_SetID_Success(t *testing.T) {
	t.Parallel()

	address, err := NewAddressBuilder().
		WithZIP("01000-000").
		WithTitle("Casa").
		WithStreet("Rua das Flores").
		WithNumber(123).
		WithNeighborhood("Centro").
		WithCity("São Paulo").
		WithState("SP").
		WithCountry("Brasil").
		Build()
	assert.NoError(t, err)

	err = address.SetID(1)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), address.ID().Value())
}

func TestAddress_SetID_AlreadySet(t *testing.T) {
	t.Parallel()

	address, err := NewAddressBuilder().
		WithID(1).
		WithZIP("01000-000").
		WithTitle("Casa").
		WithStreet("Rua das Flores").
		WithNumber(123).
		WithNeighborhood("Centro").
		WithCity("São Paulo").
		WithState("SP").
		WithCountry("Brasil").
		Build()
	assert.NoError(t, err)

	err = address.SetID(2)

	assert.Error(t, err)
	fieldErr, ok := err.(*domain.FieldError)
	assert.True(t, ok)
	assert.Equal(t, "id", fieldErr.Field)
}

func TestAddress_Delete_Success(t *testing.T) {
	t.Parallel()

	address, err := NewAddressBuilder().
		WithID(1).
		WithZIP("01000-000").
		WithTitle("Casa").
		WithStreet("Rua das Flores").
		WithNumber(123).
		WithNeighborhood("Centro").
		WithCity("São Paulo").
		WithState("SP").
		WithCountry("Brasil").
		Build()
	assert.NoError(t, err)

	oldUpdatedAt := address.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	err = address.Delete()

	assert.NoError(t, err)
	assert.True(t, address.IsDeleted())
	assert.NotNil(t, address.DeletedAt())
	assert.True(t, address.UpdatedAt().After(oldUpdatedAt))
}

func TestAddress_Delete_AlreadyDeleted(t *testing.T) {
	t.Parallel()

	address, err := NewAddressBuilder().
		WithID(1).
		WithZIP("01000-000").
		WithTitle("Casa").
		WithStreet("Rua das Flores").
		WithNumber(123).
		WithNeighborhood("Centro").
		WithCity("São Paulo").
		WithState("SP").
		WithCountry("Brasil").
		Build()
	assert.NoError(t, err)

	err = address.Delete()
	assert.NoError(t, err)

	err = address.Delete()

	assert.Error(t, err)
	fieldErr, ok := err.(*domain.FieldError)
	assert.True(t, ok)
	assert.Equal(t, "deleted_at", fieldErr.Field)
}

func TestAddress_Restore_Success(t *testing.T) {
	t.Parallel()

	address, err := NewAddressBuilder().
		WithID(1).
		WithZIP("01000-000").
		WithTitle("Casa").
		WithStreet("Rua das Flores").
		WithNumber(123).
		WithNeighborhood("Centro").
		WithCity("São Paulo").
		WithState("SP").
		WithCountry("Brasil").
		Build()
	assert.NoError(t, err)

	err = address.Delete()
	assert.NoError(t, err)
	assert.True(t, address.IsDeleted())

	oldUpdatedAt := address.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	address.Restore()

	assert.False(t, address.IsDeleted())
	assert.Nil(t, address.DeletedAt())
	assert.True(t, address.UpdatedAt().After(oldUpdatedAt))
}

func TestAddress_Restore_NotDeleted(t *testing.T) {
	t.Parallel()

	address, err := NewAddressBuilder().
		WithID(1).
		WithZIP("01000-000").
		WithTitle("Casa").
		WithStreet("Rua das Flores").
		WithNumber(123).
		WithNeighborhood("Centro").
		WithCity("São Paulo").
		WithState("SP").
		WithCountry("Brasil").
		Build()
	assert.NoError(t, err)
	assert.False(t, address.IsDeleted())

	address.Restore()

	assert.False(t, address.IsDeleted())
	assert.Nil(t, address.DeletedAt())
}

func TestAddress_UpdateZIP_Success(t *testing.T) {
	t.Parallel()

	address, err := NewAddressBuilder().
		WithID(1).
		WithZIP("01000-000").
		WithTitle("Casa").
		WithStreet("Rua das Flores").
		WithNumber(123).
		WithNeighborhood("Centro").
		WithCity("São Paulo").
		WithState("SP").
		WithCountry("Brasil").
		Build()
	assert.NoError(t, err)

	oldUpdatedAt := address.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	err = address.UpdateZIP("20000-000")

	assert.NoError(t, err)
	assert.Equal(t, "20000-000", address.ZIP().Value())
	assert.True(t, address.UpdatedAt().After(oldUpdatedAt))
}

func TestAddress_UpdateStreet_Success(t *testing.T) {
	t.Parallel()

	address, err := NewAddressBuilder().
		WithID(1).
		WithZIP("01000-000").
		WithTitle("Casa").
		WithStreet("Rua das Flores").
		WithNumber(123).
		WithNeighborhood("Centro").
		WithCity("São Paulo").
		WithState("SP").
		WithCountry("Brasil").
		Build()
	assert.NoError(t, err)

	oldUpdatedAt := address.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	err = address.UpdateStreet("Avenida Paulista")

	assert.NoError(t, err)
	assert.Equal(t, "Avenida Paulista", address.Street().Value())
	assert.True(t, address.UpdatedAt().After(oldUpdatedAt))
}

func TestAddress_UpdateNumber(t *testing.T) {
	t.Parallel()

	address, err := NewAddressBuilder().
		WithID(1).
		WithZIP("01000-000").
		WithTitle("Casa").
		WithStreet("Rua das Flores").
		WithNumber(123).
		WithNeighborhood("Centro").
		WithCity("São Paulo").
		WithState("SP").
		WithCountry("Brasil").
		Build()
	assert.NoError(t, err)

	oldUpdatedAt := address.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	address.UpdateNumber(456)

	assert.Equal(t, uint(456), address.Number())
	assert.True(t, address.UpdatedAt().After(oldUpdatedAt))
}

func TestAddress_FullAddress(t *testing.T) {
	t.Parallel()

	address, err := NewAddressBuilder().
		WithID(1).
		WithZIP("01000-000").
		WithTitle("Casa").
		WithStreet("Rua das Flores").
		WithNumber(123).
		WithComplement("Apto 101").
		WithNeighborhood("Centro").
		WithCity("São Paulo").
		WithState("SP").
		WithCountry("Brasil").
		Build()
	assert.NoError(t, err)

	full := address.FullAddress()

	assert.Contains(t, full, "Rua das Flores")
	assert.Contains(t, full, "Centro")
	assert.Contains(t, full, "São Paulo")
	assert.Contains(t, full, "SP")
	assert.Contains(t, full, "01000-000")
}
