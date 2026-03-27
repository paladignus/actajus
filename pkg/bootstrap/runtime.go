// Package bootstrap provides shared application bootstrap for entrypoints.
package bootstrap

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paladignus/actajus/internal/module/company"
	companypg "github.com/paladignus/actajus/internal/module/company/infrastructure/persistence/database/postgres"
	"github.com/paladignus/actajus/internal/module/identity"
	identitypg "github.com/paladignus/actajus/internal/module/identity/infrastructure/persistence/database/postgres"
	"github.com/paladignus/actajus/internal/module/person"
	"github.com/paladignus/actajus/internal/shared/application/messaging"
	"github.com/paladignus/actajus/internal/shared/application/repository"
	sharedservice "github.com/paladignus/actajus/internal/shared/application/service"
	"github.com/paladignus/actajus/internal/shared/application/uow"
	"github.com/paladignus/actajus/internal/shared/infrastructure/config"
	"github.com/paladignus/actajus/internal/shared/infrastructure/logger"
	sharedpostgres "github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
	sharedserialization "github.com/paladignus/actajus/internal/shared/infrastructure/serialization"
	shareduuid "github.com/paladignus/actajus/internal/shared/infrastructure/service"
	"github.com/redis/go-redis/v9"
)

// Runtime holds the shared infrastructure and modules used by multiple entrypoints.
type Runtime struct {
	Config  config.Config
	Logger  repository.Logger
	DB      *pgxpool.Pool
	RDB     *redis.Client
	UoW     uow.UnitOfWork
	Outbox  messaging.OutboxFactory
	Codec   sharedservice.MessageSerializer
	IDGen   sharedservice.IDGenerator
	Modules Modules
}

// Modules holds the application modules initialized by the shared bootstrap.
type Modules struct {
	Company  company.Module
	Identity identity.Module
	Person   person.Module
}

// NewRuntime loads config and initializes the shared runtime.
func NewRuntime(ctx context.Context) (*Runtime, error) {
	return NewRuntimeWithConfig(ctx, config.Load())
}

// NewRuntimeWithConfig initializes the shared runtime using the provided config.
func NewRuntimeWithConfig(ctx context.Context, cfg config.Config) (*Runtime, error) {
	log := logger.NewDefaultLogger()

	db, err := sharedpostgres.NewConnection(ctx, &cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("initialize database connection: %w", err)
	}

	rdb, err := setupRedis(ctx, cfg.Redis)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("initialize redis connection: %w", err)
	}

	uow := sharedpostgres.NewUnitOfWork(db)
	outbox := sharedpostgres.NewOutboxFactory(db)
	codec := sharedserialization.NewJSONSerializer()
	idGen := shareduuid.NewUUIDGenerator()
	modules, err := initializeModules(db, log, uow, rdb, cfg, outbox, codec, idGen)
	if err != nil {
		_ = rdb.Close()
		db.Close()
		return nil, fmt.Errorf("initialize shared modules: %w", err)
	}

	return &Runtime{
		Config:  cfg,
		Logger:  log,
		DB:      db,
		RDB:     rdb,
		UoW:     uow,
		Outbox:  outbox,
		Codec:   codec,
		IDGen:   idGen,
		Modules: modules,
	}, nil
}

// Close releases the shared resources owned by the runtime.
func (r *Runtime) Close() {
	if r.RDB != nil {
		_ = r.RDB.Close()
		r.RDB = nil
	}
	if r.DB != nil {
		r.DB.Close()
		r.DB = nil
	}
}

func initializeModules(
	db *pgxpool.Pool,
	log repository.Logger,
	unitOfWork uow.UnitOfWork,
	rdb *redis.Client,
	cfg config.Config,
	outbox messaging.OutboxFactory,
	codec sharedservice.MessageSerializer,
	idGen sharedservice.IDGenerator,
) (Modules, error) {
	companyFactory := companypg.NewFactory(db)
	companyRead := companypg.NewCompanyReadRepository(db)
	companyModule, err := company.NewModule(company.Dependencies{
		DB:             db,
		Logger:         log,
		UoW:            unitOfWork,
		Repository:     companyFactory,
		ReadRepository: companyRead,
	})
	if err != nil {
		return Modules{}, fmt.Errorf("initialize company module: %w", err)
	}

	identityFactory := identitypg.NewFactory(db)
	identityModule, err := identity.NewModule(identity.Dependencies{
		Logger:        log,
		DB:            db,
		RDB:           rdb,
		Config:        cfg.Auth,
		UoW:           unitOfWork,
		Repository:    identityFactory,
		OutboxFactory: outbox,
		Serializer:    codec,
		IDGenerator:   idGen,
	})
	if err != nil {
		return Modules{}, fmt.Errorf("initialize identity module: %w", err)
	}

	return Modules{
		Company:  companyModule,
		Identity: identityModule,
		Person:   person.NewModule(db, log),
	}, nil
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
