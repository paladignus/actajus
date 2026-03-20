// Package domain
package domain

import (
	"testing"
	"time"

	"github.com/paladignus/actajus/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

func TestCompanyBuilder_Build_Success(t *testing.T) {
	t.Parallel()

	now := time.Now()
	company, err := NewCompanyBuilder().
		WithID(1).
		WithRegisteredBy(123).
		WithRegisteredByName("João Silva").
		WithName("Empresa X LTDA").
		WithTradeName("Empresa X").
		WithCNPJ("10.123.456/0001-00").
		WithIDAddress(10).
		WithIDPhone(20).
		WithIDEmail(30).
		WithIDSocialMedia(40).
		WithCreatedAt(now).
		WithUpdatedAt(now).
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, company)
	assert.Equal(t, int64(1), company.ID().Value())
	assert.Equal(t, int64(123), company.RegisteredBy().Value())
	assert.Equal(t, "Empresa X LTDA", company.Name().Value())
	assert.Equal(t, "Empresa X", company.TradeName().Value())
	assert.Equal(t, "10.123.456/0001-00", company.CNPJ().Value())
	assert.False(t, company.IsDeleted())
}

func TestCompanyBuilder_Build_InvalidCNPJ(t *testing.T) {
	t.Parallel()

	company, err := NewCompanyBuilder().
		WithID(1).
		WithRegisteredBy(123).
		WithName("Empresa X").
		WithCNPJ("invalid").
		Build()

	assert.Error(t, err)
	assert.Nil(t, company)
	fieldErr, ok := err.(*domain.FieldError)
	assert.True(t, ok)
	assert.Equal(t, "cnpj", fieldErr.Field)
}

func TestCompanyBuilder_Build_MissingRegisteredBy(t *testing.T) {
	t.Parallel()

	// registeredBy=0 e id=0 (não definido) = erro
	company, err := NewCompanyBuilder().
		WithName("Empresa X").
		WithCNPJ("10.123.456/0001-00").
		Build()

	assert.Error(t, err)
	assert.Nil(t, company)
	fieldErr, ok := err.(*domain.FieldError)
	assert.True(t, ok)
	assert.Equal(t, "registered_by", fieldErr.Field)
}

func TestCompany_SetID_Success(t *testing.T) {
	t.Parallel()

	company, err := NewCompanyBuilder().
		WithRegisteredBy(123).
		WithName("Empresa X").
		WithCNPJ("10.123.456/0001-00").
		Build()
	assert.NoError(t, err)

	err = company.SetID(1)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), company.ID().Value())
}

func TestCompany_SetID_AlreadySet(t *testing.T) {
	t.Parallel()

	company, err := NewCompanyBuilder().
		WithID(1).
		WithRegisteredBy(123).
		WithName("Empresa X").
		WithCNPJ("10.123.456/0001-00").
		Build()
	assert.NoError(t, err)

	err = company.SetID(2)

	assert.Error(t, err)
	fieldErr, ok := err.(*domain.FieldError)
	assert.True(t, ok)
	assert.Equal(t, "id", fieldErr.Field)
}

func TestCompany_Delete_Success(t *testing.T) {
	t.Parallel()

	company, err := NewCompanyBuilder().
		WithID(1).
		WithRegisteredBy(123).
		WithName("Empresa X").
		WithCNPJ("10.123.456/0001-00").
		Build()
	assert.NoError(t, err)

	oldUpdatedAt := company.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	err = company.Delete()

	assert.NoError(t, err)
	assert.True(t, company.IsDeleted())
	assert.NotNil(t, company.DeletedAt())
	assert.True(t, company.UpdatedAt().After(oldUpdatedAt))
}

func TestCompany_Delete_AlreadyDeleted(t *testing.T) {
	t.Parallel()

	company, err := NewCompanyBuilder().
		WithID(1).
		WithRegisteredBy(123).
		WithName("Empresa X").
		WithCNPJ("10.123.456/0001-00").
		Build()
	assert.NoError(t, err)

	err = company.Delete()
	assert.NoError(t, err)

	err = company.Delete()

	assert.Error(t, err)
	fieldErr, ok := err.(*domain.FieldError)
	assert.True(t, ok)
	assert.Equal(t, "deleted_at", fieldErr.Field)
}

func TestCompany_Restore_Success(t *testing.T) {
	t.Parallel()

	company, err := NewCompanyBuilder().
		WithID(1).
		WithRegisteredBy(123).
		WithName("Empresa X").
		WithCNPJ("10.123.456/0001-00").
		Build()
	assert.NoError(t, err)

	err = company.Delete()
	assert.NoError(t, err)
	assert.True(t, company.IsDeleted())

	oldUpdatedAt := company.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	company.Restore()

	assert.False(t, company.IsDeleted())
	assert.Nil(t, company.DeletedAt())
	assert.True(t, company.UpdatedAt().After(oldUpdatedAt))
}

