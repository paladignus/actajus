// Package company
package company

import (
	"net/http"

	"connectrpc.com/connect"
	address "github.com/paladignus/actajus/internal/module/address/application/mapper"
	"github.com/paladignus/actajus/internal/module/company/application/mapper"
	"github.com/paladignus/actajus/internal/module/company/application/repository"
	"github.com/paladignus/actajus/internal/module/company/application/usecase"
	"github.com/paladignus/actajus/internal/module/company/presentation/grpc/handler"
	email "github.com/paladignus/actajus/internal/module/email/application/mapper"
	phone "github.com/paladignus/actajus/internal/module/phone/application/mapper"
	socialMedia "github.com/paladignus/actajus/internal/module/social_media/application/mapper"
	sharedRepo "github.com/paladignus/actajus/internal/shared/application/repository"
	"github.com/paladignus/actajus/internal/shared/application/uow"
	sharedPostgres "github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
	"github.com/paladignus/actajus/proto/company/v1/companyv1connect"
)

type Dependencies struct {
	DB         sharedPostgres.Executor
	Logger     sharedRepo.Logger
	UoW        uow.UnitOfWork
	Repository repository.Factory
}

type Module struct {
	handlerImpl *handler.CompanyHandler
	db          sharedPostgres.Executor
	logger      sharedRepo.Logger
}

func NewModule(d Dependencies) (Module, error) {
	addrProject := address.NewAddressProjectionMapper()
	phoneProject := phone.NewPhoneProjectionMapper()
	emailProject := email.NewEmailProjectionMapper()
	socialMediaProject := socialMedia.NewSocialMediaProjectionMapper()
	projection := mapper.NewCompanyProjectionMapper(
		addrProject,
		phoneProject,
		emailProject,
		socialMediaProject,
	)
	address := address.NewAddressMapper()
	phone := phone.NewPhoneMapper()
	email := email.NewEmailMapper()
	socialMedia := socialMedia.NewSocialMediaMapper()
	mapper := mapper.NewCompanyMapper(
		address,
		phone,
		email,
		socialMedia,
	)
	createUC := usecase.NewCreateCompany(
		d.UoW,
		d.Repository,
		*mapper,
		projection,
	)
	updateUC := usecase.NewUpdateCompany(
		d.UoW,
		d.Repository,
		*mapper,
		projection,
	)
	handlerImpl := handler.NewCompanyHandler(
		createUC,
		updateUC,
	)
	return Module{
		&handlerImpl,
		d.DB,
		d.Logger,
	}, nil
}

func (m Module) Mount(mux *http.ServeMux, opts ...connect.HandlerOption) {
	mux.Handle(companyv1connect.NewCompanyServiceHandler(m.handlerImpl, opts...))
}

// type Module struct {
// 	Handler handler.CompanyHandler
// }
//
// func NewModule(pool *pgxpool.Pool) Module {
// 	uow := sharedPostgres.NewUnitOfWork(pool)
// 	repository := postgres.NewFactory(pool)
// 	companyRepository := postgres.NewCompany(pool)
// 	companyReadRepository := postgres.NewCompanyReadRepository(pool)
// 	addressRepository := addressDB.NewAddress(pool)
// 	phoneRepository := phoneDB.NewPhone(pool)
// 	emailRepository := emailDB.NewEmail(pool)
// 	socialMediaRepository := socialMediaDB.NewSocialMedia(pool)
// 	addrProject := address.NewAddressProjectionMapper()
// 	address := address.NewAddressMapper()
// 	phoneProject := phone.NewPhoneProjectionMapper()
// 	phone := phone.NewPhoneMapper()
// 	emailProject := email.NewEmailProjectionMapper()
// 	email := email.NewEmailMapper()
// 	socialMediaProject := socialMedia.NewSocialMediaProjectionMapper()
// 	socialMedia := socialMedia.NewSocialMediaMapper()
// 	projection := mapper.NewCompanyProjectionMapper(
// 		addrProject,
// 		phoneProject,
// 		emailProject,
// 		socialMediaProject,
// 	)
// 	mapper := mapper.NewCompanyMapper(
// 		address,
// 		phone,
// 		email,
// 		socialMedia,
// 	)
// 	createUC := usecase.NewCreateCompany(uow, repository, *mapper, projection)
// 	updateUC := usecase.NewUpdateCompany(uow, repository, *mapper, projection)
// 	deleteUC := usecase.NewDeleteCompany(uow, repository)
// 	listUC := usecase.NewListCompanies(companyReadRepository)
// 	findByCNPJ := usecase.NewFindByCNPJ(
// 		companyRepository,
// 		addressRepository,
// 		phoneRepository,
// 		emailRepository,
// 		socialMediaRepository,
// 		projection,
// 	)
// 	findByID := usecase.NewFindByID(
// 		companyRepository,
// 		addressRepository,
// 		phoneRepository,
// 		emailRepository,
// 		socialMediaRepository,
// 		projection,
// 	)
// 	handler := handler.NewCompanyHandler(
// 		createUC,
// 		updateUC,
// 		deleteUC,
// 		listUC,
// 		findByCNPJ,
// 		findByID,
// 	)
// 	return Module{handler}
// }
