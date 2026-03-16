// Package main
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"github.com/paladignus/actajus/internal/module/identity"
	identitypg "github.com/paladignus/actajus/internal/module/identity/infrastructure/persistence/database/postgres"
	"github.com/paladignus/actajus/internal/module/notification"
	"github.com/paladignus/actajus/internal/module/notification/infrastructure/email"
	"github.com/paladignus/actajus/internal/module/person"
	"github.com/paladignus/actajus/internal/shared/infrastructure/config"
	"github.com/paladignus/actajus/internal/shared/infrastructure/logger"
	sharednats "github.com/paladignus/actajus/internal/shared/infrastructure/messaging/nats"
	sharedoutbox "github.com/paladignus/actajus/internal/shared/infrastructure/messaging/outbox"
	postgresShared "github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
	sharedserialization "github.com/paladignus/actajus/internal/shared/infrastructure/serialization"
	shareduuid "github.com/paladignus/actajus/internal/shared/infrastructure/service"
	"github.com/paladignus/actajus/internal/shared/presentation/interceptor"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	appLogger := logger.NewDefaultLogger()
	db, err := postgresShared.NewConnection(ctx, &cfg.Database)
	if err != nil {
		appLogger.Error(ctx, "❌ error initializing the database connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	appLogger.Info(ctx, "✅ Database connected successfully")
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       0,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		appLogger.Error(ctx, "❌ redis ping failed", "error", err)
		os.Exit(1)
	}
	defer rdb.Close()
	appLogger.Info(ctx, "✅ Cache connected successfully")
	serializer := sharedserialization.NewJSONSerializer()
	idGen := shareduuid.NewUUIDGenerator()
	uow := postgresShared.NewUnitOfWork(db)
	identityFactory := identitypg.NewFactory(db)
	outboxFactory := postgresShared.NewOutboxFactory(db)
	natsBoot, err := sharednats.NewBootstrap(cfg.NATS.URL)
	if err != nil {
		appLogger.Error(ctx, "❌ nats connect failed", "error", err)
		os.Exit(1)
	}
	defer natsBoot.Conn.Close()
	if err := natsBoot.EnsureNotificationStream(ctx); err != nil {
		appLogger.Error(ctx, "❌ ensure notification stream failed", "error", err)
		os.Exit(1)
	}
	emailConsumer, err := natsBoot.EnsureEmailWorkerConsumer(ctx)
	if err != nil {
		appLogger.Error(ctx, "❌ ensure email consumer failed", "error", err)
		os.Exit(1)
	}
	publisher := sharednats.NewJetStreamPublisher(natsBoot.JS)
	dispatcher := sharedoutbox.NewDispatcher(db, publisher, 200)
	go func() {
		if err := dispatcher.Run(ctx, 2*time.Second); err != nil && err != context.Canceled {
			appLogger.Error(context.Background(), "outbox dispatcher stopped", "error", err)
		}
	}()
	smtpPort, err := strconv.Atoi(cfg.SMTP.Port)
	if err != nil {
		appLogger.Error(ctx, "❌ invalid SMTP_PORT", "value", cfg.SMTP.Port, "error", err)
		os.Exit(1)
	}
	emailSender := email.NewSMTPSender(
		cfg.SMTP.Host,
		smtpPort,
		cfg.SMTP.User,
		cfg.SMTP.Pass,
		cfg.SMTP.From,
	)
	notificationMod, err := notification.NewModule(notification.Dependencies{
		Logger:        appLogger,
		DB:            db,
		Serializer:    serializer,
		Consumer:      emailConsumer,
		EmailSender:   emailSender,
		PublicBaseURL: cfg.SMTP.PublicBaseURL,
	})
	if err != nil {
		appLogger.Error(ctx, "❌ failed to init notification module", "error", err)
		os.Exit(1)
	}
	notificationMod.Start(ctx)
	mux := http.NewServeMux()
	personMod := person.NewModule(db, appLogger)
	// userRepo := identitypg.NewUser(db) // mantém por compatibilidade com Login/Refresh
	identityMod, err := identity.NewModule(identity.Dependencies{
		Logger:        appLogger,
		DB:            db,
		RDB:           rdb,
		Config:        cfg.Auth,
		UoW:           uow,
		Repository:    identityFactory,
		OutboxFactory: outboxFactory,
		Serializer:    serializer,
		IDGenerator:   idGen,
	})
	if err != nil {
		appLogger.Error(ctx, "❌ failed to init identity module", "error", err)
		os.Exit(1)
	}
	recoverI := interceptor.NewRecoverInterceptor(appLogger)
	loggingI := interceptor.NewLoggingInterceptor(appLogger)
	authI := interceptor.NewAuthInterceptor(
		identityMod.ValidateAccess,
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
	rbacI := interceptor.NewRBACInterceptor(identityMod.RBACChecker, rbacRules)
	opts := connect.WithInterceptors(recoverI, loggingI, authI, rbacI)
	identityMod.Mount(mux, opts)
	personMod.Mount(mux, opts)
	appLogger.Info(ctx, "✅ Modules initialized successfully")
	handler := corsMiddleware(mux)
	srv := &http.Server{
		Addr:              cfg.Server.GRPCPort,
		Handler:           h2c.NewHandler(handler, &http2.Server{}),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	go func() {
		appLogger.Info(ctx, "🚀 Server starting", "addr", cfg.Server.GRPCPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error(ctx, "❌ server error", "error", err)
			os.Exit(1)
		}
	}()
	<-ctx.Done()
	appLogger.Info(context.Background(), "🛑 Server is shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		appLogger.Error(context.Background(), "❌ server forced to shutdown", "error", err)
		os.Exit(1)
	}
	appLogger.Info(context.Background(), "✅ Server stopped")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Connect-Protocol-Version")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Expose-Headers", "Connect-Protocol-Version")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// func main() {
// 	cfg := config.Load()
// 	ctx := context.Background()
// 	appLogger := logger.NewDefaultLogger()
// 	db, err := postgresShared.NewConnection(ctx, &cfg.Database)
// 	if err != nil {
// 		appLogger.Error(ctx, "❌ error initializing the database connection", "error", err)
// 		os.Exit(1)
// 	}
// 	defer db.Close()
// 	appLogger.Info(ctx, "✅ Database connected successfully")
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     cfg.Redis.Addr,
// 		Password: cfg.Redis.Password,
// 		DB:       0,
// 	})
// 	if err := rdb.Ping(ctx).Err(); err != nil {
// 		appLogger.Error(ctx, "❌ redis ping failed", "error", err)
// 		os.Exit(1)
// 	}
// 	defer rdb.Close()
// 	appLogger.Info(ctx, "✅ Cache connected successfully")
//
// 	mux := http.NewServeMux()
// 	personMod := person.NewModule(db, appLogger)
// 	// userMod := user.NewModule(db)
// 	userRepo := postgres.NewUser(db)
// 	identityMod, err := identity.NewModule(identity.Dependencies{
// 		Logger: appLogger,
// 		DB:     db,
// 		RDB:    rdb,
// 		Config: cfg.Auth,
// 		// SkipRBACRebuild: true, // Descomente para desabilitar o rebuild automático
// 	})
// 	if err != nil {
// 		appLogger.Error(ctx, "❌ failed to init identity module", "error", err)
// 		os.Exit(1)
// 	}
// 	// O rebuild dos índices RBAC é feito automaticamente dentro de identity.NewModule
// 	recoverI := interceptor.NewRecoverInterceptor(appLogger)
// 	loggingI := interceptor.NewLoggingInterceptor(appLogger)
// 	authI := interceptor.NewAuthInterceptor(
// 		identityMod.ValidateAccess,
// 		interceptor.WithWhitelistProcedures(
// 			"/identity.v1.AuthService/Login",
// 			"/identity.v1.AuthService/Refresh",
// 			"/identity.v1.AuthService/RequestPasswordReset",
// 			"/identity.v1.AuthService/ConfirmPasswordReset",
// 			"/grpc.health.v1.Health/Check",
// 		),
// 	)
// 	rbacRules := map[string]string{
// 		"/identity.v1.RbacAdminService/AssignRoleToUser":         "rbac:admin",
// 		"/identity.v1.RbacAdminService/RemoveRoleFromUser":       "rbac:admin",
// 		"/identity.v1.RbacAdminService/GrantPermissionToRole":    "rbac:admin",
// 		"/identity.v1.RbacAdminService/RevokePermissionFromRole": "rbac:admin",
// 		"/person.v1.PersonService/CreatePerson":                  "person:create",
// 		"/person.v1.PersonService/DeletePerson":                  "person:delete",
// 	}
// 	rbacI := interceptor.NewRBACInterceptor(identityMod.RBACChecker, rbacRules)
// 	opts := connect.WithInterceptors(
// 		recoverI,
// 		loggingI,
// 		authI,
// 		rbacI,
// 	)
// 	identityMod.Mount(mux, opts)
// 	personMod.Mount(mux, opts)
// 	// personPath, personHandler := personMod.Route(opts)
// 	// mux.Handle(personPath, personHandler)
// 	appLogger.Info(ctx, "✅ Modules initialized successfully")
// 	handler := corsMiddleware(mux)
// 	srv := &http.Server{
// 		Addr:              cfg.Server.GRPCPort,
// 		Handler:           h2c.NewHandler(handler, &http2.Server{}),
// 		ReadHeaderTimeout: 5 * time.Second,
// 		ReadTimeout:       10 * time.Second,
// 		WriteTimeout:      10 * time.Second,
// 		IdleTimeout:       120 * time.Second,
// 	}
// 	go func() {
// 		appLogger.Info(ctx, "🚀 Server starting", "addr", cfg.Server.GRPCPort)
// 		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			appLogger.Error(ctx, "❌ server error", "error", err)
// 			os.Exit(1)
// 			// log.Fatalf("❌ server error: %v", err)
// 		}
// 	}()
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 	<-quit
// 	appLogger.Info(ctx, "🛑 Server is shutting down...")
// 	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	if err := srv.Shutdown(shutdownCtx); err != nil {
// 		appLogger.Error(ctx, "❌ server forced to shutdown", "error", err)
// 		os.Exit(1)
// 	}
// 	appLogger.Info(ctx, "✅ Server stopped")
// }
//
// func corsMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Connect-Protocol-Version")
// 		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
// 		w.Header().Set("Access-Control-Expose-Headers", "Connect-Protocol-Version")
// 		if r.Method == http.MethodOptions {
// 			w.WriteHeader(http.StatusNoContent)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }
