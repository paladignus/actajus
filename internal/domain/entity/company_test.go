// Package entity
package entity

// func TestCompany(t *testing.T) {
// 	input := dto.Company{
// 		RegisteredBy: 1, Name: "Actajus", TradeName: "Actajus", CNPJ: "10.123.456/0001-00",
// 	}
// 	t.Run("should return an Company", func(t *testing.T) {
// 		Company, err := NewCompany(input)
// 		assert.NoError(t, err)
// 		assert.Equal(t, "Actajus", Company.Name.Value())
// 		assert.Equal(t, "Actajus", Company.TradeName.Value())
// 		assert.Equal(t, "10.123.456/0001-00", Company.CNPJ.Value())
// 		assert.True(t, Company.CreatedAt.Equal(Company.UpdatedAt))
// 		assert.True(t, time.Now().After(Company.CreatedAt))
// 		assert.True(t, time.Now().After(Company.UpdatedAt))
// 		assert.Nil(t, Company.DeletedAt)
// 	})
//
// 	t.Run("should return false if the registered by is invalid", func(t *testing.T) {
// 		input.RegisteredBy = 0
// 		_, err := NewCompany(input)
// 		assert.Error(t, err)
// 		assert.Equal(t, exception.ErrInvalidRegisteredBy, err)
// 	})
//
// 	t.Run("should return false if the name is invalid", func(t *testing.T) {
// 		input.RegisteredBy = 1
// 		input.Name = "11"
// 		_, err := NewCompany(input)
// 		assert.Error(t, err)
// 		assert.Equal(t, exception.ErrInvalidName, err)
// 	})
//
// 	t.Run("should return false if the trade name is invalid", func(t *testing.T) {
// 		input.Name = "Actajus"
// 		input.TradeName = "11"
// 		_, err := NewCompany(input)
// 		assert.Error(t, err)
// 		assert.Equal(t, exception.ErrInvalidTradeName, err)
// 	})
//
// 	t.Run("should return false if the cnpj is invalid", func(t *testing.T) {
// 		input.TradeName = "Trade Actajus"
// 		input.CNPJ = "10.123.456/0001"
// 		_, err := NewCompany(input)
// 		assert.Error(t, err)
// 		assert.Equal(t, exception.ErrInvalidCNPJ, err)
// 	})
//
// 	t.Run("should return false if Company is not deleted", func(t *testing.T) {
// 		input.CNPJ = "10.123.456/0001-00"
// 		Company, err := NewCompany(input)
// 		assert.NoError(t, err)
// 		assert.False(t, Company.IsDeleted())
// 	})
// }
