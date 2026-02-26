// Package company
package company

import (
	"github.com/jackc/pgx/v5/pgxpool"
	address "github.com/paladignus/actajus/internal/module/address/application/mapper"
	addressDB "github.com/paladignus/actajus/internal/module/address/infrastructure/persistence/database/postgres"
	"github.com/paladignus/actajus/internal/module/company/application/mapper"
	"github.com/paladignus/actajus/internal/module/company/application/usecase"
	"github.com/paladignus/actajus/internal/module/company/infrastructure/persistence/database/postgres"
	"github.com/paladignus/actajus/internal/module/company/presentation/handler"
	email "github.com/paladignus/actajus/internal/module/email/application/mapper"
	emailDB "github.com/paladignus/actajus/internal/module/email/infrastructure/persistence/database/postgres"
	phone "github.com/paladignus/actajus/internal/module/phone/application/mapper"
	phoneDB "github.com/paladignus/actajus/internal/module/phone/infrastructure/persistence/database/postgres"
	socialMedia "github.com/paladignus/actajus/internal/module/social_media/application/mapper"
	socialMediaDB "github.com/paladignus/actajus/internal/module/social_media/infrastructure/persistence/database/postgres"
)

type Module struct {
	Handler handler.CompanyHandler
}

func NewModule(pool *pgxpool.Pool) Module {
	uow := postgres.NewCompanyUnitOfWork(pool)
	companyRepository := postgres.NewCompany(pool)
	companyReadRepository := postgres.NewCompanyReadRepository(pool)
	addressRepository := addressDB.NewAddress(pool)
	phoneRepository := phoneDB.NewPhone(pool)
	emailRepository := emailDB.NewEmail(pool)
	socialMediaRepository := socialMediaDB.NewSocialMedia(pool)
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
	updateUC := usecase.NewUpdateCompany(&uow, *mapper, projection)
	deleteUC := usecase.NewDeleteCompany(&uow)
	listUC := usecase.NewListCompanies(companyReadRepository)
	findByCNPJ := usecase.NewFindByCNPJ(
		companyRepository,
		addressRepository,
		phoneRepository,
		emailRepository,
		socialMediaRepository,
		projection,
	)
	findByID := usecase.NewFindByID(
		companyRepository,
		addressRepository,
		phoneRepository,
		emailRepository,
		socialMediaRepository,
		projection,
	)
	handler := handler.NewCompanyHandler(
		createUC,
		updateUC,
		deleteUC,
		listUC,
		findByCNPJ,
		findByID,
	)
	return Module{handler}
}
