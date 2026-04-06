// Package handler contains shared WEB HTTP handlers.
package handler

import (
	"net/http"

	sharedRepo "github.com/paladignus/actajus/internal/shared/application/repository"
	accessguard "github.com/paladignus/actajus/internal/shared/presentation/web/app/handler/access"
	supporthandler "github.com/paladignus/actajus/internal/shared/presentation/web/app/handler/support"
	webhtml "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/html"
	navigation "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/http/navigation"
)

type AppHandler struct {
	logger    sharedRepo.Logger
	assets    http.Handler
	html      *webhtml.Responder
	navigator *navigation.Navigator
	guard     accessguard.Guard
	support   supporthandler.Service
	auth      authDependencies
	identity  identityDependencies
	company   companyDependencies
}
