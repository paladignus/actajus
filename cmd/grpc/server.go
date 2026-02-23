package main

// Package main
func main() {
	// authI := interceptor.NewAuthInterceptor(
	// 	validateAccessUsecase,
	// 	interceptor.WithWhitelist(
	// 		"/identity.v1.AuthService/Login",
	// 		"/identity.v1.AuthService/Refresh",
	// 		"/grpc.health.v1.Health/Check",
	// 	),
	// )
	//
	// mux := http.NewServeMux()
	// mux.Handle(identityv1connect.NewAuthServiceHandler(authHandler, connect.WithInterceptors(authI)))
}

// Dica: o valor exato de Procedure depende do serviço gerado.
// Se você logar req.Spec().Procedure, você pega os nomes certinhos para whitelist.

// authI := interceptor.NewAuthInterceptor(
// 	validateAccessUsecase,
// 	interceptor.WithLogger(logger),
// 	interceptor.WithProcedureLogging(true),
//
// 	// whitelist por prefixo é MUITO mais fácil de manter:
// 	interceptor.WithWhitelistPrefixes(
// 		"/identity.v1.AuthService/",     // login/refresh/logout/logout-all etc (se quiser liberar só alguns, use procedures)
// 		"/grpc.health.v1.Health/",       // health checks
// 	),
//
// 	// ou whitelist pontual:
// 	interceptor.WithWhitelistProcedures(
// 		"/identity.v1.AuthService/Login",
// 		"/identity.v1.AuthService/Refresh",
// 	),
// )

// Dica: Eu recomendo prefixo só para health e procedures específicos para auth.
// Ex: você normalmente quer proteger /identity.v1.AuthService/Logout e
// /LogoutAll com auth, então não libere o prefixo inteiro do AuthService.