func TestCompany_UpdateName_Success(t *testing.T) {
	t.Parallel()

	company, err := NewCompanyBuilder().
		WithID(1).
		WithRegisteredBy(123).
		WithName("Empresa X").
		WithCNPJ("10.123.456/0001-00").
		Build()
	assert.NoError(t, err)

	oldUpdatedAt := company.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	err = company.UpdateName("Empresa Y LTDA")

	assert.NoError(t, err)
	assert.Equal(t, "Empresa Y LTDA", company.Name().Value())
	assert.True(t, company.UpdatedAt().After(oldUpdatedAt))
}

func TestCompany_UpdateName_Empty(t *testing.T) {
	t.Parallel()

	company, err := NewCompanyBuilder().
		WithID(1).
		WithRegisteredBy(123).
		WithName("Empresa X").
		WithCNPJ("10.123.456/0001-00").
		Build()
	assert.NoError(t, err)

	err = company.UpdateName("")

	assert.Error(t, err)
	fieldErr, ok := err.(*domain.FieldError)
	assert.True(t, ok)
	assert.Equal(t, "name", fieldErr.Field)
}

func TestCompany_UpdateTradeName_Success(t *testing.T) {
	t.Parallel()

	company, err := NewCompanyBuilder().
		WithID(1).
		WithRegisteredBy(123).
		WithTradeName("Empresa X").
		WithCNPJ("10.123.456/0001-00").
		Build()
	assert.NoError(t, err)

	oldUpdatedAt := company.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	err = company.UpdateTradeName("Empresa Y")

	assert.NoError(t, err)
	assert.Equal(t, "Empresa Y", company.TradeName().Value())
	assert.True(t, company.UpdatedAt().After(oldUpdatedAt))
}

func TestCompany_UpdateCNPJ_Success(t *testing.T) {
	t.Parallel()

	company, err := NewCompanyBuilder().
		WithID(1).
		WithRegisteredBy(123).
		WithName("Empresa X").
		WithCNPJ("10.123.456/0001-00").
		Build()
	assert.NoError(t, err)

	oldUpdatedAt := company.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	err = company.UpdateCNPJ("25.461.378/0001-13")

	assert.NoError(t, err)
	assert.Equal(t, "25.461.378/0001-13", company.CNPJ().Value())
	assert.True(t, company.UpdatedAt().After(oldUpdatedAt))
}

func TestCompany_SetAddress_Success(t *testing.T) {
	t.Parallel()

	company, err := NewCompanyBuilder().
		WithID(1).
		WithRegisteredBy(123).
		WithName("Empresa X").
		WithCNPJ("10.123.456/0001-00").
		Build()
	assert.NoError(t, err)

	assert.False(t, company.HasAddress())

	oldUpdatedAt := company.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	company.SetAddress(100)

	assert.True(t, company.HasAddress())
	assert.Equal(t, int64(100), company.IDAddress().Value())
	assert.True(t, company.UpdatedAt().After(oldUpdatedAt))
}

func TestCompany_SetPhone_Success(t *testing.T) {
	t.Parallel()

	company, err := NewCompanyBuilder().
		WithID(1).
		WithRegisteredBy(123).
		WithName("Empresa X").
		WithCNPJ("10.123.456/0001-00").
		Build()
	assert.NoError(t, err)

	assert.False(t, company.HasPhone())

	oldUpdatedAt := company.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	company.SetPhone(200)

	assert.True(t, company.HasPhone())
	assert.Equal(t, int64(200), company.IDPhone().Value())
	assert.True(t, company.UpdatedAt().After(oldUpdatedAt))
}

func TestCompany_SetEmail_Success(t *testing.T) {
	t.Parallel()

	company, err := NewCompanyBuilder().
		WithID(1).
		WithRegisteredBy(123).
		WithName("Empresa X").
		WithCNPJ("10.123.456/0001-00").
		Build()
	assert.NoError(t, err)

	assert.False(t, company.HasEmail())

	oldUpdatedAt := company.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	company.SetEmail(300)

	assert.True(t, company.HasEmail())
	assert.Equal(t, int64(300), company.IDEmail().Value())
	assert.True(t, company.UpdatedAt().After(oldUpdatedAt))
}

func TestCompany_SetSocialMedia_Success(t *testing.T) {
	t.Parallel()

	company, err := NewCompanyBuilder().
		WithID(1).
		WithRegisteredBy(123).
		WithName("Empresa X").
		WithCNPJ("10.123.456/0001-00").
		Build()
	assert.NoError(t, err)

	assert.False(t, company.HasSocialMedia())

	oldUpdatedAt := company.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	company.SetSocialMedia(400)

	assert.True(t, company.HasSocialMedia())
	assert.Equal(t, int64(400), company.IDSocialMedia().Value())
	assert.True(t, company.UpdatedAt().After(oldUpdatedAt))
}
