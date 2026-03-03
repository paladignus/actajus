// Package main
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"github.com/paladignus/actajus/internal/module/identity"
	"github.com/paladignus/actajus/internal/module/identity/infrastructure/persistence/database/postgres"
	"github.com/paladignus/actajus/internal/module/person"

	"github.com/paladignus/actajus/internal/shared/infrastructure/config"
	"github.com/paladignus/actajus/internal/shared/infrastructure/logger"
	postgresShared "github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
	"github.com/paladignus/actajus/internal/shared/presentation/interceptor"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()
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
	mux := http.NewServeMux()
	personMod := person.NewModule(db, appLogger)
	// userMod := user.NewModule(db)
	userRepo := postgres.NewUser(db)
	identityMod, err := identity.NewModule(identity.Dependencies{
		Logger: appLogger,
		DB:     db,
		RDB:    rdb,
		Config: cfg.Auth,
		Users:  userRepo,
	})
	if err != nil {
		appLogger.Error(ctx, "❌ failed to init identity module", "error", err)
		os.Exit(1)
	}
	if err := identity.RebuildRBACIndexes(ctx, identity.Dependencies{
		DB:  db,
		RDB: rdb,
	}); err != nil {
		appLogger.Error(ctx, "⚠️ failed to rebuild RBAC indexes", "error", err)
	} else {
		appLogger.Info(ctx, "✅ RBAC indexes rebuilt successfully")
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
	opts := connect.WithInterceptors(
		recoverI,
		loggingI,
		authI,
		rbacI,
	)
	identityMod.Mount(mux, opts)
	personMod.Mount(mux, opts)
	// personPath, personHandler := personMod.Route(opts)
	// mux.Handle(personPath, personHandler)
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
			// log.Fatalf("❌ server error: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	appLogger.Info(ctx, "🛑 Server is shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		appLogger.Error(ctx, "❌ server forced to shutdown", "error", err)
		os.Exit(1)
	}
	appLogger.Info(ctx, "✅ Server stopped")
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

// // Package main
// package main

// import (
// 	"context"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"connectrpc.com/connect"
// 	"github.com/paladignus/actajus/internal/module/identity"
// 	"github.com/paladignus/actajus/internal/module/identity/application/mapper"
// 	"github.com/paladignus/actajus/internal/module/identity/application/usecase"
// 	"github.com/paladignus/actajus/internal/module/identity/infrastructure/persistence/database/postgres"
// 	"github.com/paladignus/actajus/internal/module/identity/infrastructure/security"
// 	"github.com/paladignus/actajus/internal/module/identity/presentation/grpc/handler"
// 	"github.com/paladignus/actajus/internal/module/identity/presentation/rbac"
// 	"github.com/paladignus/actajus/internal/module/person"
// 	"github.com/paladignus/actajus/internal/shared/infrastructure/config"
// 	"github.com/paladignus/actajus/internal/shared/infrastructure/logger"
// 	postgresShared "github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
// 	"github.com/paladignus/actajus/internal/shared/presentation/interceptor"
// 	"github.com/paladignus/actajus/internal/shared/presentation/validation"
// 	"github.com/paladignus/actajus/proto/identity/v1/identityv1connect"
// 	"github.com/redis/go-redis/v9"
// 	"golang.org/x/net/http2"
// 	"golang.org/x/net/http2/h2c"
// )

// func main() {
// 	config := config.Load()
// 	ctx := context.Background()
// 	logger := logger.NewDefaultLogger()
// 	db, err := postgresShared.NewConnection(ctx, &config.Database)
// 	if err != nil {
// 		logger.Error(ctx, "❌ error initializing the database connection.", "error", err)
// 		os.Exit(1)
// 	}
// 	defer db.Close()
// 	logger.Info(ctx, "✅ Database connected successfully")
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     config.Redis.Addr,
// 		Password: config.Redis.Password,
// 		DB:       0,
// 	})
// 	if err := rdb.Ping(ctx).Err(); err != nil {
// 		logger.Error(ctx, "redis ping failed", "err", err)
// 		os.Exit(1)
// 	}
// 	defer rdb.Close()
// 	logger.Info(ctx, "✅ Cache connected successfully")
// 	person := person.NewModule(db, logger)
// 	identity, err := identity.NewModule(db, logger, config.Auth, rdb)
// 	authzRepo := postgres.NewAuthorization(db)
// 	authzSvc := security.NewAuthorizationService(
// 		authzRepo,
// 		rdb,
// 		security.WithAuthzPrefix("rbac:"),     // opcional
// 		security.WithAuthzTTL(10*time.Minute), // ajuste conforme quiser
// 	)
// 	roleUserRepo := postgres.NewRoleUserAdminRepository(db)
// 	permRoleRepo := postgres.NewPermissionRoleAdminRepository(db)

