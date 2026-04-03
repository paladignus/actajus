// Package handler provides compatibility wrappers for the shared WEB handlers.
package handler

import sharedhandler "github.com/paladignus/actajus/internal/shared/presentation/web/app/handler"

type AppHandler = sharedhandler.AppHandler

var NewAppHandler = sharedhandler.NewAppHandler
