package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/paladignus/actajus/internal/shared/infrastructure/config"
	"github.com/paladignus/actajus/internal/shared/presentation/interceptor"
	"github.com/paladignus/actajus/pkg/bootstrap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type app struct {
	runtime      *bootstrap.Runtime
	notification *bootstrap.NotificationRuntime
	httpServer   *http.Server
}

func newApp(ctx context.Context, cfg config.Config) (*app, error) {
	runtime, err := bootstrap.NewRuntimeWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("initialize shared runtime: %w", err)
	}

	notificationRuntime, err := bootstrap.NewNotificationRuntime(ctx, runtime)
	if err != nil {
		runtime.Close()
		return nil, fmt.Errorf("initialize notification runtime: %w", err)
	}

	mux := http.NewServeMux()
	recoverI := interceptor.NewRecoverInterceptor(runtime.Logger)
	loggingI := interceptor.NewLoggingInterceptor(runtime.Logger)
	authI := interceptor.NewAuthInterceptor(
		runtime.Modules.Identity.ValidateAccess,
		interceptor.WithWhitelistProcedures(
			"/identity.v1.AuthService/Login",
			"/identity.v1.AuthService/Refresh",
			"/identity.v1.AuthService/RequestPasswordReset",
			"/identity.v1.AuthService/ConfirmPasswordReset",
			"/grpc.health.v1.Health/Check",
		),
	)
	rbacRules := map[string]string{
		"/identity.v1.RbacAdminService/AssignRoleToUser":         "rbac:admin",
		"/identity.v1.RbacAdminService/RemoveRoleFromUser":       "rbac:admin",
		"/identity.v1.RbacAdminService/GrantPermissionToRole":    "rbac:admin",
		"/identity.v1.RbacAdminService/RevokePermissionFromRole": "rbac:admin",
		"/person.v1.PersonService/CreatePerson":                  "person:create",
		"/person.v1.PersonService/DeletePerson":                  "person:delete",
	}
	rbacI := interceptor.NewRBACInterceptor(runtime.Modules.Identity.RBACChecker, rbacRules)
	opts := connect.WithInterceptors(recoverI, loggingI, authI, rbacI)
	runtime.Modules.Identity.Mount(mux, opts)
	runtime.Modules.Person.Mount(mux, opts)
	runtime.Modules.Company.Mount(mux, opts)

	httpServer := &http.Server{
		Addr:              cfg.Server.GRPCPort,
		Handler:           h2c.NewHandler(corsMiddleware(mux), &http2.Server{}),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	return &app{
		runtime:      runtime,
		notification: notificationRuntime,
		httpServer:   httpServer,
	}, nil
}

func (a *app) run(ctx context.Context) error {
	a.runtime.Logger.Info(ctx, "Database connected successfully")
	a.runtime.Logger.Info(ctx, "Cache connected successfully")
	a.runtime.Logger.Info(ctx, "Modules initialized successfully")

	go func() {
		a.runtime.Logger.Info(ctx, "Server starting", "addr", a.httpServer.Addr)
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.runtime.Logger.Error(context.Background(), "server error", "error", err)
		}
	}()

	<-ctx.Done()
	a.runtime.Logger.Info(context.Background(), "Server is shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}
	a.runtime.Logger.Info(context.Background(), "Server stopped")
	return nil
}

func (a *app) close() {
	if a.notification != nil {
		a.notification.Close()
		a.notification = nil
	}
	if a.runtime != nil {
		a.runtime.Close()
	}
}
