// Package company
package company

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/paladignus/actajus/internal/module/company/application/mapper"
	"github.com/paladignus/actajus/internal/module/company/application/repository"
	"github.com/paladignus/actajus/internal/module/company/application/usecase"
	"github.com/paladignus/actajus/internal/module/company/presentation/grpc/handler"
	sharedRepo "github.com/paladignus/actajus/internal/shared/application/repository"
	"github.com/paladignus/actajus/internal/shared/application/uow"
	sharedPostgres "github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
	"github.com/paladignus/actajus/proto/company/v1/companyv1connect"
)

type Dependencies struct {
	DB             sharedPostgres.Executor
	Logger         sharedRepo.Logger
	UoW            uow.UnitOfWork
	Repository     repository.Factory
	ReadRepository repository.CompanyReadRepository
}

type Module struct {
	Create      usecase.CreateCompany
	List        usecase.ListCompanies
	Update      usecase.UpdateCompany
	Delete      usecase.DeleteCompany
	FindByID    usecase.FindByID
	FindByCNPJ  usecase.FindByCNPJ
	handlerImpl *handler.CompanyHandler
	db          sharedPostgres.Executor
	logger      sharedRepo.Logger
}

func NewModule(d Dependencies) (Module, error) {
	projection := mapper.NewCompanyProjectionMapper()
	mapper := mapper.NewCompanyMapper()
	createUC := usecase.NewCreateCompany(
		d.UoW,
		d.Repository,
		*mapper,
		projection,
	)
	listUC := usecase.NewListCompanies(
		d.ReadRepository,
	)
	updateUC := usecase.NewUpdateCompany(
		d.UoW,
		d.Repository,
		*mapper,
		projection,
	)
	deleteUC := usecase.NewDeleteCompany(
		d.UoW,
		d.Repository,
	)
	findByID := usecase.NewFindByID(
		d.Repository,
		projection,
	)
	findByCNPJ := usecase.NewFindByCNPJ(
		d.Repository,
		projection,
	)
	handlerImpl := handler.NewCompanyHandler(
		createUC,
		listUC,
		updateUC,
		deleteUC,
		findByID,
		findByCNPJ,
	)
	return Module{
		Create:      createUC,
		List:        listUC,
		Update:      updateUC,
		Delete:      deleteUC,
		FindByID:    findByID,
		FindByCNPJ:  findByCNPJ,
		handlerImpl: &handlerImpl,
		db:          d.DB,
		logger:      d.Logger,
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
