// Package company
package company

import (
	"github.com/jackc/pgx/v5/pgxpool"
	address "github.com/paladignus/actajus/internal/module/address/application/mapper"
	"github.com/paladignus/actajus/internal/module/company/application/mapper"
	"github.com/paladignus/actajus/internal/module/company/application/usecase"
	"github.com/paladignus/actajus/internal/module/company/infrastructure/persistence/database"
	"github.com/paladignus/actajus/internal/module/company/presentation/handler"
	email "github.com/paladignus/actajus/internal/module/email/application/mapper"
	phone "github.com/paladignus/actajus/internal/module/phone/application/mapper"
	socialMedia "github.com/paladignus/actajus/internal/module/social_media/application/mapper"
)

type Module struct {
	Handler handler.CompanyHandler
}

func NewModule(pool *pgxpool.Pool) Module {
	uow := database.NewCompanyUnitOfWork(pool)
	addrProject := address.NewAddressProjectionMapper()
	address := address.NewAddressMapper()
	phoneProject := phone.NewPhoneProjectionMapper()
	phone := phone.NewPhoneMapper()
	emailProject := email.NewEmailProjectionMapper()
	email := email.NewEmailMapper()
	socialMediaProject := socialMedia.NewSocialMediaProjectionMapper()
	socialMedia := socialMedia.NewSocialMediaMapper()
	projection := mapper.NewCompanyProjectionMapper(
		addrProject,
		phoneProject,
		emailProject,
		socialMediaProject,
	)
	mapper := mapper.NewCompanyMapper(
		address,
		phone,
		email,
		socialMedia,
	)
	createUC := usecase.NewCreateCompany(&uow, *mapper, projection)
	handler := handler.NewCompanyHandler(createUC)
	return Module{handler}
}
