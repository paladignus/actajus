// Package web keeps backward-compatible wiring for the shared WEB presentation layer.
package web

import (
	"embed"
	"os"
	"strings"

	webapp "github.com/paladignus/actajus/internal/shared/presentation/web/app"
	webtemplate "github.com/paladignus/actajus/internal/shared/presentation/web/template"
)

//go:embed presentation/http/handler/content/templates/layouts/*.gohtml presentation/http/handler/content/templates/pages/*.gohtml presentation/http/handler/content/templates/partials/*.gohtml presentation/http/handler/content/templates/components/*.gohtml presentation/http/handler/content/dist/.vite/manifest.json presentation/http/handler/content/dist/assets/* presentation/http/handler/content/dist/*.svg
var contentFS embed.FS

type Dependencies = webapp.Dependencies

type Module = webapp.Module

func NewModule(dep Dependencies) (Module, error) {
	devServerURL := strings.TrimSpace(os.Getenv("ACTAJUS_WEB_DEV_SERVER_URL"))
	if devServerURL == "" && strings.EqualFold(strings.TrimSpace(os.Getenv("ACTAJUS_WEB_ENV")), "development") {
		devServerURL = "http://localhost:5173"
	}
	dep.RendererSource = webtemplate.Source{
		TemplateFS: contentFS,
		TemplatePatterns: []string{
			"presentation/http/handler/content/templates/layouts/*.gohtml",
			"presentation/http/handler/content/templates/pages/*.gohtml",
			"presentation/http/handler/content/templates/partials/*.gohtml",
			"presentation/http/handler/content/templates/components/*.gohtml",
		},
		AssetFS:      contentFS,
		AssetRoot:    "presentation/http/handler/content/dist",
		ManifestPath: "presentation/http/handler/content/dist/.vite/manifest.json",
		PublicPath:   "/assets",
		DevServerURL: devServerURL,
	}
	return webapp.NewModule(dep)
}
