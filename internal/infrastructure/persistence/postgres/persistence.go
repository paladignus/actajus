// Package postgres
package postgres

type Persistence struct {
	db PgxPool
}

func NewPersistence(db PgxPool) Persistence {
	return Persistence{db}
}

func (p Persistence) User() User {
	return NewUser(p.db)
}

func (p Persistence) PasswordResetToken() PasswordResetToken {
	return NewPasswordResetToken(p.db)
}

func (p Persistence) Enterprise() Enterprise {
	return *NewEnterprise(p.db)
}

func (p Persistence) Address() Address {
	return NewAddress(p.db)
}
