// Package di provides dependency injection container
package di

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/company"
	"github.com/paladignus/actajus/internal/module/identity"
	"github.com/paladignus/actajus/internal/module/notification"
	"github.com/paladignus/actajus/internal/module/person"
	"github.com/paladignus/actajus/internal/shared/application/repository"
	"github.com/paladignus/actajus/internal/shared/application/uow"
	"github.com/paladignus/actajus/internal/shared/infrastructure/config"
	sharedPostgres "github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
	website "github.com/paladignus/actajus/internal/shared/presentation/web/site"
	"github.com/paladignus/actajus/pkg/bootstrap"
)

// Container represents the dependency injection container
type Container struct {
	Config  config.Config
	DB      sharedPostgres.Executor
	Logger  repository.Logger
	UoW     uow.UnitOfWork
	RDB     any // Redis client - checked at runtime for health check
	Modules Modules

	runtime *bootstrap.Runtime
}

// Modules holds all initialized modules
type Modules struct {
	Company      company.Module
	Identity     identity.Module
	Notification notification.Module
	Person       person.Module
	Web          website.Module
}

// NewContainer creates a new DI container
func NewContainer() (*Container, error) {
	runtime, err := bootstrap.NewRuntime(context.Background())
	if err != nil {
		return nil, fmt.Errorf("shared runtime initialization failed: %w", err)
	}

	return &Container{
		Config: runtime.Config,
		DB:     runtime.DB,
		Logger: runtime.Logger,
		UoW:    runtime.UoW,
		RDB:    runtime.RDB,
		Modules: Modules{
			Company:  runtime.Modules.Company,
			Identity: runtime.Modules.Identity,
			Person:   runtime.Modules.Person,
			Web:      runtime.Modules.Web,
		},
		runtime: runtime,
	}, nil
}

// Close releases all resources
func (c *Container) Close() {
	if c.runtime != nil {
		c.runtime.Close()
		c.runtime = nil
	}
}
