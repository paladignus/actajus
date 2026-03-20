// Package di provides dependency injection container
package di

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/address/application/mapper"
	"github.com/paladignus/actajus/internal/module/company"
	companyMapper "github.com/paladignus/actajus/internal/module/company/application/mapper"
	companyRepo "github.com/paladignus/actajus/internal/module/company/application/repository"
	companypg "github.com/paladignus/actajus/internal/module/company/infrastructure/persistence/database/postgres"
	emailMapper "github.com/paladignus/actajus/internal/module/email/application/mapper"
	"github.com/paladignus/actajus/internal/module/identity"
	identitypg "github.com/paladignus/actajus/internal/module/identity/infrastructure/persistence/database/postgres"
	"github.com/paladignus/actajus/internal/module/notification"
	"github.com/paladignus/actajus/internal/module/person"
	phoneMapper "github.com/paladignus/actajus/internal/module/phone/application/mapper"
	socialMediaMapper "github.com/paladignus/actajus/internal/module/social_media/application/mapper"
	"github.com/paladignus/actajus/internal/shared/application/repository"
	"github.com/paladignus/actajus/internal/shared/application/uow"
	"github.com/paladignus/actajus/internal/shared/infrastructure/config"
	sharedPostgres "github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
	"github.com/paladignus/actajus/internal/shared/infrastructure/logger"
	"github.com/redis/go-redis/v9"
)

// Container represents the dependency injection container
type Container struct {
	Config  config.Config
	DB      sharedPostgres.Executor
	Logger  repository.Logger
	UoW     uow.UnitOfWork
	Modules Modules
}

// Modules holds all initialized modules
type Modules struct {
	Company      company.Module
	Identity     identity.Module
	Notification notification.Module
	Person       person.Module
}

// NewContainer creates a new DI container
func NewContainer() (*Container, error) {
	cfg := config.Load()
	ctx := context.Background()
	log := logger.NewDefaultLogger()

	db, err := sharedPostgres.NewConnection(ctx, &cfg.Database)
	if err != nil {
		log.Error(ctx, "error initializing the database connection.", "error", err)
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	uow := sharedPostgres.NewUnitOfWork(db)

	// Setup Redis
	rdb, err := setupRedis(ctx, cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	modules, err := initializeModules(db, log, uow, rdb, cfg)
	if err != nil {
		rdb.Close()
		return nil, fmt.Errorf("module initialization failed: %w", err)
	}

	return &Container{
		Config:  cfg,
		DB:      db,
		Logger:  log,
		UoW:     uow,
		Modules: modules,
	}, nil
}

// Close releases all resources
func (c *Container) Close() {
	// DB close is handled by the database package if needed
}

func initializeModules(db sharedPostgres.Executor, log repository.Logger, uow uow.UnitOfWork, rdb *redis.Client, cfg config.Config) (Modules, error) {
	var modules Modules
	var err error

	// Company module
	modules.Company, err = initializeCompanyModule(db, log, uow)
	if err != nil {
		return modules, fmt.Errorf("company module initialization failed: %w", err)
	}

	// Identity module
	modules.Identity, err = initializeIdentityModule(db, log, uow, rdb, cfg)
	if err != nil {
		return modules, fmt.Errorf("identity module initialization failed: %w", err)
	}

	// Person module
	modules.Person = initializePersonModule(db, log)

	// Notification module will be initialized separately due to complex dependencies

	return modules, nil
}

func initializeCompanyModule(db sharedPostgres.Executor, log repository.Logger, uow uow.UnitOfWork) (company.Module, error) {
	companyFactory := companypg.NewFactory(db)
	companyRead := companypg.NewCompanyReadRepository(db)

	return company.NewModule(company.Dependencies{
		DB:             db,
		Logger:         log,
		UoW:            uow,
		Repository:     companyFactory,
		ReadRepository: companyRead,
	})
}

// CompanyModuleDeps holds company module dependencies for usecases
type CompanyModuleDeps struct {
	UoW            uow.UnitOfWork
	Repository     companyRepo.Factory
	ReadRepository companyRepo.CompanyReadRepository
	Mapper         *companyMapper.CompanyMapper
	Projection     companyMapper.CompanyProjectionMapper
}

// NewCompanyModuleDeps creates company module dependencies
func NewCompanyModuleDeps(uow uow.UnitOfWork, repo companyRepo.Factory, readRepo companyRepo.CompanyReadRepository) CompanyModuleDeps {
	addrProject := mapper.NewAddressProjectionMapper()
	phoneProject := phoneMapper.NewPhoneProjectionMapper()
	emailProject := emailMapper.NewEmailProjectionMapper()
	socialMediaProject := socialMediaMapper.NewSocialMediaProjectionMapper()
	projection := companyMapper.NewCompanyProjectionMapper(
		addrProject,
		phoneProject,
		emailProject,
		socialMediaProject,
	)
	addr := mapper.NewAddressMapper()
	phone := phoneMapper.NewPhoneMapper()
	email := emailMapper.NewEmailMapper()
	socialMedia := socialMediaMapper.NewSocialMediaMapper()
	mapper := companyMapper.NewCompanyMapper(
		addr,
		phone,
		email,
		socialMedia,
	)
	return CompanyModuleDeps{
		UoW:            uow,
		Repository:     repo,
		ReadRepository: readRepo,
		Mapper:         mapper,
		Projection:     projection,
	}
}

func setupRedis(ctx context.Context, cfg config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}

func initializeIdentityModule(db sharedPostgres.Executor, log repository.Logger, uow uow.UnitOfWork, rdb *redis.Client, cfg config.Config) (identity.Module, error) {
	identityFactory := identitypg.NewFactory(db)
	
	return identity.NewModule(identity.Dependencies{
		Logger:     log,
		DB:         db,
		RDB:        rdb,
		Config:     cfg.Auth,
		UoW:        uow,
		Repository: identityFactory,
	})
}

func initializePersonModule(db sharedPostgres.Executor, log repository.Logger) person.Module {
	return person.NewModule(db, log)
}