// 	// RBAC cache invalidator (você já tem authzSvc)
// 	// cacheInvalidator := authzSvc // se ele já expõe InvalidateUser(ctx, IDUser)
// 	v := validation.New()
// 	rbacMapper := mapper.NewRBACAdminMapper(v)
// 	roleUsersIndex := security.NewRBACRoleUsersIndex(
// 		rdb,
// 		security.WithRoleUsersIndexPrefix("rbac:"),
// 	)
// 	pairs, err := postgres.ListAllRoleUsers(ctx, db)
// 	if err != nil {
// 		logger.Error(ctx, "failed to load role_user snapshot", "error", err)
// 	} else {
// 		rolePairs := make([]security.RoleUserPair, 0, len(pairs))
// 		for _, p := range pairs {
// 			rolePairs = append(rolePairs, security.RoleUserPair{
// 				IDRole: p.IDRole,
// 				IDUser: p.IDUser,
// 			})
// 		}
// 		if err := security.RebuildRoleUsersIndex(ctx, roleUsersIndex, rolePairs); err != nil {
// 			logger.Error(ctx, "failed to rebuild role_users index", "error", err)
// 		} else {
// 			logger.Info(ctx, "role_users index rebuilt successfully", "count", len(rolePairs))
// 		}
// 	}
// 	// usecases
// 	assignUC := usecase.NewAssignRoleToUser(roleUserRepo, authzSvc, roleUsersIndex, rbacMapper)
// 	removeUC := usecase.NewRemoveRoleFromUser(roleUserRepo, authzSvc, roleUsersIndex, rbacMapper)
// 	grantUC := usecase.NewGrantPermissionToRole(roleUserRepo, permRoleRepo, authzSvc, roleUsersIndex, rbacMapper)
// 	revokeUC := usecase.NewRevokePermissionFromRole(roleUserRepo, permRoleRepo, authzSvc, roleUsersIndex, rbacMapper)

// 	// 2) Handler connect
// 	rbacAdminHandler := handler.NewRbacAdminHandler(assignUC, removeUC, grantUC, revokeUC)
// 	rbacChecker := rbac.NewChecker(authzSvc)
// 	logger.Info(ctx, "✅ Modules initialized successfully")
// 	mux := http.NewServeMux()
// 	recoverI := interceptor.NewRecoverInterceptor(logger)
// 	loggingI := interceptor.NewLoggingInterceptor(logger)
// 	authI := interceptor.NewAuthInterceptor(
// 		identity.ValidateAccess,
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
// 		// "/identity.v1.AuthService/Refresh":                       "auth:refresh",
// 		"/person.v1.PersonService/CreatePerson": "person:create",
// 		"/person.v1.PersonService/DeletePerson": "person:delete",
// 	}
// 	rbacI := interceptor.NewRBACInterceptor(rbacChecker, rbacRules)
// 	interceptors := connect.WithInterceptors(
// 		recoverI,
// 		loggingI,
// 		authI,
// 		rbacI,
// 	)
// 	if err != nil {
// 		logger.Error(ctx, "❌ error initializing the access token service.", "error", err)
// 		os.Exit(1)
// 	}
// 	rbacPath, rbacHTTPHandler := identityv1connect.NewRbacAdminServiceHandler(rbacAdminHandler)
// 	mux.Handle(rbacPath, rbacHTTPHandler)
// 	mux.Handle(identity.Route(interceptors))
// 	mux.Handle(person.Route(interceptors))
// 	// mux.Handle(identity.RbacAdminRoute(interceptors))
// 	handler := corsMiddleware(mux)
// 	srv := &http.Server{
// 		Addr:              config.Server.GRPCPort,
// 		Handler:           h2c.NewHandler(handler, &http2.Server{}),
// 		ReadHeaderTimeout: 5 * time.Second,
// 		ReadTimeout:       10 * time.Second,
// 		WriteTimeout:      10 * time.Second,
// 		IdleTimeout:       120 * time.Second,
// 	}
// 	go func() {
// 		logger.Info(ctx, "🚀 Server starting on", "addr", config.Server.GRPCPort)
// 		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			log.Fatalf("❌ Server error: %v", err)
// 		}
// 	}()
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 	<-quit
// 	logger.Info(ctx, "❌ Server is shutting down...")
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	if err := srv.Shutdown(ctx); err != nil {
// 		logger.Error(ctx, "❌ Server forced to shutdown", "error", err)
// 		os.Exit(1)
// 	}
// 	logger.Info(ctx, "❌ Server stopped")
// }

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
