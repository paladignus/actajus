// Package unitofwork
package unitofwork

import "context"

type UnitOfWork interface {
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// type UnitOfWorkCompany interface {
// 	UnitOfWork
// 	Company() ICompany
// 	Address() IAddress
// 	CompanyAddress() ICompanyAddress
// 	Phone() IPhone
// 	CompanyPhone() ICompanyPhone
// 	Email() IEmail
// 	CompanyEmail() ICompanyEmail
// 	SocialMedia() ISocialMedia
// }

// type UnitOfWorkDefault interface {
// 	Company() UnitOfWorkCompany
// }
